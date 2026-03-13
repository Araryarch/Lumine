package ui

import (
	"context"
	"fmt"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"lumine/internal/config"
	"lumine/internal/docker"
)

type view int

const (
	servicesView view = iota
	projectsView
	databasesView
	runtimesView
	addServiceView
	addProjectView
	settingsView
)

type panel int

const (
	sidebarPanel panel = iota
	mainPanel
	detailPanel
)

type serviceStatus struct {
	name        string
	running     bool
	containerId string
}

type statusMsg map[string]serviceStatus

type versionListMsg struct {
	serviceType string
	versions    []string
}

type projectCreationMsg struct {
	success bool
	message string
}

type portConflictInfo struct {
	Port         int
	ServiceName  string
	Alternatives []int
}

type model struct {
	config              *config.Config
	docker              *docker.Manager
	cursor              int
	sidebarCursor       int
	versionCursor       int
	projectTypeCursor   int
	cleanupCursor       int
	portConflictCursor  int
	selected            map[int]bool
	activePanel         panel
	currentView         view
	width               int
	height              int
	err                 error
	statusMessage       string
	serviceStatus       map[string]serviceStatus
	logs                []string
	ready               bool
	showVersionList     bool
	showProjectCreate   bool
	showCleanupDialog   bool
	showConfirmDialog   bool
	showPortConflict    bool
	availableVersions   []string
	selectedService     *config.Service
	selectedProject     *config.Project
	portConflict        *portConflictInfo
	searchQuery         string
	sidebarItems        []string
	newProjectName      string
	newProjectPath      string
	confirmInput        string
	customPortInput     string
}

func NewModel() model {
	cfg, _ := config.LoadConfig()
	dockerMgr, _ := docker.NewManager()

	sidebarItems := []string{
		"Services",
		"Projects",
		"Databases",
		"Runtimes",
		"New Project",
		"Logs",
		"Refresh",
		"Quit",
	}

	return model{
		config:           cfg,
		docker:           dockerMgr,
		selected:         make(map[int]bool),
		activePanel:      mainPanel,
		currentView:      servicesView,
		serviceStatus:    make(map[string]serviceStatus),
		logs:             []string{},
		ready:            false,
		showVersionList:  false,
		availableVersions: []string{},
		sidebarItems:     sidebarItems,
		sidebarCursor:    0,
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		m.checkStatus(),
		tea.EnterAltScreen,
	)
}

func (m model) checkStatus() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		containers, err := m.docker.ListContainers(ctx)
		if err != nil {
			return statusMsg{}
		}

		status := make(map[string]serviceStatus)
		for _, container := range containers {
			for _, name := range container.Names {
				if strings.HasPrefix(name, "/lumine-") {
					serviceName := strings.TrimPrefix(name, "/lumine-")
					status[serviceName] = serviceStatus{
						name:        serviceName,
						running:     container.State == "running",
						containerId: container.ID,
					}
				}
			}
		}
		return statusMsg(status)
	}
}

