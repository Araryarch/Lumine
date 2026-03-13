package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"lumine/internal/config"
)

func (m model) renderSidebarFixed(width, height int) string {
	var s strings.Builder

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == sidebarPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColor).
		Background(bgColor).
		Padding(1, 0).
		Width(width).
		Height(height)

	if m.activePanel == sidebarPanel {
		style = style.BorderForeground(primaryColor)
	}

	s.WriteString(lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		MarginBottom(1).
		Render("󰘤 Navigation") + "\n\n")

	menuItems := []struct {
		name string
		icon string
	}{
		{"Services", "󰘦"},
		{"Projects", "󰉋"},
		{"Databases", "󱆟"},
		{"Runtimes", "󰌠"},
		{"Logs", "󰌱"},
		{"Tasks", "󰘦"},
	}

	for i, item := range menuItems {
		var line strings.Builder

		if m.sidebarCursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("  "))
			line.WriteString(item.icon + " ")
			nameStyle := lipgloss.NewStyle().
				Foreground(bgColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1).
				Render(item.name)
			line.WriteString(nameStyle)
		} else {
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("   "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(infoColor).
				Render(item.icon + " "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(fgColor).
				Render(item.name))
		}

		lineStr := line.String()
		s.WriteString(lineStr + "\n")
	}

	s.WriteString("\n")

	dividerWidth := width - 4
	divider := lipgloss.NewStyle().
		Foreground(surface1).
		Render(strings.Repeat("─", dividerWidth))
	s.WriteString("  " + divider + "\n\n")

	shortcutItems := []struct {
		name string
		icon string
	}{
		{"New Project", "󰈙"},
		{"Settings", "󰒓"},
	}

	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Bold(true).
		Padding(0, 1).
		Render("Shortcuts") + "\n\n")

	for i, item := range shortcutItems {
		idx := i + len(menuItems)
		var line strings.Builder

		if m.sidebarCursor == idx {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("  "))
			line.WriteString(item.icon + " ")
			nameStyle := lipgloss.NewStyle().
				Foreground(bgColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1).
				Render(item.name)
			line.WriteString(nameStyle)
		} else {
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("   "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(infoColor).
				Render(item.icon + " "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(fgColor).
				Render(item.name))
		}

		lineStr := line.String()
		s.WriteString(lineStr + "\n")
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < height-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderMainContent() string {
	switch m.currentView {
	case servicesView:
		return m.renderServicesPanel()
	case projectsView:
		return m.renderProjectsPanel()
	case databasesView:
		return m.renderDatabasePanel()
	case runtimesView:
		return m.renderRuntimesPanel()
	case logsView:
		return m.renderLogsPanel()
	case backgroundTasksView:
		return m.renderBackgroundTasksPanel()
	case addServiceView:
		return m.renderAddServicePanel()
	case settingsView:
		return m.renderSettingsPanel()
	default:
		return m.renderServicesPanel()
	}
}

func (m model) renderServicesPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == mainPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == mainPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰘦  Services")

	s.WriteString(titleStyle + "\n\n")

	if len(m.config.Services) == 0 {
		emptyIcon := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("󰉋")

		emptyTitle := lipgloss.NewStyle().
			Foreground(subtleColor).
			Bold(true).
			Render("No services configured")

		emptyDesc := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Press 'n' to add a new service")

		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(2, 4).
			Width(panelWidth - 8).
			Align(lipgloss.Center).
			Render(lipgloss.JoinVertical(lipgloss.Center,
				emptyIcon,
				"",
				emptyTitle,
				emptyDesc,
			))

		s.WriteString(emptyBox)
	} else {
		maxVisibleItems := panelHeight - 5
		if maxVisibleItems < 3 {
			maxVisibleItems = 3
		}

		startIdx := 0
		endIdx := len(m.config.Services)
		showScrollTop := false
		showScrollBottom := false

		if len(m.config.Services) > maxVisibleItems {
			startIdx = m.cursor - (maxVisibleItems / 2)
			if startIdx < 0 {
				startIdx = 0
			}
			endIdx = startIdx + maxVisibleItems
			if endIdx > len(m.config.Services) {
				endIdx = len(m.config.Services)
				startIdx = endIdx - maxVisibleItems
				if startIdx < 0 {
					startIdx = 0
				}
			}

			showScrollTop = startIdx > 0
			showScrollBottom = endIdx < len(m.config.Services)
		}

		if showScrollTop {
			scrollIndicator := lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Render(fmt.Sprintf("  ↑ %d more", startIdx))
			s.WriteString(scrollIndicator + "\n")
		}

		for i := startIdx; i < endIdx; i++ {
			service := m.config.Services[i]
			var line strings.Builder

			icon := getIconForService(service.Type)
			if m.cursor == i {
				line.WriteString(lipgloss.NewStyle().
					Foreground(primaryColor).
					Bold(true).
					Render(" "))
			} else {
				line.WriteString(lipgloss.NewStyle().
					Foreground(mutedColor).
					Render("  "))
			}

			badge := getServiceBadge(service.Type).Render(strings.ToUpper(service.Type))
			line.WriteString(badge + " ")

			iconStyle := lipgloss.NewStyle().Foreground(fgColor)
			line.WriteString(iconStyle.Render(icon) + " ")

			nameStyle := lipgloss.NewStyle().Bold(true).Foreground(fgColor)
			line.WriteString(nameStyle.Render(service.Name))

			status := m.serviceStatus[service.Name]
			if status.running {
				statusBadge := lipgloss.NewStyle().
					Foreground(bgColor).
					Background(successColor).
					Padding(0, 1).
					Bold(true).
					Render(" 󰀄 ")
				line.WriteString(" " + statusBadge)
			} else {
				statusBadge := lipgloss.NewStyle().
					Foreground(bgColor).
					Background(mutedColor).
					Padding(0, 1).
					Render(" 󰀊 ")
				line.WriteString(" " + statusBadge)
			}

			portBadge := lipgloss.NewStyle().
				Foreground(warningColor).
				Background(surface0).
				Padding(0, 1).
				Render(fmt.Sprintf(":%d", service.Port))
			line.WriteString(" " + portBadge)

			lineStr := line.String()
			if m.cursor == i {
				lineStr = lipgloss.NewStyle().
					Background(surface0).
					Width(panelWidth-4).
					Padding(0, 1).
					Render(lineStr)
			} else {
				lineStr = lipgloss.NewStyle().
					Foreground(fgColor).
					Padding(0, 1).
					Render(lineStr)
			}

			s.WriteString(lineStr + "\n")
		}

		if showScrollBottom {
			scrollIndicator := lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Render(fmt.Sprintf("  ↓ %d more", len(m.config.Services)-endIdx))
			s.WriteString(scrollIndicator + "\n")
		}
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderAddServicePanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == mainPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == mainPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰈙  Add New Service")

	s.WriteString(titleStyle + "\n\n")

	availableServices := []struct {
		name string
		desc string
		icon string
	}{
		{"PHP", "PHP-FPM or Apache with PHP", "󰌞"},
		{"MySQL", "MySQL Database Server", "󱆟"},
		{"MariaDB", "MariaDB Database Server", "󱆟"},
		{"PostgreSQL", "PostgreSQL Database", "󱆢"},
		{"Nginx", "Nginx Web Server", "󰖟"},
		{"Apache", "Apache HTTP Server", "󰖟"},
		{"Caddy", "Caddy Web Server", "󰖟"},
		{"Redis", "Redis Cache Server", "󰝚"},
		{"MongoDB", "MongoDB NoSQL Database", "󱆦"},
		{"phpMyAdmin", "MySQL Web Interface", "󰖶"},
		{"Adminer", "Database Management Tool", "󱆦"},
		{"Elasticsearch", "Search Engine", "󰉋"},
		{"RabbitMQ", "Message Queue", "󰘦"},
		{"Memcached", "Memory Cache", "󰘦"},
	}

	subHeader := lipgloss.NewStyle().
		Foreground(infoColor).
		Render("Select a service type to add:")
	s.WriteString(subHeader + "\n\n")

	maxVisibleItems := panelHeight - 8
	if maxVisibleItems < 5 {
		maxVisibleItems = 5
	}

	startIdx := 0
	endIdx := len(availableServices)
	showScrollTop := false
	showScrollBottom := false

	if len(availableServices) > maxVisibleItems {
		startIdx = m.cursor - (maxVisibleItems / 2)
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + maxVisibleItems
		if endIdx > len(availableServices) {
			endIdx = len(availableServices)
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}

		showScrollTop = startIdx > 0
		showScrollBottom = endIdx < len(availableServices)
	}

	if showScrollTop {
		scrollInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↑ %d more", startIdx))
		s.WriteString(scrollInfo + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		svc := availableServices[i]
		var line strings.Builder

		if m.cursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(" "))
		} else {
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("  "))
		}

		badge := getServiceBadge(strings.ToLower(svc.name)).Render(svc.name)
		line.WriteString(badge + " ")

		iconStyle := lipgloss.NewStyle().Foreground(infoColor)
		line.WriteString(iconStyle.Render(svc.icon) + " ")

		desc := lipgloss.NewStyle().
			Foreground(fgColor).
			Render(svc.desc)
		line.WriteString(desc)

		lineStr := line.String()

		if m.cursor == i {
			lineStr = lipgloss.NewStyle().
				Background(surface0).
				Foreground(fgColor).
				Width(panelWidth-4).
				Padding(0, 1).
				Render(lineStr)
		} else {
			lineStr = lipgloss.NewStyle().
				Foreground(fgColor).
				Padding(0, 1).
				Render(lineStr)
		}

		s.WriteString(lineStr + "\n")
	}

	if showScrollBottom {
		scrollInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↓ %d more", len(availableServices)-endIdx))
		s.WriteString(scrollInfo + "\n")
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderSettingsPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == mainPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == mainPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰒓  Settings")

	s.WriteString(titleStyle + "\n\n")

	serverBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var serverContent strings.Builder

	serverContent.WriteString(lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Render("󰖟  Web Server") + "\n\n")

	servers := []struct {
		name   string
		active bool
	}{
		{"Nginx", true},
		{"Apache", false},
		{"Caddy", false},
	}

	for _, srv := range servers {
		var line strings.Builder
		if srv.active {
			line.WriteString(lipgloss.NewStyle().
				Foreground(successColor).
				Render("󰀄 "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(successColor).
				Padding(0, 1).
				Bold(true).
				Render(srv.name + "  ACTIVE"))
		} else {
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("󰀊 "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(srv.name))
		}
		serverContent.WriteString(line.String() + "\n")
	}

	s.WriteString(serverBox.Render(serverContent.String()) + "\n\n")

	runtimeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(infoColor).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var runtimeContent strings.Builder

	runtimeContent.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true).
		Render("󰌠  Runtime Versions") + "\n\n")

	runtimes := []struct {
		name    string
		version string
		icon    string
	}{
		{"PHP", m.config.Runtimes.PHP, "󰌞"},
		{"Node.js", m.config.Runtimes.Node, "󰛦"},
		{"Python", m.config.Runtimes.Python, "󰌠"},
		{"Go", m.config.Runtimes.Go, "󰟓"},
		{"Rust", m.config.Runtimes.Rust, "󱘘"},
		{"Bun", m.config.Runtimes.Bun, "󰛦"},
		{"Deno", m.config.Runtimes.Deno, "󰛦"},
	}

	for _, rt := range runtimes {
		var line strings.Builder
		line.WriteString(rt.icon + " ")
		line.WriteString(lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true).
			Render(rt.name))
		line.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(": "))
		line.WriteString(lipgloss.NewStyle().
			Foreground(successColor).
			Background(surface0).
			Padding(0, 1).
			Render("v" + rt.version))
		runtimeContent.WriteString(line.String() + "\n")
	}

	s.WriteString(runtimeBox.Render(runtimeContent.String()) + "\n\n")

	helpBox := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Width(panelWidth - 8).
		Align(lipgloss.Center).
		Render("Press 'v' on Runtimes page to change versions")

	s.WriteString(helpBox)

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderDetailPanelDynamic() string {
	switch m.currentView {
	case servicesView:
		return m.renderServiceDetailPanel()
	case projectsView:
		return m.renderProjectDetailPanel()
	case databasesView:
		return m.renderDatabaseDetailPanel()
	case runtimesView:
		return m.renderRuntimeDetailPanel()
	case logsView:
		return m.renderLogsDetailPanel()
	case backgroundTasksView:
		return m.renderTasksDetailPanel()
	case settingsView:
		return m.renderSettingsDetailPanel()
	default:
		return m.renderServiceDetailPanel()
	}
}

