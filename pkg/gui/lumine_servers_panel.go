package gui

import (
	"fmt"
	"strings"

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
			statusColor := color.FgRed
			if service.Status == "running" {
				statusColor = color.FgGreen
			} else if service.Status == "stopped" {
				statusColor = color.FgYellow
			}

			return []string{
				utils.ColoredString(service.Name, color.FgCyan),
				utils.ColoredString(service.Status, statusColor),
				fmt.Sprintf(":%d", service.Port),
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
		if service.LogPath == "" {
			return "No log file configured"
		}

		content, err := gui.OSCommand.RunCommandWithOutput(fmt.Sprintf("tail -n 50 %s", service.LogPath))
		if err != nil {
			return fmt.Sprintf("Error reading logs: %v", err)
		}

		return content
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
