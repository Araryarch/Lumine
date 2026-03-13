# Lumine Cleanup Tool

Cross-platform cleanup utility for Lumine containers and volumes.

## Features

- ✅ Cross-platform (Windows, Linux, macOS)
- ✅ Interactive menu
- ✅ Auto-backup before destructive operations
- ✅ Color-coded output
- ✅ Safe confirmation prompts
- ✅ Detailed summary

## Installation

### From Source

```bash
# Build
make build-cleanup

# Install
make install-cleanup

# Or install everything
make install-all
```

### From Release

Download the appropriate binary for your platform:

**Linux:**
```bash
# AMD64
wget https://github.com/Araryarch/lumine/releases/latest/download/lumine-cleanup-linux-amd64
chmod +x lumine-cleanup-linux-amd64
sudo mv lumine-cleanup-linux-amd64 /usr/local/bin/lumine-cleanup

# ARM64
wget https://github.com/Araryarch/lumine/releases/latest/download/lumine-cleanup-linux-arm64
chmod +x lumine-cleanup-linux-arm64
sudo mv lumine-cleanup-linux-arm64 /usr/local/bin/lumine-cleanup
```

**macOS:**
```bash
# Intel
curl -L https://github.com/Araryarch/lumine/releases/latest/download/lumine-cleanup-darwin-amd64 -o lumine-cleanup
chmod +x lumine-cleanup
sudo mv lumine-cleanup /usr/local/bin/

# Apple Silicon
curl -L https://github.com/Araryarch/lumine/releases/latest/download/lumine-cleanup-darwin-arm64 -o lumine-cleanup
chmod +x lumine-cleanup
sudo mv lumine-cleanup /usr/local/bin/
```

**Windows:**
```powershell
Invoke-WebRequest -Uri "https://github.com/Araryarch/lumine/releases/latest/download/lumine-cleanup-windows-amd64.exe" -OutFile "lumine-cleanup.exe"
Move-Item lumine-cleanup.exe C:\Windows\System32\
```

## Usage

### Interactive Mode

```bash
lumine-cleanup
```

This will show a menu:

```
🌟 Lumine Cleanup Tool

Select cleanup option:
  1) Stop containers only
  2) Remove containers (keep data)
  3) Remove containers + volumes (DELETE DATA)
  4) Nuclear cleanup (REMOVE EVERYTHING)
  5) Cancel

Enter choice [1-5]:
```

### Via Makefile

```bash
# Interactive
make cleanup-interactive

# Stop containers
make cleanup-stop

# Remove containers
make cleanup-remove

# Nuclear cleanup
make cleanup-nuclear
```

### Programmatic

```bash
# Stop containers
echo "1" | lumine-cleanup

# Remove containers
echo "2" | lumine-cleanup

# Remove with volumes
echo "3" | lumine-cleanup

# Nuclear
echo "4" | lumine-cleanup
```

## Options

### 1. Stop Containers Only

- Stops all running Lumine containers
- Keeps all data
- Safe operation

### 2. Remove Containers (Keep Data)

- Stops all containers
- Removes containers
- Keeps volumes (data preserved)
- Requires confirmation

### 3. Remove Containers + Volumes

- Creates backup first
- Stops all containers
- Removes containers
- Removes volumes (DATA LOSS!)
- Requires typing 'yes'

### 4. Nuclear Cleanup

- Creates backup first
- Removes all containers
- Removes all volumes
- Removes network
- Requires typing 'DESTROY'

## Safety Features

### Confirmation Prompts

- Simple operations: y/N confirmation
- Dangerous operations: Type 'yes'
- Nuclear option: Type 'DESTROY'

### Auto-Backup

Before destructive operations, the tool attempts to:
1. Create MySQL backup
2. Save to timestamped file
3. Continue only if backup succeeds or user confirms

### Color-Coded Output

- 🔵 Cyan: Information
- 🟢 Green: Success
- 🟡 Yellow: Warning
- 🔴 Red: Error/Danger

## Examples

### Safe Cleanup

```bash
# Stop everything
lumine-cleanup
# Choose option 1

# Remove containers (keep data)
lumine-cleanup
# Choose option 2
# Confirm with 'y'
```

### Complete Reset

```bash
# Remove everything
lumine-cleanup
# Choose option 3
# Type 'yes' to confirm

# Backup is created automatically
# All data is removed
```

### Emergency Cleanup

```bash
# Nuclear option
lumine-cleanup
# Choose option 4
# Type 'DESTROY' to confirm

# Everything is removed
# Fresh start
```

## Troubleshooting

### "Failed to connect to Docker"

```bash
# Check Docker is running
docker ps

# Start Docker
# macOS: Open Docker Desktop
# Linux: sudo systemctl start docker
# Windows: Start Docker Desktop
```

### "Permission denied"

```bash
# Linux/macOS: Run with sudo
sudo lumine-cleanup

# Or add user to docker group
sudo usermod -aG docker $USER
newgrp docker
```

### "Backup failed"

The tool will ask if you want to continue without backup.

Options:
1. Cancel and backup manually
2. Continue without backup (risky)

## Integration

### CI/CD

```yaml
# GitHub Actions
- name: Cleanup Lumine
  run: echo "2" | lumine-cleanup
```

### Scripts

```bash
#!/bin/bash
# cleanup-and-restart.sh

# Stop containers
echo "1" | lumine-cleanup

# Wait
sleep 2

# Start fresh
make db-setup
```

### Cron Jobs

```bash
# Daily cleanup of stopped containers
0 2 * * * echo "2" | /usr/local/bin/lumine-cleanup
```

## Development

### Build

```bash
# Build for current platform
go build -o lumine-cleanup ./cmd/cleanup

# Build for all platforms
make build-all
```

### Test

```bash
# Run
./lumine-cleanup

# Test with Docker
docker ps -a --filter "name=lumine-"
```

## See Also

- [Main Documentation](../../README.md)
- [Database Guide](../../docs/DATABASE.md)
- [Cleanup Guide](../../docs/CLEANUP.md)
- [Makefile Commands](../../docs/MAKEFILE.md)
