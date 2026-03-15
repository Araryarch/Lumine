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
						Key:    "versions",
						Title:  "Available Versions",
						Render: gui.renderLumineLanguageVersions,
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
			statusColor := color.FgRed
			if service.Status == "running" {
				statusColor = color.FgGreen
			} else if service.Status == "stopped" {
				statusColor = color.FgYellow
			}

			return []string{
				utils.ColoredString(service.Name, color.FgCyan),
				service.Version,
				utils.ColoredString(service.Status, statusColor),
			}
		},
	}
}

func (gui *Gui) renderLumineLanguageInfo(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(service.Name, color.FgCyan) + "\n"
		output += utils.WithPadding("Status: ", 15) + gui.getColoredStatus(service.Status) + "\n"
		output += utils.WithPadding("Version: ", 15) + service.Version + "\n"
		output += utils.WithPadding("Type: ", 15) + string(service.Type) + "\n"
		output += utils.WithPadding("Port: ", 15) + fmt.Sprintf("%d", service.Port) + "\n"

		return output
	})
}

func (gui *Gui) renderLumineLanguageVersions(service *lumine.Service) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := utils.ColoredString("Available Versions:\n\n", color.FgYellow)
		
		var versions []string
		serviceType := string(service.Type)
		
		if strings.Contains(serviceType, "php") {
			versions = []string{"7.4", "8.0", "8.1", "8.2", "8.3"}
			output += "PHP Versions:\n"
		} else if strings.Contains(serviceType, "node") {
			versions = []string{"16", "18", "20", "21"}
			output += "Node.js Versions:\n"
		} else {
			return "Version switching not available for this runtime"
		}
		
		for _, v := range versions {
			indicator := "  "
			if v == service.Version {
				indicator = utils.ColoredString("● ", color.FgGreen)
			}
			output += fmt.Sprintf("%s%s\n", indicator, v)
		}
		
		output += "\nPress 'v' to switch version"
		
		return output
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

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Stop %s?", service.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping runtime...", func() error {
			if err := gui.Orchestrator.StopService(service.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineLanguages()
		})
	}, nil)
}

func (gui *Gui) handleLumineLanguageVersionSwitch(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	var versions []string
	var switchFunc func(string) error
	serviceType := string(service.Type)

	if strings.Contains(serviceType, "php") {
		versions = []string{"7.4", "8.0", "8.1", "8.2", "8.3"}
		switchFunc = gui.Orchestrator.SwitchPHPVersion
	} else if strings.Contains(serviceType, "node") {
		versions = []string{"16", "18", "20", "21"}
		switchFunc = gui.Orchestrator.SwitchNodeVersion
	} else {
		return gui.createErrorPanel("Version switching not supported for this runtime")
	}

	menuItems := make([]*types.MenuItem, len(versions))
	for i, version := range versions {
		v := version
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{fmt.Sprintf("Switch to %s", v)},
			OnPress: func() error {
				return gui.WithWaitingStatus("Switching version...", func() error {
					if err := switchFunc(v); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					return gui.refreshLumineLanguages()
				})
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Switch %s Version", service.Name),
		Items: menuItems,
	})
}

func (gui *Gui) handleLumineLanguageAdd(g *gocui.Gui, v *gocui.View) error {
	runtimeTypes := []string{"PHP", "Node.js", "Python", "Ruby", "Go"}
	
	menuItems := make([]*types.MenuItem, len(runtimeTypes))
	for i, rt := range runtimeTypes {
		runtimeType := rt
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{runtimeType},
			OnPress: func() error {
				return gui.createPromptPanel(fmt.Sprintf("%s Version (e.g., 8.3)", runtimeType), func(g *gocui.Gui, v *gocui.View) error {
					version := gui.trimmedContent(v)
					if version == "" {
						return gui.createErrorPanel("Version cannot be empty")
					}

					serviceName := fmt.Sprintf("lumine-%s-%s", strings.ToLower(runtimeType), version)
					image := fmt.Sprintf("%s:%s-fpm-alpine", strings.ToLower(runtimeType), version)
					
					if runtimeType == "Node.js" {
						image = fmt.Sprintf("node:%s-alpine", version)
					}

					customService := &lumine.CustomService{
						Name:         serviceName,
						Type:         "language",
						Image:        image,
						Port:         9000,
						InternalPort: 9000,
						Enabled:      true,
						Environment:  make(map[string]string),
						Volumes:      make(map[string]string),
					}

					return gui.WithWaitingStatus("Adding runtime...", func() error {
						if err := gui.Orchestrator.AddCustomService(customService); err != nil {
							return gui.createErrorPanel(err.Error())
						}
						return gui.refreshLumineLanguages()
					})
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Select Runtime Type",
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

// Handler for editing language runtime settings
func (gui *Gui) handleLumineLanguageEdit(g *gocui.Gui, v *gocui.View) error {
	service, err := gui.Panels.LumineLanguages.GetSelectedItem()
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
					return gui.refreshLumineLanguages()
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
					return gui.refreshLumineLanguages()
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
					return gui.refreshLumineLanguages()
				})
			},
		},
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Edit %s Settings", service.Name),
		Items: menuItems,
	})
}
