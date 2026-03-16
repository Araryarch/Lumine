package gui

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/gui/types"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/tasks"
	"github.com/Araryarch/Lumine/pkg/utils"
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
						Key:    "network",
						Title:  "Network",
						Render: gui.renderLumineProjectNetwork,
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
			
			networkIndicator := ""
			if project.NetworkName != "" {
				networkIndicator = utils.ColoredString("🔗", color.FgYellow)
			}
			
			return []string{
				utils.ColoredString(project.Name, color.FgCyan),
				string(project.Type),
				project.URL,
				sslIndicator + tunnelIndicator + networkIndicator,
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

func (gui *Gui) renderLumineProjectNetwork(project *lumine.Project) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if project.NetworkName == "" {
			return utils.ColoredString("⚠ No isolated network configured\n\n", color.FgYellow) +
				"This project is not using an isolated Docker network.\n" +
				"Services can communicate with all other containers.\n\n" +
				"Press 'N' to create an isolated network for better security."
		}
		
		output := utils.ColoredString("🔗 Isolated Network Active\n\n", color.FgGreen)
		output += utils.WithPadding("Network Name: ", 20) + utils.ColoredString(project.NetworkName, color.FgCyan) + "\n"
		output += utils.WithPadding("Project: ", 20) + project.Name + "\n\n"
		
		output += utils.ColoredString("Benefits:\n", color.FgYellow)
		output += "  • Services isolated from other projects\n"
		output += "  • Better security and resource management\n"
		output += "  • Services can communicate using container names\n"
		output += "  • No port conflicts between projects\n\n"
		
		// Get network info
		networkInfo, err := gui.Orchestrator.GetProjectNetworkInfo(project.Name)
		if err == nil && networkInfo != "" {
			output += utils.ColoredString("\nNetwork Details:\n", color.FgYellow)
			output += networkInfo
		}
		
		output += "\n\nPress 'C' to connect services to this network"
		output += "\nPress 'D' to disconnect services from this network"
		
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
						if err := gui.Orchestrator.CreateProjectWithNetwork(projectName, projectType, config.DefaultPHPVersion, config.DefaultNodeVersion); err != nil {
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
	
	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Delete project '%s'? This will remove all files and network.", project.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Deleting project...", func() error {
			if err := gui.Orchestrator.DeleteProjectWithNetwork(project.Name); err != nil {
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

// Handler for editing project settings
func (gui *Gui) handleLumineProjectEdit(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}

	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"Edit PHP Version", fmt.Sprintf("Current: %s", project.PHPVersion)},
			OnPress: func() error {
				versions := []string{"7.4", "8.0", "8.1", "8.2", "8.3"}
				versionItems := make([]*types.MenuItem, len(versions))
				
				for i, version := range versions {
					v := version
					versionItems[i] = &types.MenuItem{
						LabelColumns: []string{fmt.Sprintf("PHP %s", v)},
						OnPress: func() error {
							project.PHPVersion = v
							gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("PHP version set to %s", v))
							return gui.refreshLumineProjects()
						},
					}
				}
				
				return gui.Menu(CreateMenuOptions{
					Title: "Select PHP Version",
					Items: versionItems,
				})
			},
		},
		{
			LabelColumns: []string{"Edit Node Version", fmt.Sprintf("Current: %s", project.NodeVersion)},
			OnPress: func() error {
				versions := []string{"16", "18", "20", "21"}
				versionItems := make([]*types.MenuItem, len(versions))
				
				for i, version := range versions {
					v := version
					versionItems[i] = &types.MenuItem{
						LabelColumns: []string{fmt.Sprintf("Node.js %s", v)},
						OnPress: func() error {
							project.NodeVersion = v
							gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Node.js version set to %s", v))
							return gui.refreshLumineProjects()
						},
					}
				}
				
				return gui.Menu(CreateMenuOptions{
					Title: "Select Node.js Version",
					Items: versionItems,
				})
			},
		},
		{
			LabelColumns: []string{"Toggle SSL", fmt.Sprintf("Current: %v", project.SSLEnabled)},
			OnPress: func() error {
				project.SSLEnabled = !project.SSLEnabled
				status := "disabled"
				if project.SSLEnabled {
					status = "enabled"
				}
				gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("SSL %s", status))
				return gui.refreshLumineProjects()
			},
		},
		{
			LabelColumns: []string{"Edit Path", project.Path},
			OnPress: func() error {
				return gui.createPromptPanel("Project Path", func(g *gocui.Gui, v *gocui.View) error {
					path := gui.trimmedContent(v)
					if path == "" {
						return gui.createErrorPanel("Path cannot be empty")
					}

					project.Path = path
					gui.Orchestrator.NotificationMgr.ShowSuccess("Path updated")
					return gui.refreshLumineProjects()
				})
			},
		},
	}

	return gui.Menu(CreateMenuOptions{
		Title: fmt.Sprintf("Edit %s Settings", project.Name),
		Items: menuItems,
	})
}

