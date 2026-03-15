package gui

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazydocker/pkg/gui/panels"
	"github.com/jesseduffield/lazydocker/pkg/gui/types"
	"github.com/jesseduffield/lazydocker/pkg/lumine"
	"github.com/jesseduffield/lazydocker/pkg/tasks"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func (gui *Gui) getLumineServicesPanel() *panels.SideListPanel[*lumine.Service] {
	return &panels.SideListPanel[*lumine.Service]{
		ContextState: &panels.ContextState[*lumine.Service]{
			GetMainTabs: func() []panels.MainTab[*lumine.Service] {
				return []panels.MainTab[*lumine.Service]{
					{
						Key:    "info",
						Title:  "Service Info",
						Render: gui.renderLumineServiceInfo,
					},
					{
						Key:    "config",
						Title:  "Configuration",
						Render: gui.renderLumineServiceConfig,
					},
					{
						Key:    "health",
						Title:  "Health Status",
						Render: gui.renderLumineServiceHealth,
					},
				}
			},
			GetItemContextCacheKey: func(service *lumine.Service) string {
				return "lumine-service-" + service.Name + "-" + service.Status
			},
		},
		ListPanel: panels.ListPanel[*lumine.Service]{
			List: panels.NewFilteredList[*lumine.Service](),
			View: gui.Views.LumineServices,
		},
		NoItemsMessage: "No Lumine services",
		Gui:            gui.intoInterface(),
		Sort: func(a *lumine.Service, b *lumine.Service) bool {
			// Sort by status (running first) then by name
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
				service.Version,
			}
		},
	}
}

func (gui *Gui) renderLumineServiceInfo(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(service.Name, color.FgCyan) + "\n"
		output += utils.WithPadding("Status: ", 15) + gui.getColoredStatus(service.Status) + "\n"
		output += utils.WithPadding("Port: ", 15) + fmt.Sprintf("%d", service.Port) + "\n"
		output += utils.WithPadding("Version: ", 15) + service.Version + "\n"
		output += utils.WithPadding("Type: ", 15) + string(service.Type) + "\n"
		output += utils.WithPadding("Command: ", 15) + service.Command + "\n"
		output += utils.WithPadding("PID: ", 15) + fmt.Sprintf("%d", service.PID) + "\n"
		
		if service.ConfigPath != "" {
			output += utils.WithPadding("Config: ", 15) + service.ConfigPath + "\n"
		}
		
		if service.LogPath != "" {
			output += utils.WithPadding("Logs: ", 15) + service.LogPath + "\n"
		}
		
		return output
	})
}

func (gui *Gui) renderLumineServiceConfig(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if service.ConfigPath == "" {
			return "No configuration file"
		}
		
		// Read config file content
		content, err := gui.OSCommand.RunCommandWithOutput(fmt.Sprintf("cat %s", service.ConfigPath))
		if err != nil {
			return fmt.Sprintf("Error reading config: %v", err)
		}
		
		return content
	})
}

func (gui *Gui) renderLumineServiceHealth(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		health := gui.Orchestrator.ServiceManager.CheckHealth(service.Name)
		
		output := ""
		output += utils.WithPadding("Service: ", 15) + service.Name + "\n"
		output += utils.WithPadding("Healthy: ", 15) + fmt.Sprintf("%v", health.Healthy) + "\n"
		output += utils.WithPadding("Uptime: ", 15) + health.Uptime.String() + "\n"
		output += utils.WithPadding("Last Check: ", 15) + health.LastCheck.Format("2006-01-02 15:04:05") + "\n"
		
		if health.Error != "" {
			output += "\n" + utils.ColoredString("Error: "+health.Error, color.FgRed) + "\n"
		}
		
		return output
	})
}

func (gui *Gui) getColoredStatus(status string) string {
	switch status {
	case "running":
		return utils.ColoredString(status, color.FgGreen)
	case "stopped":
		return utils.ColoredString(status, color.FgYellow)
	case "error":
		return utils.ColoredString(status, color.FgRed)
	default:
		return status
	}
}

// Keybinding handlers for Lumine services
func (gui *Gui) handleLumineServiceStart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServices.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	if service.Status == "running" {
		return gui.createErrorPanel("Service is already running")
	}
	
	return gui.WithWaitingStatus("Starting service...", func() error {
		if err := gui.Orchestrator.StartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineServices()
	})
}

func (gui *Gui) handleLumineServiceStop(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServices.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	if service.Status != "running" {
		return gui.createErrorPanel("Service is not running")
	}
	
	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Stop %s?", service.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping service...", func() error {
			if err := gui.Orchestrator.StopService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineServices()
		})
	}, nil)
}

func (gui *Gui) handleLumineServiceRestart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServices.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	return gui.WithWaitingStatus("Restarting service...", func() error {
		if err := gui.Orchestrator.RestartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineServices()
	})
}

func (gui *Gui) handleLumineServiceVersionSwitch(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServices.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	var versions []string
	var switchFunc func(string) error
	
	switch service.Type {
	case lumine.ServiceTypePHP:
		versions = []string{"7.4", "8.0", "8.1", "8.2", "8.3"}
		switchFunc = gui.Orchestrator.SwitchPHPVersion
	case lumine.ServiceTypeNode:
		versions = []string{"16", "18", "20", "22"}
		switchFunc = gui.Orchestrator.SwitchNodeVersion
	default:
		return gui.createErrorPanel("Version switching not supported for this service")
	}
	
	menuItems := make([]*types.MenuItem, len(versions))
	for i, version := range versions {
		v := version // capture for closure
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{fmt.Sprintf("Switch to %s", v)},
			OnPress: func() error {
				return gui.WithWaitingStatus("Switching version...", func() error {
					if err := switchFunc(v); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					return gui.refreshLumineServices()
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Switch %s Version", service.Name),
		Items: menuItems,
	})
}

func (gui *Gui) handleLumineServiceHealth(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineServices.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	health := gui.Orchestrator.ServiceManager.CheckHealth(service.Name)
	
	message := fmt.Sprintf("Service: %s\nHealthy: %v\nUptime: %s", 
		service.Name, health.Healthy, health.Uptime)
	
	if health.Error != "" {
		message += fmt.Sprintf("\nError: %s", health.Error)
	}
	
	return gui.createConfirmationPanel("Health Status", message, func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}, nil)
}

func (gui *Gui) refreshLumineServices() error {
	if gui.Orchestrator == nil || gui.Panels.LumineServices == nil {
		return nil
	}
	
	services := gui.Orchestrator.ServiceManager.ListServices()
	gui.Panels.LumineServices.SetItems(services)
	
	return gui.Panels.LumineServices.RerenderList()
}
