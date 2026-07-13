use crate::service::ServiceManager;
use crate::settings_manager::SettingsManager;
use tauri::{State, Manager};

#[tauri::command]
pub fn get_services(manager: State<'_, ServiceManager>) -> Vec<crate::service::ServiceInfo> {
    manager.get_all()
}

#[tauri::command]
pub fn get_service(
    manager: State<'_, ServiceManager>,
    id: String,
) -> Option<crate::service::ServiceInfo> {
    manager.get(&id)
}

#[tauri::command]
pub fn start_service(manager: State<'_, ServiceManager>, id: String) -> Result<String, String> {
    manager.start(&id)?;
    Ok(format!("{} started", id))
}

#[tauri::command]
pub fn stop_service(manager: State<'_, ServiceManager>, id: String) -> Result<String, String> {
    manager.stop(&id)?;
    Ok(format!("{} stopped", id))
}

#[tauri::command]
pub fn restart_service(manager: State<'_, ServiceManager>, id: String) -> Result<String, String> {
    manager.restart(&id)?;
    Ok(format!("{} restarted", id))
}

#[tauri::command]
pub fn set_port(manager: State<'_, ServiceManager>, id: String, port: u16) -> Result<(), String> {
    manager.update_port(&id, port)
}

#[tauri::command]
pub fn clear_service_log(manager: State<'_, ServiceManager>, id: String) -> Result<(), String> {
    manager.clear_log(&id)
}

#[tauri::command]
pub fn check_port(port: u16) -> bool {
    std::net::TcpListener::bind(format!("127.0.0.1:{}", port)).is_ok()
}

#[tauri::command]
pub fn start_all(manager: State<'_, ServiceManager>) -> Vec<String> {
    let all_services = manager.get_all();
    all_services
        .iter()
        .map(|svc| match manager.start(&svc.id) {
            Ok(_) => format!("{} started", svc.id),
            Err(e) => format!("{} failed: {}", svc.id, e),
        })
        .collect()
}

#[tauri::command]
pub fn stop_all(manager: State<'_, ServiceManager>) -> Vec<String> {
    let all_services = manager.get_all();
    all_services
        .iter()
        .map(|svc| match manager.stop(&svc.id) {
            Ok(_) => format!("{} stopped", svc.id),
            Err(e) => format!("{} failed: {}", svc.id, e),
        })
        .collect()
}

#[tauri::command]
pub fn add_service(
    manager: State<'_, ServiceManager>,
    name: String,
    service_type: String,
    executable_path: String,
    arguments: String,
    port: u16,
    container_port: Option<u16>,
    volume_path: Option<String>,
    env: Option<std::collections::HashMap<String, String>>,
    auto_start: Option<bool>,
) -> Result<(), String> {
    let id = name.to_lowercase().replace(" ", "_");
    let info = crate::service::ServiceInfo {
        id,
        name,
        description: format!("Custom {} Service", service_type),
        status: crate::service::ServiceStatus::Stopped,
        config: crate::service::ServiceConfig {
            port,
            executable_path,
            arguments,
            service_type,
            runner: "docker".to_string(),
            container_port,
            volume_path,
            env,
            auto_start,
        },
        log: Vec::new(),
    };
    manager.add_service(info)
}

#[tauri::command]
pub fn edit_service(
    manager: State<'_, ServiceManager>,
    id: String,
    name: String,
    service_type: String,
    executable_path: String,
    arguments: String,
    port: u16,
    container_port: Option<u16>,
    volume_path: Option<String>,
    env: Option<std::collections::HashMap<String, String>>,
    auto_start: Option<bool>,
) -> Result<(), String> {
    let info = crate::service::ServiceInfo {
        id: id.clone(),
        name,
        description: format!("Custom {} Service", service_type),
        status: crate::service::ServiceStatus::Stopped,
        config: crate::service::ServiceConfig {
            port,
            executable_path,
            arguments,
            service_type,
            runner: "docker".to_string(),
            container_port,
            volume_path,
            env,
            auto_start,
        },
        log: Vec::new(),
    };
    manager.edit_service(&id, info)
}

#[tauri::command]
pub fn delete_service(manager: State<'_, ServiceManager>, id: String) -> Result<(), String> {
    manager.delete_service(&id)
}

// ---- Hosts Commands ----

fn sync_proxy(
    proxy_manager: &State<'_, crate::proxy_manager::ProxyManager>,
    hosts_manager: &State<'_, crate::hosts_manager::HostsManager>,
    service_manager: &State<'_, crate::service::ServiceManager>,
    project_manager: &State<'_, crate::project_manager::ProjectManager>,
) {
    if let Ok(hosts) = hosts_manager.get_hosts() {
        let services = service_manager.get_all();
        let projects = project_manager.get_all();
        let _ = proxy_manager.sync_config(&hosts, &services, &projects);
    }
}

#[tauri::command]
pub fn get_hosts(
    manager: State<'_, crate::hosts_manager::HostsManager>,
) -> Result<Vec<crate::hosts_manager::HostEntry>, String> {
    manager.get_hosts()
}

#[tauri::command]
pub fn add_host(
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    name: String,
    php: Option<String>,
    comment: String,
    is_enable: bool,
    is_ipv6: bool,
) -> Result<(), String> {
    let result = hosts_manager.add_host(crate::hosts_manager::HostEntry {
        name,
        php,
        comment,
        is_enable,
        is_ipv6,
    });
    if result.is_ok() {
        sync_proxy(&proxy_manager, &hosts_manager, &service_manager, &project_manager);
    }
    result
}

