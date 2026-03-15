# GUI Integration Guide

## Overview

Dokumen ini menjelaskan cara mengintegrasikan Lumine components dengan TUI (gocui).

## Architecture

```
┌─────────────────────────────────────────┐
│           GUI Layer (gocui)             │
│  - Views & Panels                       │
│  - Keyboard bindings                    │
│  - Rendering                            │
└──────────────┬──────────────────────────┘
               │
┌──────────────▼──────────────────────────┐
│         Orchestrator                    │
│  - Coordinates all managers             │
│  - Handles notifications                │
│  - Manages lifecycle                    │
└──────────────┬──────────────────────────┘
               │
    ┌──────────┼──────────┬──────────┐
    │          │          │          │
┌───▼───┐ ┌───▼───┐ ┌───▼───┐ ┌───▼───┐
│Service│ │Version│ │Project│ │Database│
│Manager│ │Manager│ │Manager│ │Manager │
└───────┘ └───────┘ └───────┘ └───────┘
```

## Integration Steps

### 1. Initialize Orchestrator in GUI

```go
// In pkg/gui/gui.go

type Gui struct {
    g             *gocui.Gui
    // ... existing fields
    Orchestrator  *lumine.Orchestrator  // Add this
}

// In NewGui function
func NewGui(...) (*Gui, error) {
    // ... existing code
    
    orchestrator, err := lumine.NewOrchestrator()
    if err != nil {
        return nil, err
    }
    
    gui := &Gui{
        // ... existing fields
        Orchestrator: orchestrator,
    }
    
    return gui, nil
}
```

### 2. Add Lumine Panels

```go
// In pkg/gui/panels.go

type Panels struct {
    // ... existing panels
    LumineServices   *panels.SideListPanel[*lumine.Service]
    LumineProjects   *panels.SideListPanel[*lumine.Project]
    LumineDatabases  *panels.SideListPanel[string]
}
```

### 3. Create Service Panel

```go
// In pkg/gui/lumine_services_panel.go

func (gui *Gui) getLumineServicesPanel() *panels.SideListPanel[*lumine.Service] {
    return panels.NewSideListPanel(
        panels.SideListPanelOpts[*lumine.Service]{
            GetTableCells: gui.getLumineServiceTableCells,
            OnSelect:      gui.onLumineServiceSelect,
            OnRender:      gui.onLumineServicesRender,
            // ... other options
        },
    )
}

func (gui *Gui) getLumineServiceTableCells(service *lumine.Service) []string {
    status := "stopped"
    statusColor := "red"
    
    if running, _ := gui.Orchestrator.ServiceManager.GetServiceStatus(service.Name); running {
        status = "running"
        statusColor = "green"
    }
    
    // Check health
    if health, exists := gui.Orchestrator.ServiceManager.GetHealthStatus(service.Name); exists {
        if !health.Healthy {
            status = "unhealthy"
            statusColor = "yellow"
        }
    }
    
    return []string{
        service.Name,
        fmt.Sprintf("[%s]%s[-]", statusColor, status),
        fmt.Sprintf("%d", service.Port),
        service.Image,
    }
}
```

### 4. Add Keybindings

```go
// In pkg/gui/keybindings.go

func (gui *Gui) keybindings(g *gocui.Gui) error {
    // ... existing keybindings
    
    // Lumine service keybindings
    if err := g.SetKeybinding("lumineServices", nil, 's', gocui.ModNone, gui.handleStartLumineService); err != nil {
        return err
    }
    
    if err := g.SetKeybinding("lumineServices", nil, 'S', gocui.ModNone, gui.handleStopLumineService); err != nil {
        return err
    }
    
    if err := g.SetKeybinding("lumineServices", nil, 'r', gocui.ModNone, gui.handleRestartLumineService); err != nil {
        return err
    }
    
    if err := g.SetKeybinding("lumineServices", nil, 'v', gocui.ModNone, gui.handleSwitchVersion); err != nil {
        return err
    }
    
    // ... more keybindings
    
    return nil
}
```

### 5. Implement Handlers

```go
// In pkg/gui/lumine_services_panel.go

func (gui *Gui) handleStartLumineService(g *gocui.Gui, v *gocui.View) error {
    service := gui.Panels.LumineServices.List.GetSelectedItem()
    if service == nil {
        return nil
    }
    
    go func() {
        if err := gui.Orchestrator.StartService(service.Name); err != nil {
            gui.ErrorChan <- err
        }
        gui.g.Update(func(*gocui.Gui) error {
            return gui.refreshLumineServices()
        })
    }()
    
    return nil
}

func (gui *Gui) handleStopLumineService(g *gocui.Gui, v *gocui.View) error {
    service := gui.Panels.LumineServices.List.GetSelectedItem()
    if service == nil {
        return nil
    }
    
    return gui.createConfirmationPanel(
        "Stop Service",
        fmt.Sprintf("Are you sure you want to stop %s?", service.Name),
        func(g *gocui.Gui, v *gocui.View) error {
            go func() {
                if err := gui.Orchestrator.StopService(service.Name); err != nil {
                    gui.ErrorChan <- err
                }
                gui.g.Update(func(*gocui.Gui) error {
                    return gui.refreshLumineServices()
                })
            }()
            return nil
        },
        nil,
    )
}

func (gui *Gui) handleSwitchVersion(g *gocui.Gui, v *gocui.View) error {
    service := gui.Panels.LumineServices.List.GetSelectedItem()
    if service == nil {
        return nil
    }
    
    // Show version selection menu
    if service.Type == lumine.ServiceTypePHPFPM {
        return gui.showPHPVersionMenu()
    } else if service.Name == "nodejs" {
        return gui.showNodeVersionMenu()
    }
    
    return nil
}
```

