package ui

import "github.com/charmbracelet/lipgloss"

var (
	bgColor   = lipgloss.Color("#080808")
	surfaceBg = lipgloss.Color("#0c0c0c")
	surface0  = lipgloss.Color("#141414")
	surface1  = lipgloss.Color("#1a1a1a")
	surface2  = lipgloss.Color("#222222")
	surface3  = lipgloss.Color("#2a2a2a")

	fgColor  = lipgloss.Color("#e4e4e4")
	fgMuted  = lipgloss.Color("#a0a0a0")
	fgSubtle = lipgloss.Color("#666666")
	fgDim    = lipgloss.Color("#444444")

	primaryColor = lipgloss.Color("#d4d4d4")
	primaryBold  = lipgloss.Color("#ffffff")

	accentColor  = lipgloss.Color("#c0c0c0")
	successColor = lipgloss.Color("#b0b0b0")
	warningColor = lipgloss.Color("#909090")
	infoColor    = lipgloss.Color("#808080")
	errorColor   = lipgloss.Color("#ffffff")

	borderColor     = lipgloss.Color("#2a2a2a")
	borderHighlight = lipgloss.Color("#404040")

	dividerColor = lipgloss.Color("#1a1a1a")
	dividerStyle = lipgloss.NewStyle().Foreground(dividerColor)
)

var (
	brandStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryBold).
			Background(surface2).
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(surface3).
			BorderBackground(surface2)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1).
			Background(surface1).
			Width(0)

	sidebarHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(fgMuted).
				Padding(1, 0, 0, 0).
				MarginBottom(1)

	menuItemStyle = lipgloss.NewStyle().
			Foreground(fgMuted).
			Padding(0, 1)

	menuItemActiveStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(fgMuted).
				Bold(true).
				Padding(0, 1)

	menuIconStyle = lipgloss.NewStyle().
			Foreground(fgSubtle).
			Render

	menuIconActiveStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Bold(true).
				Render

	panelBorderStyle  = lipgloss.NormalBorder()
	panelActiveBorder = lipgloss.ThickBorder()

	panelStyle = lipgloss.NewStyle().
			Border(panelBorderStyle).
			BorderForeground(borderColor).
			Background(surfaceBg).
			Padding(1, 1)

	panelActiveStyle = lipgloss.NewStyle().
				Border(panelActiveBorder).
				BorderForeground(borderHighlight).
				Background(surfaceBg).
				Padding(1, 1)

	panelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(primaryColor).
			Padding(0, 1).
			Background(surface0)

	listItemStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Padding(0, 1)

	listItemSelectedStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(fgMuted).
				Bold(true).
				Padding(0, 1)

	listItemHoverStyle = lipgloss.NewStyle().
				Foreground(fgColor).
				Background(surface1).
				Padding(0, 1)

	statusRunningStyle = lipgloss.NewStyle().
				Foreground(primaryBold).
				Bold(true)

	statusStoppedStyle = lipgloss.NewStyle().
				Foreground(fgSubtle)

	tagStyle = lipgloss.NewStyle().
			Foreground(bgColor).
			Background(fgMuted).
			Padding(0, 1).
			Bold(true)

	tagAltStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface2).
			Padding(0, 1)

	tagMutedStyle = lipgloss.NewStyle().
			Foreground(fgMuted).
			Background(surface1).
			Padding(0, 1)

	infoLabelStyle = lipgloss.NewStyle().
			Foreground(fgSubtle)

	infoValueStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Bold(true)

	infoBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(borderColor).
			Padding(1, 2).
			Background(surfaceBg)

	infoBoxActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(borderHighlight).
				Padding(1, 2).
				Background(surfaceBg)

	dividerStyleBox = lipgloss.NewStyle().
			Foreground(dividerColor).
			Render

	statusBarStyle = lipgloss.NewStyle().
			Background(surface0).
			Foreground(fgMuted).
			Padding(0, 2).
			Bold(true)

	statusPillStyle = lipgloss.NewStyle().
			Padding(0, 2).
			Bold(true)

	helpKeyStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface2).
			Padding(0, 1).
			MarginRight(1)

	helpDescStyle = lipgloss.NewStyle().
			Foreground(fgSubtle)

	emptyStateStyle = lipgloss.NewStyle().
			Foreground(fgSubtle).
			Italic(true).
			Align(lipgloss.Center)

	scrollIndicatorStyle = lipgloss.NewStyle().
				Foreground(fgSubtle).
				Italic(true)

	buttonStyle = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface2).
			Padding(0, 2).
			Border(lipgloss.NormalBorder()).
			BorderForeground(borderColor)

	buttonActiveStyle = lipgloss.NewStyle().
				Foreground(bgColor).
				Background(fgMuted).
				Bold(true).
				Padding(0, 2).
				Border(lipgloss.NormalBorder()).
				BorderForeground(fgMuted)

	phpBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#999999")).
			Padding(0, 1).
			Bold(true)

	mysqlBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#888888")).
			Padding(0, 1).
			Bold(true)

	nginxBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#777777")).
			Padding(0, 1).
			Bold(true)

	redisBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#666666")).
			Padding(0, 1).
			Bold(true)

	postgresBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#aaaaaa")).
			Padding(0, 1).
			Bold(true)

	mariadbBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#999999")).
			Padding(0, 1).
			Bold(true)

	apacheBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#888888")).
			Padding(0, 1).
			Bold(true)

	mongoBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#777777")).
			Padding(0, 1).
			Bold(true)

	defaultBadge = lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface2).
			Padding(0, 1).
			Bold(true)

	successBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#b0b0b0")).
			Padding(0, 1).
			Bold(true)

	errorBadge = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#ffffff")).
			Padding(0, 1).
			Bold(true)

	mutedBadge = lipgloss.NewStyle().
			Foreground(fgSubtle).
			Background(surface1).
			Padding(0, 1)
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

