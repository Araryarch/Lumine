package ui

import "github.com/charmbracelet/lipgloss"

var (
	// Catppuccin Mocha color scheme
	primaryColor   = lipgloss.Color("#89b4fa")  // Blue
	secondaryColor = lipgloss.Color("#cba6f7")  // Mauve
	successColor   = lipgloss.Color("#a6e3a1")  // Green
	errorColor     = lipgloss.Color("#f38ba8")  // Red
	warningColor   = lipgloss.Color("#f9e2af")  // Yellow
	infoColor      = lipgloss.Color("#94e2d5")  // Teal
	mutedColor     = lipgloss.Color("#6c7086")  // Overlay0
	bgColor        = lipgloss.Color("#1e1e2e")  // Base
	fgColor        = lipgloss.Color("#cdd6f4")  // Text
	borderColor    = lipgloss.Color("#45475a")  // Surface1
	surfaceColor   = lipgloss.Color("#313244")  // Surface0
	
	// Additional Catppuccin colors
	peachColor     = lipgloss.Color("#fab387")  // Peach
	pinkColor      = lipgloss.Color("#f5c2e7")  // Pink
	lavenderColor  = lipgloss.Color("#b4befe")  // Lavender
	sapphireColor  = lipgloss.Color("#74c7ec")  // Sapphire

	// Base styles
	baseStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(bgColor)

	// Title bar
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#1a1b26")).
			Background(primaryColor).
			Padding(0, 2).
			Width(100)

	// Panel styles with better borders
	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginRight(1)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.ThickBorder()).
				BorderForeground(primaryColor).
				Padding(1, 2).
				MarginRight(1)

	// List item styles
	selectedItemStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 1)

	// Status styles with icons
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
			Underline(true).
			MarginBottom(1)

	subHeaderStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Italic(true)

	// Help styles - more prominent
	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(surfaceColor).
			Padding(0, 1).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lavenderColor).
			Bold(true)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(fgColor)

	// Info box
	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(infoColor).
			Padding(1, 2).
			MarginTop(1)

	warningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningColor).
			Padding(1, 2).
			MarginTop(1)

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1)

	// Status bar - more prominent
	statusBarStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 2).
			Bold(true)

	// Divider
	dividerStyle = lipgloss.NewStyle().
			Foreground(borderColor)

	// Service type badges with Catppuccin colors
	phpBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(lavenderColor).
			Padding(0, 1).
			Bold(true)

	mysqlBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(sapphireColor).
			Padding(0, 1).
			Bold(true)

	nginxBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(successColor).
			Padding(0, 1).
			Bold(true)

	redisBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(errorColor).
			Padding(0, 1).
			Bold(true)

	postgresBadge = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(primaryColor).
				Padding(0, 1).
				Bold(true)

	mariadbBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(infoColor).
			Padding(0, 1).
			Bold(true)

	apacheBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(peachColor).
			Padding(0, 1).
			Bold(true)

	mongoBadge = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(successColor).
			Padding(0, 1).
			Bold(true)

	defaultBadge = lipgloss.NewStyle().
			Foreground(fgColor).
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