#[tauri::command]
pub fn edit_host(
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    old_name: String,
    name: String,
    php: Option<String>,
    comment: String,
    is_enable: bool,
    is_ipv6: bool,
) -> Result<(), String> {
    let result = hosts_manager.edit_host(&old_name, crate::hosts_manager::HostEntry {
        name,
        php,
        comment,
        is_enable,
        is_ipv6,
    });
    if result.is_ok() {
        sync_proxy(&proxy_manager, &hosts_manager, &service_manager, &project_manager);
    }
    result
}

#[tauri::command]
pub fn toggle_host(
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    name: String,
    is_enable: bool,
) -> Result<(), String> {
    let result = hosts_manager.toggle_host(&name, is_enable);
    if result.is_ok() {
        sync_proxy(&proxy_manager, &hosts_manager, &service_manager, &project_manager);
    }
    result
}

#[tauri::command]
pub fn delete_host(
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    name: String,
) -> Result<(), String> {
    let result = hosts_manager.delete_host(&name);
    if result.is_ok() {
        sync_proxy(&proxy_manager, &hosts_manager, &service_manager, &project_manager);
    }
    result
}

#[derive(serde::Deserialize)]
struct DockerTagResult {
    name: String,
}

#[derive(serde::Deserialize)]
struct DockerTagsResponse {
    results: Vec<DockerTagResult>,
}

#[tauri::command]
pub async fn fetch_docker_tags(image_name: String) -> Result<Vec<String>, String> {
    let repo = if image_name.contains('/') {
        image_name.clone()
    } else {
        format!("library/{}", image_name)
    };

    let url = format!("https://hub.docker.com/v2/repositories/{}/tags/?page_size=100", repo);
    
    let client = reqwest::Client::builder()
        .user_agent("Lumine")
        .timeout(std::time::Duration::from_secs(10))
        .build()
        .map_err(|e| format!("Failed to create client: {}", e))?;

    let response = client
        .get(&url)
        .send()
        .await
        .map_err(|e| format!("Network error: {}", e))?;

    if !response.status().is_success() {
        return Err(format!("Docker API error: {}", response.status()));
    }

    let data: DockerTagsResponse = response.json()
        .await
        .map_err(|e| format!("Failed to parse response: {}", e))?;

    let tags = data.results.into_iter().map(|r| r.name).collect();
    Ok(tags)
}

#[tauri::command]
pub fn open_url(url: String) -> Result<(), String> {
    opener::open(&url).map_err(|e| format!("Failed to open URL: {}", e))
}

#[tauri::command]
pub fn open_hosts_file() -> Result<(), String> {
    let path = crate::hosts_manager::HostsManager::get_system_hosts_path();
    #[cfg(target_os = "windows")]
    {
        crate::utils::create_command("notepad")
            .arg(path)
            .spawn()
            .map_err(|e| format!("Failed to open notepad: {}", e))?;
    }
    #[cfg(target_os = "macos")]
    {
        crate::utils::create_command("open")
            .arg("-e")
            .arg(path)
            .spawn()
            .map_err(|e| format!("Failed to open editor: {}", e))?;
    }
    #[cfg(target_os = "linux")]
    {
        crate::utils::create_command("xdg-open")
            .arg(path)
            .spawn()
            .map_err(|e| format!("Failed to open editor: {}", e))?;
    }
    Ok(())
}

// ---- MkCert Commands ----

#[tauri::command]
pub fn check_mkcert(
    manager: State<'_, crate::mkcert_manager::MkCertManager>,
) -> Result<bool, String> {
    Ok(manager.check_installed())
}

#[tauri::command]
pub async fn download_mkcert(
    manager: State<'_, crate::mkcert_manager::MkCertManager>,
) -> Result<String, String> {
    manager.download_mkcert().await
}

#[tauri::command]
pub fn install_root_ca(
    manager: State<'_, crate::mkcert_manager::MkCertManager>,
) -> Result<String, String> {
    manager.install_root_ca()
}

#[tauri::command]
pub fn generate_cert(
    manager: State<'_, crate::mkcert_manager::MkCertManager>,
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    domain: String,
) -> Result<crate::mkcert_manager::CertEntry, String> {
    let result = manager.generate_cert(&domain);
    if result.is_ok() {
        sync_proxy(&proxy_manager, &hosts_manager, &service_manager, &project_manager);
    }
    result
}

#[tauri::command]
pub fn get_certs(
    manager: State<'_, crate::mkcert_manager::MkCertManager>,
) -> Result<Vec<crate::mkcert_manager::CertEntry>, String> {
    manager.get_certs()
}

#[tauri::command]
pub fn delete_cert(
    manager: State<'_, crate::mkcert_manager::MkCertManager>,
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    domain: String,
) -> Result<(), String> {
    let result = manager.delete_cert(&domain);
    if result.is_ok() {
        sync_proxy(&proxy_manager, &hosts_manager, &service_manager, &project_manager);
    }
    result
}

// ---- Settings Commands ----

#[tauri::command]
pub fn get_settings(
    manager: State<'_, crate::settings_manager::SettingsManager>,
) -> Result<crate::settings_manager::AppSettings, String> {
    manager.get_settings()
}

