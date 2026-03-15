package gui

import (
	"context"
	"fmt"
	"time"

	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/gui/types"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/tasks"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getLumineDatabasesPanel() *panels.SideListPanel[*lumine.Service] {
	return &panels.SideListPanel[*lumine.Service]{
		ContextState: &panels.ContextState[*lumine.Service]{
			GetMainTabs: func() []panels.MainTab[*lumine.Service] {
				return []panels.MainTab[*lumine.Service]{
					{
						Key:    "info",
						Title:  "Database Info",
						Render: gui.renderLumineDatabaseServiceInfo,
					},
					{
						Key:    "config",
						Title:  "Configuration",
						Render: gui.renderLumineDatabaseServiceConfig,
					},
					{
						Key:    "health",
						Title:  "Health Status",
						Render: gui.renderLumineDatabaseServiceHealth,
					},
					{
						Key:    "logs",
						Title:  "Logs",
						Render: gui.renderLumineDatabaseServiceLogs,
					},
				}
			},
			GetItemContextCacheKey: func(service *lumine.Service) string {
				return "lumine-database-service-" + service.Name + "-" + service.Status
			},
		},
		ListPanel: panels.ListPanel[*lumine.Service]{
			List: panels.NewFilteredList[*lumine.Service](),
			View: gui.Views.LumineDatabases,
		},
		NoItemsMessage: "No database services",
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
				string(service.Type),
				utils.ColoredString(statusText, statusColor),
			}
		},
	}
}

func (gui *Gui) renderLumineDatabaseServiceInfo(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(service.Name, color.FgCyan) + "\n"
		output += utils.WithPadding("Type: ", 15) + string(service.Type) + "\n"
		output += utils.WithPadding("Status: ", 15) + gui.getColoredStatus(service.Status) + "\n"
		output += utils.WithPadding("Port: ", 15) + fmt.Sprintf("%d", service.Port) + "\n"
		output += utils.WithPadding("Image: ", 15) + service.Image + "\n"
		output += utils.WithPadding("Version: ", 15) + service.Version + "\n"

		if service.Status == "running" {
			output += "\n" + utils.ColoredString("● Service is running", color.FgGreen) + "\n"
		}

		return output
	})
}

func (gui *Gui) renderLumineDatabaseServiceConfig(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := utils.ColoredString("Environment Variables:\n\n", color.FgYellow)
		
		if len(service.Environment) == 0 {
			return "No environment variables configured"
		}

		for key, value := range service.Environment {
			output += fmt.Sprintf("%s = %s\n", key, value)
		}

		return output
	})
}

func (gui *Gui) renderLumineDatabaseServiceHealth(service *lumine.Service) tasks.TaskFunc {
	return gui.NewTickerTask(TickerTaskOpts{
		Func: func(ctx context.Context, notifyStopped chan struct{}) {
			health := gui.Orchestrator.ServiceManager.CheckHealth(service.Name)

			output := ""
			output += utils.WithPadding("Service: ", 15) + service.Name + "\n"
			output += utils.WithPadding("Healthy: ", 15) + fmt.Sprintf("%v", health.Healthy) + "\n"
			output += utils.WithPadding("Uptime: ", 15) + health.Uptime.String() + "\n"
			output += utils.WithPadding("Last Check: ", 15) + health.LastCheck.Format("2006-01-02 15:04:05") + "\n"

			if health.Error != "" {
				output += "\n" + utils.ColoredString("Error: "+health.Error, color.FgRed) + "\n"
			}

			gui.reRenderStringMain(output)
		},
		Duration:   time.Second * 2,
		Before:     func(ctx context.Context) { gui.clearMainView() },
		Wrap:       true,
		Autoscroll: false,
	})
}

