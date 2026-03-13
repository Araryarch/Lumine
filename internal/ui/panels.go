package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"lumine/internal/config"
)

func (m model) renderSidebar() string {
	var s strings.Builder

	sidebarWidth := 20
	sidebarHeight := m.height - 6

	style := panelStyle
	if m.activePanel == sidebarPanel {
		style = activePanelStyle
	}

	for i, item := range m.sidebarItems {
		var line string
		if m.sidebarCursor == i {
			line = selectedItemStyle.Width(sidebarWidth - 4).Render("> " + item)
		} else {
			line = normalItemStyle.Render("  " + item)
		}
		s.WriteString(line + "\n")
	}

	// Add spacing
	for i := len(m.sidebarItems); i < sidebarHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Width(sidebarWidth).Height(sidebarHeight).Render(s.String())
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
	case addServiceView:
		return m.renderAddServicePanel()
	case addProjectView:
		return m.renderAddProjectPanel()
	case settingsView:
		return m.renderSettingsPanel()
	default:
		return m.renderServicesPanel()
	}
}

func (m model) renderServicesPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("Services")
	s.WriteString(header + "\n")

	for i, service := range m.config.Services {
		var line strings.Builder

		// Cursor
		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		// Checkbox
		checkbox := "[ ]"
		if m.selected[i] {
			checkbox = "[x]"
		}

		// Status indicator
		status := m.serviceStatus[service.Name]
		statusIcon := "*"
		statusStyle := stoppedStatusStyle
		statusText := "stopped"

		if status.running {
			statusStyle = runningStatusStyle
			statusText = "running"
		}

		// Service badge
		badge := getServiceBadge(service.Type).Render(strings.ToUpper(service.Type))

		// Build line
		line.WriteString(cursor)
		line.WriteString(checkbox + " ")
		line.WriteString(badge + " ")

		serviceName := lipgloss.NewStyle().Bold(true).Render(service.Name)
		line.WriteString(serviceName)

		// Status
		statusStr := fmt.Sprintf(" %s %s", statusStyle.Render(statusIcon), statusStyle.Render(statusText))
		line.WriteString(statusStr)

		// Port info
		portInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(fmt.Sprintf(" :%d", service.Port))
		line.WriteString(portInfo)

		// Version
		versionInfo := lipgloss.NewStyle().
			Foreground(infoColor).
			Render(fmt.Sprintf(" v%s", service.Version))
		line.WriteString(versionInfo)

		lineStr := line.String()
		if m.cursor == i {
			lineStr = selectedItemStyle.Width(panelWidth - 4).Render(lineStr)
		} else {
			lineStr = normalItemStyle.Render(lineStr)
		}

		s.WriteString(lineStr + "\n")
	}

	// Add some spacing
	for i := len(m.config.Services); i < panelHeight-5; i++ {
		s.WriteString("\n")
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderAddServicePanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("Add New Service")
	s.WriteString(header + "\n\n")

	availableServices := []struct {
		name string
		desc string
	}{
		{"PHP", "PHP-FPM or Apache with PHP"},
		{"MySQL", "MySQL Database Server"},
		{"MariaDB", "MariaDB Database Server"},
		{"PostgreSQL", "PostgreSQL Database"},
		{"Nginx", "Nginx Web Server"},
		{"Apache", "Apache HTTP Server"},
		{"Caddy", "Caddy Web Server"},
		{"Redis", "Redis Cache Server"},
		{"MongoDB", "MongoDB NoSQL Database"},
		{"phpMyAdmin", "MySQL Web Interface"},
		{"Adminer", "Database Management Tool"},
		{"Elasticsearch", "Search Engine"},
		{"RabbitMQ", "Message Queue"},
		{"Memcached", "Memory Cache"},
	}

	s.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Select a service type to add:\n\n"))

	for _, svc := range availableServices {
		badge := getServiceBadge(strings.ToLower(svc.name)).Render(svc.name)
		s.WriteString(fmt.Sprintf("  %s - %s\n", badge, svc.desc))
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render("Press 'n' to add a service (coming soon)"))

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderSettingsPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("Settings")
	s.WriteString(header + "\n\n")

	s.WriteString(lipgloss.NewStyle().Bold(true).Render("Configuration File:\n"))
	s.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render(config.ConfigFile + "\n\n"))

	s.WriteString(lipgloss.NewStyle().Bold(true).Render("Docker Status:\n"))
	s.WriteString(lipgloss.NewStyle().Foreground(successColor).Render("* Connected\n\n"))

	s.WriteString(lipgloss.NewStyle().Bold(true).Render("Total Services:\n"))
	s.WriteString(fmt.Sprintf("%d configured\n\n", len(m.config.Services)))

	runningCount := 0
	for _, status := range m.serviceStatus {
		if status.running {
			runningCount++
		}
	}
	s.WriteString(lipgloss.NewStyle().Bold(true).Render("Running Services:\n"))
	s.WriteString(fmt.Sprintf("%d / %d\n", runningCount, len(m.config.Services)))

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderDetailPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == detailPanel {
		style = activePanelStyle
	}

	// Service details section
	if m.cursor < len(m.config.Services) {
		service := m.config.Services[m.cursor]
		status := m.serviceStatus[service.Name]

		header := headerStyle.Render("📋 Service Details")
		s.WriteString(header + "\n\n")

		// Service name
		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Name: "))
		s.WriteString(service.Name + "\n")

		// Type
		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Type: "))
		badge := getServiceBadge(service.Type).Render(service.Type)
		s.WriteString(badge + "\n")

		// Version
		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Version: "))
		s.WriteString(service.Version + "\n")

		// Port
		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Port: "))
		s.WriteString(fmt.Sprintf("%d\n", service.Port))

		// Status
		s.WriteString(lipgloss.NewStyle().Bold(true).Render("Status: "))
		if status.running {
			s.WriteString(runningStatusStyle.Render("* Running") + "\n")
		} else {
			s.WriteString(stoppedStatusStyle.Render("* Stopped") + "\n")
		}

		// Container ID
		if status.containerId != "" {
			s.WriteString(lipgloss.NewStyle().Bold(true).Render("Container: "))
			s.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(status.containerId[:12]) + "\n")
		}

		// Environment variables
		if len(service.Env) > 0 {
			s.WriteString("\n")
			s.WriteString(lipgloss.NewStyle().Bold(true).Render("Environment:\n"))
			for k, v := range service.Env {
				s.WriteString(fmt.Sprintf("  %s=%s\n",
					lipgloss.NewStyle().Foreground(infoColor).Render(k),
					v))
			}
		}

		s.WriteString("\n")
		divider := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(strings.Repeat("─", panelWidth-6))
		s.WriteString(divider + "\n\n")
	}

	// Logs section
	logsHeader := headerStyle.Render("📝 Activity Logs")
	s.WriteString(logsHeader + "\n\n")

	// Show last logs
	startIdx := 0
	maxLogs := 8
	if len(m.logs) > maxLogs {
		startIdx = len(m.logs) - maxLogs
	}

	if len(m.logs) == 0 {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No activity yet...") + "\n")
	} else {
		for i := startIdx; i < len(m.logs); i++ {
			logLine := lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("• " + m.logs[i])
			s.WriteString(logLine + "\n")
		}
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderStatusBar() string {
	var parts []string

	// Selected count
	selectedCount := 0
	for _, selected := range m.selected {
		if selected {
			selectedCount++
		}
	}

	if selectedCount > 0 {
		parts = append(parts, fmt.Sprintf("Selected: %d", selectedCount))
	}

	// Running services count
	runningCount := 0
	for _, status := range m.serviceStatus {
		if status.running {
			runningCount++
		}
	}
	parts = append(parts, fmt.Sprintf("Running: %d/%d", runningCount, len(m.config.Services)))

	// Status message
	if m.statusMessage != "" {
		parts = append(parts, m.statusMessage)
	}

	statusText := strings.Join(parts, " • ")
	return statusBarStyle.Width(m.width - 2).Render(statusText)
}

func (m model) renderHelp() string {
	var helps []string

	if m.activePanel == mainPanel && m.currentView == servicesView {
		helps = []string{
			helpKeyStyle.Render("↑/↓,j/k") + " navigate",
			helpKeyStyle.Render("h/l") + " panels",
			helpKeyStyle.Render("space") + " select",
			helpKeyStyle.Render("s") + " start",
			helpKeyStyle.Render("x") + " stop",
			helpKeyStyle.Render("r") + " restart",
			helpKeyStyle.Render("v") + " change version",
			helpKeyStyle.Render("delete") + " remove",
			helpKeyStyle.Render("c") + " cleanup all",
			helpKeyStyle.Render("ctrl+c") + " quit",
		}
	} else if m.activePanel == sidebarPanel {
		helps = []string{
			helpKeyStyle.Render("↑/↓,j/k") + " navigate",
			helpKeyStyle.Render("enter") + " select",
			helpKeyStyle.Render("l") + " go to main",
			helpKeyStyle.Render("ctrl+c") + " quit",
		}
	} else {
		helps = []string{
			helpKeyStyle.Render("h/l") + " switch panel",
			helpKeyStyle.Render("tab") + " next panel",
			helpKeyStyle.Render("ctrl+c") + " quit",
		}
	}

	helpText := strings.Join(helps, " • ")
	return helpStyle.Render(helpText)
}
