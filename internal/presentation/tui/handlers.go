package tui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (c *Controller) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (c *Controller) cursorDown(g *gocui.Gui, v *gocui.View) error {
	switch c.currentView {
	case ViewMain:
		maxIdx := len(c.menuItems) - 1
		if c.selectedIdx < maxIdx {
			c.selectedIdx++
			menuView, _ := g.View("menu")
			if menuView != nil {
				c.renderMenu(menuView)
			}
		}
	case ViewProjects:
		if len(c.projectList) > 0 && c.selectedIdx < len(c.projectList)-1 {
			c.selectedIdx++
		}
	case ViewCreateProject:
		if c.selectedIdx < 5 {
			c.selectedIdx++
		}
	case ViewLogs:
		if len(c.serviceList) > 0 && c.selectedIdx < len(c.serviceList)-1 {
			c.selectedIdx++
		}
	}
	return nil
}

func (c *Controller) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if c.selectedIdx > 0 {
		c.selectedIdx--
		if c.currentView == ViewMain {
			menuView, _ := g.View("menu")
			if menuView != nil {
				c.renderMenu(menuView)
			}
		}
	}
	return nil
}

func (c *Controller) selectItem(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewMain {
		c.selectMenuItem()
	}
	return nil
}

func (c *Controller) selectMenuItem() {
	views := []ViewMode{ViewServices, ViewProjects, ViewDatabases, ViewCreateProject, ViewLogs, ViewSettings}
	if c.selectedIdx < len(views) {
		c.currentView = views[c.selectedIdx]
		c.selectedIdx = 0
	}
}

func (c *Controller) goBack(g *gocui.Gui, v *gocui.View) error {
	if c.currentView != ViewMain {
		c.currentView = ViewMain
		c.selectedIdx = 0
		c.message = ""
	}
	return nil
}

func (c *Controller) startServices(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewServices {
		go func() {
			if err := c.serviceApp.StartAll(); err != nil {
				c.showMessage(fmt.Sprintf("Error: %v", err), "error")
			} else {
				c.showMessage("All services started successfully", "success")
			}
			g.Update(func(g *gocui.Gui) error { return nil })
		}()
	}
	return nil
}

func (c *Controller) stopServices(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewServices {
		go func() {
			if err := c.serviceApp.StopAll(); err != nil {
				c.showMessage(fmt.Sprintf("Error: %v", err), "error")
			} else {
				c.showMessage("All services stopped successfully", "success")
			}
			g.Update(func(g *gocui.Gui) error { return nil })
		}()
	}
	return nil
}

func (c *Controller) restartServices(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewServices {
		go func() {
			if err := c.serviceApp.RestartAll(); err != nil {
				c.showMessage(fmt.Sprintf("Error: %v", err), "error")
			} else {
				c.showMessage("All services restarted successfully", "success")
			}
			g.Update(func(g *gocui.Gui) error { return nil })
		}()
	}
	return nil
}

func (c *Controller) deleteProject(g *gocui.Gui, v *gocui.View) error {
	if c.currentView == ViewProjects && len(c.projectList) > 0 && c.selectedIdx < len(c.projectList) {
		project := c.projectList[c.selectedIdx]
		if err := c.projectApp.Delete(project.Name); err != nil {
			c.showMessage(fmt.Sprintf("Error: %v", err), "error")
		} else {
			c.showMessage(fmt.Sprintf("Project '%s' deleted", project.Name), "success")
		}
	}
	return nil
}

func (c *Controller) showHelp(g *gocui.Gui, v *gocui.View) error {
	c.showMessage("j/k: navigate | Enter: select | b: back | q: quit | ?: help", "info")
	return nil
}

func (c *Controller) showMessage(msg string, msgType string) {
	c.message = msg
	c.messageType = msgType
}
