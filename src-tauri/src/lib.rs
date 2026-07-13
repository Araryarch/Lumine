mod commands;
mod service;
mod hosts_manager;
mod mkcert_manager;
mod settings_manager;
mod setup_manager;
pub mod stack_manager;
pub mod project_manager;
pub mod proxy_manager;
pub mod utils;

#[cfg_attr(mobile, tauri::mobile_entry_point)]
pub fn run() {
  use tauri::Manager;
  tauri::Builder::default()
    .plugin(tauri_plugin_autostart::init(
      tauri_plugin_autostart::MacosLauncher::LaunchAgent,
      Some(vec![]),
    ))
    .on_window_event(|window, event| match event {
        tauri::WindowEvent::CloseRequested { api, .. } => {
            let state = window.state::<settings_manager::SettingsManager>();
            if let Ok(settings) = state.get_settings() {
                if settings.minimize_to_tray {
                    let _ = window.hide();
                    api.prevent_close();
                }
            }
        }
        _ => {}
    })
    .setup(|app| {
      use tauri::Manager;
      
      let app_data_dir = app.path().app_data_dir().unwrap_or_else(|_| std::env::current_dir().unwrap());
      
      if let Err(e) = setup_manager::initialize_default_packages(app) {
        log::error!("Failed to initialize default packages: {}", e);
      }
      let proxy_manager = proxy_manager::ProxyManager::new(app_data_dir.clone());

      let service_mgr = service::ServiceManager::new(app_data_dir.clone());
      let hosts_mgr = hosts_manager::HostsManager::new(app.handle());
      let stack_mgr = stack_manager::StackManager::new(app_data_dir.clone());
      
      app.manage(mkcert_manager::MkCertManager::new(app_data_dir.clone()));
      app.manage(settings_manager::SettingsManager::new(app_data_dir.clone()));
      
      let workspace_dir = {
          #[cfg(target_os = "windows")]
          {
              let windir = std::env::var("USERPROFILE").unwrap_or_else(|_| "C:\\Users\\Default".to_string());
              std::path::Path::new(&windir).join("Documents\\Lumine\\www")
          }
          #[cfg(not(target_os = "windows"))]
          {
              let home = std::env::var("HOME").unwrap_or_else(|_| "/".to_string());
              std::path::Path::new(&home).join("Lumine/www")
          }
      };
      let _ = std::fs::create_dir_all(&workspace_dir);
      
      let proj_mgr = project_manager::ProjectManager::new(app_data_dir.clone(), workspace_dir);
      
      app.manage(proxy_manager.clone());
      app.manage(service_mgr);
      app.manage(hosts_mgr);
      app.manage(stack_mgr);
      app.manage(proj_mgr);
      
      app.handle().plugin(tauri_plugin_dialog::init())?;
      app.handle().plugin(tauri_plugin_fs::init())?;
      
      if cfg!(debug_assertions) {
        app.handle().plugin(
          tauri_plugin_log::Builder::default()
            .level(log::LevelFilter::Info)
            .build(),
        )?;
      }

      // Spawn background thread for Docker operations so we don't block the UI
      let app_handle = app.handle().clone();
      std::thread::spawn(move || {
          let proxy_manager = app_handle.state::<proxy_manager::ProxyManager>();
          let hosts_mgr = app_handle.state::<hosts_manager::HostsManager>();
          let service_mgr = app_handle.state::<service::ServiceManager>();
          let proj_mgr = app_handle.state::<project_manager::ProjectManager>();
          
          proxy_manager.ensure_proxy_running();

          if let Ok(hosts) = hosts_mgr.get_hosts() {
              let services = service_mgr.get_all();
              let projects = proj_mgr.get_all();
              let _ = proxy_manager.sync_config(&hosts, &services, &projects);
          }

          // Auto-start Services
          let settings_mgr = app_handle.state::<settings_manager::SettingsManager>();
          let global_auto_start = settings_mgr.get_settings().map(|s| s.auto_start_services).unwrap_or(false);
          let docker_running = commands::check_docker();
          
          for service in service_mgr.get_all() {
              if service.config.runner == "docker" && !docker_running {
                  continue;
              }
              if global_auto_start || service.config.auto_start == Some(true) {
                  let _ = service_mgr.start(&service.id);
              }
          }
      });

      // System Tray
      use tauri::menu::{Menu, MenuItem};
      use tauri::tray::{TrayIconBuilder, MouseButton, MouseButtonState, TrayIconEvent};

      let show_i = MenuItem::with_id(app, "show", "Show Lumine", true, None::<&str>)?;
      let quit_i = MenuItem::with_id(app, "quit", "Quit", true, None::<&str>)?;
      let menu = Menu::with_items(app, &[&show_i, &quit_i])?;

      let _tray = TrayIconBuilder::new()
          .icon(app.default_window_icon().unwrap().clone())
          .menu(&menu)
          .show_menu_on_left_click(false)
          .on_menu_event(|app, event| match event.id.as_ref() {
              "quit" => {
                  app.exit(0);
              }
              "show" => {
                  if let Some(window) = app.get_webview_window("main") {
                      let _ = window.show();
                      let _ = window.set_focus();
                  }
              }
              _ => {}
          })
          .on_tray_icon_event(|tray, event| {
              if let TrayIconEvent::Click {
                  button: MouseButton::Left,
                  button_state: MouseButtonState::Up,
                  ..
              } = event
              {
                  let app = tray.app_handle();
                  if let Some(window) = app.get_webview_window("main") {
                      let _ = window.show();
                      let _ = window.set_focus();
                  }
              }
          })
          .build(app)?;

      Ok(())
    })
    .invoke_handler(tauri::generate_handler![
      commands::get_services,
      commands::get_service,
      commands::start_service,
      commands::stop_service,
      commands::restart_service,
      commands::clear_service_log,
      commands::set_port,
      commands::check_port,
      commands::start_all,
      commands::stop_all,
      commands::add_service,
      commands::edit_service,
      commands::delete_service,
      commands::get_hosts,
      commands::add_host,
      commands::edit_host,
      commands::toggle_host,
      commands::delete_host,
      commands::open_hosts_file,
      commands::fetch_docker_tags,
      commands::check_mkcert,
      commands::download_mkcert,
      commands::install_root_ca,
      commands::generate_cert,
      commands::get_certs,
      commands::delete_cert,
      commands::get_settings,
      commands::save_settings,
      commands::check_docker,
      commands::exit_app,
      commands::open_log_file,
      commands::get_projects,
      commands::get_workspace_dir,
      commands::add_project,
      commands::edit_project,
      commands::start_project,
      commands::stop_project,
      commands::delete_project,
      commands::clear_project_log,
      commands::open_project_terminal,
      commands::create_new_project,
      commands::open_url,
      commands::open_in_explorer,
      commands::open_service_terminal,
      commands::open_language_editor,
      commands::open_in_editor,
      commands::get_available_tools,
      commands::pull_docker_image,
      commands::check_docker_images,
      commands::get_language_container_status,
      commands::stop_language_container,
      commands::export_config,
      commands::import_config,
      commands::factory_reset,
      commands::get_docker_stats,
      commands::stop_container,
      commands::remove_container,
      commands::open_docker_shell,
      commands::start_tunnel,
      commands::get_tunnel_url,
      commands::stop_tunnel,
      commands::get_stacks,
      commands::add_stack,
      commands::edit_stack,
      commands::delete_stack,
      commands::start_stack,
      commands::stop_stack,
    ])
    .build(tauri::generate_context!())
    .expect("error while running tauri application")
    .run(|app_handle, event| match event {
        tauri::RunEvent::Exit => {
            let proxy_manager = app_handle.state::<crate::proxy_manager::ProxyManager>();
            proxy_manager.stop_proxy();
            
            let service_manager = app_handle.state::<crate::service::ServiceManager>();
            service_manager.stop_all();

            let project_manager = app_handle.state::<crate::project_manager::ProjectManager>();
            project_manager.stop_all();

            let hosts_manager = app_handle.state::<crate::hosts_manager::HostsManager>();
            let _ = hosts_manager.disable_all_hosts();

            // Forcefully cleanup ANY zombie containers left behind (e.g. from disposable terminals)
            #[cfg(target_os = "windows")]
            {
                use std::os::windows::process::CommandExt;
                if let Ok(output) = crate::utils::create_command("docker")
                    .args(["ps", "-a", "-q", "-f", "name=lumine_"])
                    .creation_flags(0x08000000)
                    .output() 
                {
                    let ids = String::from_utf8_lossy(&output.stdout);
                    for id in ids.lines() {
                        let id = id.trim();
                        if !id.is_empty() {
                            let _ = crate::utils::create_command("docker")
                                .args(["rm", "-f", id])
                                .creation_flags(0x08000000)
                                .output();
                        }
                    }
                }
            }
            #[cfg(not(target_os = "windows"))]
            {
                if let Ok(output) = crate::utils::create_command("docker")
                    .args(["ps", "-a", "-q", "-f", "name=lumine_"])
                    .output() 
                {
                    let ids = String::from_utf8_lossy(&output.stdout);
                    for id in ids.lines() {
                        let id = id.trim();
                        if !id.is_empty() {
                            let _ = crate::utils::create_command("docker")
                                .args(["rm", "-f", id])
                                .output();
                        }
                    }
                }
            }
        }
        _ => {}
    });
}
