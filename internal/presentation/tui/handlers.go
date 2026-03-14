package tui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (c *Controller) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Controller) cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	switch v.Name() {
	case "services":
		if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList)-1 {
			c.selectedService++
		}
	case "projects":
		if len(c.projectList) > 0 && c.selectedProject < len(c.projectList)-1 {
			c.selectedProject++
		}
	case "createproject":
		if c.selectedProjectType < len(c.projectTypes)-1 {
			c.selectedProjectType++
		}
	}
	return nil
}

func (c *Controller) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v == nil {
		return nil
	}

	switch v.Name() {
	case "services":
		if c.selectedService > 0 {
			c.selectedService--
		}
	case "projects":
		if c.selectedProject > 0 {
			c.selectedProject--
		}
	case "createproject":
		if c.selectedProjectType > 0 {
			c.selectedProjectType--
		}
	}
	return nil
}

func (c *Controller) nextPanel(g *gocui.Gui, v *gocui.View) error {
	return c.CyclePanels(g, true)
}

func (c *Controller) focusServices(g *gocui.Gui, v *gocui.View) error {
	// If in logs view, close logs first
	if c.currentView == ViewLogs {
		c.currentView = ViewServices
	}
	return c.FocusPanel(g, "services")
}

func (c *Controller) focusMain(g *gocui.Gui, v *gocui.View) error {
	return c.FocusPanel(g, "main")
}

func (c *Controller) selectService(g *gocui.Gui, v *gocui.View) error {
	if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		// Auto focus to main panel to perform actions
		c.FocusPanel(g, "main")
		c.message = fmt.Sprintf("Selected: %s - Use s/x/r to control", c.serviceList[c.selectedService].Name)
		c.messageType = "info"
	}
	return nil
}

func (c *Controller) selectProject(g *gocui.Gui, v *gocui.View) error {
	if len(c.projectList) > 0 && c.selectedProject < len(c.projectList) {
		// Show project details in main panel
		c.message = fmt.Sprintf("Viewing project %s", c.projectList[c.selectedProject].Name)
		c.messageType = "info"
	}
	return nil
}

func (c *Controller) startService(g *gocui.Gui, v *gocui.View) error {
	if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		svc := c.serviceList[c.selectedService]
		c.message = fmt.Sprintf("Starting %s...", svc.Name)
		c.messageType = "info"

		go func() {
			if err := c.serviceApp.Start(svc.Name); err != nil {
				g.Execute(func(g *gocui.Gui) error {
					c.message = fmt.Sprintf("✗ Failed to start %s: %v", svc.Name, err)
					c.messageType = "error"
					return nil
				})
			} else {
				g.Execute(func(g *gocui.Gui) error {
					c.message = fmt.Sprintf("✓ %s started successfully", svc.Name)
					c.messageType = "success"
					// Refresh the service list
					serviceList, _ := c.serviceApp.GetAllStatuses()
					c.serviceList = serviceList
					return nil
				})
			}
		}()
	}
	return nil
}

func (c *Controller) stopService(g *gocui.Gui, v *gocui.View) error {
	if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		svc := c.serviceList[c.selectedService]
		c.message = fmt.Sprintf("Stopping %s...", svc.Name)
		c.messageType = "info"

		go func() {
			if err := c.serviceApp.Stop(svc.Name); err != nil {
				g.Execute(func(g *gocui.Gui) error {
					c.message = fmt.Sprintf("✗ Failed to stop %s: %v", svc.Name, err)
					c.messageType = "error"
					return nil
				})
			} else {
				g.Execute(func(g *gocui.Gui) error {
					c.message = fmt.Sprintf("✓ %s stopped successfully", svc.Name)
					c.messageType = "success"
					// Refresh the service list
					serviceList, _ := c.serviceApp.GetAllStatuses()
					c.serviceList = serviceList
					return nil
				})
			}
		}()
	}
	return nil
}

func (c *Controller) restartService(g *gocui.Gui, v *gocui.View) error {
	if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		svc := c.serviceList[c.selectedService]
		c.message = fmt.Sprintf("Restarting %s...", svc.Name)
		c.messageType = "info"

		go func() {
			if err := c.serviceApp.Restart(svc.Name); err != nil {
				g.Execute(func(g *gocui.Gui) error {
					c.message = fmt.Sprintf("✗ Failed to restart %s: %v", svc.Name, err)
					c.messageType = "error"
					return nil
				})
			} else {
				g.Execute(func(g *gocui.Gui) error {
					c.message = fmt.Sprintf("✓ %s restarted successfully", svc.Name)
					c.messageType = "success"
					// Refresh the service list
					serviceList, _ := c.serviceApp.GetAllStatuses()
					c.serviceList = serviceList
					return nil
				})
			}
		}()
	}
	return nil
}

