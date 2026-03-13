# Port Management Guide

Lumine automatically handles port conflicts and finds alternative ports when needed.

## Automatic Port Detection

When starting a service, Lumine will:

1. ✅ Check if the configured port is available
2. ✅ If taken, automatically find an alternative port
3. ✅ Show notification about port change
4. ✅ Update service configuration with new port

## Default Ports

### Databases
- **MySQL**: 3306 → Alternatives: 3307, 3308, 3309, 33060
- **PostgreSQL**: 5432 → Alternatives: 5433, 5434, 5435, 54320
- **MariaDB**: 3307 → Alternatives: 3308, 3309, 3310
- **MongoDB**: 27017 → Alternatives: 27018, 27019, 27020
- **Redis**: 6379 → Alternatives: 6380, 6381, 6382
- **Elasticsearch**: 9200 → Alternatives: 9201, 9202, 9203

### Web Servers
- **Nginx/Apache**: 80 → Alternatives: 8080, 8000, 8888, 3000
- **HTTP Alt**: 8080 → Alternatives: 8081, 8082, 8083, 8084

### Admin Panels
- **phpMyAdmin**: 8080
- **Adminer**: 8081
- **Mongo Express**: 8082
- **Redis Commander**: 8083
- **pgAdmin**: 8084

## Port Conflict Resolution

### Automatic Resolution

```bash
# Start service with port conflict
make dev

# Lumine will:
# 1. Detect port 3306 is in use
# 2. Try alternative: 3307
# 3. If 3307 is free, use it
# 4. Show: "⚠️  Port 3306 is in use, using alternative port 3307 for mysql"
```

### Manual Port Selection

In TUI:
1. Start service
2. If port conflict detected, dialog appears
3. Choose from alternative ports
4. Or enter custom port
5. Service starts with selected port

### Configuration

Edit `~/.lumine/config.yaml`:

```yaml
services:
  - name: mysql
    type: mysql
    version: "8.0"
    port: 3306  # Will auto-change if taken
```

## Port Checking

### Check Port Availability

```bash
# Using netstat
netstat -tuln | grep :3306

# Using lsof
lsof -i :3306

# Using ss
ss -tuln | grep :3306
```

### Find What's Using a Port

```bash
# Linux
sudo lsof -i :3306

# macOS
sudo lsof -i :3306

# Windows
netstat -ano | findstr :3306
```

### Kill Process on Port

```bash
# Linux/macOS
sudo kill -9 $(lsof -t -i:3306)

# Windows
# Find PID first
netstat -ano | findstr :3306
# Then kill
taskkill /PID <PID> /F
```

## Port Ranges

### Reserved Ranges
- **1-1023**: System/privileged ports (requires root)
- **1024-49151**: Registered ports
- **49152-65535**: Dynamic/private ports

### Lumine Ranges
- **3000-3999**: Web applications
- **5000-5999**: API services
- **8000-8999**: Development servers
- **9000-9999**: Microservices

## Custom Port Configuration

### Set Custom Port

```yaml
services:
  - name: mysql-custom
    type: mysql
    version: "8.0"
    port: 13306  # Custom port
```

### Port Mapping

```yaml
services:
  - name: nginx
    type: nginx
    version: latest
    port: 8080  # Host port
    # Container always uses standard port (80 for nginx)
```

## Troubleshooting

### Port Already in Use

```bash
# Check what's using the port
lsof -i :3306

# Options:
# 1. Stop the service using the port
# 2. Let Lumine use alternative port
# 3. Change Lumine config to use different port
```

### Permission Denied (Port < 1024)

```bash
# Option 1: Use port >= 1024
port: 8080

# Option 2: Run with sudo (not recommended)
sudo lumine

# Option 3: Use authbind (Linux)
authbind --deep lumine
```

### Port Conflict After Restart

```bash
# Clean up old containers
make containers-stop
make containers-remove

# Restart
make db-setup
```

### Multiple Services Same Port

```yaml
# BAD: Both use 3306
services:
  - name: mysql1
    port: 3306
  - name: mysql2
    port: 3306  # Will auto-change to 3307

# GOOD: Different ports
services:
  - name: mysql1
    port: 3306
  - name: mysql2
    port: 3307
```