func tickCmd() tea.Cmd {
	return tea.Tick(time.Second*2, func(t time.Time) tea.Msg {
		return t
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.ready = true
		return m, nil

	case statusMsg:
		m.serviceStatus = msg
		return m, tickCmd()

	case time.Time:
		return m, m.checkStatus()

	case versionListMsg:
		m.availableVersions = msg.versions
		m.showVersionList = true
		return m, nil

	case tea.KeyMsg:
		// Version list navigation
		if m.showVersionList {
			switch msg.String() {
			case "esc":
				m.showVersionList = false
				m.versionCursor = 0
			case "up", "k":
				if m.versionCursor > 0 {
					m.versionCursor--
				}
			case "down", "j":
				if m.versionCursor < len(m.availableVersions)-1 {
					m.versionCursor++
				}
			case "enter":
				// Apply selected version
				if m.selectedService != nil && m.versionCursor < len(m.availableVersions) {
					m.selectedService.Version = m.availableVersions[m.versionCursor]
					config.SaveConfig(m.config)
					m.statusMessage = fmt.Sprintf("Updated %s to version %s", m.selectedService.Name, m.selectedService.Version)
					m.showVersionList = false
					m.versionCursor = 0
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab":
			m.activePanel = (m.activePanel + 1) % 3

		case "shift+tab":
			if m.activePanel == 0 {
				m.activePanel = 2
			} else {
				m.activePanel--
			}

		case "h", "left":
			m.activePanel = sidebarPanel

		case "l", "right":
			if m.activePanel == sidebarPanel {
				m.activePanel = mainPanel
			} else if m.activePanel == mainPanel {
				m.activePanel = detailPanel
			}

		case "up", "k":
			if m.activePanel == sidebarPanel && m.sidebarCursor > 0 {
				m.sidebarCursor--
			} else if m.activePanel == mainPanel && m.cursor > 0 {
				m.cursor--
			}

		case "down", "j":
			if m.activePanel == sidebarPanel && m.sidebarCursor < len(m.sidebarItems)-1 {
				m.sidebarCursor++
			} else if m.activePanel == mainPanel && m.cursor < len(m.config.Services)-1 {
				m.cursor++
			}

		case "enter":
			if m.activePanel == sidebarPanel {
				switch m.sidebarCursor {
				case 0: // Services
					m.currentView = servicesView
					m.activePanel = mainPanel
				case 1: // Projects
					m.currentView = projectsView
					m.activePanel = mainPanel
				case 2: // Databases
					m.currentView = databasesView
					m.activePanel = mainPanel
				case 3: // Runtimes
					m.currentView = runtimesView
					m.activePanel = mainPanel
				case 4: // New Project
					m.currentView = addProjectView
					m.activePanel = mainPanel
				case 5: // Logs
					m.activePanel = detailPanel
				case 6: // Refresh
					return m, m.checkStatus()
				case 7: // Quit
					return m, tea.Quit
				}
			} else if m.activePanel == mainPanel {
				if m.currentView == servicesView {
					// Start service
					m.startServices()
					m.statusMessage = "Starting services..."
				} else if m.currentView == projectsView {
					// Start project
					m.statusMessage = "Starting project..."
				} else if m.currentView == databasesView {
					// Start database
					m.statusMessage = "Starting database..."
				}
			}

		case "s":
			// Start selected or current service
			if m.activePanel == mainPanel {
				m.startServices()
				m.statusMessage = "Starting services..."
			}

		case "x":
			// Stop selected or current service
			if m.activePanel == mainPanel {
				m.stopServices()
				m.statusMessage = "Stopping services..."
			}

		case "delete", "backspace":
			// Show cleanup dialog
			if m.activePanel == mainPanel && (m.currentView == servicesView || m.currentView == databasesView) {
				if m.cursor < len(m.config.Services) {
					m.selectedService = &m.config.Services[m.cursor]
					m.showCleanupDialog = true
				}
			}

		case "r":
			// Restart services
			if m.activePanel == mainPanel {
				m.stopServices()
				time.Sleep(time.Second)
				m.startServices()
				m.statusMessage = "Restarting services..."
			}

		case "v":
			// Show version selector
			if m.activePanel == mainPanel && m.cursor < len(m.config.Services) {
				m.selectedService = &m.config.Services[m.cursor]
				return m, m.fetchVersions(m.selectedService.Type)
			}

		case " ":
			if m.activePanel == mainPanel {
				m.selected[m.cursor] = !m.selected[m.cursor]
			}

		case "c":
			// Show cleanup dialog for all containers
			if m.activePanel == mainPanel {
				m.showCleanupDialog = true
			}

		case "a":
			// Select all
			if m.activePanel == mainPanel {
				for i := range m.config.Services {
					m.selected[i] = true
				}
			}

		case "d":
			// Deselect all
			if m.activePanel == mainPanel {
				m.selected = make(map[int]bool)
			}

		case "n":
			// Add new service
			m.currentView = addServiceView
		}
	}

	return m, nil
}

func (m *model) startServices() {
	hasSelection := false
	for i := range m.selected {
		if m.selected[i] {
			hasSelection = true
			break
		}
	}

	if !hasSelection {
		m.selected[m.cursor] = true
	}

	for i, service := range m.config.Services {
		if m.selected[i] {
			go func(s config.Service) {
				ctx := context.Background()
				if err := m.docker.StartService(ctx, s); err != nil {
					m.logs = append(m.logs, fmt.Sprintf("Error starting %s: %v", s.Name, err))
				} else {
					m.logs = append(m.logs, fmt.Sprintf("Started %s successfully", s.Name))
				}
			}(service)
		}
	}
}

func (m *model) stopServices() {
	hasSelection := false
	for i := range m.selected {
		if m.selected[i] {
			hasSelection = true
			break
		}
	}

	if !hasSelection {
		m.selected[m.cursor] = true
	}

	for i, service := range m.config.Services {
		if m.selected[i] {
			go func(s config.Service) {
				ctx := context.Background()
				if err := m.docker.StopService(ctx, s.Name); err != nil {
					m.logs = append(m.logs, fmt.Sprintf("Error stopping %s: %v", s.Name, err))
				} else {
					m.logs = append(m.logs, fmt.Sprintf("Stopped %s successfully", s.Name))
				}
			}(service)
		}
	}
}

func (m model) View() string {
	if !m.ready {
		return "Initializing..."
	}

	// Show version selector overlay if active
	if m.showVersionList {
		return m.renderVersionSelector()
	}

	// Show cleanup dialog if active
	if m.showCleanupDialog {
		return m.renderCleanupDialog()
	}

	// Show confirm dialog if active
	if m.showConfirmDialog {
		return m.renderConfirmDialog("This action cannot be undone!")
	}

	// Show port conflict dialog if active
	if m.showPortConflict {
		return m.renderPortConflictDialog()
	}

	var s strings.Builder

	// Title bar
	title := titleStyle.Width(m.width - 2).Render("LUMINE - Docker Development Environment Manager")
	s.WriteString(title + "\n")

	// Main layout: Sidebar | Main Content | Detail Panel
	sidebar := m.renderSidebar()
	mainContent := m.renderMainContent()
	detailPanel := m.renderDetailPanel()

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		mainContent,
		detailPanel,
	)
	s.WriteString(content + "\n")

	// Status bar
	statusBar := m.renderStatusBar()
	s.WriteString(statusBar)

	// Help
	help := m.renderHelp()
	s.WriteString(help)

	return baseStyle.Render(s.String())
}
