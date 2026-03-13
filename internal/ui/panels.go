package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"lumine/internal/config"
)

func (m model) renderSidebarFixed(width, height int) string {
	var s strings.Builder

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(bgColor).
		Padding(1, 1).
		Width(width).
		Height(height)
	
	if m.activePanel == sidebarPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	// Sidebar header
	sidebarHeader := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Navigation")
	s.WriteString(sidebarHeader + "\n\n")

	// Menu items without icons
	menuItems := []string{
		"Services",
		"Projects",
		"Databases",
		"Runtimes",
		"Logs",
		"Tasks",
		"Settings",
		"Refresh",
		"Quit",
	}

	for i, item := range menuItems {
		var line strings.Builder

		// Cursor and selection
		if m.sidebarCursor == i {
			cursor := lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("> ")
			line.WriteString(cursor)
		} else {
			line.WriteString("  ")
		}

		// Name
		name := item
		if m.sidebarCursor == i {
			name = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(name)
		} else {
			name = lipgloss.NewStyle().
				Foreground(fgColor).
				Render(name)
		}
		line.WriteString(name)

		lineStr := line.String()
		
		// Highlight selected
		if m.sidebarCursor == i {
			lineStr = lipgloss.NewStyle().
				Background(surfaceColor).
				Width(width - 4).
				Padding(0, 1).
				Render(lineStr)
		}

		s.WriteString(lineStr + "\n")
	}

	// Fill remaining space
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < height-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

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

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == mainPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Services")
	s.WriteString(header + "\n\n")
	
	if len(m.config.Services) == 0 {
		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(2, 4).
			Width(panelWidth - 8).
			Align(lipgloss.Center)
		
		var empty strings.Builder
		empty.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No services configured yet") + "\n\n")
		empty.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Render("Press ") +
			lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("'n'") +
			lipgloss.NewStyle().Foreground(infoColor).Render(" to add a service"))
		
		s.WriteString(emptyBox.Render(empty.String()))
	} else {
		// Calculate visible items with scroll
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

		// Scroll indicator top
		if showScrollTop {
			s.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n")
		}

		for i := startIdx; i < endIdx; i++ {
			service := m.config.Services[i]
			var line strings.Builder

			// Cursor
			cursor := "  "
			if m.cursor == i {
				cursor = lipgloss.NewStyle().
					Foreground(primaryColor).
					Bold(true).
					Render("▶ ")
			}

			// Checkbox
			checkbox := "☐ "
			if m.selected[i] {
				checkbox = lipgloss.NewStyle().
					Foreground(successColor).
					Render("☑ ")
			}

			// Status indicator with icon
			status := m.serviceStatus[service.Name]
			statusIcon := "●"
			statusStyle := stoppedStatusStyle
			statusText := "stopped"

			if status.running {
				statusIcon = "●"
				statusStyle = runningStatusStyle
				statusText = "running"
			}

			// Service badge
			badge := getServiceBadge(service.Type).Render(strings.ToUpper(service.Type))

			// Build line
			line.WriteString(cursor)
			line.WriteString(checkbox)
			line.WriteString(badge + " ")

			serviceName := lipgloss.NewStyle().Bold(true).Render(service.Name)
			line.WriteString(serviceName)

			// Status
			statusStr := fmt.Sprintf(" %s %s", statusStyle.Render(statusIcon), statusStyle.Render(statusText))
			line.WriteString(statusStr)

			// Port info
			portInfo := lipgloss.NewStyle().
				Foreground(warningColor).
				Render(fmt.Sprintf(" :%d", service.Port))
			line.WriteString(portInfo)

			// Version
			versionInfo := lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(fmt.Sprintf(" v%s", service.Version))
			line.WriteString(versionInfo)

			lineStr := line.String()
			if m.cursor == i {
				lineStr = lipgloss.NewStyle().
					Background(surfaceColor).
					Width(panelWidth - 4).
					Padding(0, 1).
					Render(lineStr)
			} else {
				lineStr = normalItemStyle.Render(lineStr)
			}

			s.WriteString(lineStr + "\n")
		}

		// Scroll indicator bottom
		if showScrollBottom {
			s.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Render(fmt.Sprintf("  ↓ %d more below", len(m.config.Services)-endIdx)) + "\n")
		}
	}

	// Add some spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderAddServicePanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == mainPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Add New Service")
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

	subHeader := lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true).
		Render("Select a service type to add:")
	s.WriteString(subHeader + "\n\n")

	// Calculate how many items can fit
	maxVisibleItems := panelHeight - 8
	if maxVisibleItems < 5 {
		maxVisibleItems = 5
	}

	// Calculate scroll offset
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

	// Show scroll indicator at top
	if showScrollTop {
		scrollInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↑ %d more above", startIdx))
		s.WriteString(scrollInfo + "\n")
	} else {
		s.WriteString("\n") // Empty line for consistent spacing
	}

	for i := startIdx; i < endIdx; i++ {
		svc := availableServices[i]
		var line strings.Builder
		
		// Cursor
		cursor := "  "
		if m.cursor == i {
			cursor = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("> ")
		}
		
		line.WriteString(cursor)
		
		badge := getServiceBadge(strings.ToLower(svc.name)).Render(svc.name)
		line.WriteString(badge + " ")
		
		desc := lipgloss.NewStyle().
			Foreground(fgColor).
			Render(svc.desc)
		line.WriteString(desc)
		
		lineStr := line.String()
		
		if m.cursor == i {
			lineStr = lipgloss.NewStyle().
				Background(surfaceColor).
				Foreground(fgColor).
				Width(panelWidth - 4).
				Padding(0, 1).
				Render(lineStr)
		} else {
			lineStr = normalItemStyle.Render(lineStr)
		}
		
		s.WriteString(lineStr + "\n")
	}

	// Show scroll indicator at bottom
	if showScrollBottom {
		scrollInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↓ %d more below", len(availableServices)-endIdx))
		s.WriteString(scrollInfo + "\n")
	} else {
		s.WriteString("\n") // Empty line for consistent spacing
	}

	// Add spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderSettingsPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == mainPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("⚙️  Settings")
	s.WriteString(header + "\n\n")

	// Web Server Configuration
	serverBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
	var serverContent strings.Builder
	
	serverContent.WriteString(lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Render("🌐 Web Server") + "\n\n")
	
	// Server options
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
				Render("● "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(fgColor).
				Bold(true).
				Render(srv.name))
			line.WriteString(lipgloss.NewStyle().
				Foreground(successColor).
				Background(surfaceColor).
				Padding(0, 1).
				Render(" ACTIVE"))
		} else {
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("○ "))
			line.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(srv.name))
		}
		serverContent.WriteString(line.String() + "\n")
	}
	
	s.WriteString(serverBox.Render(serverContent.String()) + "\n\n")

	// Runtime Versions
	runtimeBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(infoColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
	var runtimeContent strings.Builder
	
	runtimeContent.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true).
		Render("🔧 Runtime Versions") + "\n\n")
	
	// Runtime versions
	runtimes := []struct {
		name    string
		version string
		icon    string
	}{
		{"PHP", m.config.Runtimes.PHP, "🐘"},
		{"Node.js", m.config.Runtimes.Node, "⬢"},
		{"Python", m.config.Runtimes.Python, "🐍"},
		{"Go", m.config.Runtimes.Go, "🔷"},
		{"Rust", m.config.Runtimes.Rust, "🦀"},
		{"Bun", m.config.Runtimes.Bun, "🥟"},
		{"Deno", m.config.Runtimes.Deno, "🦕"},
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
			Render("v" + rt.version))
		runtimeContent.WriteString(line.String() + "\n")
	}
	
	s.WriteString(runtimeBox.Render(runtimeContent.String()) + "\n\n")

	// Help text
	helpBox := lipgloss.NewStyle().
		Foreground(mutedColor).
		Italic(true).
		Width(panelWidth - 8).
		Align(lipgloss.Center)
	
	helpText := "Press 'v' on Runtimes page to change versions"
	s.WriteString(helpBox.Render(helpText))

	// Add spacing
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

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	// Service details section
	if m.cursor < len(m.config.Services) {
		service := m.config.Services[m.cursor]
		status := m.serviceStatus[service.Name]

		header := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Underline(true).
			Render("Service Details")
		s.WriteString(header + "\n\n")

		// Info box
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(panelWidth - 8)
		
		var info strings.Builder
		
		// Service name
		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Name: "))
		info.WriteString(lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Render(service.Name) + "\n")

		// Type with badge
		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Type: "))
		badge := getServiceBadge(service.Type).Render(strings.ToUpper(service.Type))
		info.WriteString(badge + "\n")

		// Version
		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Version: "))
		info.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Render(service.Version) + "\n")

		// Port
		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Port: "))
		info.WriteString(lipgloss.NewStyle().
			Foreground(warningColor).
			Bold(true).
			Render(fmt.Sprintf(":%d", service.Port)) + "\n")

		// Status
		info.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Status: "))
		if status.running {
			info.WriteString(lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true).
				Render("● Running") + "\n")
		} else {
			info.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("● Stopped") + "\n")
		}

		// Container ID
		if status.containerId != "" {
			info.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render("Container: "))
			info.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(status.containerId[:12]) + "\n")
		}
		
		s.WriteString(infoBox.Render(info.String()) + "\n\n")

		// Environment variables
		if len(service.Env) > 0 {
			envHeader := lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				Render("Environment Variables:")
			s.WriteString(envHeader + "\n\n")
			
			envBox := lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderColor).
				Padding(1, 2).
				Width(panelWidth - 8)
			
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
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No service selected") + "\n")
	}

	// Fill remaining space
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderProjectDetailPanel() string {
	var s strings.Builder
	
	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Project Details")
	s.WriteString(header + "\n\n")

	if m.cursor < len(m.config.Projects) {
		project := m.config.Projects[m.cursor]
		
		// Info box
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(panelWidth - 8)
		
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
			info.WriteString(lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("Running") + "\n")
		} else {
			info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Stopped") + "\n")
		}
		
		s.WriteString(infoBox.Render(info.String()))
	} else {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No project selected"))
	}

	// Fill remaining space
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderDatabaseDetailPanel() string {
	var s strings.Builder
	
	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Database Details")
	s.WriteString(header + "\n\n")

	// Hardcoded databases for now
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
		
		// Info box
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(panelWidth - 8)
		
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
			info.WriteString(lipgloss.NewStyle().Foreground(successColor).Bold(true).Render("Running") + "\n")
		} else {
			info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Stopped") + "\n")
		}
		
		if db.adminURL != "" {
			info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Admin: "))
			info.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(db.adminURL) + "\n")
		}
		
		s.WriteString(infoBox.Render(info.String()))
	} else {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No database selected"))
	}

	// Fill remaining space
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderRuntimeDetailPanel() string {
	var s strings.Builder
	
	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Runtime Details")
	s.WriteString(header + "\n\n")

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
		
		// Info box
		infoBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Width(panelWidth - 8)
		
		var info strings.Builder
		
		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Runtime: "))
		info.WriteString(lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render(rt.name) + "\n\n")
		
		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Current Version: "))
		info.WriteString(lipgloss.NewStyle().Foreground(infoColor).Bold(true).Render(rt.version) + "\n\n")
		
		info.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Italic(true).Render("Press 'v' to change version"))
		
		s.WriteString(infoBox.Render(info.String()))
	} else {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No runtime selected"))
	}

	// Fill remaining space
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderLogsDetailPanel() string {
	var s strings.Builder
	
	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("📊 Log Statistics")
	s.WriteString(header + "\n\n")
	
	// Stats
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
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
	
	stats.WriteString(lipgloss.NewStyle().Foreground(successColor).Render("✓ Success: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", successCount)) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(errorColor).Render("✗ Errors: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", errorCount)) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(warningColor).Render("⚠ Warnings: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", warningCount)) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render("ℹ Info: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", infoCount)) + "\n")
	
	s.WriteString(statsBox.Render(stats.String()))
	
	return style.Render(s.String())
}