#[tauri::command]
pub fn save_settings(
    app_handle: tauri::AppHandle,
    manager: State<'_, crate::settings_manager::SettingsManager>,
    settings: crate::settings_manager::AppSettings,
) -> Result<(), String> {
    use tauri_plugin_autostart::ManagerExt;
    
    // Configure Autostart based on setting
    let autolaunch = app_handle.autolaunch();
    if settings.start_on_boot {
        let _ = autolaunch.enable();
    } else {
        let _ = autolaunch.disable();
    }
    
    manager.save_settings(settings)
}

static DOCKER_AUTOSTART_ATTEMPTED: std::sync::atomic::AtomicBool = std::sync::atomic::AtomicBool::new(false);

#[tauri::command]
pub fn check_docker() -> bool {
    let docker_exe = crate::service::get_docker_executable();
    let mut cmd = crate::utils::create_command(&docker_exe);
    cmd.arg("info");
    cmd.stdout(std::process::Stdio::null());
    cmd.stderr(std::process::Stdio::null());
    #[cfg(target_os = "windows")]
    {
        use std::os::windows::process::CommandExt;
        cmd.creation_flags(0x08000000);
    }
    
    let is_running = cmd.status().map(|s| s.success()).unwrap_or(false);
    
    if !is_running {
        // Attempt to start only once per session
        if !DOCKER_AUTOSTART_ATTEMPTED.swap(true, std::sync::atomic::Ordering::Relaxed) {
            #[cfg(target_os = "windows")]
            {
                let docker_desktop = "C:\\Program Files\\Docker\\Docker\\Docker Desktop.exe";
                if std::path::Path::new(docker_desktop).exists() {
                    let mut start_cmd = crate::utils::create_command(docker_desktop);
                    use std::os::windows::process::CommandExt;
                    start_cmd.creation_flags(0x08000000);
                    let _ = start_cmd.spawn();
                }
            }
        }
    }
    
    is_running
}

#[tauri::command]
pub fn open_log_file(id: String, app_handle: tauri::AppHandle) -> Result<(), String> {
    use tauri::Manager;
    let app_data_dir = app_handle.path().app_data_dir().map_err(|e| e.to_string())?;
    let log_file_path = app_data_dir.join("logs").join(format!("{}.log", id));
    if log_file_path.exists() {
        if let Err(e) = opener::open(&log_file_path) {
            return Err(format!("Failed to open log file: {}", e));
        }
    } else {
        return Err("Log file does not exist yet".to_string());
    }
    Ok(())
}

#[tauri::command]
pub fn exit_app(app_handle: tauri::AppHandle) {
    app_handle.exit(0);
}

// =======================
// PROJECT COMMANDS
// =======================

use crate::project_manager::{ProjectManager, ProjectInfo};

#[tauri::command]
pub fn get_projects(manager: State<'_, ProjectManager>) -> Vec<ProjectInfo> {
    manager.get_all()
}

#[tauri::command]
pub fn get_workspace_dir(settings: State<'_, SettingsManager>) -> String {
    settings.get_settings().map(|s| s.document_root).unwrap_or_else(|_| "C:\\Lumine\\www".to_string())
}

#[tauri::command]
pub fn add_project(
    manager: State<'_, ProjectManager>, 
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    project: ProjectInfo
) -> Result<(), String> {
    let res = manager.add_project(project);
    let hosts = hosts_manager.get_hosts().unwrap_or_default();
    let services = service_manager.get_all();
    let projects = manager.get_all();
    let _ = proxy_manager.sync_config(&hosts, &services, &projects);
    res
}

#[tauri::command]
pub fn edit_project(manager: State<'_, ProjectManager>, id: String, command: String) -> Result<(), String> {
    manager.edit_project(&id, command)
}

#[tauri::command]
pub fn start_project(
    manager: State<'_, ProjectManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    id: String
) -> Result<(), String> {
    let res = manager.start(&id);
    let hosts = hosts_manager.get_hosts().unwrap_or_default();
    let services = service_manager.get_all();
    let projects = manager.get_all();
    let _ = proxy_manager.sync_config(&hosts, &services, &projects);
    res
}

#[tauri::command]
pub fn stop_project(
    manager: State<'_, ProjectManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    id: String
) -> Result<(), String> {
    let res = manager.stop(&id);
    let hosts = hosts_manager.get_hosts().unwrap_or_default();
    let services = service_manager.get_all();
    let projects = manager.get_all();
    let _ = proxy_manager.sync_config(&hosts, &services, &projects);
    res
}

#[tauri::command]
pub fn clear_project_log(manager: State<'_, ProjectManager>, id: String) -> Result<(), String> {
    manager.clear_log(&id)
}

#[tauri::command]
pub fn open_project_terminal(
    manager: State<'_, ProjectManager>,
    settings_mgr: State<'_, crate::settings_manager::SettingsManager>,
    id: String,
    image_override: Option<String>
) -> Result<(), String> {
    let settings = settings_mgr.get_settings().unwrap_or_default();
    manager.open_terminal(&id, &settings, image_override)
}