func getProjectBadge(projectType string) lipgloss.Style {
	switch projectType {
	case "laravel":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#aaaaaa")).
			Padding(0, 1).
			Bold(true)
	case "nextjs":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#999999")).
			Padding(0, 1).
			Bold(true)
	case "vue":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#888888")).
			Padding(0, 1).
			Bold(true)
	case "django":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#777777")).
			Padding(0, 1).
			Bold(true)
	case "express":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#666666")).
			Padding(0, 1).
			Bold(true)
	case "axum", "actix", "rocket":
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#080808")).
			Background(lipgloss.Color("#555555")).
			Padding(0, 1).
			Bold(true)
	default:
		return lipgloss.NewStyle().
			Foreground(fgColor).
			Background(surface2).
			Padding(0, 1).
			Bold(true)
	}
}

func getIconForService(serviceType string) string {
	icons := map[string]string{
		"php":           "PHP",
		"mysql":         "MySQL",
		"mariadb":       "MariaDB",
		"nginx":         "Nginx",
		"redis":         "Redis",
		"postgres":      "PGSQL",
		"postgresql":    "PGSQL",
		"apache":        "Apache",
		"mongodb":       "Mongo",
		"mongo":         "Mongo",
		"phpmyadmin":    "PMA",
		"adminer":       "Admin",
		"elasticsearch": "ES",
		"rabbitmq":      "RMQ",
		"memcached":     "MC",
		"caddy":         "Caddy",
	}
	if icon, ok := icons[serviceType]; ok {
		return icon
	}
	return "SVC"
}

func getIconForProject(projectType string) string {
	icons := map[string]string{
		"laravel":   "Laravel",
		"nextjs":    "Next.js",
		"vue":       "Vue",
		"django":    "Django",
		"express":   "Express",
		"fastapi":   "FastAPI",
		"nuxt":      "Nuxt",
		"sveltekit": "Svelte",
		"remix":     "Remix",
		"nestjs":    "NestJS",
		"axum":      "Axum",
		"actix":     "Actix",
		"rocket":    "Rocket",
	}
	if icon, ok := icons[projectType]; ok {
		return icon
	}
	return "App"
}

func getIconForRuntime(runtime string) string {
	icons := map[string]string{
		"php":    "PHP",
		"node":   "Node",
		"python": "Python",
		"rust":   "Rust",
		"bun":    "Bun",
		"deno":   "Deno",
		"go":     "Go",
	}
	if icon, ok := icons[runtime]; ok {
		return icon
	}
	return runtime
}

func getStatusIcon(running bool) string {
	if running {
		return "●"
	}
	return "○"
}

func getLevelIcon(level string) string {
	switch level {
	case "success":
		return "✓"
	case "error":
		return "✗"
	case "warning":
		return "⚠"
	default:
		return "ℹ"
	}
}
