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

func (gui *Gui) getLumineFilesPanel() *panels.SideListPanel[*lumine.Service] {
	return &panels.SideListPanel[*lumine.Service]{
		ContextState: &panels.ContextState[*lumine.Service]{
			GetMainTabs: func() []panels.MainTab[*lumine.Service] {
				return []panels.MainTab[*lumine.Service]{
					{
						Key:    "info",
						Title:  "Service Info",
						Render: gui.renderLumineFileServiceInfo,
					},
					{
						Key:    "config",
						Title:  "Configuration",
						Render: gui.renderLumineFileServiceConfig,
					},
					{
						Key:    "logs",
						Title:  "Logs",
						Render: gui.renderLumineFileServiceLogs,
					},
				}
			},
			GetItemContextCacheKey: func(service *lumine.Service) string {
				return "lumine-file-" + service.Name + "-" + service.Status
			},
		},
		ListPanel: panels.ListPanel[*lumine.Service]{
			List: panels.NewFilteredList[*lumine.Service](),
			View: gui.Views.LumineFiles,
		},
		NoItemsMessage: "No file services",
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

func (gui *Gui) renderLumineFileServiceInfo(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(service.DisplayName, color.FgCyan) + "\n"
		output += utils.WithPadding("Status: ", 15) + gui.getColoredStatus(service.Status) + "\n"
		output += utils.WithPadding("Type: ", 15) + string(service.Type) + "\n"
		output += utils.WithPadding("Image: ", 15) + service.Image + "\n"
		output += utils.WithPadding("Port: ", 15) + fmt.Sprintf("%d", service.Port) + "\n"

		if service.Status == "running" {
			output += "\n" + utils.ColoredString("Access Information:\n", color.FgYellow)
			
			serviceType := string(service.Type)
			if serviceType == "ftp" || serviceType == "sftp" {
				output += fmt.Sprintf("  Host: localhost\n")
				output += fmt.Sprintf("  Port: %d\n", service.Port)
				output += fmt.Sprintf("  User: Check environment variables\n")
			} else if serviceType == "webdav" {
				output += fmt.Sprintf("  URL: http://localhost:%d\n", service.Port)
			}
		}

		return output
	})
}

func (gui *Gui) renderLumineFileServiceConfig(service *lumine.Service) tasks.TaskFunc {
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

func (gui *Gui) renderLumineFileServiceLogs(service *lumine.Service) tasks.TaskFunc {
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
func (gui *Gui) handleLumineFileServiceStart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineFiles.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status == "running" {
		return gui.createErrorPanel("Service is already running")
	}

	return gui.WithWaitingStatus("Starting file service...", func() error {
		if err := gui.Orchestrator.StartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineFiles()
	})
}

func (gui *Gui) handleLumineFileServiceStop(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineFiles.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status != "running" {
		return gui.createErrorPanel("Service is not running")
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Stop %s?", service.DisplayName), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping file service...", func() error {
			if err := gui.Orchestrator.StopService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineFiles()
		})
	}, nil)
}

func (gui *Gui) handleLumineFileServiceRestart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineFiles.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.WithWaitingStatus("Restarting file service...", func() error {
		if err := gui.Orchestrator.RestartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineFiles()
	})
}

func (gui *Gui) handleLumineFileServiceAdd(g *gocui.Gui, v *gocui.View) error {
	fileServices := []string{"FTP Server", "SFTP Server", "WebDAV", "Samba/SMB"}
	
	menuItems := make([]*types.MenuItem, len(fileServices))
	for i, fs := range fileServices {
		fileService := fs
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{fileService},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("%s Name", fileService), func(g *gocui.Gui, v *gocui.View) error {
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

						var image string
						var serviceType string
						env := make(map[string]string)
						
						switch fileService {
						case "FTP Server":
							image = "fauria/vsftpd"
							serviceType = "ftp"
							env["FTP_USER"] = "lumine"
							env["FTP_PASS"] = "lumine"
							env["PASV_ADDRESS"] = "localhost"
						case "SFTP Server":
							image = "atmoz/sftp"
							serviceType = "sftp"
							env["SFTP_USERS"] = "lumine:lumine:1001"
						case "WebDAV":
							image = "bytemark/webdav"
							serviceType = "webdav"
							env["AUTH_TYPE"] = "Basic"
							env["USERNAME"] = "lumine"
							env["PASSWORD"] = "lumine"
						case "Samba/SMB":
							image = "dperson/samba"
							serviceType = "samba"
							env["USER"] = "lumine;lumine"
						}

						customService := &lumine.CustomService{
							Name:         serviceName,
							Type:         serviceType,
							Image:        image,
							Port:         port,
							InternalPort: port,
							Enabled:      true,
							Environment:  env,
							Volumes:      make(map[string]string),
						}

						return gui.WithWaitingStatus("Adding file service...", func() error {
							if err := gui.Orchestrator.AddCustomService(customService); err != nil {
								return gui.createErrorPanel(err.Error())
							}
							return gui.refreshLumineFiles()
						})
					})
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Select File Service Type",
		Items: menuItems,
	})
}

