package tui

import (
	"fmt"
	"strings"

	"github.com/Araryarch/lumine/internal/domain/project"
	"github.com/jesseduffield/gocui"
)

func (c *Controller) renderTitle(v *gocui.View) {
	v.Clear()
	fmt.Fprintf(v, "\033[1;32m⚡ Lumine\033[0m - Docker Development Manager (Laragon for Linux)")
}

func (c *Controller) getMainTitle() string {
	if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		return fmt.Sprintf(" %s - Details ", c.serviceList[c.selectedService].Name)
	}
	if len(c.projectList) > 0 && c.selectedProject < len(c.projectList) {
		return fmt.Sprintf(" %s - Details ", c.projectList[c.selectedProject].Name)
	}
	return " Info "
}

func (c *Controller) getLogsTitle() string {
	if len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		return c.serviceList[c.selectedService].Name
	}
	return ""
}

func (c *Controller) renderServicesList(v *gocui.View) {
	v.Clear()

	serviceList, err := c.serviceApp.GetAllStatuses()
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError: %v\033[0m\n", err)
		return
	}
	c.serviceList = serviceList

	if len(serviceList) == 0 {
		fmt.Fprintln(v, "\033[1;33mNo services configured\033[0m")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "Add services in ~/.lumine/config.yaml")
		return
	}

	for i, svc := range serviceList {
		prefix := "  "
		if i == c.selectedService {
			prefix = "\033[1;32m▶\033[0m "
		}

		statusIcon := "●"
		statusColor := "\033[1;31m"
		if svc.Running {
			statusColor = "\033[1;32m"
		}

		fmt.Fprintf(v, "%s%s%s\033[0m %s\n", prefix, statusColor, statusIcon, svc.Name)
		if i == c.selectedService {
			fmt.Fprintf(v, "  \033[2m%s:%d\033[0m\n", svc.Image, svc.Port)
		}
	}
}

func (c *Controller) renderProjectsList(v *gocui.View) {
	v.Clear()

	projectList, err := c.projectApp.List()
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError: %v\033[0m\n", err)
		return
	}
	c.projectList = projectList

	if len(projectList) == 0 {
		fmt.Fprintln(v, "\033[1;37mNo projects found\033[0m")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "Press 'n' to create a new project")
		return
	}

	for i, proj := range projectList {
		prefix := "  "
		if i == c.selectedProject {
			prefix = "\033[1;32m▶\033[0m "
		}

		typeIcon := "📁"
		switch proj.Type {
		case project.TypePHP:
			typeIcon = "🐘"
		case project.TypeLaravel:
			typeIcon = "🔺"
		case project.TypeWordPress:
			typeIcon = "📝"
		case project.TypeNodeJS:
			typeIcon = "🟢"
		case project.TypeStatic:
			typeIcon = "📄"
		}

		fmt.Fprintf(v, "%s%s %s\n", prefix, typeIcon, proj.Name)
		if i == c.selectedProject {
			fmt.Fprintf(v, "  \033[2m%s\033[0m\n", proj.URL)
		}
	}
}

func (c *Controller) renderMainPanel(v *gocui.View) {
	v.Clear()

	// Show service logs or project details
	if c.currentView == ViewServices && len(c.serviceList) > 0 && c.selectedService < len(c.serviceList) {
		svc := c.serviceList[c.selectedService]

		fmt.Fprintln(v, "\033[1;32mService Information\033[0m")
		fmt.Fprintln(v, strings.Repeat("─", 60))
		fmt.Fprintf(v, "Name:   %s\n", svc.Name)
		fmt.Fprintf(v, "Image:  %s\n", svc.Image)
		fmt.Fprintf(v, "Port:   %d\n", svc.Port)

		statusColor := "\033[1;31m"
		statusText := "Stopped"
		if svc.Running {
			statusColor = "\033[1;32m"
			statusText = "Running"
		}
		fmt.Fprintf(v, "Status: %s%s\033[0m\n", statusColor, statusText)
		fmt.Fprintln(v, "")

		fmt.Fprintln(v, "\033[1;33mActions:\033[0m")
		fmt.Fprintln(v, "  s - Start service")
		fmt.Fprintln(v, "  x - Stop service")
		fmt.Fprintln(v, "  r - Restart service")
		fmt.Fprintln(v, "  S - Start all services")
		fmt.Fprintln(v, "  X - Stop all services")

	} else if c.currentView == ViewProjects && len(c.projectList) > 0 && c.selectedProject < len(c.projectList) {
		proj := c.projectList[c.selectedProject]

		fmt.Fprintln(v, "\033[1;32mProject Information\033[0m")
		fmt.Fprintln(v, strings.Repeat("─", 60))
		fmt.Fprintf(v, "Name: %s\n", proj.Name)
		fmt.Fprintf(v, "Type: %s\n", proj.Type)
		fmt.Fprintf(v, "Path: %s\n", proj.Path)
		fmt.Fprintf(v, "URL:  %s\n", proj.URL)
		fmt.Fprintln(v, "")

		fmt.Fprintln(v, "\033[1;33mActions:\033[0m")
		fmt.Fprintln(v, "  d - Delete project")
		fmt.Fprintln(v, "  n - Create new project")

	} else {
		fmt.Fprintln(v, "\033[1;32mWelcome to Lumine!\033[0m")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "A Docker-based development environment manager for Linux")
		fmt.Fprintln(v, "Inspired by Laragon for Windows")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "\033[1;33mQuick Start:\033[0m")
		fmt.Fprintln(v, "  • Use Tab to switch between Services and Projects")
		fmt.Fprintln(v, "  • Use j/k or arrow keys to navigate")
		fmt.Fprintln(v, "  • Press Enter to view details")
		fmt.Fprintln(v, "  • Press ? for help")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "\033[1;33mServices:\033[0m")
		fmt.Fprintln(v, "  Manage Docker containers (Nginx, MySQL, Redis, etc.)")
		fmt.Fprintln(v, "")
		fmt.Fprintln(v, "\033[1;33mProjects:\033[0m")
		fmt.Fprintln(v, "  Manage your web projects (PHP, Laravel, WordPress, Node.js)")
	}
}