func (m model) renderServiceDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	if m.cursor < len(m.config.Services) {
		service := m.config.Services[m.cursor]
		status := m.serviceStatus[service.Name]

		titleStyle := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Padding(0, 1).
			Background(surface0).
			Width(panelWidth - 4).
			Render("󰡨  Service Details")

		s.WriteString(titleStyle + "\n\n")

		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(1, 2).
			Width(panelWidth - 8).
			Background(surfaceBg)

		var info strings.Builder

		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Name: "))
		info.WriteString(lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Render(service.Name) + "\n")

		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Type: "))
		badge := getServiceBadge(service.Type).Render(strings.ToUpper(service.Type))
		info.WriteString(badge + "\n")

		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Version: "))
		info.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Background(surface0).
			Padding(0, 1).
			Render("v"+service.Version) + "\n")

		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Port: "))
		info.WriteString(lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Render(fmt.Sprintf(":%d", service.Port)) + "\n")

		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Status: "))
		if status.running {
			info.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(successColor).
				Padding(0, 1).
				Bold(true).
				Render(" 󰀄 Running") + "\n")
		} else {
			info.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(mutedColor).
				Padding(0, 1).
				Render(" 󰀊 Stopped") + "\n")
		}

		if status.containerId != "" {
			info.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("Container: "))
			info.WriteString(lipgloss.NewStyle().
				Foreground(subtleColor).
				Render(status.containerId[:12]) + "\n")
		}

		s.WriteString(infoBox.Render(info.String()) + "\n\n")

		if len(service.Env) > 0 {
			envHeader := lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				Render("󰌮  Environment Variables")
			s.WriteString(envHeader + "\n\n")

			envBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(surface1).
				Padding(1, 2).
				Width(panelWidth - 8).
				Background(surfaceBg)

			var envContent strings.Builder
			for k, v := range service.Env {
				envContent.WriteString(lipgloss.NewStyle().
					Foreground(infoColor).
					Render(k))
				envContent.WriteString(lipgloss.NewStyle().
					Foreground(mutedColor).
					Render("="))
				envContent.WriteString(lipgloss.NewStyle().
					Foreground(fgColor).
					Render(v) + "\n")
			}

			s.WriteString(envBox.Render(envContent.String()) + "\n")
		}
	} else {
		emptyState := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("󰉋  No service selected")
		s.WriteString(emptyState)
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderProjectDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰉋  Project Details")

	s.WriteString(titleStyle + "\n\n")

	if m.cursor < len(m.config.Projects) {
		project := m.config.Projects[m.cursor]

		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(1, 2).
			Width(panelWidth - 8).
			Background(surfaceBg)

		var info strings.Builder

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Name: "))
		info.WriteString(lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(project.Name) + "\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Type: "))
		badge := getProjectBadge(project.Type).Render(strings.ToUpper(project.Type))
		info.WriteString(badge + "\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Path: "))
		info.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(project.Path) + "\n")

		if project.Domain != "" {
			info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Domain: "))
			info.WriteString(lipgloss.NewStyle().Foreground(warningColor).Render(project.Domain) + "\n")
		}

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Status: "))
		if project.Status == "running" {
			info.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(successColor).
				Padding(0, 1).
				Bold(true).
				Render(" 󰀄 Running") + "\n")
		} else {
			info.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(mutedColor).
				Padding(0, 1).
				Render(" 󰀊 Stopped") + "\n")
		}

		s.WriteString(infoBox.Render(info.String()))
	} else {
		emptyState := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("󰉋  No project selected")
		s.WriteString(emptyState)
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderDatabaseDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󱆟  Database Details")

	s.WriteString(titleStyle + "\n\n")

	databases := []struct {
		name     string
		type_    string
		port     int
		status   string
		adminURL string
	}{
		{"MySQL", "mysql", 3306, "running", "http://localhost:8080"},
		{"PostgreSQL", "postgres", 5432, "running", "http://localhost:8084"},
		{"MariaDB", "mariadb", 3307, "stopped", "http://localhost:8080"},
		{"MongoDB", "mongodb", 27017, "running", "http://localhost:8082"},
		{"Redis", "redis", 6379, "running", "http://localhost:8083"},
		{"Elasticsearch", "elasticsearch", 9200, "stopped", ""},
	}

	if m.cursor < len(databases) {
		db := databases[m.cursor]

		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(1, 2).
			Width(panelWidth - 8).
			Background(surfaceBg)

		var info strings.Builder

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Name: "))
		info.WriteString(lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(db.name) + "\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Type: "))
		badge := getServiceBadge(db.type_).Render(strings.ToUpper(db.type_))
		info.WriteString(badge + "\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Port: "))
		info.WriteString(lipgloss.NewStyle().Foreground(warningColor).Bold(true).Render(fmt.Sprintf(":%d", db.port)) + "\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Status: "))
		if db.status == "running" {
			info.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(successColor).
				Padding(0, 1).
				Bold(true).
				Render(" 󰀄 Running") + "\n")
		} else {
			info.WriteString(lipgloss.NewStyle().
				Foreground(bgColor).
				Background(mutedColor).
				Padding(0, 1).
				Render(" 󰀊 Stopped") + "\n")
		}

		if db.adminURL != "" {
			info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Admin: "))
			info.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(db.adminURL) + "\n")
		}

		s.WriteString(infoBox.Render(info.String()))
	} else {
		emptyState := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("󰉋  No database selected")
		s.WriteString(emptyState)
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderRuntimeDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰌠  Runtime Details")

	s.WriteString(titleStyle + "\n\n")

	runtimes := []struct {
		name    string
		version string
	}{
		{"PHP", m.config.Runtimes.PHP},
		{"Node.js", m.config.Runtimes.Node},
		{"Python", m.config.Runtimes.Python},
		{"Rust", m.config.Runtimes.Rust},
		{"Bun", m.config.Runtimes.Bun},
		{"Deno", m.config.Runtimes.Deno},
		{"Go", m.config.Runtimes.Go},
	}

	if m.cursor < len(runtimes) {
		rt := runtimes[m.cursor]

		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(1, 2).
			Width(panelWidth - 8).
			Background(surfaceBg)

		var info strings.Builder

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Runtime: "))
		info.WriteString(lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(rt.name) + "\n\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Current Version: "))
		info.WriteString(lipgloss.NewStyle().Foreground(infoColor).Bold(true).
			Background(surface0).Padding(0, 1).Render("v"+rt.version) + "\n\n")

		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Italic(true).
			Render("Press 'v' to change version"))

		s.WriteString(infoBox.Render(info.String()))
	} else {
		emptyState := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("󰉋  No runtime selected")
		s.WriteString(emptyState)
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderLogsDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰌱  Log Statistics")

	s.WriteString(titleStyle + "\n\n")

	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(surface1).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var stats strings.Builder

	successCount := 0
	errorCount := 0
	warningCount := 0
	infoCount := 0

	for _, log := range m.logs {
		switch log.level {
		case "success":
			successCount++
		case "error":
			errorCount++
		case "warning":
			warningCount++
		default:
			infoCount++
		}
	}

	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Total Logs: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Bold(true).Render(fmt.Sprintf("%d", len(m.logs))) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(successColor).Padding(0, 1).Render(" 󰸞 Success: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", successCount)) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(errorColor).Padding(0, 1).Render(" 󰚌 Errors: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", errorCount)) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(warningColor).Padding(0, 1).Render(" 󰀦 Warnings: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", warningCount)) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(infoColor).Padding(0, 1).Render(" 󰋽 Info: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", infoCount)))

	s.WriteString(statsBox.Render(stats.String()))

	return style.Render(s.String())
}

