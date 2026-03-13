package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderProjectsPanel() string {
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
		Render("Projects")
	s.WriteString(header + "\n\n")

	if len(m.config.Projects) == 0 {
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
			Render("No projects yet") + "\n\n")
		empty.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Render("Press ") +
			lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("'n'") +
			lipgloss.NewStyle().Foreground(infoColor).Render(" to create one"))
		
		s.WriteString(emptyBox.Render(empty.String()))
	} else {
		// Calculate visible items with scroll
		maxVisibleItems := panelHeight - 5
		if maxVisibleItems < 3 {
			maxVisibleItems = 3
		}

		startIdx := 0
		endIdx := len(m.config.Projects)
		showScrollTop := false
		showScrollBottom := false
		
		if len(m.config.Projects) > maxVisibleItems {
			startIdx = m.cursor - (maxVisibleItems / 2)
			if startIdx < 0 {
				startIdx = 0
			}
			endIdx = startIdx + maxVisibleItems
			if endIdx > len(m.config.Projects) {
				endIdx = len(m.config.Projects)
				startIdx = endIdx - maxVisibleItems
				if startIdx < 0 {
					startIdx = 0
				}
			}
			
			showScrollTop = startIdx > 0
			showScrollBottom = endIdx < len(m.config.Projects)
		}

		// Scroll indicator top
		if showScrollTop {
			s.WriteString(lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n")
		}

		for i := startIdx; i < endIdx; i++ {
			project := m.config.Projects[i]
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
				Render(fmt.Sprintf("  ↓ %d more below", len(m.config.Projects)-endIdx)) + "\n")
		}
	}

	// Add spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderRuntimesPanel() string {
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
		Render("Runtimes")
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

	// Calculate visible items with scroll
	maxVisibleItems := panelHeight - 6
	if maxVisibleItems < 3 {
		maxVisibleItems = 3
	}

	startIdx := 0
	endIdx := len(runtimes)
	showScrollTop := false
	showScrollBottom := false
	
	if len(runtimes) > maxVisibleItems {
		startIdx = m.cursor - (maxVisibleItems / 2)
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + maxVisibleItems
		if endIdx > len(runtimes) {
			endIdx = len(runtimes)
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}
		
		showScrollTop = startIdx > 0
		showScrollBottom = endIdx < len(runtimes)
	}

	// Scroll indicator top
	if showScrollTop {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↑ %d more above", startIdx)) + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		rt := runtimes[i]
		var line strings.Builder

		cursor := "  "
		if m.cursor == i {
			cursor = lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("> ")
		}

		line.WriteString(cursor)
		line.WriteString(lipgloss.NewStyle().Bold(true).Render(rt.name))
		line.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Render(fmt.Sprintf(" v%s", rt.version)))

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
			Render(fmt.Sprintf("  ↓ %d more below", len(runtimes)-endIdx)) + "\n")
	}

	s.WriteString("\n")
	s.WriteString(lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Press ") +
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("'v'") +
		lipgloss.NewStyle().Foreground(mutedColor).Render(" to change version"))

	// Add spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
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
	case "axum", "actix", "rocket":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#CE422B")).
			Padding(0, 1).
			Bold(true)
	default:
		return defaultBadge
	}
}
