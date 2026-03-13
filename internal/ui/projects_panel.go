package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderProjectsPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("Projects")
	s.WriteString(header + "\n\n")

	if len(m.config.Projects) == 0 {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No projects yet. Press 'n' to create one.") + "\n")
	} else {
		for i, project := range m.config.Projects {
			var line strings.Builder

			cursor := "  "
			if m.cursor == i {
				cursor = "> "
			}

			// Status
			statusIcon := "*"
			statusStyle := stoppedStatusStyle
			if project.Status == "running" {
				statusStyle = runningStatusStyle
			}

			// Project type badge
			badge := getProjectBadge(project.Type).Render(strings.ToUpper(project.Type))

			line.WriteString(cursor)
			line.WriteString(badge + " ")
			line.WriteString(lipgloss.NewStyle().Bold(true).Render(project.Name))
			line.WriteString(fmt.Sprintf(" %s %s", statusStyle.Render(statusIcon), statusStyle.Render(project.Status)))

			// Domain
			if project.Domain != "" {
				domainInfo := lipgloss.NewStyle().
					Foreground(infoColor).
					Render(fmt.Sprintf(" → %s", project.Domain))
				line.WriteString(domainInfo)
			}

			lineStr := line.String()
			if m.cursor == i {
				lineStr = selectedItemStyle.Width(panelWidth - 4).Render(lineStr)
			} else {
				lineStr = normalItemStyle.Render(lineStr)
			}

			s.WriteString(lineStr + "\n")
		}
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderRuntimesPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("Runtimes")
	s.WriteString(header + "\n\n")

	runtimes := []struct {
		name    string
		version string
		icon    string
	}{
		{"PHP", m.config.Runtimes.PHP, "PHP"},
		{"Node.js", m.config.Runtimes.Node, "NODE"},
		{"Python", m.config.Runtimes.Python, "PY"},
		{"Bun", m.config.Runtimes.Bun, "BUN"},
		{"Deno", m.config.Runtimes.Deno, "DENO"},
		{"Go", m.config.Runtimes.Go, "GO"},
		{"Rust", m.config.Runtimes.Rust, "RUST"},
	}

	for i, rt := range runtimes {
		var line strings.Builder

		cursor := "  "
		if m.cursor == i {
			cursor = "> "
		}

		line.WriteString(cursor)
		line.WriteString(rt.icon + " ")
		line.WriteString(lipgloss.NewStyle().Bold(true).Render(rt.name))
		line.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Render(fmt.Sprintf(" v%s", rt.version)))

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
		Foreground(mutedColor).
		Render("Press 'v' to change version"))

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func getProjectBadge(projectType string) lipgloss.Style {
	switch projectType {
	case "laravel":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#FF2D20")).
			Padding(0, 1).
			Bold(true)
	case "nextjs":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	case "vue":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#42B883")).
			Padding(0, 1).
			Bold(true)
	case "django":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#092E20")).
			Padding(0, 1).
			Bold(true)
	case "express":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#000000")).
			Padding(0, 1).
			Bold(true)
	default:
		return defaultBadge
	}
}


	case "axum", "actix", "rocket":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#CE422B")).
			Padding(0, 1).
			Bold(true)
