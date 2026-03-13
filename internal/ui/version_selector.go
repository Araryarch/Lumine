package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderVersionSelector() string {
	// Modal dimensions
	modalWidth := 70
	modalHeight := 22

	var content strings.Builder

	// Title with icon
	titleStyle := lipgloss.NewStyle().
		Foreground(secondaryColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	var title string
	if m.selectedService != nil {
		title = titleStyle.Render(fmt.Sprintf("🔧 Select Version for %s", m.selectedService.Name))
	} else if m.selectedRuntimeType != "" {
		title = titleStyle.Render(fmt.Sprintf("🔧 Select %s Version", strings.Title(m.selectedRuntimeType)))
	} else {
		title = titleStyle.Render("🔧 Select Version")
	}
	content.WriteString(title + "\n")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat("─", modalWidth-4))
	content.WriteString(divider + "\n\n")

	// Current version info
	if m.selectedService != nil {
		currentBox := lipgloss.NewStyle().
			Foreground(infoColor).
			Background(surfaceColor).
			Padding(0, 2).
			Align(lipgloss.Center).
			Width(modalWidth - 8).
			Render(fmt.Sprintf("Current: %s", m.selectedService.Version))
		content.WriteString(currentBox + "\n\n")
	}

	// Subheader
	subheader := lipgloss.NewStyle().
		Foreground(mutedColor).
		Bold(true).
		Render("Available Versions:")
	content.WriteString(subheader + "\n")

	// Calculate visible items
	maxVisibleItems := modalHeight - 12
	if maxVisibleItems < 5 {
		maxVisibleItems = 5
	}

	startIdx := 0
	endIdx := len(m.availableVersions)
	showScrollTop := false
	showScrollBottom := false

	if len(m.availableVersions) > maxVisibleItems {
		startIdx = m.versionCursor - (maxVisibleItems / 2)
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx = startIdx + maxVisibleItems
		if endIdx > len(m.availableVersions) {
			endIdx = len(m.availableVersions)
			startIdx = endIdx - maxVisibleItems
			if startIdx < 0 {
				startIdx = 0
			}
		}

		showScrollTop = startIdx > 0
		showScrollBottom = endIdx < len(m.availableVersions)
	}

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

	// Version list with better styling
	for i := startIdx; i < endIdx; i++ {
		version := m.availableVersions[i]
		var line strings.Builder

		// Selection indicator
		if i == m.versionCursor {
			line.WriteString(lipgloss.NewStyle().
				Foreground(secondaryColor).
				Bold(true).
				Render("▶ "))
		} else {
			line.WriteString("  ")
		}

		// Version badge
		versionBadge := lipgloss.NewStyle().
			Foreground(successColor).
			Bold(true).
			Render("v" + version)
		line.WriteString(versionBadge)

		// Tag for special versions
		if strings.Contains(version, "latest") {
			tag := lipgloss.NewStyle().
				Foreground(warningColor).
				Background(surfaceColor).
				Padding(0, 1).
				Render("LATEST")
			line.WriteString(" " + tag)
		} else if strings.Contains(version, "alpine") {
			tag := lipgloss.NewStyle().
				Foreground(infoColor).
				Background(surfaceColor).
				Padding(0, 1).
				Render("ALPINE")
			line.WriteString(" " + tag)
		}

		lineStr := line.String()

		// Highlight selected
		if i == m.versionCursor {
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
			Render(fmt.Sprintf("↓ %d more below ↓", len(m.availableVersions)-endIdx))
		content.WriteString(scrollBottom + "\n")
	} else {
		content.WriteString("\n")
	}

	// Counter
	counter := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		Width(modalWidth - 4).
		Render(fmt.Sprintf("[%d/%d]", m.versionCursor+1, len(m.availableVersions)))
	content.WriteString(counter + "\n")

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

	// Modal box
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
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

