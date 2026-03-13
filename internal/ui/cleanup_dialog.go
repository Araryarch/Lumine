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
	// Modal dimensions
	modalWidth := 75

	var content strings.Builder

	// Title with icon
	titleStyle := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	title := titleStyle.Render("🗑️  Cleanup Options")
	content.WriteString(title + "\n")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat("─", modalWidth-4))
	content.WriteString(divider + "\n\n")

	// Warning box
	warningBox := lipgloss.NewStyle().
		Foreground(warningColor).
		Background(surfaceColor).
		Padding(0, 2).
		Align(lipgloss.Center).
		Width(modalWidth - 8).
		Bold(true).
		Render("⚠️  Select what to remove  ⚠️")
	content.WriteString(warningBox + "\n\n")

	// Options with better styling
	options := []struct {
		name        string
		description string
		icon        string
		danger      bool
	}{
		{"Remove Container", "Stop and remove this container only", "📦", false},
		{"Remove with Volume", "Remove container and its data volume", "💾", true},
		{"Remove All Containers", "Remove all Lumine containers", "🗂️", true},
		{"Nuclear Cleanup", "Remove EVERYTHING (containers, volumes, networks)", "💣", true},
	}

	for i, opt := range options {
		var line strings.Builder

		// Selection indicator
		if m.cleanupCursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true).
				Render("▶ "))
		} else {
			line.WriteString("  ")
		}

		// Icon
		line.WriteString(opt.icon + " ")

		// Option name
		nameStyle := lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true)
		if opt.danger {
			nameStyle = nameStyle.Foreground(errorColor)
		}
		line.WriteString(nameStyle.Render(opt.name))

		lineStr := line.String()

		// Highlight selected
		if m.cleanupCursor == i {
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

		// Description
		descStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Width(modalWidth - 10).
			Padding(0, 0, 0, 4)
		content.WriteString(descStyle.Render(opt.description) + "\n\n")
	}

	// Current selection info
	if m.selectedService != nil {
		content.WriteString("\n")
		infoBox := lipgloss.NewStyle().
			Foreground(infoColor).
			Background(surfaceColor).
			Padding(0, 2).
			Align(lipgloss.Center).
			Width(modalWidth - 8).
			Render(fmt.Sprintf("Selected: %s", m.selectedService.Name))
		content.WriteString(infoBox + "\n")
	}

	// Bottom divider
	content.WriteString("\n" + divider + "\n")

	// Help text with icons
	helpItems := []string{
		lipgloss.NewStyle().Foreground(infoColor).Render("↑↓") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" navigate"),
		lipgloss.NewStyle().Foreground(errorColor).Render("enter") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" confirm"),
		lipgloss.NewStyle().Foreground(successColor).Render("esc") +
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
		BorderForeground(errorColor).
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

func (m model) renderConfirmDialog(message string) string {
	// Modal dimensions
	modalWidth := 65

	var content strings.Builder

	// Title with icon
	titleStyle := lipgloss.NewStyle().
		Foreground(errorColor).
		Bold(true).
		Align(lipgloss.Center).
		Width(modalWidth - 4)

	title := titleStyle.Render("⚠️  WARNING  ⚠️")
	content.WriteString(title + "\n")

	// Divider
	divider := lipgloss.NewStyle().
		Foreground(borderColor).
		Render(strings.Repeat("─", modalWidth-4))
	content.WriteString(divider + "\n\n")

	// Message box
	messageBox := lipgloss.NewStyle().
		Foreground(fgColor).
		Background(surfaceColor).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(modalWidth - 8).
		Render(message)
	content.WriteString(messageBox + "\n\n")

	// Confirmation prompt
	promptStyle := lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true)
	content.WriteString(promptStyle.Render("Type 'yes' to confirm:") + "\n\n")

	// Input field with box
	inputBox := lipgloss.NewStyle().
		Foreground(infoColor).
		Background(lipgloss.Color("#313244")).
		Padding(0, 2).
		Width(modalWidth - 8).
		Align(lipgloss.Center).
		Render(m.confirmInput + "█")
	content.WriteString(inputBox + "\n\n")

	// Bottom divider
	content.WriteString(divider + "\n")

	// Help text
	helpText := lipgloss.NewStyle().
		Foreground(mutedColor).
		Align(lipgloss.Center).
		Width(modalWidth - 4).
		Render("Type 'yes' and press enter to confirm  •  esc to cancel")
	content.WriteString(helpText)

	// Modal box
	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
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

