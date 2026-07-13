use serde::{Deserialize, Serialize};
use std::collections::HashMap;
use std::process::{Child, Stdio};
use std::sync::{Arc, Mutex};
#[cfg(target_os = "windows")]
use std::os::windows::process::CommandExt;
use std::io::{BufRead, BufReader};
use std::thread;

#[cfg(target_os = "windows")]
const CREATE_NO_WINDOW: u32 = 0x08000000;

#[derive(Debug, Clone, Serialize, Deserialize, PartialEq)]
pub enum ProjectStatus {
    Running,
    Stopped,
    Error(String),
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ProjectInfo {
    pub id: String,
    pub name: String,
    pub framework: String,
    pub path: String,
    pub start_command: String,
    pub port: u16,
    #[serde(default = "default_status")]
    pub status: ProjectStatus,
    #[serde(default)]
    pub log: Vec<String>,
    #[serde(default)]
    pub env: Option<std::collections::HashMap<String, String>>,
}

fn default_status() -> ProjectStatus {
    ProjectStatus::Stopped
}

pub struct ManagedProject {
    pub info: ProjectInfo,
    pub process: Option<Child>,
    pub log: Arc<Mutex<Vec<String>>>,
}

pub struct ProjectManager {
    projects: Mutex<HashMap<String, ManagedProject>>,
    data_file: std::path::PathBuf,
    workspace_dir: std::path::PathBuf,
}

impl ProjectManager {
    pub fn new(app_data_dir: std::path::PathBuf, workspace_dir: std::path::PathBuf) -> Self {
        let data_file = app_data_dir.join("projects.json");
        
        // Create custom bashrc for beautiful terminal prompt and modern features
        let bashrc_path = app_data_dir.join("lumine_bashrc");
        let bashrc_content = r#"
if [ -f ~/.bashrc ]; then
    source ~/.bashrc
fi

# Modern Autocomplete & History Search Bindings
if [ -n "$BASH_VERSION" ]; then
    bind 'set completion-ignore-case on'
    bind 'set show-all-if-ambiguous on'
    bind 'set colored-stats on'
    bind 'set visible-stats on'
    bind '"\e[A": history-search-backward'
    bind '"\e[B": history-search-forward'
    
    # Try to load bash_completion if available
    if [ -f /etc/bash_completion ]; then
        . /etc/bash_completion
    elif [ -f /usr/share/bash-completion/bash_completion ]; then
        . /usr/share/bash-completion/bash_completion
    fi

    # Persistent history configuration
    export HISTFILE=/tmp/lumine_bash_history
    export HISTSIZE=10000
    export HISTFILESIZE=10000
    shopt -s histappend
    export PROMPT_COMMAND="history -a; history -c; history -r; $PROMPT_COMMAND"
fi

export PS1="\[\e[38;2;140;170;238m\]╭─\[\e[38;2;186;187;241m\]\u\[\e[38;2;172;176;190m\]@\[\e[38;2;244;184;228m\]lumine \[\e[38;2;140;170;238m\]in \[\e[38;2;186;187;241m\]\w\n\[\e[38;2;140;170;238m\]╰─❯ \[\e[0m\]"
"#;
        if !bashrc_path.exists() {
            let _ = std::fs::write(&bashrc_path, bashrc_content);
        }

        // Ensure history file exists
        let hist_path = app_data_dir.join("lumine_bash_history");
        if !hist_path.exists() {
            let _ = std::fs::write(&hist_path, "");
        }

        let manager = ProjectManager {
            projects: Mutex::new(HashMap::new()),
            data_file,
            workspace_dir,
        };
        manager.load_data();
        manager
    }

    pub fn get_workspace_dir(&self) -> String {
        self.workspace_dir.to_string_lossy().to_string()
    }

    fn load_data(&self) {
        if let Ok(data) = std::fs::read_to_string(&self.data_file) {
            if let Ok(mut saved_projects) = serde_json::from_str::<Vec<ProjectInfo>>(&data) {
                let mut map = self.projects.lock().unwrap();
                for proj in saved_projects.iter_mut() {
                    proj.status = ProjectStatus::Stopped; // Always default to stopped on boot
                    map.insert(proj.id.clone(), ManagedProject {
                        info: proj.clone(),
                        process: None,
                        log: Arc::new(Mutex::new(proj.log.clone())),
                    });
                }
            }
        }
    }

