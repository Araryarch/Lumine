# Changelog

All notable changes to Lumine will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/), and this project adheres to [Semantic Versioning](https://semver.org/).

## [0.1.0] - 2026-07-13

### Added

- **17 pre-configured services** — Nginx, Apache, Caddy, PHP-FPM, Node.js, Bun, Python, Go, Deno, MySQL, MariaDB, PostgreSQL, MongoDB, Redis, phpMyAdmin, Adminer
- **One-click start/stop** — individual or batch, with live status indicators
- **Per-service auto-start** — configure individual services to start on boot independently of the global setting
- **Full-page service editor** — add and edit services with an inline editor
- **Live logs** — real-time terminal output from every running service with log rotation
- **Built-in reverse proxy** — Nginx on ports 80/443, auto-routes domains to services with `^~` location matching
- **Instant HTTPS** — generate and trust SSL certs for any domain via mkcert
- **Hosts editor** — visual `/etc/hosts` management with enable/disable toggles, IPv6 support, and auto-sync
- **Cloudflare Tunnel** — expose local services to the internet via `cloudflared`
- **Docker stats dashboard** — live CPU and memory monitoring with color-coded bar graphs
- **Stop & remove containers** — manage containers directly from the stats view
- **Image pulling UI** — pre-download Docker images from the dashboard (pull all or selected)
- **Docker Hub browser** — search and pick image versions when adding services
- **Stacks** — group services into named stacks and start/stop them together
- **10 pre-configured stacks** — XAMPP Migration, MERN Stack, Laravel Local, LEMP, LAMP, Python Data, Go Microservices, Spring Boot, Ruby on Rails, Rust Backend
- **Custom stacks** — create, edit, and delete your own stack configurations
- **Project scaffolding** — create Laravel, React, Vue, or Next.js projects from a wizard
- **Project terminals** — open disposable Docker terminals with custom bashrc and themed prompts
- **System tray** — runs quietly in the background, start on boot
- **Factory reset** — reset all app data, services, and settings to defaults
- **Import/Export config** — backup and restore service configurations
- **Settings panel** — configure preferred terminal, code editor, and file explorer with auto-detection
- **Dark mode UI** — built with React 19, Tailwind CSS v4, Space Grotesk + JetBrains Mono
- **ErrorBoundary** — graceful crash recovery with reload option
- Docker Desktop auto-start on Windows when Docker is not running
- Zombie Docker container cleanup on application exit
- Document root mounted into Nginx reverse proxy container for serving static files

### Technical

- **Backend:** Rust + Tauri v2 with async I/O (`tokio`)
- **Frontend:** React 19 + TypeScript + Vite + Tailwind CSS v4
- **State:** Zustand store for shared services, projects, and docker stats
- **Cross-platform:** Windows (.exe, .msi), macOS (.dmg), Linux (.deb)
