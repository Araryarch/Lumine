package gui

import (
	"fmt"
	"strings"

	"github.com/Araryarch/Lumine/pkg/gui/types"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

// ShowFirstSetupWizard shows the first setup wizard
func (gui *Gui) ShowFirstSetupWizard() error {
	wizard := lumine.NewFirstSetupWizard(gui.Orchestrator)
	
	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"1. Check Docker Installation", "Verify Docker is ready"},
			OnPress: func() error {
				return gui.handleCheckDocker()
			},
		},
		{
			LabelColumns: []string{"2. Check Development Tools", "PHP, Node.js, Composer, npm, etc."},
			OnPress: func() error {
				return gui.handleCheckDevTools()
			},
		},
		{
			LabelColumns: []string{"3. Configure Default Settings", "Set versions and preferences"},
			OnPress: func() error {
				return gui.handleConfigureDefaults()
			},
		},
		{
			LabelColumns: []string{"4. Choose Stack Template", "LAMP, LEMP, MEAN, JAMstack"},
			OnPress: func() error {
				return gui.handleChooseStack(wizard)
			},
		},
		{
			LabelColumns: []string{"5. View Project Templates", "See available project types"},
			OnPress: func() error {
				return gui.handleViewProjectTemplates()
			},
		},
		{
			LabelColumns: []string{"Quick Start (Recommended)", "Auto-detect and setup"},
			OnPress: func() error {
				return gui.handleQuickStart(wizard)
			},
		},
		{
			LabelColumns: []string{"Skip Setup", "Continue without setup"},
			OnPress: func() error {
				gui.Orchestrator.ConfigManager.MarkSetupComplete()
				return nil
			},
		},
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Welcome to Lumine - First Setup",
		Items: menuItems,
	})
}

func (gui *Gui) handleCheckDocker() error {
	return gui.WithWaitingStatus("Checking Docker...", func() error {
		output, err := gui.OSCommand.RunCommandWithOutput("docker --version")
		if err != nil {
			return gui.createErrorPanel("Docker is not installed or not running.\n\nPlease install Docker:\n  Ubuntu: sudo apt install docker.io\n  Fedora: sudo dnf install docker\n  Or visit: https://docs.docker.com/get-docker/")
		}
		
		// Check if Docker daemon is running
		_, err = gui.OSCommand.RunCommandWithOutput("docker ps")
		if err != nil {
			return gui.createErrorPanel("Docker daemon is not running.\n\nStart Docker:\n  sudo systemctl start docker\n  sudo systemctl enable docker")
		}
		
		return gui.createInfoPanel("Docker Check", fmt.Sprintf("✓ Docker is installed and running\n\n%s", output))
	})
}

func (gui *Gui) handleCheckDevTools() error {
	return gui.WithWaitingStatus("Checking development tools...", func() error {
		gui.Orchestrator.ToolManager.CheckAllInstallations()
		
		tools := gui.Orchestrator.ToolManager.GetAllTools()
		
		output := utils.ColoredString("Development Tools Status:\n\n", color.FgYellow)
		
		installedCount := 0
		for _, tool := range tools {
			status := "✗ Not Installed"
			statusColor := color.FgRed
			if tool.Installed {
				status = "✓ Installed"
				statusColor = color.FgGreen
				installedCount++
			}
			
			output += fmt.Sprintf("%s %-15s %s\n", 
				utils.ColoredString(status, statusColor),
				tool.DisplayName,
				tool.Version)
		}
		
		output += fmt.Sprintf("\n%d/%d tools installed", installedCount, len(tools))
		
		return gui.createInfoPanel("Development Tools", output)
	})
}

func (gui *Gui) handleConfigureDefaults() error {
	config := gui.Orchestrator.ConfigManager.Get()
	
	menuItems := []*types.MenuItem{
		{
			LabelColumns: []string{"Default PHP Version", config.DefaultPHPVersion},
			OnPress: func() error {
				return gui.handleSettingPHPVersion()
			},
		},
		{
			LabelColumns: []string{"Default Node.js Version", config.DefaultNodeVersion},
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
			LabelColumns: []string{"Projects Directory", config.ProjectsDirectory},
			OnPress: func() error {
				return gui.handleSettingProjectsDir()
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
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Configure Default Settings",
		Items: menuItems,
	})
}

func (gui *Gui) handleChooseStack(wizard *lumine.FirstSetupWizard) error {
	stacks := gui.Orchestrator.TemplateManager.ListStacks()
	
	menuItems := make([]*types.MenuItem, 0, len(stacks))
	
	for stackName, stack := range stacks {
		sn := stackName
		s := stack
		
		recommended := ""
		if sn == wizard.GetRecommendedStack() {
			recommended = " " + utils.ColoredString("(Recommended)", color.FgGreen)
		}
		
		menuItems = append(menuItems, &types.MenuItem{
			LabelColumns: []string{s.Name + recommended, s.Description},
			OnPress: func() error {
				return gui.handleInstallStack(wizard, sn, s)
			},
		})
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Choose Stack Template",
		Items: menuItems,
	})
}