func (c *Controller) renderStatusBar(v *gocui.View) {
	v.Clear()

	runningServices := 0
	for _, svc := range c.serviceList {
		if svc.Running {
			runningServices++
		}
	}

	// First line: Status info
	statusLine := fmt.Sprintf("\033[1;32m●\033[0m %d/%d services | %d projects",
		runningServices, len(c.serviceList), len(c.projectList))

	// Message
	if c.message != "" {
		color := "\033[1;37m"
		switch c.messageType {
		case "success":
			color = "\033[1;32m"
		case "error":
			color = "\033[1;31m"
		case "warning":
			color = "\033[1;33m"
		}
		statusLine += fmt.Sprintf(" | %s%s\033[0m", color, c.message)
	}

	fmt.Fprintln(v, statusLine)

	// Second line: Shortcuts based on current panel
	var shortcuts string
	switch c.currentView {
	case ViewServices:
		shortcuts = "\033[1;33m[Services]\033[0m j/k:navigate Enter:select l:logs Tab:switch R:refresh ?:help q:quit"
	case ViewProjects:
		shortcuts = "\033[1;33m[Projects]\033[0m j/k:navigate Enter:select d:delete n:new Tab:switch R:refresh ?:help q:quit"
	case ViewMain:
		shortcuts = "\033[1;33m[Main]\033[0m s:start x:stop r:restart S:start-all X:stop-all ←:back R:refresh ?:help q:quit"
	case ViewLogs:
		shortcuts = "\033[1;33m[Logs]\033[0m l:toggle-logs ←/Esc:back R:refresh q:quit"
	case ViewCreateProject:
		shortcuts = "\033[1;33m[Create]\033[0m j/k:navigate Enter:select b/Esc:back"
	default:
		shortcuts = "Tab:switch j/k:navigate Enter:select R:refresh ?:help q:quit"
	}

	fmt.Fprintln(v, shortcuts)
}

func (c *Controller) renderLogsPanel(v *gocui.View) {
	v.Clear()

	if len(c.serviceList) == 0 || c.selectedService >= len(c.serviceList) {
		fmt.Fprintln(v, "\033[1;33mNo service selected\033[0m")
		return
	}

	svc := c.serviceList[c.selectedService]

	fmt.Fprintf(v, "\033[1;36mLogs for %s\033[0m\n", svc.Name)
	fmt.Fprintln(v, strings.Repeat("─", 60))

	logs, err := c.serviceApp.GetLogs(svc.Name, "100")
	if err != nil {
		fmt.Fprintf(v, "\033[1;31mError fetching logs: %v\033[0m\n", err)
		return
	}

	if logs == "" {
		fmt.Fprintln(v, "\033[1;33mNo logs available\033[0m")
		return
	}

	// Clean up docker logs output (remove header bytes)
	cleanLogs := c.cleanDockerLogs(logs)
	fmt.Fprintln(v, cleanLogs)
}

func (c *Controller) cleanDockerLogs(logs string) string {
	// Docker logs may have extra bytes at the beginning, let's try to clean it
	// This is a simple approach - just return as is for now
	return logs
}

func (c *Controller) renderCreateProjectPanel(v *gocui.View) {
	v.Clear()

	fmt.Fprintln(v, "\033[1;36mCreate New Project\033[0m")
	fmt.Fprintln(v, strings.Repeat("─", 60))
	fmt.Fprintln(v, "")

	projectTypes := c.projectTypes
	for i, ptype := range projectTypes {
		prefix := "  "
		if i == c.selectedProjectType {
			prefix = "\033[1;32m▶\033[0m "
		}
		fmt.Fprintf(v, "%s%s\n", prefix, ptype)
	}

	fmt.Fprintln(v, "")
	fmt.Fprintln(v, "\033[1;33mInstructions:\033[0m")
	fmt.Fprintln(v, "  j/k - Navigate")
	fmt.Fprintln(v, "  Enter - Select project type")
	fmt.Fprintln(v, "  b/Esc - Back to menu")
}