func (m model) renderTasksDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰘦  Task Statistics")

	s.WriteString(titleStyle + "\n\n")

	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(surface1).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var stats strings.Builder

	runningCount := 0
	completedCount := 0
	failedCount := 0

	for _, task := range m.backgroundTasks {
		switch task.status {
		case "running":
			runningCount++
		case "completed":
			completedCount++
		case "failed":
			failedCount++
		}
	}

	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Total Tasks: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Bold(true).Render(fmt.Sprintf("%d", len(m.backgroundTasks))) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(infoColor).Padding(0, 1).Render(" ⟳ Running: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", runningCount)) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(successColor).Padding(0, 1).Render(" 󰸞 Completed: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", completedCount)) + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(errorColor).Padding(0, 1).Render(" 󰚌 Failed: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf(" %d", failedCount)))

	s.WriteString(statsBox.Render(stats.String()))

	return style.Render(s.String())
}

func (m model) renderSettingsDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == detailPanel {
		borderStyle = lipgloss.DoubleBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == detailPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("󰒓  System Info")

	s.WriteString(titleStyle + "\n\n")

	aboutBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var about strings.Builder

	about.WriteString(lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render("󰘦  Lumine") + "\n")
	about.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Docker Development Manager") + "\n\n")
	about.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Version: "))
	about.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Background(surface0).
		Padding(0, 1).
		Render("v1.0.0") + "\n")
	about.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Theme: "))
	about.WriteString(lipgloss.NewStyle().
		Foreground(secondaryColor).
		Background(surface0).
		Padding(0, 1).
		Render("Catppuccin Mocha") + "\n")

	s.WriteString(aboutBox.Render(about.String()) + "\n\n")

	dockerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(infoColor).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var dockerInfo strings.Builder

	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true).
		Render("󰡨  Docker") + "\n\n")

	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(bgColor).
		Background(successColor).
		Padding(0, 1).
		Render(" 󰀄  Connected") + "\n\n")

	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Config: "))
	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(fgColor).
		Render(config.ConfigFile) + "\n")

	s.WriteString(dockerBox.Render(dockerInfo.String()) + "\n\n")

	runningCount := 0
	for _, status := range m.serviceStatus {
		if status.running {
			runningCount++
		}
	}

	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 2).
		Width(panelWidth - 8).
		Background(surfaceBg)

	var stats strings.Builder

	stats.WriteString(lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		Render("󰡨  Statistics") + "\n\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Services: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Bold(true).Render(fmt.Sprintf("%d", len(m.config.Services))) + "\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Running: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(bgColor).
		Background(successColor).Padding(0, 1).
		Render(fmt.Sprintf(" %d ", runningCount)) + "\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Projects: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Bold(true).Render(fmt.Sprintf("%d", len(m.config.Projects))) + "\n")

	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Logs: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(fmt.Sprintf("%d", len(m.logs))) + "\n")

	s.WriteString(statsBox.Render(stats.String()))

	return style.Render(s.String())
}

