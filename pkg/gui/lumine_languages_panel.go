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

func (gui *Gui) getLumineLanguagesPanel() *panels.SideListPanel[*lumine.Service] {
	return &panels.SideListPanel[*lumine.Service]{
		ContextState: &panels.ContextState[*lumine.Service]{
			GetMainTabs: func() []panels.MainTab[*lumine.Service] {
				return []panels.MainTab[*lumine.Service]{
					{
						Key:    "info",
						Title:  "Runtime Info",
						Render: gui.renderLumineLanguageInfo,
					},
					{
						Key:    "tools",
						Title:  "Available Tools",
						Render: gui.renderLumineLanguageTools,
					},
					{
						Key:    "logs",
						Title:  "Logs",
						Render: gui.renderLumineLanguageLogs,
					},
				}
			},
			GetItemContextCacheKey: func(service *lumine.Service) string {
				return "lumine-language-" + service.Name + "-" + service.Version
			},
		},
		ListPanel: panels.ListPanel[*lumine.Service]{
			List: panels.NewFilteredList[*lumine.Service](),
			View: gui.Views.LumineLanguages,
		},
		NoItemsMessage: "No language runtimes",
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
				service.Version,
				utils.ColoredString(statusText, statusColor),
			}
		},
	}
}

func (gui *Gui) renderLumineLanguageInfo(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(service.DisplayName, color.FgCyan) + "\n"
		output += utils.WithPadding("Status: ", 15) + gui.getColoredStatus(service.Status) + "\n"
		output += utils.WithPadding("Version: ", 15) + service.Version + "\n"
		output += utils.WithPadding("Type: ", 15) + string(service.Type) + "\n"
		output += utils.WithPadding("Image: ", 15) + service.Image + "\n"
		
		if service.Status == "running" {
			output += "\n" + utils.ColoredString("Container is running", color.FgGreen) + "\n"
			output += "\nYou can execute commands inside this container.\n"
		}

		return output
	})
}

func (gui *Gui) renderLumineLanguageTools(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := utils.ColoredString("Available Tools in Container:\n\n", color.FgYellow)
		
		serviceType := string(service.Type)
		
		if serviceType == "php-fpm" {
			output += utils.ColoredString("PHP Tools:\n", color.FgCyan)
			output += "  • php          - PHP CLI\n"
			output += "  • composer     - Dependency manager\n"
			output += "  • pecl         - PHP extensions\n\n"
			output += utils.ColoredString("Usage:\n", color.FgYellow)
			output += "  docker exec lumine-php php -v\n"
			output += "  docker exec lumine-php composer install\n"
			output += "  docker exec lumine-php php artisan migrate\n"
		} else if serviceType == "node" {
			output += utils.ColoredString("Node.js Tools:\n", color.FgCyan)
			output += "  • node         - Node.js runtime\n"
			output += "  • npm          - Package manager\n"
			output += "  • npx          - Package runner\n\n"
			output += utils.ColoredString("Alternative Runtimes:\n", color.FgCyan)
			output += "  • Bun          - Fast all-in-one runtime\n"
			output += "  • Deno         - Secure TypeScript runtime\n\n"
			output += utils.ColoredString("Usage:\n", color.FgYellow)
			output += "  docker exec lumine-node node -v\n"
			output += "  docker exec lumine-node npm install\n"
			output += "  docker exec lumine-node npm run dev\n"
		} else if serviceType == "python" {
			output += utils.ColoredString("Python Tools:\n", color.FgCyan)
			output += "  • python3      - Python runtime\n"
			output += "  • pip          - Package installer\n"
			output += "  • venv         - Virtual environments\n\n"
			output += utils.ColoredString("Package Managers:\n", color.FgCyan)
			output += "  • Poetry       - Dependency management\n"
			output += "  • Pipenv       - Python dev workflow\n\n"
			output += utils.ColoredString("Usage:\n", color.FgYellow)
			output += "  docker exec lumine-python python3 -V\n"
			output += "  docker exec lumine-python pip install -r requirements.txt\n"
			output += "  docker exec lumine-python python3 manage.py runserver\n"
		} else {
			output += "No tools information available for this runtime.\n"
		}
		
		return output
	})
}

func (gui *Gui) renderLumineLanguageLogs(service *lumine.Service) tasks.TaskFunc {
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
func (gui *Gui) handleLumineLanguageStart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status == "running" {
		return gui.createErrorPanel("Runtime is already running")
	}

	return gui.WithWaitingStatus("Starting runtime...", func() error {
		if err := gui.Orchestrator.StartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineLanguages()
	})
}

func (gui *Gui) handleLumineLanguageStop(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	if service.Status != "running" {
		return gui.createErrorPanel("Runtime is not running")
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Stop %s?", service.DisplayName), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping runtime...", func() error {
			if err := gui.Orchestrator.StopService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineLanguages()
		})
	}, nil)
}

func (gui *Gui) handleLumineLanguageRestart(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.WithWaitingStatus("Restarting runtime...", func() error {
		if err := gui.Orchestrator.RestartService(service.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineLanguages()
	})
}

func (gui *Gui) handleLumineLanguageVersionSwitch(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	versions := gui.Orchestrator.GetAvailableVersions(string(service.Type))
	if len(versions) == 0 {
		return gui.createErrorPanel("No versions available for this runtime")
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
					return gui.refreshLumineLanguages()
				})
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Switch %s Version", service.DisplayName),
		Items: menuItems,
	})
}

func (gui *Gui) handleLumineLanguageEdit(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"Edit Image", fmt.Sprintf("Current: %s", service.Image)},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("New Image for %s", service.DisplayName), func(g *gocui.Gui, v *gocui.View) error {
					image := gui.trimmedContent(v)
					if image == "" {
						return gui.createErrorPanel("Image cannot be empty")
					}

					service.Image = image
					gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Image updated to %s", image))
					return gui.refreshLumineLanguages()
				})
			},
		},
		{
			LabelColumns: []string{"Add Volume Mount", "Mount host directory"},
			OnPress: func() error {
				return gui.createPromptPanel("Host Path (e.g., /home/user/projects)", func(g *gocui.Gui, v *gocui.View) error {
					hostPath := gui.trimmedContent(v)
					if hostPath == "" {
						return gui.createErrorPanel("Host path cannot be empty")
					}

					return gui.createPromptPanel("Container Path (e.g., /var/www)", func(g *gocui.Gui, v *gocui.View) error {
						containerPath := gui.trimmedContent(v)
						if containerPath == "" {
							return gui.createErrorPanel("Container path cannot be empty")
						}

						if service.Volumes == nil {
							service.Volumes = make(map[string]string)
						}
						service.Volumes[hostPath] = containerPath
						gui.Orchestrator.NotificationMgr.ShowSuccess("Volume mount added")
						return gui.refreshLumineLanguages()
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

func (gui *Gui) refreshLumineLanguages() error {
	if gui.Orchestrator == nil || gui.Panels.LumineLanguages == nil {
		return nil
	}

	languages := gui.Orchestrator.ServiceManager.ListLanguageServices()
	gui.Panels.LumineLanguages.SetItems(languages)
	return gui.Panels.LumineLanguages.RerenderList()
}