// Handler for creating project network
func (gui *Gui) handleLumineProjectCreateNetwork(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	if project.NetworkName != "" {
		return gui.createErrorPanel("Project already has a network")
	}
	
	return gui.WithWaitingStatus("Creating network...", func() error {
		network, err := gui.Orchestrator.NetworkManager.CreateProjectNetwork(project.Name)
		if err != nil {
			return gui.createErrorPanel(err.Error())
		}
		
		project.NetworkName = network.Name
		gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Created network: %s", network.Name))
		return gui.refreshLumineProjects()
	})
}

// Handler for connecting service to project network
func (gui *Gui) handleLumineProjectConnectService(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	if project.NetworkName == "" {
		return gui.createErrorPanel("Project has no network. Press 'N' to create one.")
	}
	
	// Get list of running services
	allServices := gui.Orchestrator.ServiceManager.ListAllServices()
	runningServices := []*lumine.Service{}
	for _, svc := range allServices {
		if svc.Running {
			runningServices = append(runningServices, svc)
		}
	}
	
	if len(runningServices) == 0 {
		return gui.createErrorPanel("No running services to connect")
	}
	
	// Create menu to select service
	menuItems := make([]*types.MenuItem, len(runningServices))
	for i, svc := range runningServices {
		service := svc // capture for closure
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{service.DisplayName, service.Name},
			OnPress: func() error {
				return gui.WithWaitingStatus("Connecting service...", func() error {
					if err := gui.Orchestrator.ConnectServiceToProject(service.Name, project.Name); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					return nil
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Select Service to Connect",
		Items: menuItems,
	})
}

// Handler for disconnecting service from project network
func (gui *Gui) handleLumineProjectDisconnectService(g *gocui.Gui, v *gocui.View) error {
	project, err := gui.Panels.LumineProjects.GetSelectedItem()
	if err != nil {
		return nil
	}
	
	if project.NetworkName == "" {
		return gui.createErrorPanel("Project has no network")
	}
	
	// Get list of running services
	allServices := gui.Orchestrator.ServiceManager.ListAllServices()
	runningServices := []*lumine.Service{}
	for _, svc := range allServices {
		if svc.Running {
			runningServices = append(runningServices, svc)
		}
	}
	
	if len(runningServices) == 0 {
		return gui.createErrorPanel("No running services to disconnect")
	}
	
	// Create menu to select service
	menuItems := make([]*types.MenuItem, len(runningServices))
	for i, svc := range runningServices {
		service := svc // capture for closure
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{service.DisplayName, service.Name},
			OnPress: func() error {
				return gui.WithWaitingStatus("Disconnecting service...", func() error {
					if err := gui.Orchestrator.DisconnectServiceFromProject(service.Name, project.Name); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					return nil
				})
			},
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Select Service to Disconnect",
		Items: menuItems,
	})
}