    fn save_data(&self) {
        let map = self.projects.lock().unwrap();
        let list: Vec<ProjectInfo> = map.values().map(|p| {
            let mut info = p.info.clone();
            info.log = p.log.lock().unwrap().clone();
            info
        }).collect();
        if let Ok(data) = serde_json::to_string_pretty(&list) {
            let _ = std::fs::write(&self.data_file, data);
        }
    }

    pub fn get_all(&self) -> Vec<ProjectInfo> {
        let mut map = self.projects.lock().unwrap();
        let mut result = Vec::new();

        for (_, mp) in map.iter_mut() {
            // Check if process exited
            if let Some(mut child) = mp.process.take() {
                match child.try_wait() {
                    Ok(Some(status)) => {
                        mp.info.status = if status.success() {
                            ProjectStatus::Stopped
                        } else {
                            ProjectStatus::Error(format!("Exited with {}", status))
                        };
                    }
                    Ok(None) => {
                        mp.process = Some(child);
                    }
                    Err(e) => {
                        mp.info.status = ProjectStatus::Error(format!("Process error: {}", e));
                    }
                }
            }

            let mut info = mp.info.clone();
            info.log = mp.log.lock().unwrap().clone();
            result.push(info);
        }
        
        result.sort_by(|a, b| a.name.cmp(&b.name));
        result
    }

    pub fn add_project(&self, mut project: ProjectInfo) -> Result<(), String> {
        let mut map = self.projects.lock().unwrap();
        if map.contains_key(&project.id) {
            return Err("Project ID already exists".into());
        }
        project.status = ProjectStatus::Stopped;
        
        map.insert(project.id.clone(), ManagedProject {
            info: project.clone(),
            process: None,
            log: Arc::new(Mutex::new(Vec::new())),
        });
        drop(map);
        self.save_data();
        Ok(())
    }

    pub fn clear_log(&self, id: &str) -> Result<(), String> {
        let map = self.projects.lock().unwrap();
        if let Some(mp) = map.get(id) {
            let mut logs = mp.log.lock().unwrap();
            logs.clear();
            drop(logs);
            drop(map);
            self.save_data();
            Ok(())
        } else {
            Err("Project not found".into())
        }
    }

    pub fn edit_project(&self, id: &str, new_command: String) -> Result<(), String> {
        let mut map = self.projects.lock().unwrap();
        if let Some(mp) = map.get_mut(id) {
            if matches!(mp.info.status, ProjectStatus::Running) {
                return Err("Cannot edit project while it is running".into());
            }
            mp.info.start_command = new_command;
        } else {
            return Err("Project not found".into());
        }
        drop(map);
        self.save_data();
        Ok(())
    }

    pub fn delete_project(&self, id: &str, delete_folder: bool) -> Result<(), String> {
        self.stop(id)?;
        let mut map = self.projects.lock().unwrap();
        if let Some(project) = map.remove(id) {
            drop(map);
            self.save_data();
            
            if delete_folder {
                if std::path::Path::new(&project.info.path).exists() {
                    let _ = std::fs::remove_dir_all(&project.info.path);
                }
            }
            Ok(())
        } else {
            Err("Project not found".into())
        }
    }

    pub fn stop(&self, id: &str) -> Result<(), String> {
        let (child_opt, container_name) = {
            let mut map = self.projects.lock().unwrap();
            let mp = map.get_mut(id).ok_or_else(|| "Project not found".to_string())?;
            
            let mut c_name = None;
            if let Ok(args) = serde_json::from_str::<Vec<String>>(&mp.info.start_command) {
                if args.get(0).map(String::as_str) == Some("docker") {
                    if let Some(pos) = args.iter().position(|x| x == "--name") {
                        if let Some(name) = args.get(pos + 1) {
                            if name.starts_with("lumine_proj_") {
                                c_name = Some(name.clone());
                            }
                        }
                    }
                }
            }
            (mp.process.take(), c_name)
        }; // Lock dropped

        if let Some(mut child) = child_opt {
            let _ = child.kill();
            let _ = child.wait();
        }

        if let Some(name) = container_name {
            #[cfg(target_os = "windows")]
            let _ = crate::utils::create_command("docker").args(["rm", "-f", &name]).creation_flags(CREATE_NO_WINDOW).output();
            #[cfg(not(target_os = "windows"))]
            let _ = crate::utils::create_command("docker").args(["rm", "-f", &name]).output();
        }

        let mut map = self.projects.lock().unwrap();
        if let Some(mp) = map.get_mut(id) {
            mp.info.status = ProjectStatus::Stopped;
            drop(map);
            self.save_data();
            Ok(())
        } else {
            Err("Project not found".into())
        }
    }

