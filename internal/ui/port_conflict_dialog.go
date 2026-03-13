package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderPortConflictDialog() string {
	var s strings.Builder

	// Create overlay
	overlayWidth := 65
	overlayHeight := 18

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
		BorderForeground(warningColor).
		Padding(1, 2).
		Width(overlayWidth).
		Height(overlayHeight).
		Background(bgColor)

	var content strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(warningColor).
		Render("WARNING: Port Conflict Detected")
	content.WriteString(header + "\n\n")

	// Conflict info
	if m.portConflict != nil {
		info := fmt.Sprintf("Port %d is already in use by another service.", m.portConflict.Port)
		content.WriteString(lipgloss.NewStyle().Foreground(fgColor).Render(info) + "\n\n")

		// Alternative ports
		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(infoColor).
			Render("Available alternative ports:") + "\n\n")

		for i, altPort := range m.portConflict.Alternatives {
			var line string
			cursor := "  "
			if m.portConflictCursor == i {
				cursor = "> "
			}

			portText := fmt.Sprintf("%sPort %d", cursor, altPort)
			if m.portConflictCursor == i {
				line = selectedItemStyle.Width(overlayWidth - 6).Render(portText)
			} else {
				line = normalItemStyle.Render(portText)
			}

			content.WriteString(line + "\n")
		}

		content.WriteString("\n")
		content.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Or enter custom port: "))
		
		if m.customPortInput != "" {
			content.WriteString(lipgloss.NewStyle().
				Foreground(infoColor).
				Render(m.customPortInput + "█"))
		}
	}

	// Help
	content.WriteString("\n\n")
	helpText := lipgloss.NewStyle().
		Foreground(infoColor).
		Render("↑/↓: Navigate  Enter: Use port  Esc: Cancel  Type: Custom port")
	content.WriteString(helpText)

	overlay := overlayStyle.Render(content.String())

	// Add left padding
	lines := strings.Split(overlay, "\n")
	for _, line := range lines {
		s.WriteString(strings.Repeat(" ", leftPadding) + line + "\n")
	}

	return s.String()
}
