# Installation Guide

## Prerequisites

Before installing Lumine, ensure you have:

- **Docker 20.10+** installed and running
- **Terminal** with true color support (recommended)
- **sudo/admin access** (for domain management)

## Quick Install

### Linux/macOS

```bash
curl -fsSL https://raw.githubusercontent.com/Araryarch/lumine/main/install.sh | bash
```

### Windows (PowerShell)

```powershell
irm https://raw.githubusercontent.com/Araryarch/lumine/main/install.ps1 | iex
```

## Manual Installation

### Linux

#### AMD64
```bash
wget https://github.com/Araryarch/lumine/releases/latest/download/lumine-linux-amd64
chmod +x lumine-linux-amd64
sudo mv lumine-linux-amd64 /usr/local/bin/lumine
```

#### ARM64
```bash
wget https://github.com/Araryarch/lumine/releases/latest/download/lumine-linux-arm64
chmod +x lumine-linux-arm64
sudo mv lumine-linux-arm64 /usr/local/bin/lumine
```

### macOS

#### Intel
```bash
curl -L https://github.com/Araryarch/lumine/releases/latest/download/lumine-darwin-amd64 -o lumine
chmod +x lumine
sudo mv lumine /usr/local/bin/
```

#### Apple Silicon (M1/M2/M3)
```bash
curl -L https://github.com/Araryarch/lumine/releases/latest/download/lumine-darwin-arm64 -o lumine
chmod +x lumine
sudo mv lumine /usr/local/bin/
```

### Windows

```powershell
Invoke-WebRequest -Uri "https://github.com/Araryarch/lumine/releases/latest/download/lumine-windows-amd64.exe" -OutFile "lumine.exe"
Move-Item lumine.exe C:\Windows\System32\
```

## Build from Source

### Requirements
- Go 1.21 or higher
- Git

### Steps

```bash
# Clone repository
git clone https://github.com/Araryarch/lumine.git
cd lumine

# Download dependencies
go mod download

# Build
go build -o lumine

# Run
./lumine
```

### Using Makefile

```bash
# Build
make build

# Install to /usr/local/bin
make install

# Build for all platforms
make build-all

# Run tests
make test

# Clean build artifacts
make clean
```

## Docker Installation

If Docker is not installed, Lumine will provide installation instructions on first run.

### Linux (Ubuntu/Debian)
```bash
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh
sudo usermod -aG docker $USER
newgrp docker
```

### macOS
```bash
# Using Homebrew
brew install --cask docker

# Or download from
# https://www.docker.com/products/docker-desktop
```

### Windows
```powershell
# Using Chocolatey
choco install docker-desktop

# Or download from
# https://www.docker.com/products/docker-desktop
```

## First Run

After installation, run:

```bash
lumine
```

Lumine will:
1. ✅ Check Docker installation
2. ✅ Validate Docker is running
3. ✅ Create configuration directory (~/.lumine)
4. ✅ Set up Docker network
5. ✅ Launch TUI

## Verification

Verify installation:

```bash
# Check version
lumine --version

# Check Docker
docker --version

# Check Docker Compose
docker compose version
```

## Troubleshooting

### Docker not running
```bash
# Linux
sudo systemctl start docker

# macOS
open -a Docker

# Windows
# Start Docker Desktop from Start Menu
```

### Permission denied
```bash
# Linux/macOS
sudo chmod +x /usr/local/bin/lumine

# Add user to docker group (Linux)
sudo usermod -aG docker $USER
newgrp docker
```

### Command not found
```bash
# Check if binary is in PATH
which lumine

# Add to PATH (Linux/macOS)
export PATH=$PATH:/usr/local/bin

# Add to PATH (Windows)
# Add C:\Windows\System32 to System Environment Variables
```

## Uninstallation

### Linux/macOS
```bash
sudo rm /usr/local/bin/lumine
rm -rf ~/.lumine
```

### Windows
```powershell
Remove-Item C:\Windows\System32\lumine.exe
Remove-Item -Recurse $env:LOCALAPPDATA\lumine
```

## Next Steps

- Read [Quick Start Guide](QUICKSTART.md)
- Check [Configuration Guide](CONFIGURATION.md)
- See [Usage Examples](EXAMPLES.md)
