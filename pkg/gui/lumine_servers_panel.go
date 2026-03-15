package gui

import (
	"fmt"

	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/gui/types"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/tasks"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getLumineServersPanel() *panels.SideListPanel[*lumine.Service] {
	return &panels.SideListPanel[*lumine.Service]{
		ContextState: &panels.ContextState[*lumine.Service]{
			GetMainTabs: func() []panels.MainTab[*lumine.Service] {
				return []panels.MainTab[*lumine.Service]{
					{
						Key:    "info",
						Title:  "Server Info",
						Render: gui.renderLumineServerInfo,
					},
					{
						Key:    "config",
						Title:  "Configuration",
						Render: gui.renderLumineServerConfig,
					},
					{
						Key:    "logs",
						Title:  "Logs",
						Render: gui.renderLumineServerLogs,
					},
				}
			},
			GetItemContextCacheKey: func(service *lumine.Service) string {
				return "lumine-server-" + service.Name + "-" + service.Status
			},
		},
		ListPanel: panels.ListPanel[*lumine.Service]{
			List: panels.NewFilteredList[*lumine.Service](),
			View: gui.Views.LumineServers,
		},
		NoItemsMessage: "No servers",
		Gui:            gui.intoInterface(),
		Sort: func(a *lumine.Service, b *lumine.Service) bool {
			if a.Status == "running" && b.Status != "running" {
				return true
			}
			if a.Status != "running" && b.Status == "running" {
				return false
			}
			return a.Name < b.Name
		},
		GetTableCells: func(service *lumine.Service) []string {
			statusText := "inactive"
			statusColor := color.FgRed
			if service.Status == "running" {
				statusText = "active"
				statusColor = color.FgGreen
			}

			displayName := service.DisplayName
			if displayName == "" {
				displayName = service.Name
			}

			return []string{
				utils.ColoredString(displayName, color.FgCyan),
				utils.ColoredString(statusText, statusColor),
			}
		},
	}
}

func (gui *Gui) renderLumineServerInfo(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(service.Name, color.FgCyan) + "\n"
		output += utils.WithPadding("Status: ", 15) + gui.getColoredStatus(service.Status) + "\n"
		output += utils.WithPadding("Port: ", 15) + fmt.Sprintf("%d", service.Port) + "\n"
		output += utils.WithPadding("Type: ", 15) + string(service.Type) + "\n"
		output += utils.WithPadding("Image: ", 15) + service.Image + "\n"

		if service.ConfigPath != "" {
			output += utils.WithPadding("Config: ", 15) + service.ConfigPath + "\n"
		}

		return output
	})
}

func (gui *Gui) renderLumineServerConfig(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if service.ConfigPath == "" {
			return "No configuration file"
		}

		content, err := gui.OSCommand.RunCommandWithOutput(fmt.Sprintf("cat %s", service.ConfigPath))
		if err != nil {
			return fmt.Sprintf("Error reading config: %v", err)
		}

		return content
	})
}

func (gui *Gui) renderLumineServerLogs(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if !service.Running {
			return "Service is not running. Start the service to view logs."
		}

		logs, err := gui.Orchestrator.ServiceManager.GetServiceLogs(service.Name, 100)
		if err != nil {
			return fmt.Sprintf("Error reading logs: %v", err)
		}

		if logs == "" {
			return "No logs available"
		}

		return logs
	})
}

// Keybinding handlers
func (gui *Gui) handleLumineServerStart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status == "running" {
		return gui.createErrorPanel("Server is already running")
	}

	return gui.WithWaitingStatus("Starting server...", func() error {
		if err := gui.Orchestrator.StartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineServers()
	})
}

func (gui *Gui) handleLumineServerStop(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status != "running" {
		return gui.createErrorPanel("Server is not running")
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Stop %s?", service.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping server...", func() error {
			if err := gui.Orchestrator.StopService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineServers()
		})
	}, nil)
}

func (gui *Gui) handleLumineServerRestart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.WithWaitingStatus("Restarting server...", func() error {
		if err := gui.Orchestrator.RestartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineServers()
	})
}

