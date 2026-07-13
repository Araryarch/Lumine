use serde::{Deserialize, Serialize};
use std::fs;
use std::path::{Path, PathBuf};
use tauri::AppHandle;
use tauri::Manager;
#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;

#[derive(Debug, Serialize, Deserialize, Clone)]
pub struct HostEntry {
    pub name: String,
    pub php: Option<String>,
    pub comment: String,
    pub is_enable: bool,
    pub is_ipv6: bool,
}

pub struct HostsManager {
    data_file: PathBuf,
}

impl HostsManager {
    pub fn new(app: &AppHandle) -> Self {
        let app_data_dir = app.path().app_data_dir().unwrap_or_else(|_| std::env::current_dir().unwrap());
        if !app_data_dir.exists() {
            let _ = fs::create_dir_all(&app_data_dir);
        }
        let data_file = app_data_dir.join("hosts.json");
        
        // Create initial default JSON if doesn't exist
        if !data_file.exists() {
            let defaults = vec![
                HostEntry {
                    name: "nginx.test".to_string(),
                    php: None,
                    comment: "Nginx (default port 80)".to_string(),
                    is_enable: false,
                    is_ipv6: false,
                },
                HostEntry {
                    name: "apache.test".to_string(),
                    php: None,
                    comment: "Apache httpd (requires port 8090)".to_string(),
                    is_enable: false,
                    is_ipv6: false,
                },
                HostEntry {
                    name: "caddy.test".to_string(),
                    php: None,
                    comment: "Caddy (requires port 8443)".to_string(),
                    is_enable: false,
                    is_ipv6: false,
                },
                HostEntry {
                    name: "phpmyadmin.test".to_string(),
                    php: None,
                    comment: "phpMyAdmin (requires port 8080)".to_string(),
                    is_enable: false,
                    is_ipv6: false,
                },
                HostEntry {
                    name: "adminer.test".to_string(),
                    php: None,
                    comment: "Adminer (requires port 8081)".to_string(),
                    is_enable: false,
                    is_ipv6: false,
                }
            ];
            let _ = fs::write(&data_file, serde_json::to_string_pretty(&defaults).unwrap());
        }

        Self { data_file }
    }

    pub fn get_hosts(&self) -> Result<Vec<HostEntry>, String> {
        let content = fs::read_to_string(&self.data_file)
            .map_err(|e| format!("Failed to read hosts data: {}", e))?;
        let hosts: Vec<HostEntry> = serde_json::from_str(&content)
            .unwrap_or_else(|_| Vec::new());
        Ok(hosts)
    }

    pub fn save_hosts(&self, hosts: &Vec<HostEntry>) -> Result<(), String> {
        let json = serde_json::to_string_pretty(hosts)
            .map_err(|e| format!("Failed to serialize hosts: {}", e))?;
        fs::write(&self.data_file, json)
            .map_err(|e| format!("Failed to save hosts data: {}", e))?;
        self.sync_system_hosts_file(hosts)?;
        Ok(())
    }

    pub fn add_host(&self, host: HostEntry) -> Result<(), String> {
        // Validate hostname to prevent injection (must be alphanumeric, dots, and hyphens only)
        if host.name.is_empty() || host.name.len() > 255 {
            return Err("Invalid hostname length".to_string());
        }
        if !host.name.chars().all(|c| c.is_ascii_alphanumeric() || c == '.' || c == '-') {
            return Err("Invalid hostname format".to_string());
        }

        let mut hosts = self.get_hosts()?;
        if hosts.iter().any(|h| h.name == host.name) {
            return Err("Host already exists".to_string());
        }

        hosts.push(host);
        self.save_hosts(&hosts)
    }

    pub fn edit_host(&self, old_name: &str, host: HostEntry) -> Result<(), String> {
        // Validate hostname to prevent injection
        if host.name.is_empty() || host.name.len() > 255 {
            return Err("Invalid hostname length".to_string());
        }
        if !host.name.chars().all(|c| c.is_ascii_alphanumeric() || c == '.' || c == '-') {
            return Err("Invalid hostname format".to_string());
        }

        let mut hosts = self.get_hosts()?;
        
        // If they are changing the name, ensure the new name doesn't already exist
        if old_name != host.name && hosts.iter().any(|h| h.name == host.name) {
            return Err("A host with the new name already exists".to_string());
        }

        if let Some(existing) = hosts.iter_mut().find(|h| h.name == old_name) {
            *existing = host;
            self.save_hosts(&hosts)
        } else {
            Err("Host not found".to_string())
        }
    }

