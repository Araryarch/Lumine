package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderVersionSelector() string {
	modalWidth := 70
	modalHeight := 22

	var content strings.Builder

	titleBox := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(primaryColor).
		Bold(true).
		Padding(1, 3).
		Render(" [R]  Select Version  ")

	content.WriteString(titleBox + "\n\n")

	divider := lipgloss.NewStyle().
		Foreground(surface1).
		Render(strings.Repeat("─", modalWidth-6))
	content.WriteString(divider + "\n\n")

	if m.selectedService != nil {
		currentBox := lipgloss.NewStyle().
			Foreground(infoColor).
			Background(surface0).
			Padding(0, 2).
			Align(lipgloss.Center).
			Width(modalWidth - 12).
			Render(fmt.Sprintf(" %s ", m.selectedService.Name))
		content.WriteString(currentBox + "\n\n")
	}

	subheader := lipgloss.NewStyle().
		Foreground(fgMuted).
		Bold(true).
		Render("Available Versions:")
	content.WriteString(subheader + "\n")

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

	if showScrollTop {
		scrollTop := lipgloss.NewStyle().
			Foreground(fgMuted).
			Italic(true).
			Align(lipgloss.Center).
			Width(modalWidth - 8).
			Render(fmt.Sprintf("↑ %d more above ↑", startIdx))
		content.WriteString(scrollTop + "\n")
	}

	for i := startIdx; i < endIdx; i++ {
		version := m.availableVersions[i]
		var line strings.Builder

		if i == m.versionCursor {
			line.WriteString(lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(" "))
		} else {
			line.WriteString("  ")
		}

		versionBadge := lipgloss.NewStyle().
			Foreground(successColor).
			Background(surface0).
			Padding(0, 1).
			Bold(true).
			Render(" v" + version + " ")
		line.WriteString(versionBadge)

		if strings.Contains(version, "latest") {
			tag := lipgloss.NewStyle().
				Foreground(bgColor).
				Background(warningColor).
				Padding(0, 1).
				Bold(true).
				Render(" LATEST ")
			line.WriteString(" " + tag)
		} else if strings.Contains(version, "alpine") {
			tag := lipgloss.NewStyle().
				Foreground(bgColor).
				Background(infoColor).
				Padding(0, 1).
				Bold(true).
				Render(" ALPINE ")
			line.WriteString(" " + tag)
		}

		lineStr := line.String()

		if i == m.versionCursor {
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

	if showScrollBottom {
		scrollBottom := lipgloss.NewStyle().
			Foreground(fgMuted).
			Italic(true).
			Align(lipgloss.Center).
			Width(modalWidth - 8).
			Render(fmt.Sprintf("↓ %d more below ↓", len(m.availableVersions)-endIdx))
		content.WriteString(scrollBottom + "\n")
	}

	counter := lipgloss.NewStyle().
		Foreground(fgMuted).
		Background(surface0).
		Align(lipgloss.Center).
		Padding(0, 1).
		Width(modalWidth - 8).
		Render(fmt.Sprintf(" [%d/%d]", m.versionCursor+1, len(m.availableVersions)))
	content.WriteString("\n" + counter + "\n")

	content.WriteString(divider + "\n")

	helpBox := lipgloss.NewStyle().
		Background(surface0).
		Padding(1, 2).
		Width(modalWidth - 8).
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

	modalBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(primaryColor).
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
