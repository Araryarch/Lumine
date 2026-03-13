package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderDatabasePanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("Databases")
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

	for i, db := range databases {
		var line strings.Builder

		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		// Status
		statusIcon := "*"
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
			Foreground(mutedColor).
			Render(fmt.Sprintf(" :%d", db.port))
		line.WriteString(portInfo)

		lineStr := line.String()
		if m.cursor == i {
			lineStr = selectedItemStyle.Width(panelWidth - 4).Render(lineStr)
		} else {
			lineStr = normalItemStyle.Render(lineStr)
		}

		s.WriteString(lineStr + "\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Render("Press 'o' to open admin panel"))

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}