func (m model) renderStatusBar() string {
	var parts []string

	viewName := ""
	switch m.currentView {
	case servicesView:
		viewName = "Services"
	case projectsView:
		viewName = "Projects"
	case databasesView:
		viewName = "Databases"
	case runtimesView:
		viewName = "Runtimes"
	case logsView:
		viewName = "Logs"
	case backgroundTasksView:
		viewName = "Tasks"
	case settingsView:
		viewName = "Settings"
	}

	if viewName != "" {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Bold(true).
			Padding(0, 2).
			Render(" 󰘦 "+viewName+" "))
	}

	selectedCount := 0
	for _, selected := range m.selected {
		if selected {
			selectedCount++
		}
	}

	if selectedCount > 0 {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(bgColor).
			Background(warningColor).
			Padding(0, 2).
			Render(fmt.Sprintf(" 󰡖 %d selected ", selectedCount)))
	}

	runningCount := 0
	for _, status := range m.serviceStatus {
		if status.running {
			runningCount++
		}
	}

	statusText := fmt.Sprintf("󰀄  %d/%d running", runningCount, len(m.config.Services))
	statusStyle := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(successColor).
		Padding(0, 2).
		Bold(true)

	if runningCount == 0 {
		statusStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(mutedColor).
			Padding(0, 2)
	}

	parts = append(parts, statusStyle.Render(statusText))

	if m.statusMessage != "" {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(bgColor).
			Background(infoColor).
			Padding(0, 2).
			Render(" 󰋽 "+m.statusMessage+" "))
	}

	statusContent := strings.Join(parts, " ")

	currentWidth := lipgloss.Width(statusContent)
	remainingWidth := m.width - currentWidth
	if remainingWidth > 0 {
		statusContent += strings.Repeat(" ", remainingWidth)
	}

	return lipgloss.NewStyle().
		Background(surface0).
		Foreground(fgColor).
		Width(m.width).
		Render(statusContent)
}

