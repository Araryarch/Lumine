package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderProjectsPanel() string {
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
		Render("󰉋  Projects")

	s.WriteString(titleStyle + "\n\n")

	if len(m.config.Projects) == 0 {
		emptyIcon := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("󰈙")

		emptyTitle := lipgloss.NewStyle().
			Foreground(subtleColor).
			Bold(true).
			Render("No projects yet")

		emptyDesc := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Press 'n' to create a new project")

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

		if showScrollTop {
			scrollIndicator := lipgloss.NewStyle().
				Foreground(mutedColor).
				Italic(true).
				Render(fmt.Sprintf("  ↑ %d more", startIdx))
			s.WriteString(scrollIndicator + "\n")
		}

		for i := startIdx; i < endIdx; i++ {
			project := m.config.Projects[i]
			var line strings.Builder

			icon := getIconForProject(project.Type)
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

			statusIcon := "󰀊"
			statusStyle := mutedColor
			if project.Status == "running" {
				statusIcon = "󰀄"
				statusStyle = successColor
			}

			badge := getProjectBadge(project.Type).Render(strings.ToUpper(project.Type))
			line.WriteString(badge + " ")

			line.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(icon + " "))
			line.WriteString(lipgloss.NewStyle().Bold(true).Foreground(fgColor).Render(project.Name))

			statusBadge := lipgloss.NewStyle().
				Foreground(bgColor).
				Background(statusStyle).
				Padding(0, 1).
				Render(" " + statusIcon + " ")
			line.WriteString(" " + statusBadge)

			if project.Domain != "" {
				domainBadge := lipgloss.NewStyle().
					Foreground(infoColor).
					Background(surface0).
					Padding(0, 1).
					Render("󰀟 " + project.Domain)
				line.WriteString(" " + domainBadge)
			}

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
				Render(fmt.Sprintf("  ↓ %d more", len(m.config.Projects)-endIdx))
			s.WriteString(scrollIndicator + "\n")
		}
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Render(s.String())
}

func (m model) renderRuntimesPanel() string {
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
		Render("󰌠  Runtimes")

	s.WriteString(titleStyle + "\n\n")

	runtimes := []struct {
		name    string
		version string
		icon    string
	}{
		{"PHP", m.config.Runtimes.PHP, "󰌞"},
		{"Node.js", m.config.Runtimes.Node, "󰛦"},
		{"Python", m.config.Runtimes.Python, "󰌠"},
		{"Rust", m.config.Runtimes.Rust, "󱘘"},
		{"Bun", m.config.Runtimes.Bun, "󰛦"},
		{"Deno", m.config.Runtimes.Deno, "󰛦"},
		{"Go", m.config.Runtimes.Go, "󰟓"},
	}

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

	if showScrollTop {
		scrollIndicator := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render(fmt.Sprintf("  ↑ %d more", startIdx))
		s.WriteString(scrollIndicator + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		rt := runtimes[i]
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

		line.WriteString(lipgloss.NewStyle().Foreground(infoColor).Render(rt.icon + " "))

		line.WriteString(lipgloss.NewStyle().Bold(true).Foreground(fgColor).Render(rt.name))

		versionBadge := lipgloss.NewStyle().
			Foreground(successColor).
			Background(surface0).
			Padding(0, 1).
			Render("v" + rt.version)
		line.WriteString(" " + versionBadge)

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
			Render(fmt.Sprintf("  ↓ %d more", len(runtimes)-endIdx))
		s.WriteString(scrollIndicator + "\n")
	}

	s.WriteString("\n")

	helpText := lipgloss.NewStyle().
		Foreground(mutedColor).
		Render("Press ") +
		lipgloss.NewStyle().Foreground(primaryColor).Bold(true).Render("'v'") +
		lipgloss.NewStyle().Foreground(mutedColor).Render(" to change version")
	s.WriteString(helpText)

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
		return lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface1).
			Padding(0, 1).
			Bold(true)
	}
}
