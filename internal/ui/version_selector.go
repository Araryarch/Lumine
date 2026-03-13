package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderVersionSelector() string {
	var s strings.Builder

	// Create overlay background
	overlayWidth := 60
	overlayHeight := 20

	// Center the overlay
	leftPadding := (m.width - overlayWidth) / 2
	topPadding := (m.height - overlayHeight) / 2

	// Add top padding
	for i := 0; i < topPadding; i++ {
		s.WriteString("\n")
	}

	// Overlay style
	overlayStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(secondaryColor).
		Padding(1, 2).
		Width(overlayWidth).
		Height(overlayHeight).
		Background(bgColor)

	var content strings.Builder

	// Header
	if m.selectedService != nil {
		header := lipgloss.NewStyle().
			Bold(true).
			Foreground(secondaryColor).
			Render(fmt.Sprintf("Select Version for %s", m.selectedService.Name))
		content.WriteString(header + "\n\n")

		currentVersion := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(fmt.Sprintf("Current: %s", m.selectedService.Version))
		content.WriteString(currentVersion + "\n\n")
	}

	// Version list
	maxDisplay := 10
	startIdx := 0
	if m.versionCursor >= maxDisplay {
		startIdx = m.versionCursor - maxDisplay + 1
	}

	endIdx := startIdx + maxDisplay
	if endIdx > len(m.availableVersions) {
		endIdx = len(m.availableVersions)
	}

	for i := startIdx; i < endIdx; i++ {
		version := m.availableVersions[i]
		var line string

		if i == m.versionCursor {
			line = selectedItemStyle.Width(overlayWidth - 6).Render(fmt.Sprintf("▶ %s", version))
		} else {
			line = normalItemStyle.Render(fmt.Sprintf("  %s", version))
		}
		content.WriteString(line + "\n")
	}

	// Scroll indicator
	if len(m.availableVersions) > maxDisplay {
		scrollInfo := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render(fmt.Sprintf("\n[%d/%d versions]", m.versionCursor+1, len(m.availableVersions)))
		content.WriteString(scrollInfo)
	}

	// Help text
	content.WriteString("\n\n")
	helpText := lipgloss.NewStyle().
		Foreground(infoColor).
		Render("↑/↓: Navigate  Enter: Select  Esc: Cancel")
	content.WriteString(helpText)

	overlay := overlayStyle.Render(content.String())

	// Add left padding
	lines := strings.Split(overlay, "\n")
	for _, line := range lines {
		s.WriteString(strings.Repeat(" ", leftPadding) + line + "\n")
	}

	return s.String()
}
