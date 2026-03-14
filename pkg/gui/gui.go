package gui

import (
	"context"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/docker/docker/client"
	"github.com/yourusername/lumine/pkg/config"
	"github.com/yourusername/lumine/pkg/projects"
	"github.com/yourusername/lumine/pkg/services"
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

type GuiController struct {
	Gui          *gocui.Gui
	Config       *config.Config
	DockerClient *client.Client
	Context      context.Context
	CurrentView  ViewMode
	SelectedIdx  int
	MenuItems    []string
	ProjectList  []projects.Project
	ServiceList  []services.ServiceStatus
	Message      string
	MessageType  string
}

func NewGuiController(g *gocui.Gui, cfg *config.Config, dockerClient *client.Client, ctx context.Context) *GuiController {
	return &GuiController{
		Gui:          g,
		Config:       cfg,
		DockerClient: dockerClient,
		Context:      ctx,
		CurrentView:  ViewMain,
		SelectedIdx:  0,
		MenuItems:    []string{"Services", "Projects", "Databases", "Create Project", "Logs", "Settings"},
	}
}

func (gc *GuiController) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// Title bar
	if v, err := g.SetView("title", 0, 0, maxX-1, 2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Frame = false
		gc.renderTitle(v)
	}

	// Menu (left sidebar)
	menuWidth := 25
	if v, err := g.SetView("menu", 0, 3, menuWidth, maxY-5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = " Menu "
		v.Highlight = true
		v.SelBgColor = gocui.ColorCyan
		v.SelFgColor = gocui.ColorBlack
		gc.renderMenu(v)
	}

	// Main content
	if v, err := g.SetView("main", menuWidth+1, 3, maxX-1, maxY-5, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Title = gc.getMainTitle()
		v.Wrap = false
		v.Autoscroll = false
		gc.renderMainView(v)
	}

	// Message bar (if any)
	if gc.Message != "" {
		if v, err := g.SetView("message", 0, maxY-4, maxX-1, maxY-2, 0); err != nil {
			if !gocui.IsUnknownView(err) {
				return err
			}
			v.Frame = false
			gc.renderMessage(v)
		}
	} else {
		g.DeleteView("message")
	}

	// Status bar
	if v, err := g.SetView("status", 0, maxY-2, maxX-1, maxY, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Frame = false
		gc.renderStatusBar(v)
	}

	// Set current view
	if _, err := g.SetCurrentView("menu"); err != nil {
		return err
	}

	return nil
}

func (gc *GuiController) renderTitle(v *gocui.View) {
	v.Clear()
	fmt.Fprintf(v, "\033[1;36m⚡ Lumine v2.0.0 - Docker Development Manager\033[0m")
}

func (gc *GuiController) renderMenu(v *gocui.View) {
	v.Clear()
	for i, item := range gc.MenuItems {
		if i == gc.SelectedIdx && gc.CurrentView == ViewMain {
			fmt.Fprintf(v, "\033[1;34m▶ %s\033[0m\n", item)
		} else {
			fmt.Fprintf(v, "  %s\n", item)
		}
	}
}

func (gc *GuiController) getMainTitle() string {
	switch gc.CurrentView {
	case ViewServices:
		return " Services "
	case ViewProjects:
		return " Projects "
	case ViewDatabases:
		return " Databases "
	case ViewCreateProject:
		return " Create New Project "
	case ViewLogs:
		return " Logs "
	case ViewSettings:
		return " Settings "
	default:
		return " Welcome "
	}
}

func (gc *GuiController) renderMainView(v *gocui.View) {
	v.Clear()

	switch gc.CurrentView {
	case ViewMain:
		gc.renderWelcome(v)
	case ViewServices:
		gc.renderServices(v)
	case ViewProjects:
		gc.renderProjects(v)
	case ViewDatabases:
		gc.renderDatabases(v)
	case ViewCreateProject:
		gc.renderCreateProject(v)
	case ViewLogs:
		gc.renderLogs(v)
	case ViewSettings:
		gc.renderSettings(v)
	}
}

