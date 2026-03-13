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
	modalWidth := 75

	var content strings.Builder

	titleBox := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(errorColor).
		Bold(true).
		Padding(1, 3).
		Render(" 󰀎  Cleanup Options  ")

	content.WriteString(titleBox + "\n\n")

	divider := lipgloss.NewStyle().
		Foreground(surface1).
		Render(strings.Repeat("─", modalWidth-6))
	content.WriteString(divider + "\n\n")

	warningBox := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(warningColor).
		Padding(0, 2).
		Align(lipgloss.Center).
		Width(modalWidth - 12).
		Bold(true).
		Render(" 󰀦  Select what to remove  ")
	content.WriteString(warningBox + "\n\n")

	options := []struct {
		name        string
		description string
		icon        string
		danger      bool
	}{
		{"Remove Container", "Stop and remove this container only", "󰡨", false},
		{"Remove with Volume", "Remove container and its data volume", "󰉋", true},
		{"Remove All Containers", "Remove all Lumine containers", "󰡨", true},
		{"Nuclear Cleanup", "Remove EVERYTHING", "󰀎", true},
	}

	for i, opt := range options {
		var line strings.Builder

		if m.cleanupCursor == i {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(" "))
		} else {
			line.WriteString("  ")
		}

		line.WriteString(opt.icon + " ")

		nameStyle := lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true)
		if opt.danger {
			nameStyle = nameStyle.Foreground(errorColor)
		}
		line.WriteString(nameStyle.Render(opt.name))

		lineStr := line.String()

		if m.cleanupCursor == i {
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

		descStyle := lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Width(modalWidth-14).
			Padding(0, 0, 0, 4)
		content.WriteString(descStyle.Render(opt.description) + "\n\n")
	}

	if m.selectedService != nil {
		content.WriteString("\n")
		infoBox := lipgloss.NewStyle().
			Foreground(bgColor).
			Background(infoColor).
			Padding(0, 2).
			Align(lipgloss.Center).
			Width(modalWidth - 12).
			Render(fmt.Sprintf(" Selected: %s ", m.selectedService.Name))
		content.WriteString(infoBox + "\n")
	}

	content.WriteString("\n" + divider + "\n")

	helpBox := lipgloss.NewStyle().
		Background(surface0).
		Padding(1, 2).
		Width(modalWidth - 8).
		Align(lipgloss.Center)

	helpItems := []string{
		lipgloss.NewStyle().Foreground(bgColor).
			Background(primaryColor).Padding(0, 2).Bold(true).Render(" ↑↓ ") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" navigate"),
		lipgloss.NewStyle().Foreground(bgColor).
			Background(errorColor).Padding(0, 2).Bold(true).Render(" enter ") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" confirm"),
		lipgloss.NewStyle().Foreground(bgColor).
			Background(successColor).Padding(0, 2).Bold(true).Render(" esc ") +
			lipgloss.NewStyle().Foreground(mutedColor).Render(" cancel"),
	}

	helpText := strings.Join(helpItems, "  "+lipgloss.NewStyle().Foreground(surface1).Render("│")+"  ")
	content.WriteString(helpBox.Render(helpText))

	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(errorColor).
		Background(bgColor).
		Padding(0, 0).
		Width(modalWidth)

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

func (m model) renderConfirmDialog(message string) string {
	modalWidth := 65

	var content strings.Builder

	titleBox := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(warningColor).
		Bold(true).
		Padding(1, 3).
		Render(" 󰀦  WARNING  ")

	content.WriteString(titleBox + "\n\n")

	divider := lipgloss.NewStyle().
		Foreground(surface1).
		Render(strings.Repeat("─", modalWidth-6))
	content.WriteString(divider + "\n\n")

	messageBox := lipgloss.NewStyle().
		Foreground(fgColor).
		Background(surface0).
		Padding(1, 2).
		Align(lipgloss.Center).
		Width(modalWidth - 12).
		Render(" " + message + " ")
	content.WriteString(messageBox + "\n\n")

	promptStyle := lipgloss.NewStyle().
		Foreground(warningColor).
		Bold(true)
	content.WriteString(promptStyle.Render("Type 'yes' to confirm:") + "\n\n")

	inputBox := lipgloss.NewStyle().
		Foreground(fgColor).
		Background(surface0).
		Padding(0, 2).
		Width(modalWidth - 12).
		Align(lipgloss.Center).
		Render(m.confirmInput + "󰍟")
	content.WriteString(inputBox + "\n\n")

	content.WriteString(divider + "\n")

	helpBox := lipgloss.NewStyle().
		Background(surface0).
		Padding(1, 2).
		Width(modalWidth - 8).
		Align(lipgloss.Center).
		Foreground(mutedColor).
		Render("Type 'yes' and press enter to confirm  •  esc to cancel")
	content.WriteString(helpBox)

	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(warningColor).
		Background(bgColor).
		Padding(0, 0).
		Width(modalWidth)

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