func (gui *Gui) handleLumineServerAdd(g *gocui.Gui, v *gocui.View) error {
	return gui.createPromptPanel("Server Name (e.g., lumine-caddy)", func(g *gocui.Gui, v *gocui.View) error {
		serviceName := gui.trimmedContent(v)
		if serviceName == "" {
			return gui.createErrorPanel("Server name cannot be empty")
		}

		return gui.createPromptPanel("Docker Image (e.g., caddy:alpine)", func(g *gocui.Gui, v *gocui.View) error {
			image := gui.trimmedContent(v)
			if image == "" {
				return gui.createErrorPanel("Image cannot be empty")
			}

			return gui.createPromptPanel("Port (e.g., 8080)", func(g *gocui.Gui, v *gocui.View) error {
				portStr := gui.trimmedContent(v)
				if portStr == "" {
					return gui.createErrorPanel("Port cannot be empty")
				}

				var port int
				if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
					return gui.createErrorPanel("Invalid port number")
				}

				customService := &lumine.CustomService{
					Name:         serviceName,
					Type:         "server",
					Image:        image,
					Port:         port,
					InternalPort: port,
					Enabled:      true,
					Environment:  make(map[string]string),
					Volumes:      make(map[string]string),
				}

				return gui.WithWaitingStatus("Adding server...", func() error {
					if err := gui.Orchestrator.AddCustomService(customService); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					return gui.refreshLumineServers()
				})
			})
		})
	})
}

func (gui *Gui) handleLumineServerRemove(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Remove server '%s'?", service.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Removing server...", func() error {
			if err := gui.Orchestrator.RemoveCustomService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineServers()
		})
	}, nil)
}

func (gui *Gui) refreshLumineServers() error {
	if gui.Orchestrator == nil || gui.Panels.LumineServers == nil {
		return nil
	}

	servers := gui.Orchestrator.ServiceManager.ListServerServices()
	gui.Panels.LumineServers.SetItems(servers)
	return gui.Panels.LumineServers.RerenderList()
}

// Handler for opening settings
func (gui *Gui) handleLumineSettings(g *gocui.Gui, v *gocui.View) error {
	config := gui.Orchestrator.ConfigManager.Get()

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"Default PHP Version", config.DefaultPHPVersion},
			OnPress: func() error {
				return gui.handleSettingPHPVersion()
			},
		},
		{
			LabelColumns: []string{"Default Node Version", config.DefaultNodeVersion},
			OnPress: func() error {
				return gui.handleSettingNodeVersion()
			},
		},
		{
			LabelColumns: []string{"Preferred Web Server", config.PreferredWebServer},
			OnPress: func() error {
				return gui.handleSettingWebServer()
			},
		},
		{
			LabelColumns: []string{"Auto Start Services", fmt.Sprintf("%v", config.AutoStartServices)},
			OnPress: func() error {
				return gui.handleToggleAutoStart()
			},
		},
		{
			LabelColumns: []string{"Enable Auto SSL", fmt.Sprintf("%v", config.EnableAutoSSL)},
			OnPress: func() error {
				return gui.handleToggleAutoSSL()
			},
		},
		{
			LabelColumns: []string{"Projects Directory", config.ProjectsDirectory},
			OnPress: func() error {
				return gui.handleSettingProjectsDir()
			},
		},
	}

	return gui.Menu(CreateMenuOptions{
		Title: "Lumine Settings",
		Items: menuItems,
	})
}

func (gui *Gui) handleSettingPHPVersion() error {
	versions := []string{"7.4", "8.0", "8.1", "8.2", "8.3"}
	menuItems := make([]*types.MenuItem, len(versions))

	for i, version := range versions {
		v := version
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{fmt.Sprintf("PHP %s", v)},
			OnPress: func() error {
				config := gui.Orchestrator.ConfigManager.Get()
				config.DefaultPHPVersion = v
				return gui.Orchestrator.UpdateConfig(config)
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: "Select Default PHP Version",
		Items: menuItems,
	})
}

func (gui *Gui) handleSettingNodeVersion() error {
	versions := []string{"16", "18", "20", "21"}
	menuItems := make([]*types.MenuItem, len(versions))

	for i, version := range versions {
		v := version
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{fmt.Sprintf("Node.js %s", v)},
			OnPress: func() error {
				config := gui.Orchestrator.ConfigManager.Get()
				config.DefaultNodeVersion = v
				return gui.Orchestrator.UpdateConfig(config)
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: "Select Default Node.js Version",
		Items: menuItems,
	})
}

func (gui *Gui) handleSettingWebServer() error {
	servers := []string{"nginx", "apache", "caddy"}
	menuItems := make([]*types.MenuItem, len(servers))

	for i, server := range servers {
		s := server
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{s},
			OnPress: func() error {
				config := gui.Orchestrator.ConfigManager.Get()
				config.PreferredWebServer = s
				return gui.Orchestrator.UpdateConfig(config)
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: "Select Preferred Web Server",
		Items: menuItems,
	})
}

