# Lumine API Documentation

## Overview

Lumine menyediakan API Go yang dapat digunakan untuk membangun aplikasi development environment manager custom. Semua komponen dapat digunakan secara independen atau melalui Orchestrator.

## Core Components

### Orchestrator

Main coordinator yang mengintegrasikan semua komponen.

```go
import "github.com/jesseduffield/lazydocker/pkg/lumine"

// Create orchestrator
orchestrator, err := lumine.NewOrchestrator()
if err != nil {
    log.Fatal(err)
}
defer orchestrator.Close()

// Access components
orchestrator.ServiceManager
orchestrator.VersionManager
orchestrator.ProjectManager
orchestrator.DatabaseManager
orchestrator.ConfigManager
orchestrator.NotificationMgr
orchestrator.DaemonMode
```

## Service Management

### ServiceManager

Mengelola Docker services lifecycle.

```go
// Create service manager
sm, err := lumine.NewServiceManager()
if err != nil {
    log.Fatal(err)
}
defer sm.Close()

// Initialize default services
sm.InitializeDefaultServices()

// Start service with port conflict resolution
port, err := sm.StartServiceWithPortCheck("nginx")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Nginx started on port %d\n", port)

// Stop service
err = sm.StopService("nginx")

// Restart service
err = sm.RestartService("nginx")

// Get service status
running, err := sm.GetServiceStatus("nginx")

// Get service info
info, err := sm.GetServiceInfo("nginx")

// Health check
err = sm.PerformHealthCheck("nginx")
health, exists := sm.GetHealthStatus("nginx")
if exists {
    fmt.Printf("Healthy: %v, Failures: %d\n", health.Healthy, health.FailureCount)
}
```

### PortManager

Mengelola alokasi port dan conflict resolution.

```go
pm := lumine.NewPortManager()

// Allocate port
port, err := pm.AllocatePort("myservice", 8080)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Allocated port: %d\n", port)

// Check if port available
available := pm.IsPortAvailable(8080)

// Release port
pm.ReleasePort(8080)

// Get all allocated ports
ports := pm.GetAllocatedPorts()
```

## Version Management

### VersionManager

Mengelola runtime versions (PHP, Node.js, Python).

```go
vm := lumine.NewVersionManager()

// Get available versions
phpVersions := vm.GetPHPVersions()
nodeVersions := vm.GetNodeVersions()

// Switch PHP version (requires ServiceManager)
err := vm.SwitchPHPVersion("8.2", serviceManager)

// Switch Node.js version
err = vm.SwitchNodeVersion("20", serviceManager)

// Get installed versions
phpVer, err := vm.GetInstalledPHPVersion()
nodeVer, err := vm.GetInstalledNodeVersion()

// Get version info
info := vm.GetVersionInfo()
fmt.Println(info)
```

## Project Management

### ProjectManager

Mengelola web development projects.

```go
pm, err := lumine.NewProjectManager("~/projects")
if err != nil {
    log.Fatal(err)
}

// Create project with virtual host
err = pm.CreateProjectWithVHost(
    "my-app",
    lumine.ProjectTypeLaravel,
    "8.3",  // PHP version
    "",     // Node version (not needed for Laravel)
    "nginx", // Web server
)

// Enable SSL
err = pm.EnableSSL("my-app")

// Expose via tunnel
publicURL, err := pm.ExposeTunnel("my-app", lumine.TunnelTypeCloudflare)
fmt.Printf("Public URL: %s\n", publicURL)

// Stop tunnel
err = pm.StopTunnel("my-app")

// List projects
projects := pm.ListProjects()
for _, p := range projects {
    fmt.Printf("%s: %s (%s)\n", p.Name, p.URL, p.Type)
}

// Scan existing projects
err = pm.ScanProjects()

// Delete project with cleanup
err = pm.DeleteProjectWithCleanup("my-app", "nginx")

// Check dependencies
deps := pm.CheckDependencies()
```

### VHostManager

Mengelola virtual host configuration.

```go
vhm := lumine.NewVHostManager()

// Create virtual host
err := vhm.CreateVirtualHost("my-app", "/path/to/project", "nginx")

// Delete virtual host
err = vhm.DeleteVirtualHost("my-app", "nginx")

// Check admin privileges
hasPrivileges := vhm.CheckAdminPrivileges()
```

### SSLManager

Mengelola SSL certificates.

```go
sm, err := lumine.NewSSLManager()
if err != nil {
    log.Fatal(err)
}

// Generate certificate
certPath, keyPath, err := sm.GenerateCertificate("my-app.test")
fmt.Printf("Cert: %s\nKey: %s\n", certPath, keyPath)

// Get certificate path
certPath, keyPath, exists := sm.GetCertificatePath("my-app.test")

// Delete certificate
err = sm.DeleteCertificate("my-app.test")
```

### TunnelManager

Mengelola localhost tunneling.

```go
tm := lumine.NewTunnelManager()

// Check if tunnel tool installed
installed := tm.CheckTunnelToolInstalled(lumine.TunnelTypeCloudflare)

// Create tunnel
tunnel, err := tm.CreateTunnel("my-app", 8080, lumine.TunnelTypeCloudflare)
if err != nil {
    log.Fatal(err)
}

// Wait for public URL
time.Sleep(5 * time.Second)
fmt.Printf("Public URL: %s\n", tunnel.PublicURL)

// Stop tunnel
err = tm.StopTunnel("my-app")

// Get tunnel info
tunnel, exists := tm.GetTunnel("my-app")
```

