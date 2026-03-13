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
	// Project types
	projectTypes := []struct {
		name    string
		desc    string
		runtime string
		icon    string
	}{
		{"Laravel", "Full-stack PHP Framework", "PHP 8.2", "⚡"},
		{"Next.js", "React Framework with SSR", "Node.js 20", "▲"},
		{"Vue", "Progressive JavaScript Framework", "Node.js 20", "V"},
		{"Django", "High-level Python Framework", "Python 3.11", "★"},
		{"Express", "Minimal Node.js Framework", "Node.js 20", "E"},
		{"FastAPI", "Modern Python API Framework", "Python 3.11", "⚡"},
		{"Nuxt", "Vue Framework with SSR", "Node.js 20", "N"},
		{"SvelteKit", "Svelte Application Framework", "Node.js 20", "S"},
		{"Remix", "Full Stack React Framework", "Node.js 20", "R"},
		{"NestJS", "Progressive Node.js Framework", "Node.js 20", "◆"},
		{"Axum", "Ergonomic Rust Web Framework", "Rust 1.75", "⚙"},
		{"Actix", "Powerful Rust Web Framework", "Rust 1.75", "⚡"},
		{"Rocket", "Simple Rust Web Framework", "Rust 1.75", "🚀"},
	}

	// Modal dimensions - larger and more spacious
	modalWidth := 80
	modalHeight := 24
	
	// Calculate visible items
	maxVisibleItems := modalHeight - 8
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
	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)
	
	title := titleStyle.Render("✨ Create New Project ✨")
	content.WriteString(title + "\n")
	
	// Divider
	divider := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat("─", modalWidth-4))
	content.WriteString(divider + "\n\n")
	
	// Subheader
	subheader := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Render("Select a project type to get started:")
	content.WriteString(subheader + "\n")
	
	// Scroll indicator top
	if showScrollTop {
		scrollTop := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Align(lipgloss.Center).
			Width(modalWidth - 4).
			Render(fmt.Sprintf("↑ %d more above ↑", startIdx))
		content.WriteString(scrollTop + "\n")
	} else {
		content.WriteString("\n")
	}
	
	// Project types list with better styling
	for i := startIdx; i < endIdx; i++ {
		pt := projectTypes[i]
		var line strings.Builder
		
		// Selection indicator
		if m.projectTypeCursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render("▶ "))
		} else {
			line.WriteString("  ")
		}
		
		// Icon
		iconStyle := lipgloss.NewStyle().
			Foreground(secondaryColor).
			Bold(true)
		line.WriteString(iconStyle.Render(pt.icon) + " ")
		
		// Badge with better styling
		badge := getProjectBadge(strings.ToLower(pt.name)).
			Padding(0, 1).
			Render(pt.name)
		line.WriteString(badge + " ")
		
		// Description
		descStyle := lipgloss.NewStyle().
			Foreground(fgColor)
		line.WriteString(descStyle.Render(pt.desc))
		
		// Runtime badge
		runtimeBadge := lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(surfaceColor).
			Padding(0, 1).
			Italic(true).
			Render(pt.runtime)
		line.WriteString(" " + runtimeBadge)
		
		lineStr := line.String()
		
		// Highlight selected with gradient-like effect
		if m.projectTypeCursor == i {
			lineStr = lipgloss.NewStyle().
				Background(lipgloss.Color("#313244")).
				Foreground(fgColor).
				Width(modalWidth - 6).
				Padding(0, 1).
				Bold(true).
				Render(lineStr)
		} else {
			lineStr = lipgloss.NewStyle().
				Width(modalWidth - 6).
				Padding(0, 1).
				Render(lineStr)
		}
		
		content.WriteString(lineStr + "\n")
	}
	
	// Scroll indicator bottom
	if showScrollBottom {
		scrollBottom := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Align(lipgloss.Center).
			Width(modalWidth - 4).
			Render(fmt.Sprintf("↓ %d more below ↓", len(projectTypes)-endIdx))
		content.WriteString(scrollBottom + "\n")
	} else {
		content.WriteString("\n")
	}
	
	// Bottom divider
	content.WriteString(divider + "\n")
	
	// Help text with icons
	helpItems := []string{
		lipgloss.NewStyle().Foreground(infoColor).Render("↑↓") + 
			lipgloss.NewStyle().Foreground(mutedColor).Render(" navigate"),
		lipgloss.NewStyle().Foreground(successColor).Render("enter") + 
			lipgloss.NewStyle().Foreground(mutedColor).Render(" select"),
		lipgloss.NewStyle().Foreground(errorColor).Render("esc") + 
			lipgloss.NewStyle().Foreground(mutedColor).Render(" cancel"),
	}
	
	helpText := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Width(modalWidth - 4).
		Render(strings.Join(helpItems, "  •  "))
	content.WriteString(helpText)
	
	// Modal box with shadow effect
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
		Background(bgColor).
		Padding(2, 2).
		Width(modalWidth)
	
	modal := modalBox.Render(content.String())
	
	// Center modal on screen with overlay
	modalWithPosition := lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars("░"),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("#45475a")),
	)
	
	return modalWithPosition
}

func (m model) renderAddProjectPanel() string {
	// This is now handled by renderProjectCreateModal
	// Keeping this for backward compatibility
	return m.renderProjectCreateModal()
}
