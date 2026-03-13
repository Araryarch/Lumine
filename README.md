# Lumine 🌟

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go)](https://go.dev/)
[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)
[![Docker](https://img.shields.io/badge/Docker-20.10+-2496ED?style=flat&logo=docker)](https://www.docker.com/)
[![Platform](https://img.shields.io/badge/Platform-Linux%20%7C%20macOS%20%7C%20Windows-lightgrey)](https://github.com/Araryarch/lumine)
[![Release](https://img.shields.io/github/v/release/Araryarch/lumine)](https://github.com/Araryarch/lumine/releases)

Aplikasi TUI (Terminal User Interface) modern untuk mengelola Docker development environment sebagai pengganti XAMPP/Laragon. Containerized, powerful, dan mudah digunakan.

## ✨ Fitur Utama

## Features

### Project Management
- Create projects: Laravel, Next.js, Vue, Django, Express, FastAPI, Nuxt, SvelteKit, Remix, NestJS, Axum, Actix, Rocket
- Project manager with status monitoring
- Auto-generate .test domains
- One-click start/stop

### Runtime Management
- Multiple runtimes: PHP, Node.js, Python, Bun, Deno, Go, Rust
- Version switching per runtime
- Per-project runtime configuration

### Docker Services
- Pre-configured: Nginx, Apache, Caddy, MySQL, PostgreSQL, MariaDB, MongoDB, Redis, Elasticsearch
- Version selector with 50+ versions
- Multi-service support
- Real-time status monitoring

### Database Management
- Support for 6 database engines
- 5 admin panels included
- Auto-initialization scripts
- Backup/restore tools

### Beautiful TUI
- 3-panel layout (Sidebar, Main, Detail)
- Minimalist design
- Vim-style navigation (h/j/k/l)
- Context-aware help

### Port Management
- Auto-detect port conflicts
- Find alternative ports automatically
- Smart port suggestions
- Manual port selection

## Usage

### Navigation

```
Arrow Keys / j/k  - Move up/down
h/l              - Switch panels
Tab              - Next panel
Enter            - Select/Execute
```

### Service Management

```
Space   - Select/Deselect
s       - Start services
x       - Stop services
r       - Restart services
v       - Change version
delete  - Remove service
c       - Cleanup all
Ctrl+C  - Quit
```

### Example Workflow

1. Start Lumine: `make dev`
2. Navigate to "Services"
3. Select MySQL and Redis (Space)
4. Press 's' to start
5. Check status in detail panel

## Documentation

- [Getting Started](GETTING_STARTED.md) - Complete setup guide
- [Quick Start](docs/QUICKSTART.md) - 5-minute tutorial
- [Database Guide](docs/DATABASE.md) - Database management
- [Web Servers](docs/WEBSERVERS.md) - Nginx, Apache, Caddy
- [Port Management](docs/PORT_MANAGEMENT.md) - Port handling
- [Makefile Commands](docs/MAKEFILE.md) - All make commands
- [Cleanup Guide](docs/CLEANUP.md) - Container cleanup

## Makefile Commands

```bash
make help          # Show all commands
make dev           # Run in development mode
make build         # Build binary
make run           # Build and run
make install       # Install to system
make test          # Run tests
make db-setup      # Setup all databases
make db-stop       # Stop databases
make db-logs       # View database logs
make clean         # Clean build artifacts
```

## 🚀 Quick Start

### Installation

#### One-line Install (Linux/macOS)
```bash
curl -fsSL https://raw.githubusercontent.com/Araryarch/lumine/main/install.sh | bash
```

#### One-line Install (Windows PowerShell)
```powershell
irm https://raw.githubusercontent.com/Araryarch/lumine/main/install.ps1 | iex
```

#### Manual Installation

**Linux:**
```bash
# AMD64
wget https://github.com/Araryarch/lumine/releases/latest/download/lumine-linux-amd64
chmod +x lumine-linux-amd64
sudo mv lumine-linux-amd64 /usr/local/bin/lumine

# ARM64
wget https://github.com/Araryarch/lumine/releases/latest/download/lumine-linux-arm64
chmod +x lumine-linux-arm64
sudo mv lumine-linux-arm64 /usr/local/bin/lumine
```

**macOS:**
```bash
# Intel
curl -L https://github.com/Araryarch/lumine/releases/latest/download/lumine-darwin-amd64 -o lumine
chmod +x lumine
sudo mv lumine /usr/local/bin/

# Apple Silicon (M1/M2/M3)
curl -L https://github.com/Araryarch/lumine/releases/latest/download/lumine-darwin-arm64 -o lumine
chmod +x lumine
sudo mv lumine /usr/local/bin/
```

**Windows:**
```powershell
# Download and add to PATH
Invoke-WebRequest -Uri "https://github.com/Araryarch/lumine/releases/latest/download/lumine-windows-amd64.exe" -OutFile "lumine.exe"
Move-Item lumine.exe C:\Windows\System32\
```

#### Build from Source
```bash
git clone https://github.com/Araryarch/lumine.git
cd lumine
go build -o lumine
./lumine
```

### First Run

Lumine will automatically:
1. ✅ Check if Docker is installed
2. ✅ Validate Docker is running
3. ✅ Attempt to start Docker if not running
4. ✅ Create configuration directory (~/.lumine)
5. ✅ Set up Docker network
6. ✅ Launch TUI

If Docker is not installed, Lumine will provide installation instructions for your OS.

## 📦 Create Your First Project

1. Start Lumine: `./lumine`
2. Navigate ke sidebar → `🚀 Projects` atau `➕ New Project`
3. Pilih framework (Laravel, Next.js, Vue, Django, dll)
4. Enter project name
5. Project akan dibuat dengan domain `.test` otomatis
6. Start project dan akses via browser!

## 🎮 Keyboard Shortcuts

### Navigation
- `↑/↓` atau `j/k` - Navigate items
- `h/l` - Switch panels (left/right)
- `tab` - Next panel
- `enter` - Select/Execute

### Service Management
- `space` - Select/Deselect service
- `s` atau `enter` - Start service(s)
- `x` - Stop service(s)
- `r` - Restart service(s)
- `v` - Change version
- `a` - Select all
- `d` - Deselect all

### Project Management
- `n` - New project
- `s` - Start project
- `x` - Stop project
- `o` - Open in browser
- `e` - Open in editor

### Views (Sidebar)
- 📦 **Services** - Manage Docker services
- 🚀 **Projects** - Manage your projects
- ⚙️ **Runtimes** - Switch runtime versions
- ➕ **New Project** - Create new project
- 📊 **Logs** - View activity logs
- 🔄 **Refresh** - Refresh status
- ❌ **Quit** - Exit

## 🛠️ Supported Frameworks

### PHP
- **Laravel** - Full-stack PHP framework
- **Symfony** - Enterprise PHP framework
- **CodeIgniter** - Lightweight PHP framework

### JavaScript/TypeScript
- **Next.js** - React framework with SSR
- **Nuxt** - Vue framework with SSR
- **SvelteKit** - Svelte framework
- **Remix** - Full-stack React framework
- **NestJS** - Progressive Node.js framework
- **Express** - Minimal Node.js framework
- **Vue** - Progressive JavaScript framework

### Python
- **Django** - High-level Python framework
- **FastAPI** - Modern Python API framework
- **Flask** - Micro Python framework

### Rust
- **Axum** - Ergonomic web framework
- **Actix-web** - Powerful, pragmatic web framework
- **Rocket** - Simple, fast web framework

## ⚙️ Supported Runtimes

### PHP
- 8.3 (fpm, apache, cli)
- 8.2 (fpm, apache, cli)
- 8.1 (fpm, apache, cli)
- 8.0, 7.4

### Node.js
- 21, 20 (LTS), 18 (LTS)
- Alpine variants

### Python
- 3.12, 3.11, 3.10, 3.9
- Slim variants

### Alternative Runtimes
- **Bun** - Fast JavaScript runtime
- **Deno** - Secure TypeScript runtime
- **Go** - For Go applications
- **Rust** - For Rust applications (Axum, Actix, Rocket)

### Package Managers
- npm, yarn, pnpm (Node.js)
- composer (PHP)
- pip, poetry (Python)
- cargo (Rust)
- go mod (Go)

## 🌐 Domain Management

Lumine automatically manages `.test` domains for your projects:

```yaml
projects:
  - name: myapp
    domain: myapp.test  # Auto-generated
    type: laravel
    port: 8000
```

Access your project at `http://myapp.test` in your browser!

## 📋 Configuration

File konfigurasi di `~/.lumine/config.yaml`:

```yaml
# Docker Services
services:
  - name: nginx
    type: nginx
    version: latest
    port: 80
  - name: mysql
    type: mysql
    version: "8.0"
    port: 3306
    env:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: lumine

# Your Projects
projects:
  - name: myapp
    type: laravel
    path: /home/user/projects/myapp
    domain: myapp.test
    runtime: php
    version: "8.2-fpm"
    port: 8000
    status: running

# Runtime Versions
runtimes:
  php: "8.2-fpm"
  node: "20-alpine"
  python: "3.11-slim"
  bun: latest
  deno: latest
  go: "1.21-alpine"
```

## 🎯 Use Cases

### Web Development
```bash
# Create Laravel project
lumine → New Project → Laravel → myapp
# Access at myapp.test
```

### API Development
```bash
# Create FastAPI project
lumine → New Project → FastAPI → myapi
# Access at myapi.test:8000
```

### Full-stack Development
```bash
# Create Next.js + Laravel
lumine → New Project → Next.js → frontend
lumine → New Project → Laravel → backend
# frontend.test + backend.test
```

## 📦 Requirements

- Go 1.21+
- Docker & Docker Compose
- Terminal dengan true color support
- sudo access (untuk domain management)

## 🔥 Why Lumine?

| Feature | XAMPP | Laragon | Lumine |
|---------|-------|---------|--------|
| Containerized | ❌ | ❌ | ✅ |
| Multiple PHP Versions | ❌ | ✅ | ✅ |
| Multiple Runtimes | ❌ | ⚠️ | ✅ |
| Project Manager | ❌ | ✅ | ✅ |
| Pretty Domains | ❌ | ✅ | ✅ |
| TUI Interface | ❌ | ❌ | ✅ |
| Cross-platform | ⚠️ | ❌ | ✅ |
| Isolated Environments | ❌ | ❌ | ✅ |

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit PRs.

## 📄 License

MIT License


## 📸 Screenshots

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ 🌟 LUMINE - Docker Development Environment Manager                         │
├──────────────┬──────────────────────────────┬──────────────────────────────┤
│              │                              │                              │
│ 📦 Services  │  ▶ ☑ PHP 8.2-fpm            │  📋 Service Details          │
│ 🚀 Projects  │    ☐ MYSQL 8.0              │                              │
│ ⚙️  Runtimes │    ☐ NGINX latest           │  Name: php                   │
│ ➕ New Proj  │    ☐ REDIS 7.2              │  Type: PHP                   │
│ 📊 Logs      │                              │  Version: 8.2-fpm            │
│ 🔄 Refresh   │  Projects:                   │  Port: 9000                  │
│ ❌ Quit      │  ▶ LARAVEL myapp.test ●     │  Status: ● Running           │
│              │    NEXTJS frontend.test ●    │                              │
│              │                              │  📝 Activity Logs            │
│              │  Runtimes:                   │  • Started php successfully  │
│              │  🐘 PHP v8.2-fpm            │  • Started mysql             │
│              │  🟢 Node.js v20-alpine      │  • Created myapp.test        │
│              │  🐍 Python v3.11-slim       │                              │
│              │  🦀 Rust vlatest            │                              │
└──────────────┴──────────────────────────────┴──────────────────────────────┘
 Running: 3/5 • Selected: 1
 ↑/↓,j/k navigate • h/l panels • space select • s start • v version • ctrl+c quit
```

## 🌟 Features Highlight

- ✅ **Zero Configuration** - Works out of the box
- ✅ **Auto Docker Setup** - Validates and starts Docker automatically
- ✅ **13+ Frameworks** - Laravel, Next.js, Django, Axum, and more
- ✅ **7 Runtimes** - PHP, Node.js, Python, Rust, Bun, Deno, Go
- ✅ **Pretty Domains** - myapp.test instead of localhost:8000
- ✅ **Cross-platform** - Linux, macOS, Windows
- ✅ **Beautiful TUI** - Vim-style navigation, color-coded badges
- ✅ **Isolated Environments** - Each project in its own container

## 🤝 Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details.

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Inspired by [Laragon](https://laragon.org/) and [Herd](https://herd.laravel.com/)
- Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Powered by [Docker](https://www.docker.com/)

## 📧 Contact

- GitHub: [@Araryarch](https://github.com/Araryarch)
- Issues: [GitHub Issues](https://github.com/Araryarch/lumine/issues)
- Discussions: [GitHub Discussions](https://github.com/Araryarch/lumine/discussions)

---

Made with ❤️ by [Araryarch](https://github.com/Araryarch)


## 🗑️ Cleanup Tool

Lumine includes a cross-platform cleanup utility:

```bash
# Interactive cleanup
lumine-cleanup

# Or via Makefile
make cleanup-interactive
```

**Options:**
1. Stop containers only
2. Remove containers (keep data)
3. Remove containers + volumes (with backup)
4. Nuclear cleanup (remove everything)

**Features:**
- ✅ Cross-platform (Windows, Linux, macOS)
- ✅ Auto-backup before destructive operations
- ✅ Safe confirmation prompts
- ✅ Color-coded output

See [Cleanup Guide](docs/CLEANUP.md) for details.

Lumine includes comprehensive database management with admin panels:

### Supported Databases
- **MySQL 8.0** - Port 3306
- **PostgreSQL 16** - Port 5432
- **MariaDB 11.2** - Port 3307
- **MongoDB 7.0** - Port 27017
- **Redis 7.2** - Port 6379
- **Elasticsearch 8.11** - Port 9200

### Admin Panels
- **phpMyAdmin** - http://localhost:8080 (MySQL/MariaDB)
- **Adminer** - http://localhost:8081 (Universal DB)
- **Mongo Express** - http://localhost:8082 (MongoDB)
- **Redis Commander** - http://localhost:8083 (Redis)
- **pgAdmin** - http://localhost:8084 (PostgreSQL)

### Quick Start

```bash
# Setup all databases
make db-setup

# Stop databases
make db-stop

# Restart databases
make db-restart

# View logs
make db-logs

# Clean database data (WARNING: destructive)
make db-clean
```

### Database Access

```bash
# MySQL
mysql -h localhost -P 3306 -u root -proot

# PostgreSQL
psql -h localhost -p 5432 -U postgres

# MongoDB
mongosh mongodb://root:root@localhost:27017

# Redis
redis-cli -h localhost -p 6379
```

### Connection Strings

**MySQL/MariaDB:**
```
root:root@tcp(localhost:3306)/lumine?charset=utf8mb4&parseTime=True
```

**PostgreSQL:**
```
host=localhost port=5432 user=postgres password=postgres dbname=lumine sslmode=disable
```

**MongoDB:**
```
mongodb://root:root@localhost:27017/lumine
```

**Redis:**
```
redis://localhost:6379
```


## 🔌 Automatic Port Management

Lumine automatically handles port conflicts:

```bash
# If port 3306 is in use:
# ⚠️  Port 3306 is in use, using alternative port 3307 for mysql
# ✓ MySQL started on port 3307
```

**Features:**
- ✅ Auto-detect port conflicts
- ✅ Find alternative ports automatically
- ✅ Smart port suggestions based on service type
- ✅ Manual port selection via TUI
- ✅ Port availability checking

**Alternative Ports:**
- MySQL: 3306 → 3307, 3308, 3309
- PostgreSQL: 5432 → 5433, 5434, 5435
- MongoDB: 27017 → 27018, 27019, 27020
- Redis: 6379 → 6380, 6381, 6382

See [Port Management Guide](docs/PORT_MANAGEMENT.md) for details.


## Web Servers

Lumine supports multiple web servers:

### Nginx
- High-performance web server
- Reverse proxy
- Load balancing
- Port: 80

### Apache
- Most popular web server
- .htaccess support
- mod_rewrite
- Port: 80

### Caddy
- Modern web server
- Automatic HTTPS
- Simple configuration
- Port: 8085 (HTTP), 8445 (HTTPS)

**Usage:**
```yaml
services:
  - name: nginx
    type: nginx
    version: latest
    port: 80
    
  - name: apache
    type: apache
    version: latest
    port: 8080
    
  - name: caddy
    type: caddy
    version: latest
    port: 8085
```

See [Web Server Guide](docs/WEBSERVERS.md) for configuration examples.
