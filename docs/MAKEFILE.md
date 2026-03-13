# Makefile Commands Reference

Quick reference for all available Makefile commands.

## 📋 Quick Commands

```bash
make help          # Show all available commands
make dev           # Run in development mode
make build         # Build the binary
make run           # Build and run
make db-setup      # Setup all databases
```

## 🚀 Development

### Running

```bash
make dev           # Run with go run (hot reload with changes)
make run           # Build and run the binary
make watch         # Watch for changes and auto-rebuild (requires entr)
```

### Building

```bash
make build         # Build for current platform
make build-all     # Build for all platforms (Linux, macOS, Windows)
```

Output location: `build/`

## 📦 Installation

```bash
make install       # Install to /usr/local/bin (requires sudo)
make uninstall     # Remove from /usr/local/bin
```

## 🧪 Testing

```bash
make test          # Run all tests
make test-coverage # Run tests with coverage report
make test-race     # Run tests with race detector
make benchmark     # Run benchmarks
```

Coverage report: `coverage.html`

## 🔍 Code Quality

```bash
make fmt           # Format code with go fmt
make vet           # Run go vet
make lint          # Run golangci-lint (requires installation)
make check         # Run fmt + vet
```

## 📚 Dependencies

```bash
make deps          # Download and tidy dependencies
make deps-update   # Update all dependencies to latest
```

## 🗄️ Database Management

### Setup & Control

```bash
make db-setup      # Start all databases and admin panels
make db-stop       # Stop all databases
make db-restart    # Restart all databases
make db-logs       # Show database logs (follow mode)
```

### Cleanup

```bash
make db-clean      # Remove all database data (WARNING: destructive!)
```

**After `make db-setup`, you'll have:**

| Service | Port | Admin Panel | URL |
|---------|------|-------------|-----|
| MySQL | 3306 | phpMyAdmin | http://localhost:8080 |
| PostgreSQL | 5432 | pgAdmin | http://localhost:8084 |
| MariaDB | 3307 | phpMyAdmin | http://localhost:8080 |
| MongoDB | 27017 | Mongo Express | http://localhost:8082 |
| Redis | 6379 | Redis Commander | http://localhost:8083 |
| Elasticsearch | 9200 | - | http://localhost:9200 |
| Adminer (All) | - | Adminer | http://localhost:8081 |

## 🐳 Docker

```bash
make docker-validate  # Check Docker installation and status
make docker-build     # Build Docker image
make docker-clean     # Clean Docker resources
```

## 🧹 Cleanup

```bash
make clean         # Clean build artifacts
make clean-all     # Clean everything (build + Docker + databases)
```

## 📦 Release

```bash
make release       # Build for all platforms + generate checksums
```

Creates:
- `build/lumine-linux-amd64`
- `build/lumine-linux-arm64`
- `build/lumine-darwin-amd64`
- `build/lumine-darwin-arm64`
- `build/lumine-windows-amd64.exe`
- `build/checksums.txt`

## ℹ️ Information

```bash
make info          # Show project and environment information
```

## 🎯 Common Workflows

### First Time Setup

```bash
# 1. Clone and setup
git clone https://github.com/Araryarch/lumine.git
cd lumine

# 2. Install dependencies
make deps

# 3. Setup databases
make db-setup

# 4. Run in dev mode
make dev
```

### Development Workflow

```bash
# Terminal 1: Run app
make dev

# Terminal 2: Watch databases
make db-logs

# Make changes, app auto-reloads
```

### Before Committing

```bash
# Format and check code
make check

# Run tests
make test

# Build to verify
make build
```

### Release Workflow

```bash
# 1. Update version
export VERSION=v1.0.0

# 2. Run tests
make test

# 3. Build release
make release

# 4. Verify builds
ls -lh build/
```

## 🔧 Customization

### Change Install Directory

```bash
make install INSTALL_DIR=/opt/lumine
```

### Change Build Directory

```bash
make build BUILD_DIR=dist
```

### Set Version

```bash
make build VERSION=v1.0.0
```

## 📝 Tips

### 1. Use Tab Completion

```bash
make <TAB><TAB>  # Shows all available commands
```

### 2. Combine Commands

```bash
make clean build run  # Clean, build, and run
```

### 3. Watch Mode

Install `entr` for auto-rebuild:

```bash
# macOS
brew install entr

# Linux
sudo apt install entr  # Debian/Ubuntu
sudo dnf install entr  # Fedora

# Then use
make watch
```

### 4. Parallel Builds

```bash
make -j4 build-all  # Use 4 parallel jobs
```

### 5. Verbose Output

```bash
make build GOFLAGS=-v  # Verbose build
```

## 🐛 Troubleshooting

### "Command not found"

```bash
# Install make
sudo apt install make  # Linux
brew install make      # macOS
```

### "Permission denied"

```bash
# For install/uninstall
sudo make install

# For database operations
sudo make db-setup
```

### "Docker not running"

```bash
# Check Docker
make docker-validate

# Start Docker
# macOS: Open Docker Desktop
# Linux: sudo systemctl start docker
```

### "Port already in use"

```bash
# Check what's using the port
lsof -i :3306

# Stop databases
make db-stop

# Or kill specific process
kill -9 <PID>
```

## 📚 Additional Resources

- [Database Guide](DATABASE.md)
- [Quick Start](QUICKSTART.md)
- [Installation Guide](INSTALLATION.md)
- [Contributing Guide](../CONTRIBUTING.md)


## 🗑️ Container & Volume Management

### List Resources

```bash
make containers-list   # List all Lumine containers
make volumes-list      # List all Lumine volumes
make network-inspect   # Inspect Lumine network
```

### Stop & Remove Containers

```bash
make containers-stop   # Stop all running containers
make containers-remove # Remove all containers (with confirmation)
make containers-clean  # Remove only stopped containers
make containers-prune  # Remove containers + volumes (DESTRUCTIVE!)
```

### Volume Management

```bash
make volumes-list      # List all volumes
make volumes-remove    # Remove all volumes (with confirmation)
make volumes-prune     # Remove unused volumes only
```

### Network Management

```bash
make network-create    # Create Lumine network
make network-remove    # Remove Lumine network
make network-inspect   # Show network details
```

### Nuclear Options ☢️

```bash
make clean-containers  # Remove all containers and volumes
make clean-everything  # DESTROY EVERYTHING (requires typing 'DESTROY')
```

**Warning Levels:**
- 🟢 Safe: `containers-list`, `volumes-list`, `network-inspect`
- 🟡 Caution: `containers-stop`, `containers-clean`, `volumes-prune`
- 🟠 Dangerous: `containers-remove`, `volumes-remove`, `db-clean`
- 🔴 DESTRUCTIVE: `containers-prune`, `clean-everything`

## 🔄 Reset & Refresh

```bash
make db-reset          # Reset databases (remove data and restart fresh)
make containers-prune  # Complete container cleanup
make network-create    # Recreate network
```

## 💡 Pro Tips

### Safe Cleanup Workflow

```bash
# 1. Stop everything
make containers-stop

# 2. List what will be removed
make containers-list
make volumes-list

# 3. Remove containers only (keep data)
make containers-remove

# 4. Or remove everything
make containers-prune
```

### Selective Cleanup

```bash
# Remove specific container
docker rm -f lumine-mysql

# Remove specific volume
docker volume rm lumine_mysql_data

# Remove stopped containers only
make containers-clean
```

### Emergency Reset

```bash
# If everything is broken
make clean-everything  # Type 'DESTROY' to confirm

# Then start fresh
make db-setup
make dev
```