func (gui *Gui) handleLumineFileServiceRemove(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineFiles.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Remove file service '%s'?", service.DisplayName), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Removing file service...", func() error {
			if err := gui.Orchestrator.RemoveCustomService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineFiles()
		})
	}, nil)
}

func (gui *Gui) handleLumineFileServiceEdit(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineFiles.GetSelectedItem()
	if err != nil {
		return nil
	}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"Edit Port", fmt.Sprintf("Current: %d", service.Port)},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("New Port for %s", service.DisplayName), func(g *gocui.Gui, v *gocui.View) error {
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
					return gui.refreshLumineFiles()
				})
			},
		},
		{
			LabelColumns: []string{"Edit Username", "Change username"},
			OnPress: func() error {
				return gui.createPromptPanel("Username", func(g *gocui.Gui, v *gocui.View) error {
					username := gui.trimmedContent(v)
					if username == "" {
						return gui.createErrorPanel("Username cannot be empty")
					}

					if service.Environment == nil {
						service.Environment = make(map[string]string)
					}
					
					serviceType := string(service.Type)
					if serviceType == "ftp" {
						service.Environment["FTP_USER"] = username
					} else if serviceType == "sftp" {
						service.Environment["SFTP_USERS"] = username + ":" + username + ":1001"
					} else if serviceType == "webdav" {
						service.Environment["USERNAME"] = username
					}

					gui.Orchestrator.NotificationMgr.ShowSuccess("Username updated")
					return gui.refreshLumineFiles()
				})
			},
		},
		{
			LabelColumns: []string{"Edit Password", "Change password"},
			OnPress: func() error {
				return gui.createPromptPanel("Password", func(g *gocui.Gui, v *gocui.View) error {
					password := gui.trimmedContent(v)
					if password == "" {
						return gui.createErrorPanel("Password cannot be empty")
					}

					if service.Environment == nil {
						service.Environment = make(map[string]string)
					}
					
					serviceType := string(service.Type)
					if serviceType == "ftp" {
						service.Environment["FTP_PASS"] = password
					} else if serviceType == "webdav" {
						service.Environment["PASSWORD"] = password
					}

					gui.Orchestrator.NotificationMgr.ShowSuccess("Password updated")
					return gui.refreshLumineFiles()
				})
			},
		},
		{
			LabelColumns: []string{"Add Volume Mount", "Mount directory"},
			OnPress: func() error {
				return gui.createPromptPanel("Host Path (e.g., /home/user/files)", func(g *gocui.Gui, v *gocui.View) error {
					hostPath := gui.trimmedContent(v)
					if hostPath == "" {
						return gui.createErrorPanel("Host path cannot be empty")
					}

					return gui.createPromptPanel("Container Path (e.g., /home/vsftpd)", func(g *gocui.Gui, v *gocui.View) error {
						containerPath := gui.trimmedContent(v)
						if containerPath == "" {
							return gui.createErrorPanel("Container path cannot be empty")
						}

						if service.Volumes == nil {
							service.Volumes = make(map[string]string)
						}
						service.Volumes[hostPath] = containerPath
						gui.Orchestrator.NotificationMgr.ShowSuccess("Volume mount added")
						return gui.refreshLumineFiles()
					})
				})
			},
		},
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Edit %s Settings", service.DisplayName),
		Items: menuItems,
	})
}

func (gui *Gui) refreshLumineFiles() error {
	if gui.Orchestrator == nil || gui.Panels.LumineFiles == nil {
		return nil
	}

	files := gui.Orchestrator.ServiceManager.ListFileServices()
	gui.Panels.LumineFiles.SetItems(files)
	return gui.Panels.LumineFiles.RerenderList()
}