func (c *Controller) startAllServices(g *gocui.Gui, v *gocui.View) error {
	c.message = "Starting all services..."
	c.messageType = "info"

	go func() {
		if err := c.serviceApp.StartAll(); err != nil {
			g.Execute(func(g *gocui.Gui) error {
				c.message = fmt.Sprintf("✗ Error: %v", err)
				c.messageType = "error"
				return nil
			})
		} else {
			g.Execute(func(g *gocui.Gui) error {
				c.message = "✓ All services started successfully"
				c.messageType = "success"
				// Refresh the service list
				serviceList, _ := c.serviceApp.GetAllStatuses()
				c.serviceList = serviceList
				return nil
			})
		}
	}()
	return nil
}

func (c *Controller) stopAllServices(g *gocui.Gui, v *gocui.View) error {
	c.message = "Stopping all services..."
	c.messageType = "info"

	go func() {
		if err := c.serviceApp.StopAll(); err != nil {
			g.Execute(func(g *gocui.Gui) error {
				c.message = fmt.Sprintf("✗ Error: %v", err)
				c.messageType = "error"
				return nil
			})
		} else {
			g.Execute(func(g *gocui.Gui) error {
				c.message = "✓ All services stopped successfully"
				c.messageType = "success"
				// Refresh the service list
				serviceList, _ := c.serviceApp.GetAllStatuses()
				c.serviceList = serviceList
				return nil
			})
		}
	}()
	return nil
}

func (c *Controller) deleteProject(g *gocui.Gui, v *gocui.View) error {
	if len(c.projectList) > 0 && c.selectedProject < len(c.projectList) {
		proj := c.projectList[c.selectedProject]
		c.message = fmt.Sprintf("Deleting project '%s'...", proj.Name)
		c.messageType = "info"

		if err := c.projectApp.Delete(proj.Name); err != nil {
			c.message = fmt.Sprintf("✗ Error: %v", err)
			c.messageType = "error"
		} else {
			c.message = fmt.Sprintf("✓ Project '%s' deleted", proj.Name)
			c.messageType = "success"
			if c.selectedProject > 0 {
				c.selectedProject--
			}
			// Refresh the project list
			projectList, _ := c.projectApp.List()
			c.projectList = projectList
		}
	}
	return nil
}

func (c *Controller) newProject(g *gocui.Gui, v *gocui.View) error {
	c.currentView = ViewCreateProject
	c.selectedProjectType = 0
	c.message = "Creating new project..."
	c.messageType = "info"
	return nil
}

func (c *Controller) selectProjectType(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewCreateProject {
		projectType := c.projectTypes[c.selectedProjectType]
		c.message = fmt.Sprintf("Selected: %s. Enter project name to create.", projectType)
		c.messageType = "info"
		// In a real implementation, we would show an input dialog here
		// For now, we'll just go back
		c.currentView = ViewProjects
	}
	return nil
}

func (c *Controller) cancelCreateProject(g *gocui.Gui, v *gocui.View) error {
	c.currentView = ViewProjects
	c.message = "Project creation cancelled"
	c.messageType = "info"
	return nil
}

func (c *Controller) toggleLogs(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewLogs {
		c.currentView = ViewServices
		c.message = "Logs closed"
		c.messageType = "info"
	} else {
		if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
			c.currentView = ViewLogs
			c.message = fmt.Sprintf("Viewing logs for %s", c.serviceList[c.selectedService].Name)
			c.messageType = "info"
		} else {
			c.message = "Select a service first to view logs"
			c.messageType = "warning"
		}
	}
	return nil
}

func (c *Controller) closeLogs(g *gocui.Gui, v *gocui.View) error {
	c.currentView = ViewServices
	c.message = "Logs closed"
	c.messageType = "info"
	return nil
}

func (c *Controller) showHelp(g *gocui.Gui, v *gocui.View) error {
	c.message = "[Services] j/k+Enter to select l:logs → [Main] s:start x:stop r:restart S:start-all X:stop-all | [Projects] j/k+Enter d:delete n:new | Tab:switch ←→:panels R:refresh"
	c.messageType = "info"
	return nil
}
