package tui

import (
	appProject "github.com/Araryarch/lumine/internal/application/project"
	appService "github.com/Araryarch/lumine/internal/application/service"
	"github.com/Araryarch/lumine/internal/domain/project"
	"github.com/Araryarch/lumine/internal/domain/service"
	"github.com/Araryarch/lumine/internal/infrastructure/config"
	"github.com/jesseduffield/gocui"
)

type ViewMode string

const (
	ViewServices      ViewMode = "services"
	ViewProjects      ViewMode = "projects"
	ViewMain          ViewMode = "main"
	ViewLogs          ViewMode = "logs"
	ViewCreateProject ViewMode = "create_project"
)

type Controller struct {
	gui                 *gocui.Gui
	config              *config.Config
	serviceApp          *appService.Service
	projectApp          *appProject.Service
	currentView         ViewMode
	selectedIdx         int
	serviceList         []service.Status
	projectList         []project.Project
	message             string
	messageType         string
	selectedService     int
	selectedProject     int
	projectTypes        []string
	selectedProjectType int
}

func NewController(
	g *gocui.Gui,
	cfg *config.Config,
	serviceSvc *appService.Service,
	projectSvc *appProject.Service,
) *Controller {
	return &Controller{
		gui:                 g,
		config:              cfg,
		serviceApp:          serviceSvc,
		projectApp:          projectSvc,
		currentView:         ViewServices,
		selectedIdx:         0,
		selectedService:     0,
		selectedProject:     0,
		projectTypes:        []string{"Static", "Node.js", "PHP", "Laravel", "WordPress"},
		selectedProjectType: 0,
	}
}

func (c *Controller) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Title bar
	if v, err := g.SetView("title", 0, 0, maxX-1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}
	if v, _ := g.View("title"); v != nil {
		c.renderTitle(v)
	}

	// Create all panels
	panels := c.GetPanels(maxX, maxY)
	for _, panel := range panels {
		if err := c.CreatePanel(g, panel); err != nil {
			return err
		}
		if err := c.RenderPanel(g, panel.Name); err != nil {
			return err
		}
	}

	// Status bar
	if v, err := g.SetView("status", 0, maxY-3, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Frame = false
	}
	if v, _ := g.View("status"); v != nil {
		c.renderStatusBar(v)
	}

	// Set initial focus
	if _, err := g.View(c.GetCurrentPanel()); err == nil {
		g.SetCurrentView(c.GetCurrentPanel())
	}

	return nil
}

func (c *Controller) SetupKeybindings() error {
	// Global keybindings
	if err := c.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, c.quit); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", 'q', gocui.ModNone, c.quit); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", '?', gocui.ModNone, c.showHelp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", 'R', gocui.ModNone, c.RefreshView); err != nil {
		return err
	}

	// Mouse support - click on views
	if err := c.gui.SetKeybinding("services", gocui.MouseLeft, gocui.ModNone, c.selectService); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", gocui.MouseLeft, gocui.ModNone, c.selectProject); err != nil {
		return err
	}

	// Navigation between panels
	if err := c.gui.SetKeybinding("", gocui.KeyTab, gocui.ModNone, c.nextPanel); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, c.focusServices); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, c.focusMain); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", 'l', gocui.ModNone, c.toggleLogs); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, c.closeLogs); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("", 'b', gocui.ModNone, c.closeLogs); err != nil {
		return err
	}

	// Services panel keybindings - ONLY navigation and selection
	if err := c.gui.SetKeybinding("services", 'j', gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("services", 'k', gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("services", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("services", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("services", gocui.KeyEnter, gocui.ModNone, c.selectService); err != nil {
		return err
	}

	// Main panel keybindings - service actions
	if err := c.gui.SetKeybinding("main", 's', gocui.ModNone, c.startService); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("main", 'x', gocui.ModNone, c.stopService); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("main", 'r', gocui.ModNone, c.restartService); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("main", 'S', gocui.ModNone, c.startAllServices); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("main", 'X', gocui.ModNone, c.stopAllServices); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("main", 'l', gocui.ModNone, c.toggleLogs); err != nil {
		return err
	}

	// Services panel keybindings - logs
	if err := c.gui.SetKeybinding("services", 'l', gocui.ModNone, c.toggleLogs); err != nil {
		return err
	}

	// Logs panel keybindings
	if err := c.gui.SetKeybinding("logs", gocui.KeyEsc, gocui.ModNone, c.closeLogs); err != nil {
		return err
	}

	// Projects panel keybindings
	if err := c.gui.SetKeybinding("projects", 'j', gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", 'k', gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", gocui.KeyEnter, gocui.ModNone, c.selectProject); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", 'd', gocui.ModNone, c.deleteProject); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("projects", 'n', gocui.ModNone, c.newProject); err != nil {
		return err
	}

	// Create project panel keybindings
	if err := c.gui.SetKeybinding("createproject", 'j', gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("createproject", 'k', gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("createproject", gocui.KeyArrowDown, gocui.ModNone, c.cursorDown); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("createproject", gocui.KeyArrowUp, gocui.ModNone, c.cursorUp); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("createproject", gocui.KeyEnter, gocui.ModNone, c.selectProjectType); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("createproject", 'b', gocui.ModNone, c.cancelCreateProject); err != nil {
		return err
	}
	if err := c.gui.SetKeybinding("createproject", gocui.KeyEsc, gocui.ModNone, c.cancelCreateProject); err != nil {
		return err
	}

	return nil
}
