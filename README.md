# Lumine

Modern Development Stack Manager - Alternative to Laragon (Windows) & Laravel Herd (macOS)

![Lumine](assets/lumine-banner.png)

## Overview

Lumine is a powerful TUI (Terminal User Interface) application for managing modern development stacks on Linux. It provides an intuitive interface for managing services, projects, and databases - similar to Laragon on Windows and Laravel Herd on macOS.

## Features

### 🚀 Service Orchestration
- **Multi-Service Management**: Nginx, Apache, PHP-FPM, MySQL, PostgreSQL, Redis, and more
- **Multi-Version Support**: Switch between PHP (7.4, 8.0, 8.1, 8.2, 8.3) and Node.js (16, 18, 20, 22) versions on-the-fly
- **Auto Port Management**: Automatic port conflict detection and resolution
- **Daemon Mode**: Background service with auto-restart on crash
- **Health Monitoring**: Real-time service health checks with uptime tracking

### 📁 Modern Project Management
- **Project Scaffolding**: Create projects with one command
  - Laravel, Symfony, CodeIgniter (PHP)
  - React, Vue, Next.js, Nuxt.js (JavaScript)
  - WordPress, Static sites
- **Automatic Virtual Hosts**: Auto-generate `.test` domains
- **SSL Certificates**: Automatic HTTPS with self-signed certificates
- **Localhost Tunneling**: Expose local projects to the internet (ngrok alternative)
- **Dependency Checking**: Detect Composer, NPM, PNPM, Yarn, Git, Docker

### 🗄️ Database Management
- **Quick Actions**: Create, drop, and backup databases
- **Multi-Database Support**: MySQL, PostgreSQL, SQLite
- **Connection Profiles**: Switch between database connections
- **Query Logging**: Real-time query monitoring
- **Slow Query Detection**: Identify performance bottlenecks

### 🎨 User Experience
- **TUI Interface**: Beautiful terminal interface with keyboard navigation
- **Toast Notifications**: Success, error, warning, and info notifications
- **Vim-style Navigation**: Familiar keyboard shortcuts
- **Real-time Updates**: Auto-refresh service status and logs

## Installation

### Quick Install (Linux)

```bash
# Clone repository
git clone https://github.com/Araryarch/Lumine.git
cd Lumine

# Fix dependencies
make tidy

# Build and run
make run
```

### Manual Installation

```bash
# Clone repository
git clone https://github.com/Araryarch/Lumine.git
cd Lumine

# Tidy dependencies (important!)
rm -rf vendor
go mod tidy

# Build
go build -mod=mod -o lumine .

# Run
./lumine

# Install to system (optional)
sudo mv lumine /usr/local/bin/
```

### Using Makefile (Recommended)

```bash
# Show available commands
make help

# Build only
make build

# Build and run
make run

# Install to system
make install

# Clean build artifacts
make clean
```

## Usage

### Start Lumine

```bash
lumine
```

### Daemon Mode

```bash
# Start daemon
lumine --daemon start

# Stop daemon
lumine --daemon stop

# Check status
lumine --daemon status
```

## Keyboard Shortcuts

### Navigation
- `1-6`: Focus Docker panels (Projects, Services, Containers, Images, Volumes, Networks)
- `7`: Focus Lumine Services
- `8`: Focus Lumine Projects
- `9`: Focus Lumine Databases
- `h/l` or `←/→`: Navigate between panels
- `j/k` or `↑/↓`: Navigate within panel
- `q`: Quit

### Lumine Services Panel (7)
- `s`: Start service
- `S`: Stop service
- `r`: Restart service
- `v`: Switch version (PHP/Node.js)
- `H`: Health check
- `Enter`: View details

### Lumine Projects Panel (8)
- `n`: New project
- `d`: Delete project
- `e`: Expose via tunnel
- `o`: Open in browser
- `t`: Open terminal
- `Enter`: View details

### Lumine Databases Panel (9)
- `c`: Create database
- `d`: Drop database
- `b`: Backup database
- `s`: Switch connection
- `Enter`: View logs

## Configuration

Configuration file: `~/.config/lumine/config.yml`

```yaml
projectsDirectory: ~/Projects
defaultPHPVersion: "8.2"
defaultNodeVersion: "20"
preferredWebServer: nginx
autoStartServices: true
enableDaemonMode: true
defaultTunnelService: ngrok
```

## Requirements

- Linux (Ubuntu, Debian, Arch, Fedora, etc.)
- Go 1.21+ (for building from source)
- Docker (optional, for containerized services)
- Root/sudo access (for editing /etc/hosts and binding to port 80)

## Architecture

Lumine is built with:
- **TUI Framework**: gocui for terminal interface
- **Service Management**: Docker SDK for containerized services
- **Project Management**: Template-based project scaffolding
- **Database Management**: Native database drivers (MySQL, PostgreSQL, SQLite)

## Contributing

Contributions are welcome! Please read [CONTRIBUTING.md](CONTRIBUTING_LUMINE.md) for details.

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Credits

- Based on [lazydocker](https://github.com/jesseduffield/lazydocker) TUI framework
- Inspired by [Laragon](https://laragon.org/) and [Laravel Herd](https://herd.laravel.com/)

## Support

- 🐛 [Report Issues](https://github.com/Araryarch/Lumine/issues)
- 💬 [Discussions](https://github.com/Araryarch/Lumine/discussions)
- 📖 [Documentation](https://github.com/Araryarch/Lumine/wiki)

---

Made with ❤️ for the Linux development community