## Port Management API

### Check Port Availability

```go
import "lumine/internal/port"

portMgr := port.NewManager()

// Check if port is available
if portMgr.IsPortAvailable(3306) {
    fmt.Println("Port 3306 is available")
}

// Find alternative port
altPort, err := portMgr.GetAlternativePort(3306)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Using port %d\n", altPort)
```

### Get Service Port

```go
// Get appropriate port for service
port, err := portMgr.GetServicePort("mysql", 0)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("MySQL will use port %d\n", port)
```

### Scan Port Range

```go
// Find all available ports in range
available := portMgr.ScanPortRange(3306, 3320)
fmt.Printf("Available ports: %v\n", available)
```

## Best Practices

### 1. Use Standard Ports When Possible

```yaml
# Good: Standard ports
mysql: 3306
postgres: 5432
redis: 6379
```

### 2. Document Custom Ports

```yaml
services:
  - name: mysql-dev
    port: 13306  # Custom: avoiding conflict with system MySQL
```

### 3. Use Port Ranges

```yaml
# Development: 3000-3999
# Testing: 4000-4999
# Staging: 5000-5999
```

### 4. Avoid Privileged Ports

```yaml
# Bad: Requires root
port: 80

# Good: No special permissions
port: 8080
```

### 5. Check Before Starting

```bash
# Always check ports before starting
make containers-list
netstat -tuln | grep -E ':(3306|5432|6379)'
```

## Port Monitoring

### Real-time Monitoring

```bash
# Watch port usage
watch -n 1 'netstat -tuln | grep -E ":(3306|5432|6379)"'

# Or with ss
watch -n 1 'ss -tuln | grep -E ":(3306|5432|6379)"'
```

### Log Port Changes

Lumine logs all port changes:

```
⚠️  Port 3306 is in use, using alternative port 3307 for mysql
✓ MySQL started on port 3307
```

### Port Usage Report

```bash
# List all Lumine ports
docker ps --filter "name=lumine-" --format "table {{.Names}}\t{{.Ports}}"
```

## Advanced Configuration

### Port Forwarding

```yaml
services:
  - name: mysql
    type: mysql
    version: "8.0"
    port: 3306
    # Forward to different host port
    host_port: 13306
```

### Multiple Interfaces

```yaml
services:
  - name: mysql-local
    port: 3306
    bind: "127.0.0.1"  # Localhost only
    
  - name: mysql-public
    port: 3307
    bind: "0.0.0.0"  # All interfaces
```

### Dynamic Port Allocation

```yaml
services:
  - name: mysql
    port: 0  # Let Lumine choose any available port
```

## Security Considerations

### 1. Bind to Localhost

```yaml
# Secure: Only accessible from localhost
bind: "127.0.0.1"

# Insecure: Accessible from network
bind: "0.0.0.0"
```

### 2. Use Firewall

```bash
# Allow only specific ports
sudo ufw allow 3306/tcp
sudo ufw enable
```

### 3. Change Default Ports

```yaml
# Change from default to avoid automated scans
mysql: 13306  # Instead of 3306
postgres: 15432  # Instead of 5432
```

## Examples

### Example 1: Multiple MySQL Instances

```yaml
services:
  - name: mysql-dev
    type: mysql
    version: "8.0"
    port: 3306
    
  - name: mysql-test
    type: mysql
    version: "8.0"
    port: 3307  # Auto-assigned if 3306 taken
    
  - name: mysql-staging
    type: mysql
    version: "8.0"
    port: 3308
```

### Example 2: Development Stack

```yaml
services:
  - name: nginx
    port: 8080  # Avoid privileged port 80
    
  - name: mysql
    port: 3306
    
  - name: redis
    port: 6379
    
  - name: elasticsearch
    port: 9200
```

### Example 3: Microservices

```yaml
services:
  - name: api-gateway
    port: 8000
    
  - name: auth-service
    port: 8001
    
  - name: user-service
    port: 8002
    
  - name: order-service
    port: 8003
```
