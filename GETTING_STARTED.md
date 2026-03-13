# Getting Started with Lumine

Complete guide to get Lumine up and running.

## Prerequisites

Before starting, make sure you have:

1. **Docker** (version 20.10 or higher)
   - Linux: `docker --version`
   - macOS: Docker Desktop
   - Windows: Docker Desktop

2. **Go** (version 1.21 or higher) - for building from source
   - Check: `go version`

3. **Make** (optional, but recommended)
   - Linux: Usually pre-installed
   - macOS: Install Xcode Command Line Tools
   - Windows: Install via Chocolatey or use Git Bash

## Quick Start (5 minutes)

### Step 1: Clone Repository

```bash
git clone https://github.com/Araryarch/lumine.git
cd lumine
```

### Step 2: Install Dependencies

```bash
go mod download
```

### Step 3: Build Lumine

```bash
# Using Make (recommended)
make build

# Or using Go directly
go build -o build/lumine .
```

### Step 4: Run Lumine

```bash
# Using Make
make run

# Or directly
./build/lumine

# Or in development mode (with auto-reload)
make dev
```

That's it! Lumine will:
1. Check if Docker is installed and running
2. Create configuration directory (~/.lumine)
3. Set up Docker network
4. Launch the TUI

## Detailed Setup

### Option 1: Development Mode (Recommended for Testing)

```bash
# 1. Clone and enter directory
git clone https://github.com/Araryarch/lumine.git
cd lumine

# 2. Download dependencies
make deps

# 3. Run in development mode
make dev
```

This will run Lumine directly with `go run`, perfect for development.

### Option 2: Build and Install

```bash
# 1. Clone repository
git clone https://github.com/Araryarch/lumine.git
cd lumine

# 2. Build
make build

# 3. Install to system (requires sudo)
make install

# 4. Run from anywhere
lumine
```

### Option 3: Build All Tools

```bash
# Build Lumine + Cleanup tool
make build-all

# Install both
make install-all

# Now you can use:
lumine          # Main TUI
lumine-cleanup  # Cleanup tool
```

## First Run

When you run Lumine for the first time:

```bash
./build/lumine
```

You'll see:

```
Checking Docker installation...
Docker is running (version 24.0.7)
Docker Compose is available (version 2.23.0)
Initializing configuration...
Configuration ready

Starting Lumine...
```

Then the TUI will appear:

```
┌─────────────────────────────────────────────────────────────┐
│ LUMINE - Docker Development Environment Manager            │
├──────────────┬──────────────────────────┬───────────────────┤
│              │                          │                   │
│ Services     │  Services                │  Service Details  │
│ Projects     │                          │                   │
│ Databases    │  > [ ] MySQL 8.0         │  Name: mysql      │
│ Runtimes     │    [ ] PostgreSQL 16     │  Type: mysql      │
│ New Project  │    [ ] Redis 7.2         │  Port: 3306       │
│ Logs         │                          │  Status: stopped  │
│ Refresh      │                          │                   │
│ Quit         │                          │                   │
│              │                          │                   │
└──────────────┴──────────────────────────┴───────────────────┘
```

## Setup Databases

### Quick Setup (All Databases)

```bash
# In another terminal
make db-setup
```

This starts:
- MySQL (port 3306)
- PostgreSQL (port 5432)
- MariaDB (port 3307)
- MongoDB (port 27017)
- Redis (port 6379)
- Elasticsearch (port 9200)

