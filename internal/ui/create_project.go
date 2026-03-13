package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type projectFormStep int

const (
	selectProjectType projectFormStep = iota
	enterProjectName
	enterProjectPath
	confirmProjectCreation
)

func (m model) renderProjectCreateModal() string {
	projectTypes := []struct {
		name    string
		desc    string
		runtime string
		icon    string
	}{
		{"Laravel", "Full-stack PHP Framework", "PHP 8.2", "󰖬"},
		{"Next.js", "React Framework with SSR", "Node.js 20", "[W]"},
		{"Vue", "Progressive JavaScript Framework", "Node.js 20", "󰡄"},
		{"Django", "High-level Python Framework", "Python 3.11", "[S]"},
		{"Express", "Minimal Node.js Framework", "Node.js 20", "[S]"},
		{"FastAPI", "Modern Python API Framework", "Python 3.11", "[S]"},
		{"Nuxt", "Vue Framework with SSR", "Node.js 20", "[W]"},
		{"SvelteKit", "Svelte Application Framework", "Node.js 20", "[W]"},
		{"Remix", "Full Stack React Framework", "Node.js 20", "[W]"},
		{"NestJS", "Progressive Node.js Framework", "Node.js 20", "[W]"},
		{"Axum", "Ergonomic Rust Web Framework", "Rust 1.75", "[S]"},
		{"Actix", "Powerful Rust Web Framework", "Rust 1.75", "[S]"},
		{"Rocket", "Simple Rust Web Framework", "Rust 1.75", "[S]"},
	}

	modalWidth := 80
	modalHeight := 26

	maxVisibleItems := modalHeight - 10
	if maxVisibleItems < 5 {
		maxVisibleItems = 5
	}

	startIdx := 0
	endIdx := len(projectTypes)
	showScrollTop := false
	showScrollBottom := false

	if len(projectTypes) > maxVisibleItems {
		startIdx = m.projectTypeCursor - (maxVisibleItems / 2)
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + maxVisibleItems
		if endIdx > len(projectTypes) {
			endIdx = len(projectTypes)
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}

		showScrollTop = startIdx > 0
		showScrollBottom = endIdx < len(projectTypes)
	}

	var content strings.Builder

	// Title with icon
	titleBox := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(primaryColor).
		Bold(true).
		Padding(1, 3).
		Render(" [+]  Create New Project  ")
	content.WriteString(titleBox + "\n\n")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(surface1).
		Render(strings.Repeat("─", modalWidth-6))
	content.WriteString(divider + "\n\n")

	// Subheader
	subheader := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Render("Select a project type to get started:")
	content.WriteString(subheader + "\n\n")

	// Scroll indicator top
	if showScrollTop {
		scrollTop := lipgloss.NewStyle().
			Foreground(fgMuted).
			Italic(true).
			Align(lipgloss.Center).
			Width(modalWidth - 6).
			Render(fmt.Sprintf("↑ %d more above ↑", startIdx))
		content.WriteString(scrollTop + "\n")
	}

	// Project types list
	for i := startIdx; i < endIdx; i++ {
		pt := projectTypes[i]
		var line strings.Builder

		// Selection indicator
		if m.projectTypeCursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(" "))
		} else {
			line.WriteString("  ")
		}

		// Icon
		line.WriteString(lipgloss.NewStyle().
			Foreground(infoColor).
			Render(pt.icon) + " ")

		// Badge with better styling
		badge := getProjectBadge(strings.ToLower(pt.name)).
			Padding(0, 1).
			Bold(true).
			Render(pt.name)
		line.WriteString(badge + " ")

		// Description
		line.WriteString(lipgloss.NewStyle().
			Foreground(fgColor).
			Render(pt.desc))

		// Runtime badge
		runtimeBadge := lipgloss.NewStyle().
			Foreground(fgMuted).
			Background(surface0).
			Padding(0, 1).
			Render(pt.runtime)
		line.WriteString(" " + runtimeBadge)

		lineStr := line.String()

		// Highlight selected
		if m.projectTypeCursor == i {
			lineStr = lipgloss.NewStyle().
				Background(surface0).
				Foreground(fgColor).
				Width(modalWidth-8).
				Padding(0, 1).
				Bold(true).
				Render(lineStr)
		} else {
			lineStr = lipgloss.NewStyle().
				Width(modalWidth-8).
				Padding(0, 1).
				Render(lineStr)
		}

		content.WriteString(lineStr + "\n")
	}

	// Scroll indicator bottom
	if showScrollBottom {
		scrollBottom := lipgloss.NewStyle().
			Foreground(fgMuted).
			Italic(true).
			Align(lipgloss.Center).
			Width(modalWidth - 6).
			Render(fmt.Sprintf("↓ %d more below ↓", len(projectTypes)-endIdx))
		content.WriteString(scrollBottom + "\n")
	}

	// Bottom divider
	content.WriteString("\n")
	content.WriteString(divider + "\n\n")

	// Help text with icons
	helpBox := lipgloss.NewStyle().
		Background(surface0).
		Padding(1, 2).
		Width(modalWidth - 6).
		Align(lipgloss.Center)

	helpItems := []string{
		lipgloss.NewStyle().Foreground(bgColor).
			Background(primaryColor).Padding(0, 2).Bold(true).Render(" ↑↓ ") +
			lipgloss.NewStyle().Foreground(fgMuted).Render(" navigate"),
		lipgloss.NewStyle().Foreground(bgColor).
			Background(successColor).Padding(0, 2).Bold(true).Render(" enter ") +
			lipgloss.NewStyle().Foreground(fgMuted).Render(" select"),
		lipgloss.NewStyle().Foreground(bgColor).
			Background(errorColor).Padding(0, 2).Bold(true).Render(" esc ") +
			lipgloss.NewStyle().Foreground(fgMuted).Render(" cancel"),
	}

	helpText := strings.Join(helpItems, "  "+lipgloss.NewStyle().Foreground(surface1).Render("│")+"  ")
	content.WriteString(helpBox.Render(helpText))

	// Modal box with shadow effect
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Background(bgColor).
		Padding(0, 0).
		Width(modalWidth)

	modal := modalBox.Render(content.String())

	// Center modal on screen with overlay
	modalWithPosition := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(surface0),
	)

	return modalWithPosition
}

func (m model) renderAddProjectPanel() string {
	return m.renderProjectCreateModal()
}
