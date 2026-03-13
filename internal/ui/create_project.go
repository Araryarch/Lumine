package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderAddProjectPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 6

	style := panelStyle
	if m.activePanel == mainPanel {
		style = activePanelStyle
	}

	header := headerStyle.Render("➕ Create New Project")
	s.WriteString(header + "\n\n")

	projectTypes := []struct {
		name    string
		desc    string
		runtime string
	}{
		{"Laravel", "PHP Framework", "PHP"},
		{"Next.js", "React Framework", "Node.js"},
		{"Vue", "Progressive Framework", "Node.js"},
		{"Django", "Python Framework", "Python"},
		{"Express", "Node.js Framework", "Node.js"},
		{"FastAPI", "Python API Framework", "Python"},
		{"Nuxt", "Vue Framework", "Node.js"},
		{"SvelteKit", "Svelte Framework", "Node.js"},
		{"Remix", "React Framework", "Node.js"},
		{"NestJS", "Node.js Framework", "Node.js"},
		{"Axum", "Rust Web Framework", "Rust"},
		{"Actix", "Rust Web Framework", "Rust"},
		{"Rocket", "Rust Web Framework", "Rust"},
	}

	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Select a project type:\n\n"))

	for i, pt := range projectTypes {
		var line strings.Builder

		cursor := "  "
		if m.projectTypeCursor == i {
			cursor = "▶ "
		}

		badge := getProjectBadge(strings.ToLower(pt.name)).Render(pt.name)
		line.WriteString(cursor)
		line.WriteString(badge + " ")
		line.WriteString(pt.desc)
		line.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(fmt.Sprintf(" (%s)", pt.runtime)))

		lineStr := line.String()
		if m.projectTypeCursor == i {
			lineStr = selectedItemStyle.Width(panelWidth - 4).Render(lineStr)
		} else {
			lineStr = normalItemStyle.Render(lineStr)
		}

		s.WriteString(lineStr + "\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(infoColor).
		Render("↑/↓: Navigate  Enter: Select  Esc: Cancel"))

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}