Plus admin panels:
- phpMyAdmin (http://localhost:8080)
- Adminer (http://localhost:8081)
- Mongo Express (http://localhost:8082)
- Redis Commander (http://localhost:8083)
- pgAdmin (http://localhost:8084)
- Caddy (http://localhost:8085)

### Selective Setup

```bash
# Start only MySQL
docker compose -f docker-compose.db.yml up -d mysql

# Start MySQL + phpMyAdmin
docker compose -f docker-compose.db.yml up -d mysql phpmyadmin

# Start all databases (no admin panels)
docker compose -f docker-compose.db.yml up -d mysql postgres mongodb redis
```

## Using Lumine TUI

### Navigation

```
Arrow Keys or j/k  - Move up/down
h/l               - Switch panels (left/right)
Tab               - Next panel
Enter             - Select/Execute
```

### Service Management

```
Space  - Select/Deselect service
s      - Start selected services
x      - Stop selected services
r      - Restart selected services
v      - Change version
delete - Remove service
c      - Cleanup all
```

### Example Workflow

1. **Start Lumine**
   ```bash
   make dev
   ```

2. **Navigate to Services**
   - Use arrow keys to select "Services" in sidebar
   - Press Enter

3. **Select Services**
   - Use arrow keys to navigate to MySQL
   - Press Space to select
   - Navigate to Redis
   - Press Space to select

4. **Start Services**
   - Press 's' or Enter
   - Services will start automatically

5. **Check Status**
   - Status updates in real-time
   - See logs in detail panel

## Common Tasks

### Start a Service

```bash
# Via TUI
1. Navigate to service
2. Press Space to select
3. Press 's' to start

# Via Docker directly
docker compose -f docker-compose.db.yml up -d mysql
```

### Create a Project

```bash
# Via TUI
1. Navigate to "New Project" in sidebar
2. Press Enter
3. Select framework (Laravel, Next.js, etc.)
4. Enter project name
5. Wait for creation

# Via command line (coming soon)
lumine create laravel myapp
```

### Change PHP Version

```bash
# Via TUI
1. Navigate to "Runtimes"
2. Select "PHP"
3. Press 'v' for version selector
4. Choose version (8.3, 8.2, 8.1, etc.)
5. Press Enter
```

### Access Databases

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

### View Logs

```bash
# Via TUI
Navigate to "Logs" in sidebar

# Via Docker
docker logs lumine-mysql
docker logs lumine-postgres

# Via Make
make db-logs
```

## Troubleshooting

### Docker Not Running

```bash
# Check Docker status
docker ps

# Start Docker
# macOS: Open Docker Desktop
# Linux: sudo systemctl start docker
# Windows: Start Docker Desktop from Start Menu
```

### Port Already in Use

Lumine will automatically find alternative ports:

```
Port 3306 is in use, using alternative port 3307 for mysql
```

Or manually check:
```bash
# Linux/macOS
lsof -i :3306

# Windows
netstat -ano | findstr :3306
```

### Build Errors

```bash
# Clean and rebuild
make clean
make deps
make build

# Or
go clean
go mod tidy
go build -o build/lumine .
```

### Permission Denied

```bash
# Linux/macOS
chmod +x build/lumine
./build/lumine

# Or install to system
sudo make install
lumine
```

## Development Workflow

### For Contributors

```bash
# 1. Fork and clone
git clone https://github.com/yourusername/lumine.git
cd lumine

# 2. Create branch
git checkout -b feature/my-feature

# 3. Install dependencies
make deps

# 4. Run in dev mode
make dev

# 5. Make changes
# Edit files...

# 6. Test
make test

# 7. Format and check
make fmt
make vet

# 8. Build
make build

# 9. Commit and push
git add .
git commit -m "Add my feature"
git push origin feature/my-feature
```

### Hot Reload (Optional)

Install `entr` for auto-reload:

```bash
# macOS
brew install entr

# Linux
sudo apt install entr  # Debian/Ubuntu
sudo dnf install entr  # Fedora

# Then use
make watch
```

## Configuration

### Default Configuration

Location: `~/.lumine/config.yaml`

```yaml
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

projects: []

runtimes:
  php: "8.2-fpm"
  node: "20-alpine"
  python: "3.11-slim"
  bun: latest
  deno: latest
  go: "1.21-alpine"
  rust: latest
```

### Custom Configuration

Edit `~/.lumine/config.yaml`:

```yaml
services:
  - name: mysql-dev
    type: mysql
    version: "8.0"
    port: 3306
    env:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: myapp_dev
      
  - name: mysql-test
    type: mysql
    version: "8.0"
    port: 3307
    env:
      MYSQL_ROOT_PASSWORD: secret
      MYSQL_DATABASE: myapp_test
```

## Next Steps

1. **Read Documentation**
   - [Quick Start Guide](docs/QUICKSTART.md)
   - [Database Guide](docs/DATABASE.md)
   - [Web Servers Guide](docs/WEBSERVERS.md)
   - [Port Management](docs/PORT_MANAGEMENT.md)

2. **Try Examples**
   - Create a Laravel project
   - Set up Next.js with PostgreSQL
   - Configure Nginx reverse proxy

3. **Customize**
   - Add your own services
   - Configure custom ports
   - Set up project templates

4. **Contribute**
   - Report bugs
   - Suggest features
   - Submit pull requests

## Quick Reference

### Makefile Commands

```bash
make help          # Show all commands
make dev           # Run in development mode
make build         # Build binary
make install       # Install to system
make test          # Run tests
make db-setup      # Setup all databases
make db-stop       # Stop databases
make clean         # Clean build artifacts
```

### Keyboard Shortcuts

```
Navigation:
  ↑/↓, j/k    - Move up/down
  h/l         - Switch panels
  Tab         - Next panel
  Enter       - Select

Actions:
  Space       - Select/Deselect
  s           - Start
  x           - Stop
  r           - Restart
  v           - Change version
  delete      - Remove
  c           - Cleanup
  Ctrl+C      - Quit
```

### Docker Commands

```bash
# List containers
docker ps -a --filter "name=lumine-"

# View logs
docker logs lumine-mysql

# Execute command
docker exec -it lumine-mysql mysql -u root -proot

# Stop all
docker stop $(docker ps -q --filter "name=lumine-")

# Remove all
docker rm $(docker ps -aq --filter "name=lumine-")
```

## Support

- GitHub Issues: https://github.com/Araryarch/lumine/issues
- Discussions: https://github.com/Araryarch/lumine/discussions
- Documentation: https://github.com/Araryarch/lumine/tree/main/docs

## License

MIT License - see [LICENSE](LICENSE) file for details.