    pub fn toggle_host(&self, name: &str, is_enable: bool) -> Result<(), String> {
        let mut hosts = self.get_hosts()?;
        if let Some(host) = hosts.iter_mut().find(|h| h.name == name) {
            host.is_enable = is_enable;
            self.save_hosts(&hosts)?;
            Ok(())
        } else {
            Err("Host not found".to_string())
        }
    }

    pub fn disable_all_hosts(&self) -> Result<(), String> {
        let mut hosts = self.get_hosts()?;
        let mut changed = false;
        for host in hosts.iter_mut() {
            if host.is_enable {
                host.is_enable = false;
                changed = true;
            }
        }
        if changed {
            self.save_hosts(&hosts)?;
        }
        Ok(())
    }
    pub fn delete_host(&self, name: &str) -> Result<(), String> {
        let mut hosts = self.get_hosts()?;
        hosts.retain(|h| h.name != name);
        self.save_hosts(&hosts)?;
        Ok(())
    }

    pub fn get_system_hosts_path() -> PathBuf {
        #[cfg(target_os = "windows")]
        {
            let windir = std::env::var("windir").unwrap_or_else(|_| "C:\\Windows".to_string());
            Path::new(&windir).join("System32\\drivers\\etc\\hosts")
        }
        #[cfg(not(target_os = "windows"))]
        {
            Path::new("/etc/hosts").to_path_buf()
        }
    }

    fn sync_system_hosts_file(&self, hosts: &Vec<HostEntry>) -> Result<(), String> {
        let hosts_path = Self::get_system_hosts_path();
        
        // Read current hosts file content
        let content = fs::read_to_string(&hosts_path)
            .map_err(|e| format!("Failed to read system hosts file: {}", e))?;

        let start_marker = "# --- Lumine Hosts Start ---";
        let end_marker = "# --- Lumine Hosts End ---";

        let mut new_content = content.clone();

        // Remove old block
        if let (Some(start_idx), Some(end_idx)) = (new_content.find(start_marker), new_content.find(end_marker)) {
            if end_idx > start_idx {
                let end_block_idx = end_idx + end_marker.len();
                let actual_end = new_content[end_block_idx..].find('\n')
                    .map(|i| end_block_idx + i + 1)
                    .unwrap_or(end_block_idx);
                new_content.replace_range(start_idx..actual_end, "");
            }
        }

        // Clean up trailing newlines
        new_content = new_content.trim_end().to_string();

        // Generate new block
        let mut new_block = String::new();
        new_block.push_str(&format!("\n\n{}\n", start_marker));
        
        for host in hosts.iter() {
            let prefix = if host.is_enable { "" } else { "# " };
            new_block.push_str(&format!("{}127.0.0.1\t{}\n", prefix, host.name));
            if host.is_ipv6 {
                new_block.push_str(&format!("{}::1\t\t{}\n", prefix, host.name));
            }
        }
        new_block.push_str(&format!("{}\n", end_marker));
        new_content.push_str(&new_block);

        // Write to a temp file first, then use elevated powershell to copy it
        let temp_dir = std::env::temp_dir();
        let temp_file = temp_dir.join("lumine_hosts_temp");
        fs::write(&temp_file, &new_content)
            .map_err(|e| format!("Failed to write temp hosts file: {}", e))?;

        #[cfg(target_os = "windows")]
        {
            let ps_script = format!(
                "Copy-Item -Path '{}' -Destination '{}' -Force",
                temp_file.display(),
                hosts_path.display()
            );
            let mut cmd = crate::utils::create_command("powershell");
            cmd.args([
                "-Command",
                &format!(
                    "Start-Process powershell -ArgumentList '-Command', '{}' -Verb RunAs -Wait -WindowStyle Hidden",
                    ps_script.replace("'", "''")
                ),
            ]);
            cmd.creation_flags(0x08000000); // CREATE_NO_WINDOW
            let output = cmd.output()
                .map_err(|e| format!("Failed to elevate: {}", e))?;
            if !output.status.success() {
                let _ = fs::remove_file(&temp_file);
                return Err("Failed to write hosts file. UAC prompt was cancelled or denied.".to_string());
            }
        }

        #[cfg(not(target_os = "windows"))]
        {
            // On macOS/Linux, try direct write first, then sudo
            if fs::write(&hosts_path, &new_content).is_err() {
                let output = crate::utils::create_command("sudo")
                    .arg("cp")
                    .arg(&temp_file)
                    .arg(&hosts_path)
                    .output()
                    .map_err(|e| format!("Failed to copy hosts file with sudo: {}", e))?;
                if !output.status.success() {
                    let _ = fs::remove_file(&temp_file);
                    return Err("Failed to write hosts file. Permission denied.".to_string());
                }
            }
        }

        let _ = fs::remove_file(&temp_file);
        Ok(())
    }
}
