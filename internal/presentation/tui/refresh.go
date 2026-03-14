package tui

import (
	"time"

	"github.com/jesseduffield/gocui"
)

// StartAutoRefresh starts automatic refresh of the UI
func (c *Controller) StartAutoRefresh(g *gocui.Gui) {
	ticker := time.NewTicker(3 * time.Second)
	go func() {
		for range ticker.C {
			g.Execute(func(g *gocui.Gui) error {
				// Refresh service list
				serviceList, err := c.serviceApp.GetAllStatuses()
				if err == nil {
					c.serviceList = serviceList
				}

				// Refresh project list
				projectList, err := c.projectApp.List()
				if err == nil {
					c.projectList = projectList
				}

				return nil
			})
		}
	}()
}

// RefreshView manually refreshes the current view
func (c *Controller) RefreshView(g *gocui.Gui, v *gocui.View) error {
	// Refresh data
	serviceList, err := c.serviceApp.GetAllStatuses()
	if err != nil {
		c.message = "Error refreshing services"
		c.messageType = "error"
		return nil
	}
	c.serviceList = serviceList

	projectList, err := c.projectApp.List()
	if err != nil {
		c.message = "Error refreshing projects"
		c.messageType = "error"
		return nil
	}
	c.projectList = projectList

	c.message = "Refreshed"
	c.messageType = "success"

	// Clear message after 2 seconds
	go func() {
		time.Sleep(2 * time.Second)
		g.Execute(func(g *gocui.Gui) error {
			c.message = ""
			return nil
		})
	}()

	return nil
}
