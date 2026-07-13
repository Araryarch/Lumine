# Lumine - Project State

This document serves as the internal context and state tracking for the Lumine project.
Lumine is a modern, fast, universal background service manager built with Tauri (Rust) and React.

## Architectural State

1. **BYOB (Bring Your Own Binary) - 100% DONE**
   - We completely abandoned the `bin/` approach and `setup.js`.
   - Users can now dynamically add custom services, languages, and databases via the "Add Service" modal.
   - The backend successfully parses user arguments, spawns executables natively, grabs stdout/stderr, and saves state to `services.json` in the Tauri AppData dir.
   - We patched the zombie process bug: `ServiceManager` now correctly uses `try_wait()` to mark services as "Stopped" if they exit externally.
   - DEADLOCK BUG FIXED: `start()` no longer calls `get_executable_path()` (which would deadlock by trying to re-lock the mutex). It now reads `executable_path` directly from the already-locked data.

2. **Hosts Management (`/etc/hosts`) - 100% DONE**
   - **Backend**: `hosts_manager.rs` is fully implemented. It can read, add, toggle, and delete host entries. It successfully requires administrative privileges.
   - **Frontend**: `HostsView.tsx` is fully wired to the backend via Tauri `invoke()` calls. It loads hosts on mount, has an Add Host modal, inline toggle switches, context menu delete, search filtering, error banners, and loading states.

3. **MkCert (Local SSL) - 100% DONE**
   - **Backend**: `mkcert_manager.rs` implements `check_installed()`, `install_root_ca()`, `generate_cert()`, `get_certs()`, and `delete_cert()`. Certs are stored in `{app_data_dir}/certs/` and metadata in `certs.json`.
   - **Frontend**: `MkCertView.tsx` shows mkcert install status, certificate table, generate modal, install CA button, and proper error/success banners.
   - **Commands**: 5 new Tauri commands registered: `check_mkcert`, `install_root_ca`, `generate_cert`, `get_certs`, `delete_cert`.

4. **Docker Integration - 100% DONE**
   - **Backend**: Implemented `get_docker_executable` with fallback paths for Windows (`C:\Program Files\Docker\Docker\resources\bin\docker.exe`).
   - **Service Manager**: Services tagged with `runner: "docker"` are successfully spawned. Stopping a Docker service explicitly executes `docker stop lumine_{id}` to prevent background zombie containers.
   
5. **Persistent Log System - 100% DONE**
   - **Backend**: `ServiceManager` automatically writes real-time logs to `logs/{id}.log` in the Tauri AppData directory via `std::fs::OpenOptions`.
   - **Frontend**: LogViewer features a "Copy All" button (clipboard API) and an "Open Full Log File" button which triggers the Tauri `open_log_file` command.
   
6. **Automated CI/CD - 100% DONE**
   - **GitHub Actions**: `.github/workflows/release.yml` uses `actions/upload-artifact@v4` to upload compiled `.exe`, `.msi`, `.dmg`, and `.deb` files directly as Action Artifacts on every push to `main`.

7. **Hidden Internal Proxy (HTTPS) - 100% DONE**
   - **Backend**: Implemented `proxy_manager.rs` that maintains a background `nginx:alpine` docker container on port 80 and 443.
   - **Dynamic Routing**: Dynamically generates `nginx.conf` linking `.test` domains and `localhost/adminpanel` subpaths to the correct service ports.
   - **SSL Support**: Auto-detects if a MkCert certificate exists for the domain, mounts it, generates port 443 SSL configs, and injects HTTP to HTTPS redirects.
   - **UI**: Added clickable URLs in the Hosts page (dynamic http/https based on cert existence) and an inline SSL generation button (Shield icon).

8. **Clean Exit Architecture - 100% DONE**
   - **Tauri Hook**: Caught `tauri::RunEvent::Exit` in `lib.rs` to trigger background cleanup.
   - **Service Termination**: `ServiceManager::stop_all()` safely halts all active Docker/binary processes before memory is released.
   - **Proxy Teardown**: `ProxyManager::stop_proxy()` executes `docker rm -f lumine_internal_proxy` to ensure port 80/443 are properly freed.

## What's Next / Pending
- **Projects Management**: The UI and logic for adding/managing user projects (virtual hosts mappings to local directories) are not fully integrated yet.

## Critical Notes for Agents

- **Tool Usage**: When writing or updating files, DO NOT use bash `cat` or `echo`. Use proper tools like `replace_file_content` or `write_to_file`.
- **Rust Backend**: All services are managed via `std::process::Child`. Since we don't have blocking `.wait()` loops to prevent UI freezing, we poll the `Child` state via `try_wait()` during `get_all()` and `get()` calls to reflect accurate status on the frontend.
- **Frontend Architecture**: React + Tailwind CSS v4 + Lucide React. Dark mode is the primary theme (`#1e1e2e` and `#181825` background colors).
- **TypeScript Config**: We strictly follow standard types. `ServiceConfig` inside `types.ts` has `executablePath`, `arguments`, and `serviceType`. DO NOT mix it up with older properties like `customPath`.
- **Log System**: Backend sends `log: Vec<String>` (capped at 100 lines FIFO). Frontend `LogViewer` accepts `log: string[]` and renders with line numbers, error highlighting.
- **StatusDot**: The `StatusDot` component is now integrated into `SidebarItem` showing live service status (green/gray/red/yellow pulse).