#[tauri::command]
pub fn delete_project(
    manager: State<'_, ProjectManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
    hosts_manager: State<'_, crate::hosts_manager::HostsManager>,
    service_manager: State<'_, crate::service::ServiceManager>,
    id: String, 
    delete_folder: bool
) -> Result<(), String> {
    let res = manager.delete_project(&id, delete_folder);
    let hosts = hosts_manager.get_hosts().unwrap_or_default();
    let services = service_manager.get_all();
    let projects = manager.get_all();
    let _ = proxy_manager.sync_config(&hosts, &services, &projects);
    res
}

#[tauri::command]
pub fn open_in_explorer(
    settings_mgr: State<'_, crate::settings_manager::SettingsManager>,
    path: String
) -> Result<(), String> {
    let settings = settings_mgr.get_settings().unwrap_or_default();
    
    // Command is basically `docker run -it --rm <executable> sh`
    
    // Fallback to opener if no specific preference or unknown
    if settings.file_explorer.is_empty() || settings.file_explorer == "explorer" {
        opener::open(&path).map_err(|e| format!("Failed to open folder: {}", e))
    } else {
        #[cfg(target_os = "windows")]
        {
            let _ = crate::utils::create_command(&settings.file_explorer)
                .arg(&path)
                .spawn();
            Ok(())
        }
        #[cfg(not(target_os = "windows"))]
        {
            opener::open(&path).map_err(|e| format!("Failed to open folder: {}", e))
        }
    }
}

#[derive(serde::Serialize)]
pub struct LanguageContainerInfo {
    id: String,
    name: String,
    status: String,
}

#[tauri::command]
pub fn get_language_container_status(executable: String) -> Result<Vec<LanguageContainerInfo>, String> {
    let safe_name = executable.replace(":", "_").replace("/", "_");
    let term_name = format!("lumine_lang_{}", safe_name);
    
    let output = crate::utils::create_command("docker")
        .args(["ps", "-a", "--format", "{{.ID}}|{{.Names}}|{{.Status}}", "-f", &format!("name={}", term_name)])
        .output()
        .map_err(|e| format!("Failed to run docker: {}", e))?;
        
    let out_str = String::from_utf8_lossy(&output.stdout);
    let mut containers = Vec::new();
    
    for line in out_str.lines() {
        let parts: Vec<&str> = line.split('|').collect();
        if parts.len() >= 3 && parts[1] == term_name {
            containers.push(LanguageContainerInfo {
                id: parts[0].to_string(),
                name: parts[1].to_string(),
                status: parts[2].to_string(),
            });
        }
    }
    
    Ok(containers)
}

#[tauri::command]
pub fn stop_language_container(executable: String) -> Result<(), String> {
    let safe_name = executable.replace(":", "_").replace("/", "_");
    let term_name = format!("lumine_lang_{}", safe_name);
    
    let status = crate::utils::create_command("docker")
        .args(["rm", "-f", &term_name])
        .status()
        .map_err(|e| format!("Failed to run docker: {}", e))?;
        
    if status.success() {
        Ok(())
    } else {
        Err("Failed to stop container".into())
    }
}