func (gui *Gui) handleInstallStack(wizard *lumine.FirstSetupWizard, stackName string, stack *lumine.StackTemplate) error {
	message := fmt.Sprintf("Install %s?\n\nThis will start:\n", stack.Name)
	for _, service := range stack.Services {
		message += fmt.Sprintf("  • %s (%s)\n", service.Name, service.Image)
	}
	
	return gui.createConfirmationPanel("Confirm Stack Installation", message, func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Installing stack...", func() error {
			if err := wizard.InstallStack(stackName); err != nil {
				return gui.createErrorPanel(fmt.Sprintf("Failed to install stack: %v", err))
			}
			
			gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("%s installed successfully!", stack.Name))
			
			// Refresh all panels
			gui.refreshLumineServers()
			gui.refreshLumineLanguages()
			gui.refreshLumineDatabases()
			
			return nil
		})
	}, nil)
}

func (gui *Gui) handleQuickStart(wizard *lumine.FirstSetupWizard) error {
	message := "Quick Start will:\n\n"
	message += "1. Check Docker installation\n"
	message += "2. Detect installed development tools\n"
	message += "3. Install recommended stack based on your tools\n"
	message += "4. Configure default settings\n\n"
	message += "Continue?"
	
	return gui.createConfirmationPanel("Quick Start", message, func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Running quick start...", func() error {
			// Step 1: Check Docker
			_, err := gui.OSCommand.RunCommandWithOutput("docker ps")
			if err != nil {
				return gui.createErrorPanel("Docker is not running. Please start Docker first.")
			}
			
			// Step 2: Check dev tools
			gui.Orchestrator.ToolManager.CheckAllInstallations()
			
			// Step 3: Get recommended stack
			recommendedStack := wizard.GetRecommendedStack()
			
			// Step 4: Install stack
			if err := wizard.InstallStack(recommendedStack); err != nil {
				return gui.createErrorPanel(fmt.Sprintf("Failed to install stack: %v", err))
			}
			
			stack, _ := gui.Orchestrator.TemplateManager.GetStack(recommendedStack)
			gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Quick start complete! %s installed.", stack.Name))
			
			// Refresh all panels
			gui.refreshLumineServers()
			gui.refreshLumineLanguages()
			gui.refreshLumineDatabases()
			
			return nil
		})
	}, nil)
}

func (gui *Gui) createInfoPanel(title, message string) error {
	return gui.createConfirmationPanel(title, message, func(g *gocui.Gui, v *gocui.View) error {
		return nil
	}, nil)
}

// Handler for showing first setup wizard
func (gui *Gui) handleShowFirstSetup(g *gocui.Gui, v *gocui.View) error {
	return gui.ShowFirstSetupWizard()
}


func (gui *Gui) handleViewProjectTemplates() error {
	templates := lumine.GetProjectTemplateInfo()
	
	menuItems := make([]*types.MenuItem, 0, len(templates))
	
	// Order: PHP frameworks, JS frameworks, Python, Static
	order := []string{"php", "laravel", "symfony", "wordpress", "react", "vue", "nextjs", "express", "python", "static"}
	
	for _, key := range order {
		if template, ok := templates[key]; ok {
			t := template
			menuItems = append(menuItems, &types.MenuItem{
				LabelColumns: []string{t.Name, t.Description},
				OnPress: func() error {
					return gui.showProjectTemplateDetails(t)
				},
			})
		}
	}
	
	return gui.Menu(CreateMenuOptions{
		Title: "Available Project Templates",
		Items: menuItems,
	})
}

func (gui *Gui) showProjectTemplateDetails(template lumine.ProjectTemplateInfo) error {
	output := utils.ColoredString(template.Name+"\n", color.FgCyan)
	output += utils.ColoredString(strings.Repeat("=", len(template.Name))+"\n\n", color.FgCyan)
	
	output += template.Description + "\n\n"
	
	output += utils.ColoredString("Requirements:\n", color.FgYellow)
	for _, req := range template.Requirements {
		output += fmt.Sprintf("  • %s\n", req)
	}
	
	output += "\n" + utils.ColoredString("Recommended Stack: ", color.FgYellow) + template.Stack + "\n\n"
	
	output += utils.ColoredString("Features:\n", color.FgYellow)
	for _, feature := range template.Features {
		output += fmt.Sprintf("  • %s\n", feature)
	}
	
	return gui.createInfoPanel(template.Name+" Template", output)
}
