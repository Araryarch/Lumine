package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type cleanupOption int

const (
	cleanupContainer cleanupOption = iota
	cleanupVolume
	cleanupNetwork
	cleanupAll
)

func (m model) renderCleanupDialog() string {
	var s strings.Builder

	// Create overlay
	overlayWidth := 70
	overlayHeight := 25

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
		BorderForeground(errorColor).
		Padding(1, 2).
		Width(overlayWidth).
		Height(overlayHeight).
		Background(bgColor)

	var content strings.Builder

	// Header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(errorColor).
		Render("⚠️  Cleanup Options")
	content.WriteString(header + "\n\n")

	// Warning
	warning := lipgloss.NewStyle().
		Foreground(warningColor).
		Render("Select what to remove:")
	content.WriteString(warning + "\n\n")

	// Options
	options := []struct {
		name        string
		description string
		danger      bool
	}{
		{"Remove Container", "Stop and remove this container only", false},
		{"Remove with Volume", "Remove container and its data volume", true},
		{"Remove All Containers", "Remove all Lumine containers", true},
		{"Nuclear Cleanup", "Remove EVERYTHING (containers, volumes, networks)", true},
	}

	for i, opt := range options {
		var line string
		cursor := "  "
		if m.cleanupCursor == i {
			cursor = "▶ "
		}

		icon := "🗑️ "
		if opt.danger {
			icon = "💣 "
		}

		optText := fmt.Sprintf("%s%s%s", cursor, icon, opt.name)
		if m.cleanupCursor == i {
			line = selectedItemStyle.Width(overlayWidth - 6).Render(optText)
		} else {
			line = normalItemStyle.Render(optText)
		}

		content.WriteString(line + "\n")
		
		// Description
		desc := lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("   " + opt.description)
		content.WriteString(desc + "\n\n")
	}

	// Current selection info
	if m.selectedService != nil {
		content.WriteString("\n")
		info := lipgloss.NewStyle().
			Foreground(infoColor).
			Render(fmt.Sprintf("Selected: %s", m.selectedService.Name))
		content.WriteString(info + "\n")
	}

	// Help
	content.WriteString("\n")
	helpText := lipgloss.NewStyle().
		Foreground(infoColor).
		Render("↑/↓: Navigate  Enter: Confirm  Esc: Cancel")
	content.WriteString(helpText)

	overlay := overlayStyle.Render(content.String())

	// Add left padding
	lines := strings.Split(overlay, "\n")
	for _, line := range lines {
		s.WriteString(strings.Repeat(" ", leftPadding) + line + "\n")
	}

	return s.String()
}

func (m model) renderConfirmDialog(message string) string {
	var s strings.Builder

	// Create overlay
	overlayWidth := 60
	overlayHeight := 12

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
		BorderForeground(errorColor).
		Padding(2, 3).
		Width(overlayWidth).
		Height(overlayHeight).
		Background(bgColor)

	var content strings.Builder

	// Icon and message
	icon := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Render("⚠️  WARNING")
	content.WriteString(icon + "\n\n")

	msg := lipgloss.NewStyle().
		Foreground(fgColor).
		Render(message)
	content.WriteString(msg + "\n\n")

	// Confirmation prompt
	prompt := lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true).
		Render("Type 'yes' to confirm: ")
	content.WriteString(prompt)

	// Input field
	input := lipgloss.NewStyle().
		Foreground(infoColor).
		Render(m.confirmInput + "█")
	content.WriteString(input)

	overlay := overlayStyle.Render(content.String())

	// Add left padding
	lines := strings.Split(overlay, "\n")
	for _, line := range lines {
		s.WriteString(strings.Repeat(" ", leftPadding) + line + "\n")
	}

	return s.String()
}
