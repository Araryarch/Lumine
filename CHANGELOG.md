# Changelog

All notable changes to Lumine will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- 🎨 Beautiful TUI with 3-panel layout
- 🚀 Project management for 13+ frameworks
- ⚙️ Runtime version switching (PHP, Node.js, Python, Bun, Deno, Go, Rust)
- 🌐 Pretty domain management with .test domains
- 🐳 Docker service management
- 🔄 Version selector for Docker images
- 📝 Activity logs and monitoring
- ✅ Docker validation and auto-setup on first run
- 🦀 Rust project support (Axum, Actix, Rocket)
- 🔧 Cross-platform support (Linux, macOS, Windows)
- 📦 One-line installation scripts
- 🤖 GitHub Actions for automated builds
- 📋 Docker Compose generation per project
- 🌐 Nginx configuration generation
- 🔐 Isolated Docker networks

### Supported Frameworks
- **PHP**: Laravel, Symfony, CodeIgniter
- **JavaScript/TypeScript**: Next.js, Nuxt, Vue, SvelteKit, Remix, NestJS, Express
- **Python**: Django, FastAPI, Flask
- **Rust**: Axum, Actix-web, Rocket

### Supported Services
- Nginx, Apache
- MySQL, MariaDB, PostgreSQL
- Redis, Memcached
- MongoDB, Elasticsearch
- phpMyAdmin, Adminer
- RabbitMQ

## [0.1.0] - TBD

### Added
- Initial release
- Basic TUI interface
- Docker integration
- Service management
- Configuration system

---

## Release Notes

### Installation

#### Linux/macOS
```bash
curl -fsSL https://raw.githubusercontent.com/Araryarch/lumine/main/install.sh | bash
```

#### Windows
```powershell
irm https://raw.githubusercontent.com/Araryarch/lumine/main/install.ps1 | iex
```

### Requirements
- Docker 20.10+
- Go 1.21+ (for building from source)
- Terminal with true color support (recommended)

### Breaking Changes
None yet - this is the initial release!

### Migration Guide
N/A - first release

### Known Issues
- Docker Desktop must be running before starting Lumine
- Windows: May require running as Administrator for domain management
- macOS: First run may require granting permissions to Docker

### Contributors
Thank you to all contributors!

---

For more information, visit: https://github.com/Araryarch/lumine
