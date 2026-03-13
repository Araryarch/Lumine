# Quick Start Guide

Get started with Lumine in 5 minutes!

## 1. Install Lumine

```bash
# Linux/macOS
curl -fsSL https://raw.githubusercontent.com/Araryarch/lumine/main/install.sh | bash

# Windows
irm https://raw.githubusercontent.com/Araryarch/lumine/main/install.ps1 | iex
```

## 2. Start Lumine

```bash
lumine
```

## 3. Create Your First Project

### Laravel Project

1. Navigate to sidebar → `➕ New Project`
2. Select `Laravel`
3. Enter project name: `myapp`
4. Wait for project creation
5. Access at `http://myapp.test`

### Next.js Project

1. Navigate to sidebar → `➕ New Project`
2. Select `Next.js`
3. Enter project name: `frontend`
4. Wait for project creation
5. Access at `http://frontend.test`

## 4. Manage Services

### Start Services

1. Navigate to `📦 Services`
2. Use `↑/↓` to select service
3. Press `space` to select
4. Press `s` to start

### Change PHP Version

1. Navigate to `⚙️ Runtimes`
2. Select `PHP`
3. Press `v` for version selector
4. Choose version (e.g., 8.3-fpm)
5. Press `enter` to apply

## 5. Keyboard Shortcuts

### Navigation
- `↑/↓` or `j/k` - Move up/down
- `h/l` - Switch panels
- `tab` - Next panel
- `enter` - Select/Execute

### Actions
- `space` - Select/Deselect
- `s` - Start
- `x` - Stop
- `r` - Restart
- `v` - Change version
- `n` - New project
- `ctrl+c` - Quit

## Common Tasks

### Create Multiple Projects

```bash
# Start Lumine
lumine

# Create backend (Laravel)
→ New Project → Laravel → "backend"

# Create frontend (Next.js)
→ New Project → Next.js → "frontend"

# Access:
# http://backend.test
# http://frontend.test
```

### Switch Runtime Versions

```bash
# Change PHP version
→ Runtimes → PHP → v → Select 8.3-fpm

# Change Node.js version
→ Runtimes → Node.js → v → Select 20-alpine
```

### Start Multiple Services

```bash
# Select multiple services
→ Services → space (on each service) → s (start all)

# Or select all
→ Services → a (select all) → s (start)
```

## Tips & Tricks

### 1. Use Vim Keys
- `j` = down
- `k` = up
- `h` = left panel
- `l` = right panel

### 2. Quick Navigation
- Press `h` to jump to sidebar
- Press `l` to jump to main panel

### 3. Batch Operations
- Press `a` to select all services
- Press `d` to deselect all

### 4. Monitor Logs
- Navigate to `📊 Logs` in sidebar
- Or check detail panel on the right

## Next Steps

- Read [Configuration Guide](CONFIGURATION.md)
- Check [Usage Examples](EXAMPLES.md)
- See [Troubleshooting](TROUBLESHOOTING.md)


## Port Conflicts

If you see a port conflict warning:

```
⚠️  Port 3306 is in use, using alternative port 3307 for mysql
```

Lumine automatically finds an alternative port. You can:

1. **Accept automatic port** - Lumine uses the alternative
2. **Choose from alternatives** - Press `v` to see options
3. **Enter custom port** - Type your preferred port
4. **Free the port** - Stop the service using it

Check what's using a port:
```bash
# Linux/macOS
lsof -i :3306

# Windows
netstat -ano | findstr :3306
```

See [Port Management Guide](PORT_MANAGEMENT.md) for details.