    pub fn start(&self, id: &str) -> Result<(), String> {
        let (path, start_command, container_name) = {
            let map = self.projects.lock().unwrap();
            let mp = map.get(id).ok_or_else(|| "Project not found".to_string())?;
            if mp.process.is_some() {
                return Ok(());
            }
            let mut c_name = None;
            if let Ok(args) = serde_json::from_str::<Vec<String>>(&mp.info.start_command) {
                if args.get(0).map(String::as_str) == Some("docker") {
                    if let Some(pos) = args.iter().position(|x| x == "--name") {
                        if let Some(name) = args.get(pos + 1) {
                            if name.starts_with("lumine_proj_") {
                                c_name = Some(name.clone());
                            }
                        }
                    }
                }
            }
            (mp.info.path.clone(), mp.info.start_command.clone(), c_name)
        };

        if let Some(name) = container_name {
            #[cfg(target_os = "windows")]
            let _ = crate::utils::create_command("docker").args(["rm", "-f", &name]).creation_flags(CREATE_NO_WINDOW).output();
            #[cfg(not(target_os = "windows"))]
            let _ = crate::utils::create_command("docker").args(["rm", "-f", &name]).output();
        }

        let mut map = self.projects.lock().unwrap();
        if let Some(mp) = map.get_mut(id) {
            if mp.process.is_some() {
                return Ok(()); // Already running
            }

            let mut cmd;
            if let Ok(args) = serde_json::from_str::<Vec<String>>(&start_command) {
                cmd = crate::utils::create_command(&args[0]);
                if args.len() > 1 {
                    cmd.args(&args[1..]);
                }
                #[cfg(target_os = "windows")]
                cmd.creation_flags(CREATE_NO_WINDOW);
            } else {
                #[cfg(target_os = "windows")]
                {
                    cmd = crate::utils::create_command("cmd");
                    cmd.args(["/C", &start_command]);
                    cmd.creation_flags(CREATE_NO_WINDOW);
                }
                #[cfg(not(target_os = "windows"))]
                {
                    cmd = crate::utils::create_command("sh");
                    cmd.args(["-c", &start_command]);
                }
            }

            cmd.current_dir(&path);
            cmd.stdout(Stdio::piped());
            cmd.stderr(Stdio::piped());

            match cmd.spawn() {
                Ok(mut child) => {
                    let stdout = child.stdout.take();
                    let stderr = child.stderr.take();
                    
                    mp.process = Some(child);
                    mp.info.status = ProjectStatus::Running;
                    
                    let log_arc = mp.log.clone();
                    
                    if let Some(out) = stdout {
                        let log_arc_clone = log_arc.clone();
                        thread::spawn(move || {
                            let reader = BufReader::new(out);
                            for line in reader.lines().filter_map(Result::ok) {
                                let mut logs = log_arc_clone.lock().unwrap();
                                logs.push(line);
                                if logs.len() > 100 {
                                    logs.remove(0);
                                }
                            }
                        });
                    }

                    if let Some(err) = stderr {
                        let log_arc_clone = log_arc.clone();
                        thread::spawn(move || {
                            let reader = BufReader::new(err);
                            for line in reader.lines().filter_map(Result::ok) {
                                let mut logs = log_arc_clone.lock().unwrap();
                                logs.push(line);
                                if logs.len() > 100 {
                                    logs.remove(0);
                                }
                            }
                        });
                    }
                    
                    drop(map);
                    self.save_data();
                    Ok(())
                }
                Err(e) => {
                    mp.info.status = ProjectStatus::Error(format!("Failed to start: {}", e));
                    drop(map);
                    self.save_data();
                    Err(e.to_string())
                }
            }
        } else {
            Err("Project not found".into())
        }
    }

