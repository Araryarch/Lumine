package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderDatabasePanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == mainPanel {
		borderStyle = lipgloss.ThickBorder()
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
		Render("󱆟  Databases")

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

	if showScrollTop {
		scrollIndicator := lipgloss.NewStyle().
			Foreground(fgMuted).
			Italic(true).
			Render(fmt.Sprintf("  ↑ %d more", startIdx))
		s.WriteString(scrollIndicator + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		db := databases[i]
		var line strings.Builder

		icon := getIconForService(db.type_)

		// Cursor indicator with arrow
		if m.cursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("▶ "))
		} else {
			line.WriteString("  ")
		}

		// Status indicator
		statusIcon := "●"
		statusStyle := fgMuted
		if db.status == "running" {
			statusIcon = "●"
			statusStyle = successColor
		}

		badge := getServiceBadge(db.type_).Render(strings.ToUpper(db.type_))
		line.WriteString(badge + " ")

		line.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(icon + " "))
		line.WriteString(lipgloss.NewStyle().Bold(true).Foreground(fgColor).Render(db.name))

		statusBadge := lipgloss.NewStyle().
			Foreground(statusStyle).
			Bold(true).
			Render(" " + statusIcon)
		line.WriteString(statusBadge)

		portBadge := lipgloss.NewStyle().
			Foreground(warningColor).
			Background(surface0).
			Padding(0, 1).
			Render(fmt.Sprintf(":%d", db.port))
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
			Foreground(fgMuted).
			Italic(true).
			Render(fmt.Sprintf("  ↓ %d more", len(databases)-endIdx))
		s.WriteString(scrollIndicator + "\n")
	}

	s.WriteString("\n")

	helpText := lipgloss.NewStyle().
		Foreground(fgMuted).
		Render("Press ") +
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("'o'") +
		lipgloss.NewStyle().Foreground(fgMuted).Render(" to open admin panel")
	s.WriteString(helpText)

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}
