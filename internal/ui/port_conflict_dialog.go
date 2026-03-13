package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderPortConflictDialog() string {
	overlayWidth := 65

	var content strings.Builder

	titleBox := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(warningColor).
		Bold(true).
		Padding(1, 3).
		Render(" 󰀦  Port Conflict Detected  ")

	content.WriteString(titleBox + "\n\n")

	divider := lipgloss.NewStyle().
		Foreground(surface1).
		Render(strings.Repeat("─", overlayWidth-6))
	content.WriteString(divider + "\n\n")

	if m.portConflict != nil {
		infoBox := lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface0).
			Padding(1, 2).
			Width(overlayWidth - 12).
			Align(lipgloss.Center).
			Render(fmt.Sprintf(" Port %d is already in use by another service ", m.portConflict.Port))
		content.WriteString(infoBox + "\n\n")

		content.WriteString(lipgloss.NewStyle().
			Bold(true).
			Foreground(infoColor).
			Render("Available alternative ports:") + "\n\n")

		for i, altPort := range m.portConflict.Alternatives {
			var line strings.Builder
			if m.portConflictCursor == i {
				line.WriteString(lipgloss.NewStyle().
					Foreground(primaryColor).
					Bold(true).
					Render(" "))
			} else {
				line.WriteString("  ")
			}

			portBadge := lipgloss.NewStyle().
				Foreground(successColor).
				Background(surface0).
				Padding(0, 1).
				Bold(true).
				Render(fmt.Sprintf(" :%d ", altPort))
			line.WriteString(portBadge)

			lineStr := line.String()
			if m.portConflictCursor == i {
				lineStr = lipgloss.NewStyle().
					Background(surface0).
					Width(overlayWidth-8).
					Padding(0, 1).
					Render(lineStr)
			} else {
				lineStr = lipgloss.NewStyle().
					Width(overlayWidth-8).
					Padding(0, 1).
					Render(lineStr)
			}

			content.WriteString(lineStr + "\n")
		}

		content.WriteString("\n")
		content.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Render("Or enter custom port: "))

		if m.customPortInput != "" {
			inputBox := lipgloss.NewStyle().
				Foreground(infoColor).
				Background(surface0).
				Padding(0, 1).
				Render(m.customPortInput + "󰍟")
			content.WriteString(inputBox)
		}
	}

	content.WriteString("\n\n" + divider + "\n")

	helpBox := lipgloss.NewStyle().
		Background(surface0).
		Padding(1, 2).
		Width(overlayWidth - 8).
		Align(lipgloss.Center)

	helpItems := []string{
		lipgloss.NewStyle().Foreground(bgColor).
			Background(primaryColor).Padding(0, 2).Bold(true).Render(" ↑↓ ") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" navigate"),
		lipgloss.NewStyle().Foreground(bgColor).
			Background(successColor).Padding(0, 2).Bold(true).Render(" enter ") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" use port"),
		lipgloss.NewStyle().Foreground(bgColor).
			Background(errorColor).Padding(0, 2).Bold(true).Render(" esc ") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" cancel"),
	}

	helpText := strings.Join(helpItems, "  "+lipgloss.NewStyle().Foreground(surface1).Render("│")+"  ")
	content.WriteString(helpBox.Render(helpText))

	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(warningColor).
		Background(bgColor).
		Padding(0, 0).
		Width(overlayWidth)

	modal := modalBox.Render(content.String())

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