    pub fn open_terminal(&self, id: &str, settings: &crate::settings_manager::AppSettings, image_override: Option<String>) -> Result<(), String> {
        let map = self.projects.lock().unwrap();
        if let Some(mp) = map.get(id) {
            let is_running = mp.info.status == ProjectStatus::Running;
            let path = mp.info.path.clone();
            
            let mut is_docker = false;
            let mut container_name = None;
            let mut volume_args = Vec::new();
            let mut workdir_args = Vec::new();
            let mut image_name = None;

            if let Ok(args) = serde_json::from_str::<Vec<String>>(&mp.info.start_command) {
                if args.get(0).map(String::as_str) == Some("docker") {
                    is_docker = true;
                    let mut iter = args.iter().peekable();
                    while let Some(arg) = iter.next() {
                        if arg == "--name" {
                            if let Some(name) = iter.next() {
                                container_name = Some(name.clone());
                            }
                        } else if arg == "-v" || arg == "--volume" {
                            if let Some(vol) = iter.next() {
                                volume_args.push("-v".to_string());
                                volume_args.push(vol.clone());
                            }
                        } else if arg == "-w" || arg == "--workdir" {
                            if let Some(wd) = iter.next() {
                                workdir_args.push("-w".to_string());
                                workdir_args.push(wd.clone());
                            }
                        } else if arg == "-p" || arg == "--publish" || arg == "-e" || arg == "--env" || arg == "--network" {
                            let _ = iter.next(); // Skip the value
                        } else if arg.starts_with("-") {
                            // Ignore other boolean flags (e.g., -d, -it, --rm)
                        } else if arg != "docker" && arg != "run" {
                            // First positional argument after `docker run` is the image
                            image_name = Some(arg.clone());
                            break; // Stop parsing after image name
                        }
                    }
                }
            }

            let term = &settings.terminal_emulator;
            
            // Helper function to spawn the terminal
            let spawn_terminal = |cmd_str: &str, ps_str: &str, bash_str: &str, start_dir: &str| {
                #[cfg(target_os = "windows")]
                {
                    use std::os::windows::process::CommandExt;
                    const CREATE_NEW_CONSOLE: u32 = 0x00000010;
                    
                    let mut cmd;
                    let lower_term = term.to_lowercase();
                    
                    if lower_term == "wt" {
                        // Windows Terminal (uses PowerShell)
                        cmd = crate::utils::create_command("wt.exe");
                        cmd.args(["-w", "0", "nt", "powershell.exe", "-NoExit", "-Command", ps_str]);
                    } else if lower_term == "powershell" || lower_term == "pwsh" {
                        // PowerShell
                        cmd = crate::utils::create_command(format!("{}.exe", lower_term));
                        cmd.creation_flags(CREATE_NEW_CONSOLE);
                        cmd.args(["-NoExit", "-Command", ps_str]);
                    } else if lower_term == "bash" || lower_term == "zsh" || lower_term == "fish" {
                        // Unix Shells
                        cmd = crate::utils::create_command(format!("{}.exe", lower_term));
                        cmd.creation_flags(CREATE_NEW_CONSOLE);
                        cmd.args(["-c", bash_str]);
                    } else if lower_term == "alacritty" || lower_term == "kitty" || lower_term == "wezterm" {
                        // Modern cross-platform terminals (using bash as inner shell)
                        cmd = crate::utils::create_command(format!("{}.exe", lower_term));
                        cmd.args(["-e", "bash", "-c", bash_str]);
                    } else {
                        // Default CMD
                        cmd = crate::utils::create_command("cmd.exe");
                        cmd.creation_flags(CREATE_NEW_CONSOLE);
                        cmd.args(["/K", cmd_str]);
                    }
                    if !start_dir.is_empty() {
                        cmd.current_dir(start_dir);
                    }
                    let _ = cmd.spawn();
                }
            };

            let cmd_str: String;
            let ps_str: String;
            let bash_str: String;
            let display_name = mp.info.name.clone();

            if is_docker {
                if is_running {
                    if let Some(name) = container_name {
                        cmd_str = format!("title Lumine Terminal - {0} & echo ======================================== & echo   Lumine Container Terminal & echo   Status: Attached to Running Server & echo   Container: {0} & echo ======================================== & echo. & docker exec -it {0} bash || docker exec -it {0} sh", name);
                        ps_str = format!("$host.ui.RawUI.WindowTitle='Lumine Terminal - {0}'\nWrite-Host '========================================'\nWrite-Host '  Lumine Container Terminal'\nWrite-Host '  Status: Attached to Running Server'\nWrite-Host '  Container: {0}'\nWrite-Host '========================================'\nWrite-Host ''\ndocker exec -it {0} bash\nif (!$?) {{ docker exec -it {0} sh }}", name);
                        bash_str = format!("echo -en \"\\033]0;Lumine Terminal - {0}\\a\"; echo \"========================================\"; echo \"  Lumine Container Terminal\"; echo \"  Status: Attached to Running Server\"; echo \"  Container: {0}\"; echo \"========================================\"; echo \"\"; docker exec -it {0} bash || docker exec -it {0} sh", name);
                        spawn_terminal(&cmd_str, &ps_str, &bash_str, "");
                        Ok(())
                    } else {
                        Err("This project lacks a --name argument.".into())
                    }
                } else {
                    let final_image = image_override.unwrap_or_else(|| image_name.clone().unwrap_or_default());
                    if !final_image.is_empty() {
                        let vols = volume_args.join(" ");
                        let wds = workdir_args.join(" ");
                        let safe_name = display_name.replace(" ", "_").replace("-", "_").to_lowercase();
                        let term_name = format!("lumine_term_{}", safe_name);
                        
                        let app_data_dir = self.data_file.parent().unwrap().to_string_lossy();
                        let rc_args = format!("-v \"{0}\\lumine_bashrc\":/tmp/lumine_bashrc:ro -v \"{0}\\lumine_bash_history\":/tmp/lumine_bash_history -e ENV=/tmp/lumine_bashrc", app_data_dir);
                        
                        cmd_str = format!("title Lumine Terminal - {0} & echo ======================================== & echo   Lumine Container Terminal & echo   Status: Disposable Session (Server is Stopped) & echo   Image: {1} & echo ======================================== & echo. & docker exec -it {4} bash --rcfile /tmp/lumine_bashrc || (docker rm -f {4} >nul 2>&1 & docker run -it --rm --name {4} {5} {2} {3} {1} bash --rcfile /tmp/lumine_bashrc || docker run -it --rm --name {4} {5} {2} {3} {1} sh)", display_name, final_image, vols, wds, term_name, rc_args);
                        ps_str = format!("$host.ui.RawUI.WindowTitle='Lumine Terminal - {0}'\nWrite-Host '========================================'\nWrite-Host '  Lumine Container Terminal'\nWrite-Host '  Status: Disposable Session (Server is Stopped)'\nWrite-Host '  Image: {1}'\nWrite-Host '========================================'\nWrite-Host ''\ndocker exec -it {4} bash --rcfile /tmp/lumine_bashrc\nif (!$?) {{\n  docker rm -f {4} 2>$null | Out-Null\n  docker run -it --rm --name {4} {5} {2} {3} {1} bash --rcfile /tmp/lumine_bashrc\n  if (!$?) {{ docker run -it --rm --name {4} {5} {2} {3} {1} sh }}\n}}", display_name, final_image, vols, wds, term_name, rc_args);
                        bash_str = format!("echo -en \"\\033]0;Lumine Terminal - {0}\\a\"; echo \"========================================\"; echo \"  Lumine Container Terminal\"; echo \"  Status: Disposable Session (Server is Stopped)\"; echo \"  Image: {1}\"; echo \"========================================\"; echo \"\"; docker exec -it {4} bash --rcfile /tmp/lumine_bashrc || (docker rm -f {4} >/dev/null 2>&1; docker run -it --rm --name {4} {5} {2} {3} {1} bash --rcfile /tmp/lumine_bashrc || docker run -it --rm --name {4} {5} {2} {3} {1} sh)", display_name, final_image, vols, wds, term_name, rc_args);
                        spawn_terminal(&cmd_str, &ps_str, &bash_str, "");
                        Ok(())
                    } else {
                        Err("Could not detect docker image name from start command.".into())
                    }
                }
            } else {
                // Non-docker project (Local)
                let drive = path.chars().next().unwrap_or('C');
                cmd_str = format!("title Lumine Terminal - {0} & echo ======================================== & echo   Lumine Local Terminal & echo   Path: {1} & echo ======================================== & echo. & {2}: & cd \"{1}\"", display_name, path, drive);
                ps_str = format!("$host.ui.RawUI.WindowTitle='Lumine Terminal - {0}'\nWrite-Host '========================================'\nWrite-Host '  Lumine Local Terminal'\nWrite-Host '  Path: {1}'\nWrite-Host '========================================'\nWrite-Host ''\nSet-Location -LiteralPath '{1}'", display_name, path);
                bash_str = format!("echo -en \"\\033]0;Lumine Terminal - {0}\\a\"; echo \"========================================\"; echo \"  Lumine Local Terminal\"; echo \"  Path: {1}\"; echo \"========================================\"; echo \"\"; cd '{1}'; exec bash", display_name, path);
                spawn_terminal(&cmd_str, &ps_str, &bash_str, &path);
                Ok(())
            }
        } else {
            Err("Project not found".into())
        }
    }
    
    pub fn stop_all(&self) {
        let ids: Vec<String> = {
            let map = self.projects.lock().unwrap();
            map.keys().cloned().collect()
        };
        for id in ids {
            let _ = self.stop(&id);
        }
    }
}
