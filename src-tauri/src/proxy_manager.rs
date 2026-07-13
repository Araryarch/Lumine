use std::fs;
use std::path::PathBuf;
#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;

#[derive(Clone)]
pub struct ProxyManager {
    app_data_dir: PathBuf,
}

impl ProxyManager {
    pub fn new(app_data_dir: PathBuf) -> Self {
        let proxy_dir = app_data_dir.join("proxy");
        if !proxy_dir.exists() {
            let _ = fs::create_dir_all(&proxy_dir);
        }
        Self { app_data_dir }
    }

    pub fn ensure_proxy_running(&self) {
        let proxy_dir = self.app_data_dir.join("proxy");
        
        // Ensure there is at least an empty conf file so nginx doesn't fail
        let conf_path = proxy_dir.join("lumine.conf");
        if !conf_path.exists() {
            let _ = fs::write(&conf_path, "");
        }

        let mut cmd = crate::utils::create_command("docker");
        #[cfg(target_os = "windows")]
        cmd.creation_flags(0x08000000);
        
        // Check if running
        let output = cmd.args(["ps", "-q", "-f", "name=lumine_internal_proxy"]).output();
        if let Ok(out) = output {
            if !out.stdout.is_empty() {
                return; // Already running
            }
        }

        // Try to remove old container if exists but stopped
        let mut rm_cmd = crate::utils::create_command("docker");
        #[cfg(target_os = "windows")]
        rm_cmd.creation_flags(0x08000000);
        let _ = rm_cmd.args(["rm", "-f", "lumine_internal_proxy"]).output();

        // Start container
        let mut start_cmd = crate::utils::create_command("docker");
        #[cfg(target_os = "windows")]
        start_cmd.creation_flags(0x08000000);
        
        let mut bind_arg = proxy_dir.to_string_lossy().to_string();
        let mut certs_arg = self.app_data_dir.join("certs").to_string_lossy().to_string();
        // format for docker
        if cfg!(target_os = "windows") {
            bind_arg = bind_arg.replace("\\", "/");
            certs_arg = certs_arg.replace("\\", "/");
        }

        // Read document root
        let settings_file = self.app_data_dir.join("settings.json");
        let doc_root = if let Ok(content) = std::fs::read_to_string(&settings_file) {
            if let Ok(json) = serde_json::from_str::<serde_json::Value>(&content) {
                json["documentRoot"].as_str().unwrap_or("C:\\Lumine\\www").to_string()
            } else {
                "C:\\Lumine\\www".to_string()
            }
        } else {
            "C:\\Lumine\\www".to_string()
        };
        let doc_root_arg = if cfg!(target_os = "windows") {
            doc_root.replace("\\", "/")
        } else {
            doc_root
        };

        let _ = start_cmd.args([
            "run", "-d",
            "--name", "lumine_internal_proxy",
            "-p", "80:80",
            "-p", "443:443",
            // Nginx needs host.docker.internal to route back to host
            "--add-host", "host.docker.internal:host-gateway",
            "-v", &format!("{}:/etc/nginx/conf.d", bind_arg),
            "-v", &format!("{}:/etc/nginx/certs", certs_arg),
            "-v", &format!("{}:/var/www/html", doc_root_arg),
            "nginx:alpine"
        ]).output();
    }

    pub fn stop_proxy(&self) {
        let mut cmd = crate::utils::create_command("docker");
        #[cfg(target_os = "windows")]
        cmd.creation_flags(0x08000000);
        let _ = cmd.args(["rm", "-f", "lumine_internal_proxy"]).output();
    }

    pub fn sync_config(
        &self, 
        hosts: &Vec<crate::hosts_manager::HostEntry>, 
        services: &Vec<crate::service::ServiceInfo>,
        projects: &Vec<crate::project_manager::ProjectInfo>
    ) -> Result<(), String> {
        let proxy_dir = self.app_data_dir.join("proxy");
        let conf_path = proxy_dir.join("lumine.conf");
        
        let mut conf_content = String::new();

        // Read certs.json to check for HTTPS
        let certs_file = self.app_data_dir.join("certs.json");
        let mut cert_domains = std::collections::HashSet::new();
        if let Ok(content) = std::fs::read_to_string(&certs_file) {
            if let Ok(certs) = serde_json::from_str::<Vec<crate::mkcert_manager::CertEntry>>(&content) {
                for c in certs {
                    cert_domains.insert(c.domain);
                }
            }
        }
        
        for host in hosts {
            if host.is_enable {
                // Find matching service port OR project port
                let mut matched_port = None;
                if let Some(svc) = services.iter().find(|s| host.comment.contains(&s.name)) {
                    matched_port = Some(svc.config.port);
                } else if let Some(proj) = projects.iter().find(|p| host.comment.contains(&p.name)) {
                    matched_port = Some(proj.port);
                }

                if let Some(port) = matched_port {
                    if cert_domains.contains(&host.name) {
                        let safe_name = host.name.replace('*', "_wildcard_").replace('.', "_");
                        conf_content.push_str(&format!(
                            "server {{\n    listen 80;\n    server_name {0};\n    return 301 https://$host$request_uri;\n}}\nserver {{\n    listen 443 ssl;\n    server_name {0};\n    ssl_certificate /etc/nginx/certs/{1}.pem;\n    ssl_certificate_key /etc/nginx/certs/{1}-key.pem;\n    location / {{\n        proxy_pass http://host.docker.internal:{2};\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n        proxy_set_header X-Forwarded-Proto https;\n    }}\n}}\n",
                            host.name, safe_name, port
                        ));
                    } else {
                        conf_content.push_str(&format!(
                            "server {{\n    listen 80;\n    server_name {};\n    location / {{\n        proxy_pass http://host.docker.internal:{};\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n    }}\n}}\n",
                            host.name, port
                        ));
                    }
                }
            }
        }

        // Add localhost routing for Admin Panels
        let mut localhost_locations = String::new();
        for svc in services {
            if svc.config.service_type == "Admin Panel" {
                let port = svc.config.port;
                let path_name = match svc.name.to_lowercase().as_str() {
                    "phpmyadmin" => "phpmyadmin",
                    "adminer" => "adminer",
                    _ => continue,
                };
                
                localhost_locations.push_str(&format!(
                    "    location = /{} {{\n        return 301 /{}/;\n    }}\n    location ^~ /{}/ {{\n        proxy_pass http://host.docker.internal:{}/;\n        proxy_set_header Host $host;\n        proxy_set_header X-Real-IP $remote_addr;\n    }}\n",
                    path_name, path_name, path_name, port
                ));
            }
        }

        if !localhost_locations.is_empty() {
            conf_content.push_str(&format!(
                "server {{\n    listen 80;\n    server_name localhost 127.0.0.1;\n{}\n}}\n",
                localhost_locations
            ));
        }

        // Write config
        fs::write(&conf_path, &conf_content).map_err(|e| e.to_string())?;

        // Reload nginx
        let mut cmd = crate::utils::create_command("docker");
        #[cfg(target_os = "windows")]
        cmd.creation_flags(0x08000000);
        
        let _ = cmd.args(["exec", "lumine_internal_proxy", "nginx", "-s", "reload"]).output();

        Ok(())
    }
}