func (m model) renderTasksDetailPanel() string {
	var s strings.Builder
	
	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("📊 Task Statistics")
	s.WriteString(header + "\n\n")
	
	// Stats
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
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
	
	stats.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render("⟳ Running: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", runningCount)) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(successColor).Render("✓ Completed: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", completedCount)) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(errorColor).Render("✗ Failed: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(fmt.Sprintf("%d", failedCount)) + "\n")
	
	s.WriteString(statsBox.Render(stats.String()))
	
	return style.Render(s.String())
}

func (m model) renderSettingsDetailPanel() string {
	var s strings.Builder
	
	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == detailPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("ℹ️  System Info")
	s.WriteString(header + "\n\n")
	
	// About box
	aboutBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
	var about strings.Builder
	
	about.WriteString(lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render("Lumine") + "\n")
	about.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Docker Development Manager") + "\n\n")
	about.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Version: "))
	about.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Render("1.0.0") + "\n")
	about.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Theme: "))
	about.WriteString(lipgloss.NewStyle().
		Foreground(secondaryColor).
		Render("Catppuccin Mocha") + "\n")
	
	s.WriteString(aboutBox.Render(about.String()) + "\n\n")

	// Docker info
	dockerBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(infoColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
	var dockerInfo strings.Builder
	
	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Bold(true).
		Render("🐳 Docker") + "\n\n")
	
	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(successColor).
		Render("● Connected") + "\n\n")
	
	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Config: "))
	dockerInfo.WriteString(lipgloss.NewStyle().
		Foreground(fgColor).
		Render(config.ConfigFile) + "\n")
	
	s.WriteString(dockerBox.Render(dockerInfo.String()) + "\n\n")

	// Statistics
	statsBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(successColor).
		Padding(1, 2).
		Width(panelWidth - 8)
	
	var stats strings.Builder
	
	stats.WriteString(lipgloss.NewStyle().
		Foreground(successColor).
		Bold(true).
		Render("📊 Statistics") + "\n\n")
	
	runningCount := 0
	for _, status := range m.serviceStatus {
		if status.running {
			runningCount++
		}
	}
	
	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Services: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Bold(true).Render(fmt.Sprintf("%d", len(m.config.Services))) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Running: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(successColor).Bold(true).Render(fmt.Sprintf("%d", runningCount)) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Projects: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(fgColor).Bold(true).Render(fmt.Sprintf("%d", len(m.config.Projects))) + "\n")
	
	stats.WriteString(lipgloss.NewStyle().Foreground(mutedColor).Render("Logs: "))
	stats.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(fmt.Sprintf("%d", len(m.logs))) + "\n")
	
	s.WriteString(statsBox.Render(stats.String()))
	
	return style.Render(s.String())
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

	// Logs section - show recent activity only
	logsHeader := headerStyle.Render("📝 Recent Activity")
	s.WriteString(logsHeader + "\n\n")

	// Show last 5 logs
	startIdx := 0
	maxLogs := 5
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
			log := m.logs[i]
			
			var icon string
			var iconStyle lipgloss.Style
			switch log.level {
			case "success":
				icon = "✓"
				iconStyle = lipgloss.NewStyle().Foreground(successColor)
			case "error":
				icon = "✗"
				iconStyle = lipgloss.NewStyle().Foreground(errorColor)
			case "warning":
				icon = "⚠"
				iconStyle = lipgloss.NewStyle().Foreground(warningColor)
			default:
				icon = "•"
				iconStyle = lipgloss.NewStyle().Foreground(infoColor)
			}
			
			logLine := fmt.Sprintf("%s %s", iconStyle.Render(icon), log.message)
			if lipgloss.Width(logLine) > panelWidth-6 {
				logLine = logLine[:panelWidth-9] + "..."
			}
			s.WriteString(logLine + "\n")
		}
		
		s.WriteString("\n")
		s.WriteString(subHeaderStyle.Render("View all logs in Logs panel"))
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderStatusBar() string {
	var parts []string

	// Current view indicator
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
		viewName = "Background Tasks"
	case settingsView:
		viewName = "Settings"
	}
	
	if viewName != "" {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Bold(true).
			Padding(0, 1).
			Render("▶ "+viewName))
	}

	// Selected count
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
			Padding(0, 1).
			Render(fmt.Sprintf("✓ %d selected", selectedCount)))
	}

	// Running services count
	runningCount := 0
	for _, status := range m.serviceStatus {
		if status.running {
			runningCount++
		}
	}
	
	statusText := fmt.Sprintf("● %d/%d running", runningCount, len(m.config.Services))
	statusStyle := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(successColor).
		Padding(0, 1)
	
	if runningCount == 0 {
		statusStyle = statusStyle.Background(mutedColor)
	}
	
	parts = append(parts, statusStyle.Render(statusText))

	// Status message
	if m.statusMessage != "" {
		parts = append(parts, lipgloss.NewStyle().
			Foreground(bgColor).
			Background(infoColor).
			Padding(0, 1).
			Render("ℹ "+m.statusMessage))
	}

	// Join all parts with spacing
	statusContent := strings.Join(parts, "  ")
	
	// Pad to full width
	padding := m.width - lipgloss.Width(statusContent) - 4
	if padding > 0 {
		statusContent += strings.Repeat(" ", padding)
	}
	
	return lipgloss.NewStyle().
		Foreground(fgColor).
		Background(surfaceColor).
		Padding(0, 1).
		Width(m.width - 2).
		Render(statusContent)
}