func (gui *Gui) handleToggleAutoStart() error {
	config := gui.Orchestrator.ConfigManager.Get()
	config.AutoStartServices = !config.AutoStartServices
	return gui.Orchestrator.UpdateConfig(config)
}

func (gui *Gui) handleToggleAutoSSL() error {
	config := gui.Orchestrator.ConfigManager.Get()
	config.EnableAutoSSL = !config.EnableAutoSSL
	return gui.Orchestrator.UpdateConfig(config)
}

func (gui *Gui) handleSettingProjectsDir() error {
	return gui.createPromptPanel("Projects Directory Path", func(g *gocui.Gui, v *gocui.View) error {
		path := gui.trimmedContent(v)
		if path == "" {
			return gui.createErrorPanel("Path cannot be empty")
		}

		config := gui.Orchestrator.ConfigManager.Get()
		config.ProjectsDirectory = path
		return gui.Orchestrator.UpdateConfig(config)
	})
}

// Handler for editing server settings
func (gui *Gui) handleLumineServerEdit(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"Edit Port", fmt.Sprintf("Current: %d", service.Port)},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("New Port for %s", service.Name), func(g *gocui.Gui, v *gocui.View) error {
					portStr := gui.trimmedContent(v)
					if portStr == "" {
						return gui.createErrorPanel("Port cannot be empty")
					}

					var port int
					if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
						return gui.createErrorPanel("Invalid port number")
					}

					// Update service port
					service.Port = port
					gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Port updated to %d", port))
					return gui.refreshLumineServers()
				})
			},
		},
		{
			LabelColumns: []string{"Edit Image", fmt.Sprintf("Current: %s", service.Image)},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("New Image for %s", service.Name), func(g *gocui.Gui, v *gocui.View) error {
					image := gui.trimmedContent(v)
					if image == "" {
						return gui.createErrorPanel("Image cannot be empty")
					}

					service.Image = image
					gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Image updated to %s", image))
					return gui.refreshLumineServers()
				})
			},
		},
		{
			LabelColumns: []string{"Edit Config Path", service.ConfigPath},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("Config Path for %s", service.Name), func(g *gocui.Gui, v *gocui.View) error {
					path := gui.trimmedContent(v)
					service.ConfigPath = path
					gui.Orchestrator.NotificationMgr.ShowSuccess("Config path updated")
					return gui.refreshLumineServers()
				})
			},
		},
		{
			LabelColumns: []string{"Edit Log Path", service.LogPath},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("Log Path for %s", service.Name), func(g *gocui.Gui, v *gocui.View) error {
					path := gui.trimmedContent(v)
					service.LogPath = path
					gui.Orchestrator.NotificationMgr.ShowSuccess("Log path updated")
					return gui.refreshLumineServers()
				})
			},
		},
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Edit %s Settings", service.Name),
		Items: menuItems,
	})
}

// Handler for switching server version
func (gui *Gui) handleLumineServerVersionSwitch(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	versions := gui.Orchestrator.GetAvailableVersions(string(service.Type))
	if len(versions) == 0 {
		return gui.createErrorPanel("No versions available for this service")
	}

	menuItems := make([]*types.MenuItem, len(versions))
	for i, version := range versions {
		v := version
		currentIndicator := ""
		if v == service.Version {
			currentIndicator = " (current)"
		}
		
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{fmt.Sprintf("%s%s", v, currentIndicator)},
			OnPress: func() error {
				return gui.WithWaitingStatus("Switching version...", func() error {
					if err := gui.Orchestrator.SwitchServiceVersion(service.Name, v); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					return gui.refreshLumineServers()
				})
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Switch %s Version", service.Name),
		Items: menuItems,
	})
}


// Handler for executing command in server container
func (gui *Gui) handleLumineServerExec(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServers.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status != "running" {
		return gui.createErrorPanel("Server is not running. Start it first.")
	}

	return gui.createPromptPanel("Command to execute", func(g *gocui.Gui, v *gocui.View) error {
		command := gui.trimmedContent(v)
		if command == "" {
			return gui.createErrorPanel("Command cannot be empty")
		}

		return gui.WithWaitingStatus(fmt.Sprintf("Executing: %s", command), func() error {
			output, err := gui.Orchestrator.ServiceManager.ExecuteCommand(service.Name, command)
			
			if err != nil {
				return gui.createErrorPanel(fmt.Sprintf("Command failed: %v\n\nOutput:\n%s", err, output))
			}
			
			return gui.createInfoPanel("Command Output", output)
		})
	})
}
