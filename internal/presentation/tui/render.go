package tui

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/Araryarch/lumine/internal/domain/project"
)

func (c *Controller) renderTitle(v *gocui.View) {
	v.Clear()
	fmt.Fprintf(v, "\033[1;36m⚡ Lumine v2.0.0 - Docker Development Manager\033[0m")
}

func (c *Controller) renderMenu(v *gocui.View) {
	v.Clear()
	for i, item := range c.menuItems {
		if i == c.selectedIdx && c.currentView == ViewMain {
			fmt.Fprintf(v, "\033[1;32m▶ %s\033[0m\n", item)
		} else {
			fmt.Fprintf(v, "  %s\n", item)
		}
	}
}

func (c *Controller) getMainTitle() string {
	titles := map[ViewMode]string{
		ViewServices:      " Services ",
		ViewProjects:      " Projects ",
		ViewDatabases:     " Databases ",
		ViewCreateProject: " Create New Project ",
		ViewLogs:          " Logs ",
		ViewSettings:      " Settings ",
	}
	if title, ok := titles[c.currentView]; ok {
		return title
	}
	return " Welcome "
}

func (c *Controller) renderMainView(v *gocui.View) {
	v.Clear()

	switch c.currentView {
	case ViewMain:
		c.renderWelcome(v)
	case ViewServices:
		c.renderServices(v)
	case ViewProjects:
		c.renderProjects(v)
	case ViewDatabases:
		c.renderDatabases(v)
	case ViewCreateProject:
		c.renderCreateProject(v)
	case ViewLogs:
		c.renderLogs(v)
	case ViewSettings:
		c.renderSettings(v)
	}
}

func (c *Controller) renderWelcome(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mWelcome to Lumine!\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "A beautiful TUI for Docker development environment management.")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mQuick Start:\033[0m")
	fmt.Fprintln(v, "  • Navigate with j/k or arrow keys")
	fmt.Fprintln(v, "  • Press Enter to select a menu item")
	fmt.Fprintln(v, "  • Press q to quit")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "Select a menu item to get started!")
}

func (c *Controller) renderServices(v *gocui.View) {
	serviceList, err := c.serviceApp.GetAllStatuses()
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError: %v\033[0m\n", err)
		return
	}
	c.serviceList = serviceList

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
		if svc.Running {
			statusColor = "\033[1;32m"
		}

		fmt.Fprintf(v, "%-20s %s● %-13s\033[0m %-10d %-30s\n",
			svc.Name, statusColor, svc.State, svc.Port, svc.Image)
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m s=start all | x=stop all | r=restart all | b=back")
}

func (c *Controller) renderProjects(v *gocui.View) {
	projectList, err := c.projectApp.List()
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError: %v\033[0m\n", err)
		return
	}
	c.projectList = projectList

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
		if i == c.selectedIdx && c.currentView == ViewProjects {
			prefix = "\033[1;32m▶\033[0m "
		}
		
		typeColor := "\033[1;32m"
		switch proj.Type {
		case project.TypeLaravel:
			typeColor = "\033[1;31m"
		case project.TypeWordPress:
			typeColor = "\033[1;36m"
		}

		fmt.Fprintf(v, "%s%-25s %s%-13s\033[0m %-40s\n",
			prefix, proj.Name, typeColor, proj.Type, proj.URL)
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m d=delete | b=back")
}

func (c *Controller) renderDatabases(v *gocui.View) {
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

func (c *Controller) renderCreateProject(v *gocui.View) {
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
		if i == c.selectedIdx && c.currentView == ViewCreateProject {
			fmt.Fprintf(v, "\033[1;32m▶ %s\033[0m\n", ptype)
		} else {
			fmt.Fprintf(v, "  %s\n", ptype)
		}
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mNavigation:\033[0m j/k=navigate | Enter=select | b=back")
}

func (c *Controller) renderLogs(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mContainer Logs\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "Select a service to view logs")
	fmt.Fprintln(v, "")

	for i, svc := range c.serviceList {
		prefix := "  "
		if i == c.selectedIdx && c.currentView == ViewLogs {
			prefix = "\033[1;32m▶\033[0m "
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

func (c *Controller) renderSettings(v *gocui.View) {
	fmt.Fprintln(v, "\033[1;36mSettings\033[0m")
	fmt.Fprintln(v, "")
	fmt.Fprintf(v, "Default PHP Version: %s\n", c.config.DefaultPHP)
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mAvailable PHP Versions:\033[0m")
	for _, ver := range c.config.PHPVersions {
		fmt.Fprintf(v, "  • PHP %s\n", ver)
	}
	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mActions:\033[0m b=back")
}

func (c *Controller) renderMessage(v *gocui.View) {
	v.Clear()
	color := "\033[1;36m"
	switch c.messageType {
	case "success":
		color = "\033[1;32m"
	case "error":
		color = "\033[1;31m"
	case "warning":
		color = "\033[1;33m"
	}
	fmt.Fprintf(v, "%s%s\033[0m", color, c.message)
}

func (c *Controller) renderStatusBar(v *gocui.View) {
	v.Clear()

	runningServices := 0
	for _, svc := range c.serviceList {
		if svc.Running {
			runningServices++
		}
	}

	fmt.Fprintf(v, "\033[1;33mServices: %d/%d | Projects: %d | j/k: nav | Enter: select | q: quit\033[0m",
		runningServices, len(c.serviceList), len(c.projectList))
}