func (m model) renderHelp() string {
	var sections []string

	if m.activePanel == mainPanel && m.currentView == servicesView {
		sections = []string{
			helpKeyStyle.Render(" ↑↓ ") + helpDescStyle.Render("navigate"),
			helpKeyStyle.Render(" space ") + helpDescStyle.Render("select"),
			helpKeyStyle.Render(" s ") + helpDescStyle.Render("start"),
			helpKeyStyle.Render(" x ") + helpDescStyle.Render("stop"),
			helpKeyStyle.Render(" r ") + helpDescStyle.Render("restart"),
			helpKeyStyle.Render(" v ") + helpDescStyle.Render("version"),
			helpKeyStyle.Render(" n ") + helpDescStyle.Render("new"),
			helpKeyStyle.Render(" q ") + helpDescStyle.Render("quit"),
		}
	} else if m.activePanel == mainPanel && m.currentView == projectsView {
		sections = []string{
			helpKeyStyle.Render(" ↑↓ ") + helpDescStyle.Render("navigate"),
			helpKeyStyle.Render(" n ") + helpDescStyle.Render("new project"),
			helpKeyStyle.Render(" , ") + helpDescStyle.Render("settings"),
			helpKeyStyle.Render(" q ") + helpDescStyle.Render("quit"),
		}
	} else if m.activePanel == mainPanel && m.currentView == runtimesView {
		sections = []string{
			helpKeyStyle.Render(" ↑↓ ") + helpDescStyle.Render("navigate"),
			helpKeyStyle.Render(" v ") + helpDescStyle.Render("change version"),
			helpKeyStyle.Render(" q ") + helpDescStyle.Render("quit"),
		}
	} else if m.activePanel == sidebarPanel {
		sections = []string{
			helpKeyStyle.Render(" ↑↓ ") + helpDescStyle.Render("navigate"),
			helpKeyStyle.Render(" enter ") + helpDescStyle.Render("select"),
			helpKeyStyle.Render(" l ") + helpDescStyle.Render("main panel"),
			helpKeyStyle.Render(" q ") + helpDescStyle.Render("quit"),
		}
	} else {
		sections = []string{
			helpKeyStyle.Render(" h/l ") + helpDescStyle.Render("switch panel"),
			helpKeyStyle.Render(" tab ") + helpDescStyle.Render("next panel"),
			helpKeyStyle.Render(" , ") + helpDescStyle.Render("settings"),
			helpKeyStyle.Render(" q ") + helpDescStyle.Render("quit"),
		}
	}

	helpText := strings.Join(sections, " "+lipgloss.NewStyle().Foreground(surface1).Render("│")+" ")

	currentWidth := lipgloss.Width(helpText)
	remainingWidth := m.width - currentWidth
	if remainingWidth > 0 {
		helpText += strings.Repeat(" ", remainingWidth)
	}

	return lipgloss.NewStyle().
		Background(surfaceBg).
		Foreground(fgColor).
		Width(m.width).
		Render(helpText)
}
