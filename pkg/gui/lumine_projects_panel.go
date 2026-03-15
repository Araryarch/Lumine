package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/jesseduffield/lazydocker/pkg/gui/panels"
	"github.com/jesseduffield/lazydocker/pkg/gui/types"
	"github.com/jesseduffield/lazydocker/pkg/lumine"
	"github.com/jesseduffield/lazydocker/pkg/tasks"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func (gui *Gui) getLumineProjectsPanel() *panels.SideListPanel[*lumine.Project] {
	return &panels.SideListPanel[*lumine.Project]{
		ContextState: &panels.ContextState[*lumine.Project]{
			GetMainTabs: func() []panels.MainTab[*lumine.Project] {
				return []panels.MainTab[*lumine.Project]{
					{
						Key:    "info",
						Title:  "Project Info",
						Render: gui.renderLumineProjectInfo,
					},
					{
						Key:    "dependencies",
						Title:  "Dependencies",
						Render: gui.renderLumineProjectDependencies,
					},
					{
						Key:    "tunnel",
						Title:  "Tunnel Status",
						Render: gui.renderLumineProjectTunnel,
					},
				}
			},
			GetItemContextCacheKey: func(project *lumine.Project) string {
				return "lumine-project-" + project.Name
			},
		},
		ListPanel: panels.ListPanel[*lumine.Project]{
			List: panels.NewFilteredList[*lumine.Project](),
			View: gui.Views.LumineProjects,
		},
		NoItemsMessage: "No Lumine projects",
		Gui:            gui.intoInterface(),
		Sort: func(a *lumine.Project, b *lumine.Project) bool {
			return a.Name < b.Name
		},
		GetTableCells: func(project *lumine.Project) []string {
			sslIndicator := ""
			if project.SSLEnabled {
				sslIndicator = utils.ColoredString("🔒", color.FgGreen)
			}
			
			tunnelIndicator := ""
			if project.TunnelActive {
				tunnelIndicator = utils.ColoredString("🌐", color.FgCyan)
			}
			
			return []string{
				utils.ColoredString(project.Name, color.FgCyan),
				string(project.Type),
				project.URL,
				sslIndicator + tunnelIndicator,
			}
		},
	}
}

func (gui *Gui) renderLumineProjectInfo(project *lumine.Project) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(project.Name, color.FgCyan) + "\n"
		output += utils.WithPadding("Type: ", 15) + string(project.Type) + "\n"
		output += utils.WithPadding("Path: ", 15) + project.Path + "\n"
		output += utils.WithPadding("URL: ", 15) + utils.ColoredString(project.URL, color.FgBlue) + "\n"
		
		if project.SSLEnabled {
			output += utils.WithPadding("HTTPS URL: ", 15) + utils.ColoredString(project.HTTPSURL, color.FgGreen) + "\n"
		}
		
		output += utils.WithPadding("PHP Version: ", 15) + project.PHPVersion + "\n"
		output += utils.WithPadding("Node Version: ", 15) + project.NodeVersion + "\n"
		output += utils.WithPadding("Created: ", 15) + project.CreatedAt.Format("2006-01-02 15:04:05") + "\n"
		
		if project.TunnelActive {
			output += "\n" + utils.ColoredString("🌐 Tunnel Active", color.FgGreen) + "\n"
			output += utils.WithPadding("Public URL: ", 15) + utils.ColoredString(project.TunnelURL, color.FgBlue) + "\n"
		}
		
		return output
	})
}

func (gui *Gui) renderLumineProjectDependencies(project *lumine.Project) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		deps := gui.Orchestrator.ProjectManager.CheckDependencies()
		
		output := utils.ColoredString("Dependency Status:\n\n", color.FgYellow)
		
		for name, status := range deps {
			statusStr := utils.ColoredString("✓ Installed", color.FgGreen)
			if !status.Installed {
				statusStr = utils.ColoredString("✗ Not Found", color.FgRed)
			}
			output += fmt.Sprintf("%s: %s (v%s)\n", name, statusStr, status.Version)
		}
		
		return output
	})
}

func (gui *Gui) renderLumineProjectTunnel(project *lumine.Project) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if !project.TunnelActive {
			return "No active tunnel for this project"
		}
		
		output := ""
		output += utils.WithPadding("Status: ", 15) + utils.ColoredString("Active", color.FgGreen) + "\n"
		output += utils.WithPadding("Public URL: ", 15) + utils.ColoredString(project.TunnelURL, color.FgBlue) + "\n"
		output += utils.WithPadding("Local Port: ", 15) + fmt.Sprintf("%d", project.TunnelPort) + "\n"
		
		return output
	})
}

// Keybinding handlers for Lumine projects
func (gui *Gui) handleLumineProjectCreate(g *gocui.Gui, v *gocui.View) error {
	projectTypes := []lumine.ProjectType{
		lumine.ProjectTypeLaravel,
		lumine.ProjectTypeWordPress,
		lumine.ProjectTypeSymfony,
		lumine.ProjectTypeCodeIgniter,
		lumine.ProjectTypeReact,
		lumine.ProjectTypeVue,
		lumine.ProjectTypeNextJS,
		lumine.ProjectTypeNuxtJS,
		lumine.ProjectTypeExpress,
		lumine.ProjectTypeStatic,
	}
	
	menuItems := make([]*types.MenuItem, len(projectTypes))
	for i, pt := range projectTypes {
		projectType := pt // capture for closure
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{string(projectType)},
			OnPress: func() error {
				return gui.createPromptPanel("Project Name", func(g *gocui.Gui, v *gocui.View) error {
					projectName := gui.trimmedContent(v)
					if projectName == "" {
						return gui.createErrorPanel("Project name cannot be empty")
					}
					
					return gui.WithWaitingStatus("Creating project...", func() error {
						config := gui.Orchestrator.ConfigManager.Get()
						if err := gui.Orchestrator.CreateProject(projectName, projectType, config.DefaultPHPVersion, config.DefaultNodeVersion); err != nil {
							return gui.createErrorPanel(err.Error())
						}
						return gui.refreshLumineProjects()
					})
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Select Project Type",
		Items: menuItems,
	})
}

func (gui *Gui) handleLumineProjectDelete(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Delete project '%s'? This will remove all files.", project.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Deleting project...", func() error {
			if err := gui.Orchestrator.DeleteProject(project.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineProjects()
		})
	}, nil)
}

func (gui *Gui) handleLumineProjectExpose(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	if project.TunnelActive {
		return gui.createConfirmationPanel("Tunnel Active", fmt.Sprintf("Public URL: %s", project.TunnelURL), func(g *gocui.Gui, v *gocui.View) error {
			return nil
		}, nil)
	}
	
	return gui.WithWaitingStatus("Creating tunnel...", func() error {
		if err := gui.Orchestrator.ExposeTunnel(project.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineProjects()
	})
}

func (gui *Gui) handleLumineProjectOpen(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	url := project.URL
	if project.SSLEnabled {
		url = project.HTTPSURL
	}
	
	return gui.OSCommand.OpenLink(url)
}

func (gui *Gui) handleLumineProjectTerminal(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	cmd := gui.OSCommand.ExecutableFromString(fmt.Sprintf("cd %s && $SHELL", project.Path))
	return gui.runSubprocess(cmd)
}

func (gui *Gui) refreshLumineProjects() error {
	if gui.Orchestrator == nil || gui.Panels.LumineProjects == nil {
		return nil
	}
	
	projects := gui.Orchestrator.ProjectManager.ListProjects()
	gui.Panels.LumineProjects.SetItems(projects)
	
	return gui.Panels.LumineProjects.RerenderList()
}