## Database Management

### DatabaseManager

Mengelola database operations.

```go
dm := lumine.NewDatabaseManager()
defer dm.Close()

// Add connection
conn := &lumine.DatabaseConnection{
    Name:     "local-mysql",
    Type:     lumine.DatabaseTypeMySQL,
    Host:     "localhost",
    Port:     3306,
    Username: "root",
    Password: "root",
    Database: "mysql",
}
dm.AddConnection(conn)

// Switch connection
err := dm.SwitchConnection("local-mysql")

// Create database
err = dm.CreateDatabase("my_app_db")

// Drop database
err = dm.DropDatabase("my_app_db")

// List databases
databases, err := dm.ListDatabases()

// Backup database
backupPath, err := dm.BackupDatabase("my_app_db")

// Get query logs
logs := dm.GetQueryLog()
for _, log := range logs {
    fmt.Printf("[%v] %s (%.2fms)\n", 
        log.Timestamp, 
        log.Query, 
        float64(log.ExecutionTime.Microseconds())/1000)
}

// Get slow queries
slowQueries := dm.GetSlowQueries(1 * time.Second)

// Get error queries
errorQueries := dm.GetErrorQueries()

// Clear logs
dm.ClearQueryLog()
```

## Daemon Mode

### DaemonMode

Background service monitoring dan auto-restart.

```go
sm, _ := lumine.NewServiceManager()
daemon, err := lumine.NewDaemonMode(sm)
if err != nil {
    log.Fatal(err)
}

// Start daemon
err = daemon.Start()

// Check if running
running := daemon.IsRunning()

// Stop daemon
err = daemon.Stop()

// Get log path
logPath := daemon.GetLogPath()
```

## Utilities

### DependencyChecker

Check development tools.

```go
dc := lumine.NewDependencyChecker()

// Check all dependencies
deps := dc.CheckAll()
for name, dep := range deps {
    if dep.Installed {
        fmt.Printf("✓ %s v%s\n", name, dep.Version)
    } else {
        fmt.Printf("✗ %s not installed\n", name)
    }
}

// Check specific tool
status := dc.Check("composer")

// Refresh status
deps = dc.Refresh()

// Quick check
installed := dc.IsInstalled("npm")
```

### NotificationManager

Toast notifications untuk TUI.

```go
nm := lumine.NewNotificationManager()

// Show notifications
nm.ShowSuccess("Operation completed")
nm.ShowError("Operation failed")
nm.ShowWarning("Port changed to 8080")
nm.ShowInfo("Service started")

// Get active notifications
notifications := nm.GetActive()
for _, n := range notifications {
    fmt.Printf("[%s] %s %s\n", n.Type, n.GetIcon(), n.Message)
}

// Dismiss notification
nm.Dismiss("notif-1")

// Clear all
nm.Clear()
```

### ConfigManager

Application configuration.

```go
cm, err := lumine.NewConfigManager()
if err != nil {
    log.Fatal(err)
}

// Get config
config := cm.Get()
fmt.Printf("PHP: %s, Node: %s\n", 
    config.DefaultPHPVersion, 
    config.DefaultNodeVersion)

// Update config
config.DefaultPHPVersion = "8.2"
err = cm.Update(config)

// Export config
err = cm.Export("/path/to/backup.yaml")

// Import config
err = cm.Import("/path/to/backup.yaml")
```

### TemplateManager

Service templates management.

```go
tm, err := lumine.NewTemplateManager()
if err != nil {
    log.Fatal(err)
}

// Create custom template
template := &lumine.ServiceTemplate{
    Name:         "my-service",
    Description:  "Custom service",
    Image:        "myimage:latest",
    Port:         8080,
    InternalPort: 8080,
    Environment: map[string]string{
        "ENV_VAR": "value",
    },
}
err = tm.CreateTemplate(template)

// Get template
template, err = tm.GetTemplate("my-service")

// List all templates
templates := tm.ListTemplates()

// Get stack template
stack, err := tm.GetStack("lamp")

// Create service from template
service, err := tm.CreateServiceFromTemplate("my-service", "my-instance")
```

## Error Handling

Semua fungsi mengembalikan error yang descriptive. Best practice:

```go
if err := orchestrator.StartService("nginx"); err != nil {
    // Error sudah di-log dan notification sudah ditampilkan
    // Handle error sesuai kebutuhan aplikasi
    log.Printf("Failed to start nginx: %v", err)
    return err
}
```

## Thread Safety

Semua managers menggunakan mutex untuk thread-safe operations:
- ServiceManager: RWMutex untuk concurrent reads
- PortManager: RWMutex untuk port allocation
- DatabaseManager: RWMutex untuk connection management
- NotificationManager: RWMutex untuk notification list

## Context Management

Orchestrator menggunakan context untuk graceful shutdown:

```go
orchestrator, _ := lumine.NewOrchestrator()

// Context akan di-cancel saat Close() dipanggil
defer orchestrator.Close()

// Semua goroutines akan berhenti dengan graceful
```

## Best Practices

1. **Always defer Close()**: Pastikan resources di-cleanup
2. **Check errors**: Jangan ignore error returns
3. **Use Orchestrator**: Untuk high-level operations, gunakan Orchestrator
4. **Direct managers**: Untuk fine-grained control, akses managers langsung
5. **Thread-safe**: Semua operations sudah thread-safe, aman untuk concurrent use

## Examples

Lihat `examples/lumine_usage.go` untuk contoh lengkap penggunaan semua fitur.
