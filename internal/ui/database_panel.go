package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderDatabasePanel() string {
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
		Render("Databases")
	s.WriteString(header + "\n\n")

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

	// Calculate visible items with scroll
	maxVisibleItems := panelHeight - 6
	if maxVisibleItems < 3 {
		maxVisibleItems = 3
	}

	startIdx := 0
	endIdx := len(databases)
	showScrollTop := false
	showScrollBottom := false
	
	if len(databases) > maxVisibleItems {
		startIdx = m.cursor - (maxVisibleItems / 2)
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + maxVisibleItems
		if endIdx > len(databases) {
			endIdx = len(databases)
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}
		
		showScrollTop = startIdx > 0
		showScrollBottom = endIdx < len(databases)
	}

	// Scroll indicator top
	if showScrollTop {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		db := databases[i]
		var line strings.Builder

		cursor := "  "
		if m.cursor == i {
			cursor = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("▶ ")
		}

		// Status
		statusIcon := "●"
		statusStyle := stoppedStatusStyle
		if db.status == "running" {
			statusStyle = runningStatusStyle
		}

		// Database badge
		badge := getServiceBadge(db.type_).Render(strings.ToUpper(db.type_))

		line.WriteString(cursor)
		line.WriteString(badge + " ")
		line.WriteString(lipgloss.NewStyle().Bold(true).Render(db.name))
		line.WriteString(fmt.Sprintf(" %s %s", statusStyle.Render(statusIcon), statusStyle.Render(db.status)))

		// Port
		portInfo := lipgloss.NewStyle().
			Foreground(warningColor).
			Render(fmt.Sprintf(" :%d", db.port))
		line.WriteString(portInfo)

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
			Render(fmt.Sprintf("  ↓ %d more below", len(databases)-endIdx)) + "\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Press ") +
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("'o'") +
		lipgloss.NewStyle().Foreground(mutedColor).Render(" to open admin panel"))

	// Add spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}
