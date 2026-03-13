package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (m model) renderLogsPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == mainPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Logs")
	s.WriteString(header + "\n\n")
	
	if len(m.logs) == 0 {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No logs yet...") + "\n")
		s.WriteString("\n")
		s.WriteString(subHeaderStyle.Render("Logs will appear here when you perform actions"))
	} else {
		// Show logs with scrolling
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
			
			// Timestamp
			timestamp := lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(log.timestamp)
			
			// Level with color
			var levelStyle lipgloss.Style
			var levelIcon string
			switch log.level {
			case "success":
				levelStyle = lipgloss.NewStyle().Foreground(successColor).Bold(true)
				levelIcon = "✓"
			case "error":
				levelStyle = lipgloss.NewStyle().Foreground(errorColor).Bold(true)
				levelIcon = "✗"
			case "warning":
				levelStyle = lipgloss.NewStyle().Foreground(warningColor).Bold(true)
				levelIcon = "⚠"
			default: // info
				levelStyle = lipgloss.NewStyle().Foreground(infoColor)
				levelIcon = "ℹ"
			}
			level := levelStyle.Render(levelIcon)
			
			// Service
			service := lipgloss.NewStyle().
				Foreground(primaryColor).
				Bold(true).
				Render(log.service)
			
			// Message
			message := lipgloss.NewStyle().
				Foreground(fgColor).
				Render(log.message)
			
			logLine := fmt.Sprintf("%s %s [%s] %s", timestamp, level, service, message)
			
			// Truncate if too long
			if lipgloss.Width(logLine) > panelWidth-6 {
				logLine = logLine[:panelWidth-9] + "..."
			}
			
			s.WriteString(logLine + "\n")
		}
		
		// Scroll indicator
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

	// Add spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}

func (m model) renderBackgroundTasksPanel() string {
	var s strings.Builder

	panelWidth := ((m.width - 24) / 2) - 2
	panelHeight := m.height - 5

	style := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Padding(1, 2).
		Width(panelWidth).
		Height(panelHeight)
	
	if m.activePanel == mainPanel {
		style = style.BorderForeground(primaryColor).Border(lipgloss.ThickBorder())
	}

	header := lipgloss.NewStyle().
		Foreground(primaryColor).
		Bold(true).
		Underline(true).
		Render("Background Tasks")
	s.WriteString(header + "\n\n")
	
	if len(m.backgroundTasks) == 0 {
		s.WriteString(lipgloss.NewStyle().
			Foreground(mutedColor).
			Italic(true).
			Render("No background tasks running...") + "\n")
		s.WriteString("\n")
		s.WriteString(subHeaderStyle.Render("Tasks will appear here when operations run"))
	} else {
		// Show recent tasks (last 20)
		startIdx := 0
		if len(m.backgroundTasks) > 20 {
			startIdx = len(m.backgroundTasks) - 20
		}

		for i := startIdx; i < len(m.backgroundTasks); i++ {
			task := m.backgroundTasks[i]
			
			// Status icon
			var statusStyle lipgloss.Style
			var statusIcon string
			switch task.status {
			case "running":
				statusStyle = lipgloss.NewStyle().Foreground(infoColor)
				statusIcon = "⟳"
			case "completed":
				statusStyle = lipgloss.NewStyle().Foreground(successColor)
				statusIcon = "✓"
			case "failed":
				statusStyle = lipgloss.NewStyle().Foreground(errorColor)
				statusIcon = "✗"
			default:
				statusStyle = lipgloss.NewStyle().Foreground(mutedColor)
				statusIcon = "○"
			}
			
			status := statusStyle.Render(statusIcon)
			
			// Task name
			name := lipgloss.NewStyle().
				Bold(true).
				Render(task.name)
			
			// Time
			time := lipgloss.NewStyle().
				Foreground(mutedColor).
				Render(task.startTime)
			
			taskLine := fmt.Sprintf("%s %s %s - %s", time, status, name, task.message)
			
			// Truncate if too long
			if lipgloss.Width(taskLine) > panelWidth-6 {
				taskLine = taskLine[:panelWidth-9] + "..."
			}
			
			s.WriteString(taskLine + "\n")
		}
	}

	// Add spacing
	currentLines := strings.Count(s.String(), "\n")
	for i := currentLines; i < panelHeight-2; i++ {
		s.WriteString("\n")
	}

	return style.Width(panelWidth).Height(panelHeight).Render(s.String())
}
