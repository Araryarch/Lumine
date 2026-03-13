package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Colors
	primaryColor   = lipgloss.Color("#7D56F4")
	secondaryColor = lipgloss.Color("#F780E2")
	successColor   = lipgloss.Color("#00D787")
	errorColor     = lipgloss.Color("#FF5F87")
	warningColor   = lipgloss.Color("#FFD700")
	infoColor      = lipgloss.Color("#5FD7FF")
	mutedColor     = lipgloss.Color("#626262")
	bgColor        = lipgloss.Color("#1a1b26")
	fgColor        = lipgloss.Color("#c0caf5")

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(bgColor)

	// Title bar
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 1)

	// Panel styles
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(1, 2).
			MarginRight(1).
			MarginTop(1)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(secondaryColor).
				Padding(1, 2).
				MarginRight(1).
				MarginTop(1)

	// List item styles
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 1)

	// Status styles
	runningStatusStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	stoppedStatusStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	errorStatusStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	// Header styles
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			MarginBottom(1)

	// Help styles
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Bold(true)

	// Info box
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(infoColor).
			Padding(0, 1).
			MarginTop(1)

	// Status bar
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(primaryColor).
			Padding(0, 1).
			MarginTop(1)

	// Service type badges
	phpBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#777BB3")).
			Padding(0, 1).
			Bold(true)

	mysqlBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#00758F")).
			Padding(0, 1).
			Bold(true)

	nginxBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#009639")).
			Padding(0, 1).
			Bold(true)

	redisBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#DC382D")).
			Padding(0, 1).
			Bold(true)

	postgresBadge = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FFFFFF")).
				Background(lipgloss.Color("#336791")).
				Padding(0, 1).
				Bold(true)

	mariadbBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#003545")).
			Padding(0, 1).
			Bold(true)

	apacheBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#D22128")).
			Padding(0, 1).
			Bold(true)

	mongoBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#47A248")).
			Padding(0, 1).
			Bold(true)

	defaultBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(mutedColor).
			Padding(0, 1).
			Bold(true)
)

func getServiceBadge(serviceType string) lipgloss.Style {
	switch serviceType {
	case "php":
		return phpBadge
	case "mysql":
		return mysqlBadge
	case "mariadb":
		return mariadbBadge
	case "nginx":
		return nginxBadge
	case "redis":
		return redisBadge
	case "postgres", "postgresql":
		return postgresBadge
	case "apache", "httpd":
		return apacheBadge
	case "mongodb", "mongo":
		return mongoBadge
	case "phpmyadmin":
		return phpBadge
	case "adminer":
		return mysqlBadge
	default:
		return defaultBadge
	}
}