// Keybinding handlers for database services
func (gui *Gui) handleLumineDatabaseServiceAdd(g *gocui.Gui, v *gocui.View) error {
	dbTypes := []string{"MySQL", "PostgreSQL", "MongoDB", "Redis", "Elasticsearch"}
	
	menuItems := make([]*types.MenuItem, len(dbTypes))
	for i, dbType := range dbTypes {
		dt := dbType
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{dt},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("%s Service Name", dt), func(g *gocui.Gui, v *gocui.View) error {
					serviceName := gui.trimmedContent(v)
					if serviceName == "" {
						return gui.createErrorPanel("Service name cannot be empty")
					}

					return gui.createPromptPanel("Port", func(g *gocui.Gui, v *gocui.View) error {
						portStr := gui.trimmedContent(v)
						if portStr == "" {
							return gui.createErrorPanel("Port cannot be empty")
						}

						var port int
						if _, err := fmt.Sscanf(portStr, "%d", &port); err != nil {
							return gui.createErrorPanel("Invalid port number")
						}

						imageMap := map[string]string{
							"MySQL":         "mysql:8.0",
							"PostgreSQL":    "postgres:alpine",
							"MongoDB":       "mongo:latest",
							"Redis":         "redis:alpine",
							"Elasticsearch": "elasticsearch:8.11.0",
						}

						customService := &lumine.CustomService{
							Name:         serviceName,
							Type:         "database",
							Image:        imageMap[dt],
							Port:         port,
							InternalPort: port,
							Enabled:      true,
							Environment:  make(map[string]string),
							Volumes:      make(map[string]string),
						}

						if dt == "MySQL" {
							customService.Environment["MYSQL_ROOT_PASSWORD"] = "root"
						} else if dt == "PostgreSQL" {
							customService.Environment["POSTGRES_PASSWORD"] = "root"
						}

						return gui.WithWaitingStatus("Adding database service...", func() error {
							if err := gui.Orchestrator.AddCustomService(customService); err != nil {
								return gui.createErrorPanel(err.Error())
							}
							return gui.refreshLumineDatabases()
						})
					})
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Select Database Type",
		Items: menuItems,
	})
}

func (gui *Gui) handleLumineDatabaseServiceStart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status == "running" {
		return gui.createErrorPanel("Database service is already running")
	}

	return gui.WithWaitingStatus("Starting database service...", func() error {
		if err := gui.Orchestrator.StartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineDatabases()
	})
}

func (gui *Gui) handleLumineDatabaseServiceStop(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status != "running" {
		return gui.createErrorPanel("Database service is not running")
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Stop %s?", service.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping database service...", func() error {
			if err := gui.Orchestrator.StopService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineDatabases()
		})
	}, nil)
}

func (gui *Gui) handleLumineDatabaseServiceRestart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.WithWaitingStatus("Restarting database service...", func() error {
		if err := gui.Orchestrator.RestartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineDatabases()
	})
}

func (gui *Gui) handleLumineDatabaseServiceRemove(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Remove database service '%s'?", service.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Removing database service...", func() error {
			if err := gui.Orchestrator.RemoveCustomService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineDatabases()
		})
	}, nil)
}

func (gui *Gui) handleLumineDatabaseServiceEdit(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineDatabases.GetSelectedItem()
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

					service.Port = port
					gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Port updated to %d", port))
					return gui.refreshLumineDatabases()
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
					return gui.refreshLumineDatabases()
				})
			},
		},
		{
			LabelColumns: []string{"Edit Root Password", "***"},
			OnPress: func() error {
				return gui.createPromptPanel("Root Password", func(g *gocui.Gui, v *gocui.View) error {
					password := gui.trimmedContent(v)
					if password == "" {
						return gui.createErrorPanel("Password cannot be empty")
					}

					if service.Environment == nil {
						service.Environment = make(map[string]string)
					}

					// Set password based on database type
					serviceType := string(service.Type)
					if serviceType == "mysql" {
						service.Environment["MYSQL_ROOT_PASSWORD"] = password
					} else if serviceType == "postgresql" {
						service.Environment["POSTGRES_PASSWORD"] = password
					}

					gui.Orchestrator.NotificationMgr.ShowSuccess("Password updated")
					return gui.refreshLumineDatabases()
				})
			},
		},
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Edit %s Settings", service.Name),
		Items: menuItems,
	})
}

func (gui *Gui) refreshLumineDatabases() error {
	if gui.Orchestrator == nil || gui.Panels.LumineDatabases == nil {
		return nil
	}

	databases := gui.Orchestrator.ServiceManager.ListDatabaseServices()
	gui.Panels.LumineDatabases.SetItems(databases)
	return gui.Panels.LumineDatabases.RerenderList()
}

// Handler for switching database version
func (gui *Gui) handleLumineDatabaseVersionSwitch(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	versions := gui.Orchestrator.GetAvailableVersions(string(service.Type))
	if len(versions) == 0 {
		return gui.createErrorPanel("No versions available for this database")
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
					return gui.refreshLumineDatabases()
				})
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Switch %s Version", service.Name),
		Items: menuItems,
	})
}

func (gui *Gui) renderLumineDatabaseServiceLogs(service *lumine.Service) tasks.TaskFunc {
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
