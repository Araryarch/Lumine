<div align="center">
<img width="150" height="150" alt="Lumine Logo" src="https://github.com/user-attachments/assets/156e848d-afd0-4849-964f-5850e3e138f6" />
  <h1>Lumine</h1>
  <p><strong>Your Local Dev Stack, Reimagined.</strong></p>
  <p>
    <img src="https://img.shields.io/badge/Status-Beta-yellow" alt="Status: Beta">
    <a href="https://github.com/Araryarch/Lumine/blob/main/LICENSE"><img src="https://img.shields.io/badge/License-AGPL--3.0-blue.svg" alt="License: AGPL-3.0"></a>
    <img src="https://img.shields.io/badge/Platform-Windows%20%7C%20macOS%20%7C%20Linux-lightgrey" alt="Platforms">
    <img src="https://img.shields.io/badge/Tech-Rust%20%2B%20Tauri%20v2-orange" alt="Tech Stack">
  </p>
</div>

---

> **Beta Version** — Lumine is currently in early beta. Features may be incomplete or subject to change. Report bugs via [Issues](https://github.com/Araryarch/Lumine/issues).

## Why Lumine?

XAMPP bundles everything into one heavy package — Apache, MySQL, PHP, all fighting for resources even when you only need one. **Lumine flips the script.** It runs your services in isolated Docker containers, giving you the flexibility of a full dev stack without the bloat.

| | **Lumine** | **XAMPP** |
|---|---|---|
| **Architecture** | Docker containers (isolated, clean) | Bundled binaries (monolithic) |
| **Services** | Nginx, Apache, Caddy, MySQL, PostgreSQL, MongoDB, Redis, Node.js, Python, PHP, Go, Deno, Bun + more | Apache, MySQL, PHP, Mercury |
| **SSL Certs** | Built-in mkcert — one click HTTPS | Manual setup, not integrated |
| **Hosts Manager** | Visual editor, auto-sync, toggle per entry | Edit `/etc/hosts` yourself |
| **Reverse Proxy** | Auto-configured Nginx proxy with HTTPS | None |
| **Project Scaffolding** | Laravel, React, Vue, Next.js in one click | None |
| **Port Conflicts** | Each service isolated in Docker | Everything fights over ports |
| **Memory Usage** | Only runs what you start (~50MB idle) | Apache + MySQL always resident |
| **Cross-Platform** | Windows, macOS, Linux (identical UX) | Windows/macOS only, different setups |
| **Modern Stack** | Rust + Tauri v2 (~10MB install) | ~400MB+ bundled |
| **System Tray** | Yes — minimize to tray, start on boot | No |

## Features

### Service Management
- **17 pre-configured services** — Nginx, Apache, Caddy, PHP-FPM, Node.js, Bun, Python, Go, Deno, MySQL, MariaDB, PostgreSQL, MongoDB, Redis, phpMyAdmin, Adminer
- **One-click start/stop** — individual or batch, with live status indicators
- **Per-service auto-start** — configure individual services to start on boot
- **Full-page service editor** — add and edit services with an inline editor (no modals)
- **Live logs** — real-time terminal output from every running service with log rotation

### Networking & Security
- **Built-in reverse proxy** — Nginx on ports 80/443, auto-routes domains to services
- **Instant HTTPS** — generate and trust SSL certs for any domain via mkcert
- **Hosts editor** — visual `/etc/hosts` management with enable/disable toggles and IPv6 support
- **Cloudflare Tunnel** — expose local services to the internet via `cloudflared`

### Docker
- **Docker stats dashboard** — live CPU and memory monitoring with color-coded bar graphs
- **Stop & remove containers** — manage containers directly from the stats view
- **Image pulling UI** — pre-download Docker images from the dashboard (pull all or selected)
- **Docker Hub browser** — search and pick image versions when adding services

### Stacks
- **Stacks** — group services into named stacks (MERN, LAMP, Laravel Local, etc.) and start/stop them together
- **10 pre-configured stacks** — XAMPP Migration, MERN Stack, Laravel Local, LEMP, LAMP, Python Data, Go Microservices, Spring Boot, Ruby on Rails, Rust Backend
- **Custom stacks** — create, edit, and delete your own stack configurations

### Projects
- **Project scaffolding** — create Laravel, React, Vue, or Next.js projects from a wizard
- **Project terminals** — open disposable Docker terminals with custom bashrc and themed prompts

### System
- **System tray** — runs quietly in the background, start on boot
- **Factory reset** — reset all app data, services, and settings to defaults
- **Import/Export config** — backup and restore service configurations
- **Settings panel** — configure preferred terminal, code editor, and file explorer with auto-detection
- **Dark mode UI** — built with React 19, Tailwind CSS v4, Space Grotesk + JetBrains Mono

## Quick Start

**Requires [Docker](https://docs.docker.com/get-docker/) installed and running.**

### Build from Source

```bash
git clone https://github.com/Araryarch/Lumine.git
cd Lumine
npm install
npm run tauri dev
```

**Prerequisites:** Node.js, Rust toolchain, Tauri CLI

## How It Works

1. Lumine starts and launches an internal Nginx reverse proxy (ports 80/443)
2. Pick services from the sidebar — Docker images run in isolated containers
3. Organize services into stacks for one-click group management
4. Monitor container performance with live CPU and memory stats
5. Map domains (`myapp.test`) to services via the Hosts tab
6. Generate SSL certs — HTTPS works instantly through the proxy
7. Expose local services to the internet via Cloudflare Tunnel
8. On exit — everything stops cleanly, no orphaned processes

## Tech Stack

| Layer | Tech |
|---|---|
| **Backend** | Rust + Tauri v2 |
| **Frontend** | React 19 + TypeScript + Vite |
| **Styling** | Tailwind CSS v4 |
| **State** | Zustand |
| **Icons** | Lucide + react-icons/si |

## Project Structure

```
Lumine/
├── src/                    # React frontend
│   ├── components/         # UI components (DockerView, StacksView, ServiceEditorView, etc.)
│   └── store/              # Zustand state management
├── src-tauri/              # Rust backend
│   └── src/
│       ├── commands.rs     # Tauri IPC commands
│       ├── service.rs      # Docker service management
│       ├── stack_manager.rs # Stack grouping logic
│       ├── project_manager.rs # Project scaffolding
│       └── utils.rs        # Shared utilities
├── CHANGELOG.md
└── LICENSE                 # AGPL-3.0
```

## Contributing

Contributions welcome! Check the [issues page](https://github.com/Araryarch/Lumine/issues), fork the repo, and open a PR.

## Support

If you find Lumine useful, consider supporting the development:

[![Buy Me a Coffee](https://img.shields.io/badge/Donate-Tako.id-ff6b6b?style=for-the-badge&logo=tako)](https://tako.id/ararya)

## License

[AGPL-3.0 License](LICENSE)

## Changelog

See [CHANGELOG.md](CHANGELOG.md) for release history.
