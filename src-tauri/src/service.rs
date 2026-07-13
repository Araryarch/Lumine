use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::process::{Child, Stdio};
use std::sync::{Arc, Mutex};
#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;
use std::io::{BufRead, BufReader, Write};
use std::fs::OpenOptions;
use std::thread;

#[cfg(target_os = "windows")]
const CREATE_NO_WINDOW: u32 = 0x08000000;

pub fn get_docker_executable() -> String {
    #[cfg(target_os = "windows")]
    {
        if let Ok(status) = crate::utils::create_command("docker").arg("--version").stdout(Stdio::null()).stderr(Stdio::null()).status() {
            if status.success() {
                return "docker".to_string();
            }
        }
        let default_path = "C:\\Program Files\\Docker\\Docker\\resources\\bin\\docker.exe";
        if std::path::Path::new(default_path).exists() {
            return default_path.to_string();
        }
        "docker".to_string()
    }
    #[cfg(not(target_os = "windows"))]
    {
        "docker".to_string()
    }
}

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum ServiceStatus {
    Running,
    Stopped,
    Error(String),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServiceConfig {
    pub port: u16,
    pub executable_path: String,
    pub arguments: String,
    pub service_type: String, // "Service", "Language", "Database"
    #[serde(default = "default_runner")]
    pub runner: String, // "binary" or "docker"
    #[serde(default)]
    pub container_port: Option<u16>,
    #[serde(default)]
    pub volume_path: Option<String>,
    #[serde(default)]
    pub env: Option<std::collections::HashMap<String, String>>,
    #[serde(default)]
    pub auto_start: Option<bool>,
}
fn default_runner() -> String {
    "binary".to_string()
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServiceInfo {
    pub id: String,
    pub name: String,
    #[serde(default)]
    pub description: String,
    #[serde(default = "default_status")]
    pub status: ServiceStatus,
    pub config: ServiceConfig,
    #[serde(default)]
    pub log: Vec<String>,
}

fn default_status() -> ServiceStatus {
    ServiceStatus::Stopped
}

pub struct ManagedService {
    pub info: ServiceInfo,
    pub process: Option<Child>,
    pub log: Arc<Mutex<Vec<String>>>,
}

pub struct ServiceManager {
    services: Mutex<HashMap<String, ManagedService>>,
    data_file: std::path::PathBuf,
}

impl ServiceManager {
    pub fn new(app_data_dir: std::path::PathBuf) -> Self {
        if !app_data_dir.exists() {
            let _ = std::fs::create_dir_all(&app_data_dir);
        }
        let logs_dir = app_data_dir.join("logs");
        if !logs_dir.exists() {
            let _ = std::fs::create_dir_all(&logs_dir);
        }
        let data_file = app_data_dir.join("services.json");
        
        let mut services_map = HashMap::new();
        
        if data_file.exists() {
            if let Ok(content) = std::fs::read_to_string(&data_file) {
                match serde_json::from_str::<Vec<ServiceInfo>>(&content) {
                    Ok(loaded_services) => {
                        for svc in loaded_services {
                            let mut info = svc.clone();
                            info.status = ServiceStatus::Stopped; // Reset status on load
                            info.log = Vec::new();
                            services_map.insert(info.id.clone(), ManagedService {
                                info,
                                process: None,
                                log: Arc::new(Mutex::new(Vec::new())),
                            });
                        }
                    }
                    Err(e) => {
                        eprintln!("[Lumine] Failed to parse services.json: {}", e);
                    }
                }
            }
        } else {
            // Create empty array
            let _ = std::fs::write(&data_file, "[]");
        }

        ServiceManager {
            services: Mutex::new(services_map),
            data_file,
        }
    }

    fn save_to_disk(&self, map: &std::sync::MutexGuard<'_, HashMap<String, ManagedService>>) {
        let list: Vec<ServiceInfo> = map.values().map(|s| s.info.clone()).collect();
        if let Ok(json) = serde_json::to_string_pretty(&list) {
            let _ = std::fs::write(&self.data_file, json);
        }
    }

    pub fn add_service(&self, info: ServiceInfo) -> Result<(), String> {
        let mut map = self.services.lock().unwrap();
        if map.contains_key(&info.id) {
            return Err("Service with this ID already exists".to_string());
        }
        map.insert(info.id.clone(), ManagedService {
            info,
            process: None,
            log: Arc::new(Mutex::new(Vec::new())),
        });
        self.save_to_disk(&map);
        Ok(())
    }

    pub fn edit_service(&self, id: &str, mut new_info: ServiceInfo) -> Result<(), String> {
        // If the service is running, stop it first before editing
        let _ = self.stop(id);
        
        let mut map = self.services.lock().unwrap();
        
        // Extract the log Arc so we don't hold the immutable borrow on map
        let existing_log = if let Some(existing) = map.get(id) {
            new_info.log = existing.info.log.clone();
            new_info.status = existing.info.status.clone();
            Some(existing.log.clone())
        } else {
            None
        };

        if let Some(log_arc) = existing_log {
            map.insert(id.to_string(), ManagedService {
                info: new_info,
                process: None,
                log: log_arc,
            });
            self.save_to_disk(&map);
            Ok(())
        } else {
            Err(format!("Service '{}' not found", id))
        }
    }

    pub fn delete_service(&self, id: &str) -> Result<(), String> {
        let _ = self.stop(id);
        let mut map = self.services.lock().unwrap();
        if map.remove(id).is_some() {
            self.save_to_disk(&map);
            Ok(())
        } else {
            Err(format!("Service '{}' not found", id))
        }
    }

    pub fn get_all(&self) -> Vec<ServiceInfo> {
        let mut map = self.services.lock().unwrap();
        let mut list: Vec<_> = map.values_mut().map(|s| {
            if let Some(mut child) = s.process.take() {
                match child.try_wait() {
                    Ok(Some(status)) => {
                        s.info.status = ServiceStatus::Stopped;
                        if let Ok(mut lg) = s.log.lock() {
                            lg.push(format!("[{}] Process exited with: {}", chrono::Local::now().format("%H:%M:%S"), status));
                        }
                    }
                    Ok(None) => {
                        s.process = Some(child); // Put it back, still running
                    }
                    Err(_) => {
                        s.info.status = ServiceStatus::Stopped;
                    }
                }
            }

            let mut info = s.info.clone();
            if let Ok(log_guard) = s.log.lock() {
                info.log = log_guard.clone();
            }
            info
        }).collect();
        list.sort_by(|a, b| a.name.cmp(&b.name));
        list
    }

    pub fn get(&self, id: &str) -> Option<ServiceInfo> {
        let mut map = self.services.lock().unwrap();
        map.get_mut(id).map(|s| {
            if let Some(mut child) = s.process.take() {
                match child.try_wait() {
                    Ok(Some(status)) => {
                        s.info.status = ServiceStatus::Stopped;
                        if let Ok(mut lg) = s.log.lock() {
                            lg.push(format!("[{}] Process exited with: {}", chrono::Local::now().format("%H:%M:%S"), status));
                        }
                    }
                    Ok(None) => {
                        s.process = Some(child);
                    }
                    Err(_) => {
                        s.info.status = ServiceStatus::Stopped;
                    }
                }
            }

            let mut info = s.info.clone();
            if let Ok(log_guard) = s.log.lock() {
                info.log = log_guard.clone();
            }
            info
        })
    }

    pub fn update_port(&self, id: &str, new_port: u16) -> Result<(), String> {
        let mut map = self.services.lock().unwrap();
        if let Some(svc) = map.get_mut(id) {
            svc.info.config.port = new_port;
            if let Ok(mut log_guard) = svc.log.lock() {
                log_guard.push(format!(
                    "[{}] Port changed to {}",
                    chrono::Local::now().format("%H:%M:%S"),
                    new_port
                ));
            }
            self.save_to_disk(&map);
            Ok(())
        } else {
            Err(format!("Service '{}' not found", id))
        }
    }



    pub fn clear_log(&self, id: &str) -> Result<(), String> {
        let mut map = self.services.lock().unwrap();
        if let Some(svc) = map.get_mut(id) {
            if let Ok(mut lg) = svc.log.lock() {
                lg.clear();
            }
            let log_file_path = self.data_file.parent().unwrap().join("logs").join(format!("{}.log", id));
            let _ = std::fs::write(&log_file_path, "");
            Ok(())
        } else {
            Err(format!("Service '{}' not found", id))
        }
    }

    pub fn stop_all(&self) {
        let ids: Vec<String> = {
            let map = self.services.lock().unwrap();
            map.keys().cloned().collect()
        };
        for id in ids {
            let _ = self.stop(&id);
        }
    }

    pub fn start(&self, id: &str) -> Result<(), String> {
        let mut map = self.services.lock().unwrap();
        let svc = map.get_mut(id).ok_or_else(|| format!("Service '{}' not found", id))?;

        if svc.info.status == ServiceStatus::Running {
            return Err(format!("{} is already running", svc.info.name));
        }

        let exe_path = svc.info.config.executable_path.clone();
        
        let docker_exe = get_docker_executable();
        let mut cmd = crate::utils::create_command(&docker_exe);
        cmd.arg("run");
        cmd.arg("--rm");
        cmd.arg(format!("--name=lumine_{}", id));
        if svc.info.config.port > 0 {
            let container_port = svc.info.config.container_port.unwrap_or(svc.info.config.port);
            cmd.arg("-p");
            cmd.arg(format!("{}:{}", svc.info.config.port, container_port));
        }
        
        if let Some(vol_path) = &svc.info.config.volume_path {
            if !vol_path.trim().is_empty() {
                cmd.arg("-v");
                cmd.arg(format!("lumine_data_{}:{}", id, vol_path.trim()));
            }
        }
        
        if let Some(env_vars) = &svc.info.config.env {
            for (key, value) in env_vars {
                cmd.arg("-e");
                cmd.arg(format!("{}={}", key, value));
            }
        }
        
        let args_str = svc.info.config.arguments.clone();
        if !args_str.trim().is_empty() {
            match shell_words::split(&args_str) {
                Ok(args) => { cmd.args(args); },
                Err(e) => {
                    let err_msg = format!("Failed to parse docker arguments: {}", e);
                    svc.info.status = ServiceStatus::Error(err_msg.clone());
                    if let Ok(mut lg) = svc.log.lock() { lg.push(err_msg.clone()); }
                    return Err(err_msg);
                }
            }
        }
        cmd.arg(exe_path.clone()); // Image name
        
        #[cfg(target_os = "windows")]
        cmd.creation_flags(CREATE_NO_WINDOW);

        cmd.stdout(Stdio::piped())
           .stderr(Stdio::piped());

        let log_file_path = self.data_file.parent().unwrap().join("logs").join(format!("{}.log", id));
        
        // Log rotation: if > 5MB, rename to .old
        if let Ok(meta) = std::fs::metadata(&log_file_path) {
            if meta.len() > 5 * 1024 * 1024 {
                let _ = std::fs::rename(&log_file_path, log_file_path.with_extension("log.old"));
            }
        }
        
        let _ = OpenOptions::new().create(true).append(true).open(&log_file_path).and_then(|mut f| writeln!(f, "--- Started service {} ---\n", id));
        match cmd.spawn() {
            Ok(mut child) => {
                let stdout = child.stdout.take();
                let stderr = child.stderr.take();
                let log_arc_out = Arc::clone(&svc.log);
                let log_arc_err = Arc::clone(&svc.log);

                // Spawn thread for stdout
                if let Some(out) = stdout {
                    let log_path = log_file_path.clone();
                    thread::spawn(move || {
                        let reader = BufReader::new(out);
                        for line in reader.lines() {
                            if let Ok(l) = line {
                                let display_line = if l.contains("failed to connect to the docker API") || l.contains("error during connect") {
                                    "Docker Error: Please open your Docker Desktop or check if the Docker service is running.".to_string()
                                } else {
                                    l.clone()
                                };
                                
                                if let Ok(mut lg) = log_arc_out.lock() {
                                    lg.push(display_line.clone());
                                    if lg.len() > 100 { lg.remove(0); } // Keep last 100 lines
                                }
                                if let Ok(mut f) = OpenOptions::new().append(true).create(true).open(&log_path) {
                                    let _ = writeln!(f, "{}", display_line);
                                }
                            }
                        }
                    });
                }

                // Spawn thread for stderr
                if let Some(err) = stderr {
                    let log_path = log_file_path.clone();
                    thread::spawn(move || {
                        let reader = BufReader::new(err);
                        for line in reader.lines() {
                            if let Ok(l) = line {
                                let display_line = if l.contains("failed to connect to the docker API") || l.contains("error during connect") {
                                    "Docker Error: Please open your Docker Desktop or check if the Docker service is running.".to_string()
                                } else {
                                    l.clone()
                                };
                                
                                if let Ok(mut lg) = log_arc_err.lock() {
                                    lg.push(display_line.clone());
                                    if lg.len() > 100 { lg.remove(0); }
                                }
                                if let Ok(mut f) = OpenOptions::new().append(true).create(true).open(&log_path) {
                                    let _ = writeln!(f, "{}", display_line);
                                }
                            }
                        }
                    });
                }

                svc.process = Some(child);
                svc.info.status = ServiceStatus::Running;
                if let Ok(mut lg) = svc.log.lock() {
                    lg.push(format!("[{}] Service {} started successfully.", chrono::Local::now().format("%H:%M:%S"), svc.info.name));
                }
                Ok(())
            }
            Err(e) => {
                let err_msg = format!("Failed to start process (expected at {}): {}", exe_path, e);
                svc.info.status = ServiceStatus::Error(err_msg.clone());
                if let Ok(mut lg) = svc.log.lock() {
                    lg.push(format!("[{}] {}", chrono::Local::now().format("%H:%M:%S"), err_msg));
                }
                Err(err_msg)
            }
        }
    }

    pub fn stop(&self, id: &str) -> Result<(), String> {
        let child_opt = {
            let mut map = self.services.lock().unwrap();
            let svc = map.get_mut(id).ok_or_else(|| format!("Service '{}' not found", id))?;

            if svc.info.status != ServiceStatus::Running {
                return Err(format!("{} is not running", svc.info.name));
            }

            svc.info.status = ServiceStatus::Stopped;
            if let Ok(mut lg) = svc.log.lock() {
                lg.push(format!("[{}] Stopping...", chrono::Local::now().format("%H:%M:%S")));
            }
            
            svc.process.take()
        }; // Lock dropped

        if let Some(mut child) = child_opt {
            let docker_exe = get_docker_executable();
            let mut stop_cmd = crate::utils::create_command(&docker_exe);
            stop_cmd.arg("stop").arg(format!("lumine_{}", id));
            stop_cmd.stdout(Stdio::null()).stderr(Stdio::null());
            #[cfg(target_os = "windows")]
            stop_cmd.creation_flags(CREATE_NO_WINDOW);
            
            let _ = stop_cmd.status(); 
            let _ = child.kill();
            child.wait().ok();
        }

        let mut map = self.services.lock().unwrap();
        if let Some(svc) = map.get_mut(id) {
            if let Ok(mut lg) = svc.log.lock() {
                lg.push(format!("[{}] Stopped", chrono::Local::now().format("%H:%M:%S")));
            }
        }
        Ok(())
    }

    pub fn restart(&self, id: &str) -> Result<(), String> {
        let _ = self.stop(id);
        self.start(id)
    }
}

impl Drop for ServiceManager {
    fn drop(&mut self) {
        if let Ok(mut map) = self.services.lock() {
            for (id, svc) in map.iter_mut() {
                if let Some(mut child) = svc.process.take() {
                    // Must explicitly stop the Docker container, not just kill the CLI
                    let docker_exe = get_docker_executable();
                    let mut stop_cmd = crate::utils::create_command(&docker_exe);
                    stop_cmd.arg("stop").arg(format!("lumine_{}", id));
                    stop_cmd.stdout(Stdio::null()).stderr(Stdio::null());
                    #[cfg(target_os = "windows")]
                    stop_cmd.creation_flags(CREATE_NO_WINDOW);
                    let _ = stop_cmd.status();
                    let _ = child.kill();
                    let _ = child.wait();
                }
            }
        }
    }
}