### 6. Render Notifications

```go
// In pkg/gui/layout.go

func (gui *Gui) layout(g *gocui.Gui) error {
    // ... existing layout code
    
    // Render notifications at top-right
    if err := gui.renderNotifications(g); err != nil {
        return err
    }
    
    return nil
}

func (gui *Gui) renderNotifications(g *gocui.Gui) error {
    notifications := gui.Orchestrator.NotificationMgr.GetActive()
    
    maxX, _ := g.Size()
    startY := 1
    
    for i, notif := range notifications {
        if i >= 3 { // Max 3 notifications
            break
        }
        
        width := 50
        height := 3
        x0 := maxX - width - 2
        y0 := startY + (i * (height + 1))
        x1 := maxX - 2
        y1 := y0 + height
        
        v, err := g.SetView(fmt.Sprintf("notification-%d", i), x0, y0, x1, y1)
        if err != nil && err != gocui.ErrUnknownView {
            return err
        }
        
        v.Frame = true
        v.Clear()
        
        // Set color based on type
        switch notif.Type {
        case lumine.NotificationSuccess:
            v.FgColor = gocui.ColorGreen
        case lumine.NotificationError:
            v.FgColor = gocui.ColorRed
        case lumine.NotificationWarning:
            v.FgColor = gocui.ColorYellow
        case lumine.NotificationInfo:
            v.FgColor = gocui.ColorCyan
        }
        
        fmt.Fprintf(v, "%s %s", notif.GetIcon(), notif.Message)
    }
    
    return nil
}
```

### 7. Add Refresh Logic

```go
// In pkg/gui/gui.go

func (gui *Gui) refreshLumineServices() error {
    services := gui.Orchestrator.ServiceManager.GetAllServices()
    
    var serviceList []*lumine.Service
    for _, service := range services {
        serviceList = append(serviceList, service)
    }
    
    gui.Panels.LumineServices.List.SetItems(serviceList)
    return nil
}

func (gui *Gui) refreshLumineProjects() error {
    projects := gui.Orchestrator.ProjectManager.ListProjects()
    gui.Panels.LumineProjects.List.SetItems(projects)
    return nil
}
```

### 8. Add to Main Loop

```go
// In pkg/gui/gui.go Run() function

go func() {
    throttledRefresh.Trigger()
    
    // ... existing refresh calls
    
    // Add Lumine refreshes
    gui.goEvery(time.Second*2, gui.refreshLumineServices)
    gui.goEvery(time.Second*5, gui.refreshLumineProjects)
    gui.goEvery(time.Second*1, gui.refreshNotifications)
}()
```

## View Names

Tambahkan view names untuk Lumine panels:

```go
const (
    // ... existing views
    ViewLumineServices  = "lumineServices"
    ViewLumineProjects  = "lumineProjects"
    ViewLumineDatabases = "lumineDatabases"
    ViewLumineLogs      = "lumineLogs"
)
```

## Menu Items

Tambahkan menu items untuk Lumine operations:

```go
func (gui *Gui) getLumineServiceMenuItems() []*types.MenuItem {
    service := gui.Panels.LumineServices.List.GetSelectedItem()
    
    return []*types.MenuItem{
        {
            Label: "Start",
            OnPress: func() error {
                return gui.handleStartLumineService(gui.g, nil)
            },
        },
        {
            Label: "Stop",
            OnPress: func() error {
                return gui.handleStopLumineService(gui.g, nil)
            },
        },
        {
            Label: "Restart",
            OnPress: func() error {
                return gui.handleRestartLumineService(gui.g, nil)
            },
        },
        {
            Label: "Change Version",
            OnPress: func() error {
                return gui.handleSwitchVersion(gui.g, nil)
            },
        },
        {
            Label: "Health Check",
            OnPress: func() error {
                return gui.handleHealthCheck(gui.g, nil)
            },
        },
    }
}
```

## Testing Integration

```go
// Test orchestrator initialization
func TestGuiWithOrchestrator(t *testing.T) {
    gui, err := NewGui(...)
    assert.NoError(t, err)
    assert.NotNil(t, gui.Orchestrator)
    
    defer gui.Orchestrator.Close()
}

// Test service operations
func TestServiceOperations(t *testing.T) {
    gui, _ := NewGui(...)
    defer gui.Orchestrator.Close()
    
    err := gui.Orchestrator.StartService("nginx")
    assert.NoError(t, err)
    
    running, _ := gui.Orchestrator.ServiceManager.GetServiceStatus("nginx")
    assert.True(t, running)
}
```

## Performance Considerations

1. **Lazy Loading**: Load services/projects on-demand
2. **Throttled Refresh**: Use throttle untuk prevent excessive updates
3. **Background Operations**: Run heavy operations di goroutines
4. **Cache Status**: Cache service status untuk reduce Docker API calls

## Next Steps

1. Implement panels di `pkg/gui/lumine_*.go`
2. Add keybindings di `pkg/gui/keybindings.go`
3. Update layout di `pkg/gui/layout.go`
4. Add menu items di `pkg/gui/menu_panel.go`
5. Test integration dengan existing lazydocker features
