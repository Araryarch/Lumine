package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderLogsPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == mainPanel {
		borderStyle = lipgloss.ThickBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == mainPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("[L]  Logs")

	s.WriteString(titleStyle + "\n\n")

	if len(m.logs) == 0 {
		emptyIcon := lipgloss.NewStyle().
			Foreground(fgMuted).
			Render("[L]")

		emptyTitle := lipgloss.NewStyle().
			Foreground(fgSubtle).
			Bold(true).
			Render("No logs yet")

		emptyDesc := lipgloss.NewStyle().
			Foreground(fgMuted).
			Render("Logs will appear when you perform actions")

		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(2, 4).
			Width(panelWidth - 8).
			Align(lipgloss.Center).
			Render(lipgloss.JoinVertical(lipgloss.Center,
				emptyIcon,
				"",
				emptyTitle,
				emptyDesc,
			))

		s.WriteString(emptyBox)
	} else {
		maxLogs := panelHeight - 6
		startIdx := len(m.logs) - maxLogs - m.logScrollOffset
		if startIdx < 0 {
			startIdx = 0
		}
		endIdx := len(m.logs) - m.logScrollOffset
		if endIdx > len(m.logs) {
			endIdx = len(m.logs)
		}

		for i := startIdx; i < endIdx; i++ {
			log := m.logs[i]

			timestamp := lipgloss.NewStyle().
				Foreground(fgMuted).
				Render(log.timestamp)

			var levelStyle lipgloss.Style
			var levelIcon string
			switch log.level {
			case "success":
				levelStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
				levelIcon = "OK"
			case "error":
				levelStyle = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
				levelIcon = "X"
			case "warning":
				levelStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
				levelIcon = "!"
			default:
				levelStyle = lipgloss.NewStyle().Foreground(infoColor)
				levelIcon = "i"
			}
			level := levelStyle.Render(levelIcon)

			service := lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(log.service)

			message := lipgloss.NewStyle().
				Foreground(fgColor).
				Render(log.message)

			logLine := fmt.Sprintf("%s %s [%s] %s", timestamp, level, service, message)

			if lipgloss.Width(logLine) > panelWidth-6 {
				logLine = logLine[:panelWidth-9] + "..."
			}

			s.WriteString(logLine + "\n")
		}

		if m.logScrollOffset > 0 {
			s.WriteString("\n")
			s.WriteString(lipgloss.NewStyle().
				Foreground(infoColor).
				Render(fmt.Sprintf("↑ Scrolled up (%d more logs below)", m.logScrollOffset)))
		}

		if startIdx > 0 {
			s.WriteString("\n")
			s.WriteString(lipgloss.NewStyle().
				Foreground(infoColor).
				Render(fmt.Sprintf("↓ %d more logs above", startIdx)))
		}
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderBackgroundTasksPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 26) / 2) - 2
	panelHeight := m.height - 5

	borderStyle := lipgloss.NormalBorder()
	if m.activePanel == mainPanel {
		borderStyle = lipgloss.ThickBorder()
	}

	borderColorStyle := borderColor
	if m.activePanel == mainPanel {
		borderColorStyle = primaryColor
	}

	style := lipgloss.NewStyle().
		Border(borderStyle).
		BorderForeground(borderColorStyle).
		Background(bgColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)

	titleStyle := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Padding(0, 1).
		Background(surface0).
		Width(panelWidth - 4).
		Render("[S]  Background Tasks")

	s.WriteString(titleStyle + "\n\n")

	if len(m.backgroundTasks) == 0 {
		emptyIcon := lipgloss.NewStyle().
			Foreground(fgMuted).
			Render("[S]")

		emptyTitle := lipgloss.NewStyle().
			Foreground(fgSubtle).
			Bold(true).
			Render("No background tasks")

		emptyDesc := lipgloss.NewStyle().
			Foreground(fgMuted).
			Render("Tasks appear when operations run")

		emptyBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Padding(2, 4).
			Width(panelWidth - 8).
			Align(lipgloss.Center).
			Render(lipgloss.JoinVertical(lipgloss.Center,
				emptyIcon,
				"",
				emptyTitle,
				emptyDesc,
			))

		s.WriteString(emptyBox)
	} else {
		startIdx := 0
		if len(m.backgroundTasks) > 20 {
			startIdx = len(m.backgroundTasks) - 20
		}

		for i := startIdx; i < len(m.backgroundTasks); i++ {
			task := m.backgroundTasks[i]

			var statusStyle lipgloss.Style
			var statusIcon string
			switch task.status {
			case "running":
				statusStyle = lipgloss.NewStyle().Foreground(infoColor)
				statusIcon = "⟳"
			case "completed":
				statusStyle = lipgloss.NewStyle().Foreground(successColor)
				statusIcon = "󰸞"
			case "failed":
				statusStyle = lipgloss.NewStyle().Foreground(errorColor)
				statusIcon = "󰚌"
			default:
				statusStyle = lipgloss.NewStyle().Foreground(fgMuted)
				statusIcon = "○"
			}

			status := statusStyle.Render(statusIcon)

			time := lipgloss.NewStyle().
				Foreground(fgMuted).
				Render(task.startTime)

			name := lipgloss.NewStyle().
				Bold(true).
				Foreground(fgColor).
				Render(task.name)

			taskLine := fmt.Sprintf("%s %s %s - %s", time, status, name, task.message)

			if lipgloss.Width(taskLine) > panelWidth-6 {
				taskLine = taskLine[:panelWidth-9] + "..."
			}

			s.WriteString(taskLine + "\n")
		}
	}

	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}