func (gc *GuiController) renderWelcome(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mWelcome to Lumine!\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "A beautiful TUI for Docker development environment management.")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mQuick Start:\033[0m")
	fmt.Fprintln(v, "  • Navigate with j/k or arrow keys")
	fmt.Fprintln(v, "  • Press Enter to select a menu item")
	fmt.Fprintln(v, "  • Press q to quit")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mFeatures:\033[0m")
	fmt.Fprintln(v, "  • Manage Docker services")
	fmt.Fprintln(v, "  • Create projects (PHP, Laravel, WordPress, Node.js)")
	fmt.Fprintln(v, "  • Database management")
	fmt.Fprintln(v, "  • View container logs")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "Select a menu item to get started!")
}

func (gc *GuiController) renderServices(v *gocui.View) {
	serviceList, err := services.GetServicesStatus(gc.Context, gc.DockerClient, gc.Config)
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError: %v\033[0m\n", err)
		return
	}
	gc.ServiceList = serviceList

	fmt.Fprintln(v, "\033[1;36mDocker Services\033[0m")
	fmt.Fprintln(v, "")
	
	if len(serviceList) == 0 {
		fmt.Fprintln(v, "\033[1;33mNo services configured\033[0m")
		return
	}

	fmt.Fprintf(v, "%-20s %-15s %-10s %-30s\n", "SERVICE", "STATUS", "PORT", "IMAGE")
	fmt.Fprintln(v, strings.Repeat("─", 75))

	for _, svc := range serviceList {
		statusColor := "\033[1;31m"
		statusIcon := "●"
		if svc.Running {
			statusColor = "\033[1;32m"
		}

		fmt.Fprintf(v, "%-20s %s%s %-13s\033[0m %-10d %-30s\n",
			svc.Name, statusColor, statusIcon, svc.Status, svc.Port, svc.Image)
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m s=start all | x=stop all | r=restart all | b=back")
}

func (gc *GuiController) renderProjects(v *gocui.View) {
	projectList, err := projects.ListProjects()
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError: %v\033[0m\n", err)
		return
	}
	gc.ProjectList = projectList

	fmt.Fprintln(v, "\033[1;36mYour Projects\033[0m")
	fmt.Fprintln(v, "")

	if len(projectList) == 0 {
		fmt.Fprintln(v, "\033[1;33mNo projects found\033[0m")
		fmt.Fprintln(v, "Create a new project from the 'Create Project' menu")
		return
	}

	fmt.Fprintf(v, "%-25s %-15s %-40s\n", "PROJECT", "TYPE", "URL")
	fmt.Fprintln(v, strings.Repeat("─", 80))

	for i, proj := range projectList {
		prefix := "  "
		if i == gc.SelectedIdx && gc.CurrentView == ViewProjects {
			prefix = "\033[1;34m▶\033[0m "
		}
		
		typeColor := "\033[1;32m"
		switch proj.Type {
		case projects.TypeLaravel:
			typeColor = "\033[1;31m"
		case projects.TypeWordPress:
			typeColor = "\033[1;34m"
		}

		fmt.Fprintf(v, "%s%-25s %s%-13s\033[0m %-40s\n",
			prefix, proj.Name, typeColor, proj.Type, proj.URL)
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m d=delete | b=back")
}

func (gc *GuiController) renderDatabases(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mDatabase Management\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;32m● MySQL\033[0m")
	fmt.Fprintln(v, "  Host: localhost | Port: 3306 | User: root | Pass: root")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;32m● Redis\033[0m")
	fmt.Fprintln(v, "  Host: localhost | Port: 6379")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;34m● phpMyAdmin\033[0m")
	fmt.Fprintln(v, "  URL: http://localhost:8080")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m b=back")
}

func (gc *GuiController) renderCreateProject(v *gocui.View) {
	projectTypes := []string{
		"PHP Project",
		"Laravel Project",
		"WordPress Site",
		"Node.js App",
		"Static HTML",
		"Back to Menu",
	}

	fmt.Fprintln(v, "\033[1;36mCreate New Project\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "Select project type:")
	fmt.Fprintln(v, "")

	for i, ptype := range projectTypes {
		if i == gc.SelectedIdx && gc.CurrentView == ViewCreateProject {
			fmt.Fprintf(v, "\033[1;34m▶ %s\033[0m\n", ptype)
		} else {
			fmt.Fprintf(v, "  %s\n", ptype)
		}
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mNavigation:\033[0m j/k=navigate | Enter=select | b=back")
}

func (gc *GuiController) renderLogs(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mContainer Logs\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "Select a service to view logs")
	fmt.Fprintln(v, "")

	for i, svc := range gc.ServiceList {
		prefix := "  "
		if i == gc.SelectedIdx && gc.CurrentView == ViewLogs {
			prefix = "\033[1;34m▶\033[0m "
		}
		
		statusColor := "\033[1;31m"
		if svc.Running {
			statusColor = "\033[1;32m"
		}

		fmt.Fprintf(v, "%s%s● %s\033[0m\n", prefix, statusColor, svc.Name)
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m Enter=view logs | b=back")
}

func (gc *GuiController) renderSettings(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mSettings\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintf(v, "Config Directory: %s\n", config.GetConfigDir())
	fmt.Fprintf(v, "Projects Directory: %s\n", config.GetProjectsDir())
	fmt.Fprintf(v, "Default PHP Version: %s\n", gc.Config.DefaultPHP)
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mAvailable PHP Versions:\033[0m")
	for _, ver := range gc.Config.PHPVersions {
		fmt.Fprintf(v, "  • PHP %s\n", ver)
	}
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m b=back")
}

func (gc *GuiController) renderMessage(v *gocui.View) {
	v.Clear()
	color := "\033[1;36m"
	switch gc.MessageType {
	case "success":
		color = "\033[1;32m"
	case "error":
		color = "\033[1;31m"
	case "warning":
		color = "\033[1;33m"
	}
	fmt.Fprintf(v, "%s%s\033[0m", color, gc.Message)
}

func (gc *GuiController) renderStatusBar(v *gocui.View) {
	v.Clear()
	
	runningServices := 0
	for _, svc := range gc.ServiceList {
		if svc.Running {
			runningServices++
		}
	}

	fmt.Fprintf(v, "\033[1;33mServices: %d/%d | Projects: %d | j/k: nav | Enter: select | q: quit\033[0m",
		runningServices, len(gc.ServiceList), len(gc.ProjectList))
}

func (gc *GuiController) SetupKeybindings() error {
	// Global keybindings
	if err := gc.Gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gc.quit); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", 'q', gocui.ModNone, gc.quit); err != nil {
		return err
	}

	// Navigation
	if err := gc.Gui.SetKeybinding("", 'j', gocui.ModNone, gc.cursorDown); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", 'k', gocui.ModNone, gc.cursorUp); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, gc.cursorDown); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, gc.cursorUp); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, gc.selectItem); err != nil {
		return err
	}

	// View-specific actions
	if err := gc.Gui.SetKeybinding("", 'b', gocui.ModNone, gc.goBack); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", 's', gocui.ModNone, gc.startServices); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", 'x', gocui.ModNone, gc.stopServices); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", 'r', gocui.ModNone, gc.restartServices); err != nil {
		return err
	}
	if err := gc.Gui.SetKeybinding("", 'd', gocui.ModNone, gc.deleteProject); err != nil {
		return err
	}

	return nil
}

func (gc *GuiController) quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func (gc *GuiController) cursorDown(g *gocui.Gui, v *gocui.View) error {
	maxIdx := len(gc.MenuItems) - 1
	
	switch gc.CurrentView {
	case ViewMain:
		if gc.SelectedIdx < maxIdx {
			gc.SelectedIdx++
		}
	case ViewProjects:
		if gc.SelectedIdx < len(gc.ProjectList)-1 {
			gc.SelectedIdx++
		}
	case ViewCreateProject:
		if gc.SelectedIdx < 5 {
			gc.SelectedIdx++
		}
	case ViewLogs:
		if gc.SelectedIdx < len(gc.ServiceList)-1 {
			gc.SelectedIdx++
		}
	}
	
	return nil
}

func (gc *GuiController) cursorUp(g *gocui.Gui, v *gocui.View) error {
	if gc.SelectedIdx > 0 {
		gc.SelectedIdx--
	}
	return nil
}

func (gc *GuiController) selectItem(g *gocui.Gui, v *gocui.View) error {
	switch gc.CurrentView {
	case ViewMain:
		gc.selectMenuItem()
	case ViewCreateProject:
		gc.selectProjectType()
	}
	return nil
}

func (gc *GuiController) selectMenuItem() {
	gc.SelectedIdx = 0
	switch gc.SelectedIdx {
	case 0:
		gc.CurrentView = ViewServices
	case 1:
		gc.CurrentView = ViewProjects
	case 2:
		gc.CurrentView = ViewDatabases
	case 3:
		gc.CurrentView = ViewCreateProject
	case 4:
		gc.CurrentView = ViewLogs
	case 5:
		gc.CurrentView = ViewSettings
	}
}

func (gc *GuiController) selectProjectType() {
	// Implementation for project creation
	gc.showMessage("Project creation coming soon!", "info")
}

func (gc *GuiController) goBack(g *gocui.Gui, v *gocui.View) error {
	if gc.CurrentView != ViewMain {
		gc.CurrentView = ViewMain
		gc.SelectedIdx = 0
		gc.Message = ""
	}
	return nil
}

func (gc *GuiController) startServices(g *gocui.Gui, v *gocui.View) error {
	if gc.CurrentView == ViewServices {
		go func() {
			if err := services.StartAllServices(gc.Context, gc.DockerClient, gc.Config); err != nil {
				gc.showMessage(fmt.Sprintf("Error: %v", err), "error")
			} else {
				gc.showMessage("All services started successfully", "success")
			}
			g.Update(func(g *gocui.Gui) error { return nil })
		}()
	}
	return nil
}

func (gc *GuiController) stopServices(g *gocui.Gui, v *gocui.View) error {
	if gc.CurrentView == ViewServices {
		go func() {
			if err := services.StopAllServices(gc.Context, gc.DockerClient, gc.Config); err != nil {
				gc.showMessage(fmt.Sprintf("Error: %v", err), "error")
			} else {
				gc.showMessage("All services stopped successfully", "success")
			}
			g.Update(func(g *gocui.Gui) error { return nil })
		}()
	}
	return nil
}

func (gc *GuiController) restartServices(g *gocui.Gui, v *gocui.View) error {
	if gc.CurrentView == ViewServices {
		go func() {
			if err := services.RestartAllServices(gc.Context, gc.DockerClient, gc.Config); err != nil {
				gc.showMessage(fmt.Sprintf("Error: %v", err), "error")
			} else {
				gc.showMessage("All services restarted successfully", "success")
			}
			g.Update(func(g *gocui.Gui) error { return nil })
		}()
	}
	return nil
}

func (gc *GuiController) deleteProject(g *gocui.Gui, v *gocui.View) error {
	if gc.CurrentView == ViewProjects && len(gc.ProjectList) > 0 && gc.SelectedIdx < len(gc.ProjectList) {
		project := gc.ProjectList[gc.SelectedIdx]
		if err := projects.DeleteProject(project.Name); err != nil {
			gc.showMessage(fmt.Sprintf("Error: %v", err), "error")
		} else {
			gc.showMessage(fmt.Sprintf("Project '%s' deleted", project.Name), "success")
		}
	}
	return nil
}

func (gc *GuiController) showMessage(msg string, msgType string) {
	gc.Message = msg
	gc.MessageType = msgType
}