func (m model) renderHelp() string {
	var sections []string

	if m.activePanel == mainPanel && m.currentView == servicesView {
		sections = []string{
			helpKeyStyle.Render("↑/↓") + helpDescStyle.Render(" navigate"),
			helpKeyStyle.Render("space") + helpDescStyle.Render(" select"),
			helpKeyStyle.Render("enter") + helpDescStyle.Render(" start"),
			helpKeyStyle.Render("s") + helpDescStyle.Render(" start"),
			helpKeyStyle.Render("x") + helpDescStyle.Render(" stop"),
			helpKeyStyle.Render("r") + helpDescStyle.Render(" restart"),
			helpKeyStyle.Render("v") + helpDescStyle.Render(" version"),
			helpKeyStyle.Render("c") + helpDescStyle.Render(" cleanup"),
			helpKeyStyle.Render("h/l") + helpDescStyle.Render(" panels"),
			helpKeyStyle.Render("q") + helpDescStyle.Render(" quit"),
		}
	} else if m.activePanel == sidebarPanel {
		sections = []string{
			helpKeyStyle.Render("↑/↓") + helpDescStyle.Render(" navigate"),
			helpKeyStyle.Render("enter") + helpDescStyle.Render(" select"),
			helpKeyStyle.Render("l") + helpDescStyle.Render(" main panel"),
			helpKeyStyle.Render("q") + helpDescStyle.Render(" quit"),
		}
	} else {
		sections = []string{
			helpKeyStyle.Render("h/l") + helpDescStyle.Render(" switch panel"),
			helpKeyStyle.Render("tab") + helpDescStyle.Render(" next panel"),
			helpKeyStyle.Render("q") + helpDescStyle.Render(" quit"),
		}
	}

	helpText := strings.Join(sections, " " + dividerStyle.Render("│") + " ")
	return helpStyle.Render(helpText)
}