#[tauri::command]
pub fn open_language_editor(executable: String) -> Result<(), String> {
    let safe_name = executable.replace(":", "_").replace("/", "_");
    let container_name = format!("lumine_lang_{}", safe_name);
    
    // Check if the container is running first
    let output = std::process::Command::new("docker")
        .args(["ps", "-q", "-f", &format!("name={}", container_name)])
        .output()
        .map_err(|e| format!("Docker Error: {}", e))?;
        
    if output.stdout.is_empty() {
        return Err("Anda harus mengklik 'Open Terminal' terlebih dahulu agar container berjalan di background!".to_string());
    }

    let json_str = format!(r#"{{"containerName":"/{}"}}"#, container_name);
    let hex_str: String = json_str.bytes().map(|b| format!("{:02x}", b)).collect();
    
    // The workspace folder is hardcoded to /workspace as per open_service_terminal logic
    let uri = format!("vscode://vscode-remote/attached-container+{}/workspace", hex_str);
    
    open::that(&uri).map_err(|e| format!("Gagal membuka VS Code: {}", e))?;
    
    Ok(())
}

#[tauri::command]
pub fn open_service_terminal(
    app: tauri::AppHandle,
    settings_mgr: State<'_, crate::settings_manager::SettingsManager>,
    executable: String
) -> Result<(), String> {
    use tauri::Manager;
    let settings = settings_mgr.get_settings().unwrap_or_default();
    let term = settings.terminal_emulator;
    let app_data_dir = app.path().app_data_dir().unwrap().to_string_lossy().to_string();
    let doc_root = settings.document_root;
    
    let mut wsl_path = String::from("/workspace");
    let mut drive_mounts = format!("-v \"{}\":/workspace", doc_root);
    
    if doc_root.len() >= 3 && &doc_root[1..3] == ":\\" {
        let drive_letter = doc_root.chars().next().unwrap().to_lowercase().to_string();
        let path_part = &doc_root[3..].replace("\\", "/");
        wsl_path = format!("/mnt/{}/{}", drive_letter, path_part);
        drive_mounts = format!("-v {}:\\:/mnt/{}", drive_letter, drive_letter);
    }

    let rc_args = format!("-v \"{0}\\lumine_bashrc\":/tmp/lumine_bashrc:ro -v \"{0}\\lumine_bash_history\":/tmp/lumine_bash_history -e ENV=/tmp/lumine_bashrc {1} -w \"{2}\"", app_data_dir, drive_mounts, wsl_path);
    
    #[cfg(target_os = "windows")]
    {
        use std::os::windows::process::CommandExt;
        const CREATE_NEW_CONSOLE: u32 = 0x00000010;
        
        let mut cmd;
        let lower_term = term.to_lowercase();
        let safe_name = executable.replace(":", "_").replace("/", "_");
        let term_name = format!("lumine_lang_{}", safe_name);
        
        let inner_cmd_cmd = "sh -c \"command -v bash >/dev/null 2>&1 && exec bash --rcfile /tmp/lumine_bashrc || exec sh\"";
        let inner_cmd_ps = "sh -c 'command -v bash >/dev/null 2>&1 && exec bash --rcfile /tmp/lumine_bashrc || exec sh'";
        
        let cmd_str = format!("title Lumine Terminal - {0} & echo ======================================== & echo   Lumine Language Terminal & echo   Status: Background Session & echo   Image: {0} & echo ======================================== & echo. & docker exec -w \"{3}\" -it {1} {4} || (docker rm -f {1} >nul 2>&1 & timeout /t 1 >nul 2>&1 & echo Pulling image if needed (this may take a while)... & docker pull {0} & docker run -d --rm --name {1} {2} {0} tail -f /dev/null >nul 2>&1 & docker exec -w \"{3}\" -it {1} {4})", executable, term_name, rc_args, wsl_path, inner_cmd_cmd);
        let ps_str = format!("$host.ui.RawUI.WindowTitle='Lumine Terminal - {0}'\nWrite-Host '========================================'\nWrite-Host '  Lumine Language Terminal'\nWrite-Host '  Status: Background Session'\nWrite-Host '  Image: {0}'\nWrite-Host '========================================'\nWrite-Host ''\ndocker exec -w '{3}' -it {1} {4} 2>$null\nif (!$?) {{\n  docker rm -f {1} 2>$null | Out-Null\n  Start-Sleep -Seconds 1\n  Write-Host 'Pulling image if needed (this may take a while)...' -ForegroundColor Cyan\n  docker pull {0}\n  docker run -d --rm --name {1} {2} {0} tail -f /dev/null 2>$null | Out-Null\n  docker exec -w '{3}' -it {1} {4}\n}}", executable, term_name, rc_args, wsl_path, inner_cmd_ps);

        if lower_term == "wt" {
            cmd = crate::utils::create_command("wt.exe");
            cmd.args(["-w", "0", "nt", "powershell.exe", "-NoExit", "-Command", &ps_str]);
        } else if lower_term == "powershell" || lower_term == "pwsh" {
            cmd = crate::utils::create_command(format!("{}.exe", lower_term));
            cmd.creation_flags(CREATE_NEW_CONSOLE);
            cmd.args(["-NoExit", "-Command", &ps_str]);
        } else {
            // Default CMD
            cmd = crate::utils::create_command("cmd.exe");
            cmd.creation_flags(CREATE_NEW_CONSOLE);
            cmd.args(["/k", &cmd_str]);
        }
        
        cmd.spawn().map_err(|e| format!("Failed to open terminal: {}", e))?;
        Ok(())
    }
    #[cfg(not(target_os = "windows"))]
    {
        // For Mac/Linux, could be implemented similarly
        Err("Terminal spawning not yet implemented for this OS.".into())
    }
}

#[tauri::command]
pub fn open_in_editor(
    settings_mgr: State<'_, crate::settings_manager::SettingsManager>,
    path: String
) -> Result<(), String> {
    let settings = settings_mgr.get_settings().unwrap_or_default();
    let editor = if settings.code_editor.is_empty() { "code".to_string() } else { settings.code_editor };
    
    #[cfg(target_os = "windows")]
    {
        use std::os::windows::process::CommandExt;
        let mut cmd = crate::utils::create_command("cmd");
        cmd.args(["/C", &editor, "."]);
        cmd.current_dir(&path);
        cmd.creation_flags(0x08000000); // CREATE_NO_WINDOW
        let _ = cmd.spawn();
    }
    #[cfg(not(target_os = "windows"))]
    {
        let _ = crate::utils::create_command(&editor)
            .arg(".")
            .current_dir(&path)
            .spawn();
    }
    
    Ok(())
}

#[tauri::command]
pub async fn create_new_project(
    app: tauri::AppHandle,
    settings: State<'_, SettingsManager>,
    _name: String,
    _framework: String,
    command: String,
) -> Result<(), String> {
    use std::process::Stdio;
    use std::io::{BufRead, BufReader};
    use tauri::Emitter;

    let workspace = settings.get_settings().map(|s| s.document_root).unwrap_or_else(|_| "C:\\Lumine\\www".to_string());
    let mut cmd;

    if let Ok(args) = serde_json::from_str::<Vec<String>>(&command) {
        cmd = crate::utils::create_command(&args[0]);
        if args.len() > 1 {
            cmd.args(&args[1..]);
        }
        #[cfg(target_os = "windows")]
        {
            use std::os::windows::process::CommandExt;
            cmd.creation_flags(0x08000000);
        }
    } else {
        #[cfg(target_os = "windows")]
        {
            cmd = crate::utils::create_command("cmd");
            cmd.args(["/C", &command]);
            use std::os::windows::process::CommandExt;
            cmd.creation_flags(0x08000000);
        }
        #[cfg(not(target_os = "windows"))]
        {
            cmd = crate::utils::create_command("sh");
            cmd.args(["-c", &command]);
        }
    }

    cmd.current_dir(&workspace);
    cmd.stdout(Stdio::piped());
    cmd.stderr(Stdio::piped());

    let mut child = cmd.spawn().map_err(|e| format!("Failed to spawn process: {}", e))?;

    let stdout = child.stdout.take();
    let stderr = child.stderr.take();

    let app_clone = app.clone();
    if let Some(out) = stdout {
        std::thread::spawn(move || {
            let reader = BufReader::new(out);
            for line in reader.lines().filter_map(Result::ok) {
                let _ = app_clone.emit("project-creation-log", line);
            }
        });
    }

    let app_clone = app.clone();
    if let Some(err) = stderr {
        std::thread::spawn(move || {
            let reader = BufReader::new(err);
            for line in reader.lines().filter_map(Result::ok) {
                let _ = app_clone.emit("project-creation-log", line);
            }
        });
    }

    let status = child.wait().map_err(|e| format!("Failed to wait: {}", e))?;
    
    if status.success() {
        let _ = app.emit("project-creation-done", "Success");
        Ok(())
    } else {
        let err_msg = format!("Process exited with status: {}", status);
        let _ = app.emit("project-creation-error", &err_msg);
        Err(err_msg)
    }
}

#[derive(serde::Serialize)]
pub struct AvailableTools {
    terminals: Vec<String>,
    editors: Vec<String>,
    explorers: Vec<String>,
}

#[tauri::command]
pub fn get_available_tools() -> Result<AvailableTools, String> {
    let mut terminals = Vec::new();
    let mut editors = Vec::new();
    let mut explorers = Vec::new();

    // Helper to check if a command exists in PATH on Windows
    let command_exists = |cmd: &str| -> bool {
        #[cfg(target_os = "windows")]
        {
            crate::utils::create_command("where")
                .arg(cmd)
                .stdout(std::process::Stdio::null())
                .stderr(std::process::Stdio::null())
                .status()
                .map(|status| status.success())
                .unwrap_or(false)
        }
        #[cfg(not(target_os = "windows"))]
        {
            crate::utils::create_command("which")
                .arg(cmd)
                .stdout(std::process::Stdio::null())
                .stderr(std::process::Stdio::null())
                .status()
                .map(|status| status.success())
                .unwrap_or(false)
        }
    };

    // --- Terminals ---
    let term_list = ["cmd", "powershell", "wt", "bash", "zsh", "fish", "gnome-terminal", "konsole", "xfce4-terminal", "terminator", "alacritty", "kitty", "wezterm"];
    for t in term_list {
        if t == "cmd" {
            #[cfg(target_os = "windows")]
            terminals.push(t.to_string());
        } else {
            if command_exists(t) { terminals.push(t.to_string()); }
        }
    }

    // --- Editors ---
    let editor_list = ["code", "cursor", "zed", "webstorm", "webstorm64", "phpstorm", "phpstorm64", "subl", "notepad", "nvim", "vim", "nano", "hx"];
    for e in editor_list {
        if e == "notepad" {
            #[cfg(target_os = "windows")]
            editors.push(e.to_string());
        } else {
            if command_exists(e) {
                // Normalize names for JetBrains on Windows
                if e == "webstorm64" {
                    if !editors.contains(&"webstorm".to_string()) { editors.push("webstorm".to_string()); }
                } else if e == "phpstorm64" {
                    if !editors.contains(&"phpstorm".to_string()) { editors.push("phpstorm".to_string()); }
                } else {
                    editors.push(e.to_string());
                }
            }
        }
    }

    // --- Explorers ---
    let explorer_list = ["explorer", "dolphin", "thunar", "nautilus", "nemo", "pcmanfm", "xdg-open", "open"];
    for exp in explorer_list {
        if exp == "explorer" {
            #[cfg(target_os = "windows")]
            explorers.push(exp.to_string());
        } else {
            if command_exists(exp) { explorers.push(exp.to_string()); }
        }
    }

    Ok(AvailableTools {
        terminals,
        editors,
        explorers,
    })
}
#[derive(Clone, serde::Serialize)]
struct PullLogPayload {
    image: String,
    line: String,
}

#[tauri::command]
pub async fn pull_docker_image(app: tauri::AppHandle, image: String) -> Result<String, String> {
    use std::process::Stdio;
    use std::io::{BufRead, BufReader};
    use tauri::Emitter;

    let mut cmd = crate::utils::create_command("docker");
    cmd.args(["pull", &image]);
    
    #[cfg(target_os = "windows")]
    {
        use std::os::windows::process::CommandExt;
        cmd.creation_flags(0x08000000);
    }
    
    cmd.stdout(Stdio::piped());
    cmd.stderr(Stdio::piped());

    let mut child = cmd.spawn().map_err(|e| format!("Failed to spawn docker pull: {}", e))?;

    let stdout = child.stdout.take();
    let stderr = child.stderr.take();

    let app_clone = app.clone();
    let image_clone = image.clone();
    if let Some(out) = stdout {
        std::thread::spawn(move || {
            let reader = BufReader::new(out);
            for line in reader.lines().filter_map(Result::ok) {
                let _ = app_clone.emit("docker-pull-log", PullLogPayload {
                    image: image_clone.clone(),
                    line,
                });
            }
        });
    }

    let app_clone = app.clone();
    let image_clone2 = image.clone();
    if let Some(err) = stderr {
        std::thread::spawn(move || {
            let reader = BufReader::new(err);
            for line in reader.lines().filter_map(Result::ok) {
                let _ = app_clone.emit("docker-pull-log", PullLogPayload {
                    image: image_clone2.clone(),
                    line,
                });
            }
        });
    }

    let status = child.wait().map_err(|e| format!("Failed to wait: {}", e))?;
    
    if status.success() {
        Ok("Success".to_string())
    } else {
        Err(format!("Docker pull exited with {}", status))
    }
}

#[tauri::command]
pub async fn check_docker_images(images: Vec<String>) -> std::collections::HashMap<String, bool> {
    let mut result = std::collections::HashMap::new();
    
    for image in images {
        let mut cmd = crate::utils::create_command("docker");
        #[cfg(target_os = "windows")]
        {
            use std::os::windows::process::CommandExt;
            cmd.creation_flags(0x08000000);
        }
        
        let output = cmd.args(["image", "inspect", &image]).output();
        
        let exists = match output {
            Ok(out) => out.status.success(),
            Err(_) => false,
        };
        result.insert(image, exists);
    }
    result
}

#[derive(serde::Serialize, serde::Deserialize)]
pub struct LumineBackup {
    services: serde_json::Value,
    projects: serde_json::Value,
    settings: serde_json::Value,
}

#[tauri::command]
pub fn export_config(app: tauri::AppHandle) -> Result<String, String> {
    let app_dir = app.path().app_data_dir().map_err(|e| e.to_string())?;
    let services_path = app_dir.join("services.json");
    let projects_path = app_dir.join("projects.json");
    let settings_path = app_dir.join("settings.json");
    
    let services_json = std::fs::read_to_string(&services_path).unwrap_or_else(|_| "[]".to_string());
    let projects_json = std::fs::read_to_string(&projects_path).unwrap_or_else(|_| "[]".to_string());
    let settings_json = std::fs::read_to_string(&settings_path).unwrap_or_else(|_| "{}".to_string());
    
    let backup = LumineBackup {
        services: serde_json::from_str(&services_json).unwrap_or(serde_json::json!([])),
        projects: serde_json::from_str(&projects_json).unwrap_or(serde_json::json!([])),
        settings: serde_json::from_str(&settings_json).unwrap_or(serde_json::json!({})),
    };
    
    serde_json::to_string(&backup).map_err(|e| e.to_string())
}

#[tauri::command]
pub fn import_config(
    app: tauri::AppHandle, 
    json_data: String,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
) -> Result<(), String> {
    let backup: LumineBackup = serde_json::from_str(&json_data).map_err(|e| e.to_string())?;
    
    // Stop all first
    let _ = service_manager.stop_all();
    let _ = project_manager.stop_all();
    proxy_manager.stop_proxy();
    
    let app_dir = app.path().app_data_dir().map_err(|e| e.to_string())?;
    
    std::fs::write(app_dir.join("services.json"), serde_json::to_string_pretty(&backup.services).unwrap()).map_err(|e| e.to_string())?;
    std::fs::write(app_dir.join("projects.json"), serde_json::to_string_pretty(&backup.projects).unwrap()).map_err(|e| e.to_string())?;
    std::fs::write(app_dir.join("settings.json"), serde_json::to_string_pretty(&backup.settings).unwrap()).map_err(|e| e.to_string())?;
    
    Ok(())
}

#[tauri::command]
pub fn factory_reset(
    app: tauri::AppHandle,
    service_manager: State<'_, crate::service::ServiceManager>,
    project_manager: State<'_, crate::project_manager::ProjectManager>,
    proxy_manager: State<'_, crate::proxy_manager::ProxyManager>,
) -> Result<(), String> {
    let _ = service_manager.stop_all();
    let _ = project_manager.stop_all();
    proxy_manager.stop_proxy();
    
    let app_dir = app.path().app_data_dir().map_err(|e| e.to_string())?;
    if app_dir.exists() {
        let _ = std::fs::remove_dir_all(&app_dir);
    }
    
    app.exit(0);
    Ok(())
}




#[derive(serde::Serialize, Clone)]
pub struct DockerStat {
    pub cpu: String,
    pub ram: String,
    pub full_name: String,
}

#[tauri::command]
pub async fn get_docker_stats() -> Result<std::collections::HashMap<String, DockerStat>, String> {
    let docker_exe = crate::service::get_docker_executable();
    let mut cmd = tokio::process::Command::new(&docker_exe);
    cmd.arg("stats");
    cmd.arg("--no-stream");
    cmd.arg("--format");
    cmd.arg("{{.Name}}|{{.CPUPerc}}|{{.MemUsage}}");

    let output = tokio::time::timeout(
        std::time::Duration::from_secs(5),
        cmd.output()
    ).await.map_err(|_| "Docker stats timed out".to_string())?
     .map_err(|e| e.to_string())?;
    
    let stdout = String::from_utf8_lossy(&output.stdout);
    
    let mut stats = std::collections::HashMap::new();
    for line in stdout.lines() {
        let parts: Vec<&str> = line.split('|').collect();
        if parts.len() == 3 {
            let full_name = parts[0].to_string();
            let display_name = if full_name.starts_with("lumine_") {
                full_name.replacen("lumine_", "", 1)
            } else {
                full_name.clone()
            };
            stats.insert(display_name, DockerStat {
                cpu: parts[1].to_string(),
                ram: parts[2].split(" /").next().unwrap_or(parts[2]).to_string(),
                full_name,
            });
        }
    }
    
    Ok(stats)
}

#[tauri::command]
pub fn stop_container(name: String) -> Result<(), String> {
    let docker_exe = crate::service::get_docker_executable();
    let output = crate::utils::create_command(&docker_exe)
        .args(["stop", &name])
        .output()
        .map_err(|e| format!("Failed to stop container: {}", e))?;
    if output.status.success() {
        Ok(())
    } else {
        Err(String::from_utf8_lossy(&output.stderr).to_string())
    }
}

#[tauri::command]
pub fn remove_container(name: String) -> Result<(), String> {
    let docker_exe = crate::service::get_docker_executable();
    let output = crate::utils::create_command(&docker_exe)
        .args(["rm", "-f", &name])
        .output()
        .map_err(|e| format!("Failed to remove container: {}", e))?;
    if output.status.success() {
        Ok(())
    } else {
        Err(String::from_utf8_lossy(&output.stderr).to_string())
    }
}

#[tauri::command]
pub fn open_docker_shell(container_name: String) -> Result<(), String> {
    #[cfg(target_os = "windows")]
    {
        let _ = std::process::Command::new("cmd")
            .args([
                "/c",
                "start",
                "cmd",
                "/c",
                &format!("docker exec -it {} sh || docker exec -it {} bash", container_name, container_name)
            ])
            .spawn()
            .map_err(|e| e.to_string())?;
        Ok(())
    }
    #[cfg(not(target_os = "windows"))]
    {
        Err("Terminal not implemented on this OS".to_string())
    }
}

#[tauri::command]
pub fn start_tunnel(id: String, port: u16) -> Result<(), String> {
    let tunnel_name = format!("lumine_tunnel_{}", id);
    
    let _ = crate::utils::create_command("docker")
        .args(["rm", "-f", &tunnel_name])
        .output();
        
    let output = crate::utils::create_command("docker")
        .args([
            "run", "-d", "--rm",
            "--name", &tunnel_name,
            "cloudflare/cloudflared",
            "tunnel", "--url", &format!("http://host.docker.internal:{}", port)
        ])
        .output()
        .map_err(|e| format!("Failed to spawn tunnel: {}", e))?;
        
    if !output.status.success() {
        return Err(String::from_utf8_lossy(&output.stderr).to_string());
    }
    
    Ok(())
}

#[tauri::command]
pub fn get_tunnel_url(id: String) -> Result<Option<String>, String> {
    let tunnel_name = format!("lumine_tunnel_{}", id);
    
    let output = crate::utils::create_command("docker")
        .args(["logs", &tunnel_name])
        .output()
        .map_err(|e| e.to_string())?;
        
    let logs = String::from_utf8_lossy(&output.stderr);
    
    for line in logs.lines() {
        if line.contains("trycloudflare.com") {
            if let Some(start) = line.find("https://") {
                let url = line[start..].split_whitespace().next().unwrap_or("").to_string();
                if !url.is_empty() {
                    return Ok(Some(url));
                }
            }
        }
    }
    
    Ok(None)
}

#[tauri::command]
pub fn stop_tunnel(id: String) -> Result<(), String> {
    let tunnel_name = format!("lumine_tunnel_{}", id);
    let _ = crate::utils::create_command("docker")
        .args(["rm", "-f", &tunnel_name])
        .output();
    Ok(())
}


// ---- Stacks Commands ----

use crate::stack_manager::{StackInfo, StackManager};

#[tauri::command]
pub fn get_stacks(manager: tauri::State<'_, StackManager>) -> Vec<StackInfo> {
    manager.get_all()
}

#[tauri::command]
pub fn add_stack(manager: tauri::State<'_, StackManager>, info: StackInfo) -> Result<(), String> {
    manager.add_stack(info)
}

#[tauri::command]
pub fn edit_stack(manager: tauri::State<'_, StackManager>, id: String, info: StackInfo) -> Result<(), String> {
    manager.edit_stack(&id, info)
}

#[tauri::command]
pub fn delete_stack(manager: tauri::State<'_, StackManager>, id: String) -> Result<(), String> {
    manager.delete_stack(&id)
}

#[tauri::command]
pub fn start_stack(manager: tauri::State<'_, StackManager>, service_manager: tauri::State<'_, crate::service::ServiceManager>, id: String) -> Result<(), String> {
    let stacks = manager.get_all();
    if let Some(stack) = stacks.iter().find(|s| s.id == id) {
        for svc_id in &stack.services {
            let _ = service_manager.start(svc_id);
        }
        Ok(())
    } else {
        Err("Stack not found".to_string())
    }
}

#[tauri::command]
pub fn stop_stack(manager: tauri::State<'_, StackManager>, service_manager: tauri::State<'_, crate::service::ServiceManager>, id: String) -> Result<(), String> {
    let stacks = manager.get_all();
    if let Some(stack) = stacks.iter().find(|s| s.id == id) {
        for svc_id in &stack.services {
            let _ = service_manager.stop(svc_id);
        }
        Ok(())
    } else {
        Err("Stack not found".to_string())
    }
}
