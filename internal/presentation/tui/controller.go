package tui

import (
	"github.com/awesome-gocui/gocui"
	appProject "github.com/Araryarch/lumine/internal/application/project"
	appService "github.com/Araryarch/lumine/internal/application/service"
	"github.com/Araryarch/lumine/internal/domain/project"
	"github.com/Araryarch/lumine/internal/domain/service"
	"github.com/Araryarch/lumine/internal/infrastructure/config"
)

type ViewMode string

const (
	ViewMain          ViewMode = "main"
	ViewServices      ViewMode = "services"
	ViewProjects      ViewMode = "projects"
	ViewDatabases     ViewMode = "databases"
	ViewCreateProject ViewMode = "create_project"
	ViewLogs          ViewMode = "logs"
	ViewSettings      ViewMode = "settings"
)

type Controller struct {
	gui            *gocui.Gui
	config         *config.Config
	serviceApp     *appService.Service
	projectApp     *appProject.Service
	currentView    ViewMode
	selectedIdx    int
	menuItems      []string
	projectList    []project.Project
	serviceList    []service.Status
	message        string
	messageType    string
}

func NewController(
	g *gocui.Gui,
	cfg *config.Config,
	serviceSvc *appService.Service,
	projectSvc *appProject.Service,
) *Controller {
	return &Controller{
		gui:         g,
		config:      cfg,
		serviceApp:  serviceSvc,
		projectApp:  projectSvc,
		currentView: ViewMain,
		selectedIdx: 0,
		menuItems:   []string{"Services", "Projects", "Databases", "Create Project", "Logs", "Settings"},
	}
}

func (c *Controller) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView("title", 0, 0, maxX-1, 2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Frame = false
		c.renderTitle(v)
	}

	menuWidth := 25
	if v, err := g.SetView("menu", 0, 3, menuWidth, maxY-5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = " Menu "
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		c.renderMenu(v)
	}

	if v, err := g.SetView("main", menuWidth+1, 3, maxX-1, maxY-5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = c.getMainTitle()
		v.Wrap = false
		c.renderMainView(v)
	}

	if c.message != "" {
		if v, err := g.SetView("message", 0, maxY-4, maxX-1, maxY-2, 0); err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}
			v.Frame = false
			c.renderMessage(v)
		}
	} else {
		g.DeleteView("message")
	}

	if v, err := g.SetView("status", 0, maxY-2, maxX-1, maxY, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Frame = false
		c.renderStatusBar(v)
	}

	if _, err := g.SetCurrentView("menu"); err != nil {
		return err
	}

	return nil
}

func (c *Controller) SetupKeybindings() error {
	bindings := []struct {
		key     interface{}
		handler func(*gocui.Gui, *gocui.View) error
	}{
		{gocui.KeyCtrlC, c.quit},
		{'q', c.quit},
		{'j', c.cursorDown},
		{'k', c.cursorUp},
		{gocui.KeyArrowDown, c.cursorDown},
		{gocui.KeyArrowUp, c.cursorUp},
		{gocui.KeyEnter, c.selectItem},
		{'b', c.goBack},
		{'s', c.startServices},
		{'x', c.stopServices},
		{'r', c.restartServices},
		{'d', c.deleteProject},
		{'?', c.showHelp},
	}

	for _, binding := range bindings {
		if err := c.gui.SetKeybinding("", binding.key, gocui.ModNone, binding.handler); err != nil {
			return err
		}
	}

	return nil
}
