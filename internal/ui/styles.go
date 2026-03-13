package ui

import "github.com/charmbracelet/lipgloss"

var (
	gradientStart   = lipgloss.Color("#89b4fa")
	gradientEnd     = lipgloss.Color("#cba6f7")
	primaryColor    = lipgloss.Color("#89b4fa")
	secondaryColor  = lipgloss.Color("#cba6f7")
	accentColor     = lipgloss.Color("#f5c2e7")
	successColor    = lipgloss.Color("#a6e3a1")
	errorColor      = lipgloss.Color("#f38ba8")
	warningColor    = lipgloss.Color("#f9e2af")
	infoColor       = lipgloss.Color("#94e2d5")
	bgColor         = lipgloss.Color("#11111b")
	surfaceBg       = lipgloss.Color("#1e1e2e")
	surface0        = lipgloss.Color("#313244")
	surface1        = lipgloss.Color("#45475a")
	surface2        = lipgloss.Color("#585b70")
	fgColor         = lipgloss.Color("#cdd6f4")
	mutedColor      = lipgloss.Color("#6c7086")
	subtleColor     = lipgloss.Color("#9399b2")
	peachColor      = lipgloss.Color("#fab387")
	pinkColor       = lipgloss.Color("#f5c2e7")
	lavenderColor   = lipgloss.Color("#b4befe")
	sapphireColor   = lipgloss.Color("#74c7ec")
	rosewaterColor  = lipgloss.Color("#f5e0dc")
	borderColor     = lipgloss.Color("#45475a")
	highlightBorder = lipgloss.Color("#89dceb")
	surfaceColor    = lipgloss.Color("#313244")
	dividerStyle    = lipgloss.NewStyle().Foreground(surface1)
)

var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(fgColor).
			Background(surfaceBg).
			Padding(0, 2)

	appTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(bgColor).
			Background(primaryColor).
			Padding(0, 3).
			MarginRight(1)

	appSubtitleStyle = lipgloss.NewStyle().
				Foreground(fgColor).
				Background(surface0).
				Padding(0, 2)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(surface1).
			Background(surfaceBg).
			Padding(1, 2).
			MarginRight(1)

	activePanelStyle = lipgloss.NewStyle().
				Border(lipgloss.DoubleBorder()).
				BorderForeground(primaryColor).
				Background(surfaceBg).
				Padding(1, 2).
				MarginRight(1)

	selectedItemStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(primaryColor).
				Bold(true).
				Padding(0, 2).
				MarginLeft(1).
				MarginRight(1)

	normalItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 2).
			MarginLeft(1).
			MarginRight(1)

	hoverItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface0).
			Padding(0, 2).
			MarginLeft(1).
			MarginRight(1)

	runningStatusStyle = lipgloss.NewStyle().
				Foreground(successColor).
				Bold(true)

	stoppedStatusStyle = lipgloss.NewStyle().
				Foreground(mutedColor)

	errorStatusStyle = lipgloss.NewStyle().
				Foreground(errorColor).
				Bold(true)

	sectionHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(primaryColor).
				MarginBottom(1).
				Padding(0, 1).
				Background(surface0).
				Border(lipgloss.NormalBorder(), false, false, true, false).
				BorderForeground(primaryColor)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Underline(true).
			MarginBottom(1)

	subHeaderStyle = lipgloss.NewStyle().
			Foreground(infoColor).
			Italic(true)

	helpStyle = lipgloss.NewStyle().
			Foreground(mutedColor).
			Background(surfaceBg).
			Padding(1, 2).
			MarginTop(1)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(lavenderColor).
			Bold(true).
			Background(surface0).
			Padding(0, 1).
			MarginRight(1)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(subtleColor)

	statusBarStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface0).
			Padding(0, 2).
			Bold(true)

	statusPillStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true)

	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			MarginTop(1).
			Background(surfaceBg)

	warningBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(warningColor).
			Padding(1, 2).
			MarginTop(1).
			Background(surfaceBg)

	errorBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(errorColor).
			Padding(1, 2).
			MarginTop(1).
			Background(surfaceBg)

	modalStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Background(bgColor).
			Padding(2, 3)

	cardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(surface1).
			Background(surfaceBg).
			Padding(1, 2).
			Margin(1, 0)

	activeCardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Background(surfaceBg).
			Padding(1, 2).
			Margin(1, 0)

	badgeStyle = lipgloss.NewStyle().
			Padding(0, 1).
			Bold(true).
			Background(surface1)

	phpBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#8892BF")).
			Padding(0, 1).
			Bold(true)

	mysqlBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color("#4479A1")).
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
			Background(lipgloss.Color("#4DB33D")).
			Padding(0, 1).
			Bold(true)

	defaultBadge = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface1).
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

func getIconForService(serviceType string) string {
	icons := map[string]string{
		"php":           "󰌞",
		"mysql":         "󱆟",
		"mariadb":       "󱆟",
		"nginx":         "󰖟",
		"redis":         "󰝚",
		"postgres":      "󱆢",
		"postgresql":    "󱆢",
		"apache":        "󰖟",
		"mongodb":       "󱆦",
		"mongo":         "󱆦",
		"phpmyadmin":    "󰖶",
		"adminer":       "󱆦",
		"elasticsearch": "󰉋",
		"rabbitmq":      "󰘦",
		"memcached":     "󰘦",
		"caddy":         "󰖟",
	}
	if icon, ok := icons[serviceType]; ok {
		return icon
	}
	return "󰘦"
}

func getIconForProject(projectType string) string {
	icons := map[string]string{
		"laravel":   "󰖬",
		"nextjs":    "󰖟",
		"vue":       "󰡄",
		"django":    "󰘦",
		"express":   "󰘦",
		"fastapi":   "󰘦",
		"nuxt":      "󰖟",
		"sveltekit": "󰖟",
		"remix":     "󰖟",
		"nestjs":    "󰖟",
		"axum":      "󰘦",
		"actix":     "󰘦",
		"rocket":    "󰘦",
	}
	if icon, ok := icons[projectType]; ok {
		return icon
	}
	return "󰘦"
}

func getIconForRuntime(runtime string) string {
	icons := map[string]string{
		"php":    "󰌞",
		"node":   "󰛦",
		"python": "󰌠",
		"rust":   "󱘘",
		"bun":    "󰛦",
		"deno":   "󰛦",
		"go":     "󰟓",
	}
	if icon, ok := icons[runtime]; ok {
		return icon
	}
	return "󰘦"
}

func getStatusIcon(running bool) string {
	if running {
		return "󰀄"
	}
	return "󰀊"
}

func getLevelIcon(level string) string {
	switch level {
	case "success":
		return "󰸞"
	case "error":
		return "󰚌"
	case "warning":
		return "󰀦"
	default:
		return "󰋽"
	}
}
