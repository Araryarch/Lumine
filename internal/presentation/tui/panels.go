package tui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

// Panel represents a UI panel
type Panel struct {
	Name      string
	Title     string
	X0, Y0    int
	X1, Y1    int
	Highlight bool
}

// GetPanels returns all panels for the layout
func (c *Controller) GetPanels(maxX, maxY int) []Panel {
	leftWidth := maxX / 3
	midHeight := maxY / 2

	panels := []Panel{
		{
			Name:      "services",
			Title:     " Services ",
			X0:        0,
			Y0:        3,
			X1:        leftWidth,
			Y1:        midHeight - 1,
			Highlight: true,
		},
		{
			Name:      "projects",
			Title:     " Projects ",
			X0:        0,
			Y0:        midHeight,
			X1:        leftWidth,
			Y1:        maxY - 4,
			Highlight: true,
		},
	}

	// Add main or create project panel
	if c.currentView == ViewCreateProject {
		panels = append(panels, Panel{
			Name:      "createproject",
			Title:     " Create Project ",
			X0:        leftWidth + 1,
			Y0:        3,
			X1:        maxX - 1,
			Y1:        maxY - 4,
			Highlight: false,
		})
	} else {
		panels = append(panels, Panel{
			Name:      "main",
			Title:     c.getMainTitle(),
			X0:        leftWidth + 1,
			Y0:        3,
			X1:        maxX - 1,
			Y1:        maxY - 4,
			Highlight: false,
		})
	}

	// Add logs panel if in logs view mode
	if c.currentView == ViewLogs {
		panels = append(panels, Panel{
			Name:      "logs",
			Title:     " Logs: " + c.getLogsTitle(),
			X0:        leftWidth + 1,
			Y0:        maxY/2 - 2,
			X1:        maxX - 1,
			Y1:        maxY - 4,
			Highlight: false,
		})
	}

	return panels
}

// CreatePanel creates or updates a panel
func (c *Controller) CreatePanel(g *gocui.Gui, panel Panel) error {
	v, err := g.SetView(panel.Name, panel.X0, panel.Y0, panel.X1, panel.Y1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = panel.Title
		v.Highlight = panel.Highlight
		if panel.Highlight {
			v.SelBgColor = gocui.ColorGreen
			v.SelFgColor = gocui.ColorBlack
		}

		// Configure panel-specific settings
		switch panel.Name {
		case "main":
			v.Wrap = true
			v.Autoscroll = false
		case "services", "projects":
			v.Wrap = false
			v.Autoscroll = false
		}
	} else {
		v.Title = panel.Title
	}
	return nil
}

// RenderPanel renders content for a specific panel
func (c *Controller) RenderPanel(g *gocui.Gui, panelName string) error {
	v, err := g.View(panelName)
	if err != nil {
		return err
	}

	switch panelName {
	case "services":
		c.renderServicesList(v)
	case "projects":
		c.renderProjectsList(v)
	case "main":
		c.renderMainPanel(v)
	case "logs":
		c.renderLogsPanel(v)
	case "createproject":
		c.renderCreateProjectPanel(v)
	}

	return nil
}

// FocusPanel sets focus to a specific panel
func (c *Controller) FocusPanel(g *gocui.Gui, panelName string) error {
	switch panelName {
	case "services":
		c.currentView = ViewServices
	case "projects":
		c.currentView = ViewProjects
	case "main":
		c.currentView = ViewMain
	case "logs":
		c.currentView = ViewLogs
	case "createproject":
		c.currentView = ViewCreateProject
	}
	return g.SetCurrentView(panelName)
}

// GetCurrentPanel returns the current active panel name
func (c *Controller) GetCurrentPanel() string {
	switch c.currentView {
	case ViewServices:
		return "services"
	case ViewProjects:
		return "projects"
	case ViewMain:
		return "main"
	case ViewLogs:
		return "logs"
	case ViewCreateProject:
		return "createproject"
	default:
		return "services"
	}
}

// CyclePanels cycles through panels
func (c *Controller) CyclePanels(g *gocui.Gui, forward bool) error {
	// If in logs view, exit logs first
	if c.currentView == ViewLogs {
		c.currentView = ViewServices
		return c.FocusPanel(g, "services")
	}

	// If in createproject view, go back to projects
	if c.currentView == ViewCreateProject {
		c.currentView = ViewProjects
		return c.FocusPanel(g, "projects")
	}

	panels := []string{"services", "projects"}
	current := c.GetCurrentPanel()

	var nextIdx int
	for i, p := range panels {
		if p == current {
			if forward {
				nextIdx = (i + 1) % len(panels)
			} else {
				nextIdx = (i - 1 + len(panels)) % len(panels)
			}
			break
		}
	}

	return c.FocusPanel(g, panels[nextIdx])
}

// ShowPopup shows a popup message
func (c *Controller) ShowPopup(g *gocui.Gui, title, message string) error {
	maxX, maxY := g.Size()
	width := 60
	height := 10
	x0 := (maxX - width) / 2
	y0 := (maxY - height) / 2

	v, err := g.SetView("popup", x0, y0, x0+width, y0+height)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = title
		v.Wrap = true
	}

	v.Clear()
	fmt.Fprintln(v, message)

	return g.SetCurrentView("popup")
}

// ClosePopup closes the popup
func (c *Controller) ClosePopup(g *gocui.Gui) error {
	g.DeleteView("popup")
	return c.FocusPanel(g, c.GetCurrentPanel())
}
