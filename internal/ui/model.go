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
	logsView
	backgroundTasksView
	addServiceView
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

type backgroundTask struct {
	id        string
	name      string
	status    string // running, completed, failed
	startTime string
	message   string
}

type logEntry struct {
	timestamp string
	level     string // info, error, warning, success
	service   string
	message   string
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
	sidebarScrollOffset int
	versionCursor       int
	projectTypeCursor   int
	cleanupCursor       int
	portConflictCursor  int
	logScrollOffset     int
	detailScrollOffset  int
	mainScrollOffset    int
	selected            map[int]bool
	activePanel         panel
	currentView         view
	width               int
	height              int
	err                 error
	statusMessage       string
	serviceStatus       map[string]serviceStatus
	logs                []logEntry
	backgroundTasks     []backgroundTask
	ready               bool
	showVersionList     bool
	showProjectCreate   bool
	showCleanupDialog   bool
	showConfirmDialog   bool
	showPortConflict    bool
	availableVersions   []string
	selectedService     *config.Service
	selectedProject     *config.Project
	selectedRuntimeType string
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
		"Logs",
		"Background Tasks",
		"New Project",
		"Settings",
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
		logs:             []logEntry{},
		backgroundTasks:  []backgroundTask{},
		ready:            false,
		showVersionList:  false,
		availableVersions: []string{},
		sidebarItems:     sidebarItems,
		sidebarCursor:    0,
		logScrollOffset:  0,
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
		m.selectedRuntimeType = msg.serviceType
		m.showVersionList = true
		return m, nil

	case tea.KeyMsg:
		// Project creation modal navigation
		if m.showProjectCreate {
			switch msg.String() {
			case "esc":
				m.showProjectCreate = false
				m.projectTypeCursor = 0
			case "up", "k":
				if m.projectTypeCursor > 0 {
					m.projectTypeCursor--
				}
			case "down", "j":
				if m.projectTypeCursor < 12 { // 13 project types (0-12)
					m.projectTypeCursor++
				}
			case "enter":
				// TODO: Create project with selected type
				m.statusMessage = fmt.Sprintf("Creating project (type: %d)", m.projectTypeCursor)
				m.showProjectCreate = false
				m.projectTypeCursor = 0
			}
			return m, nil
		}

		// Cleanup dialog navigation
		if m.showCleanupDialog {
			switch msg.String() {
			case "esc":
				m.showCleanupDialog = false
				m.cleanupCursor = 0
				m.selectedService = nil
			case "up", "k":
				if m.cleanupCursor > 0 {
					m.cleanupCursor--
				}
			case "down", "j":
				if m.cleanupCursor < 3 {
					m.cleanupCursor++
				}
			case "enter":
				return m, m.performCleanup()
			}
			return m, nil
		}

		// Version list navigation
		if m.showVersionList {
			switch msg.String() {
			case "esc":
				m.showVersionList = false
				m.versionCursor = 0
				m.selectedService = nil
				m.selectedRuntimeType = ""
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
				if m.versionCursor < len(m.availableVersions) {
					selectedVersion := m.availableVersions[m.versionCursor]
					
					// Check if we're updating a service or a runtime
					if m.selectedService != nil {
						// Update service version
						m.selectedService.Version = selectedVersion
						m.statusMessage = fmt.Sprintf("Updated %s to version %s", m.selectedService.Name, selectedVersion)
					} else if m.selectedRuntimeType != "" {
						// Update runtime version
						switch m.selectedRuntimeType {
						case "php":
							m.config.Runtimes.PHP = selectedVersion
						case "node":
							m.config.Runtimes.Node = selectedVersion
						case "python":
							m.config.Runtimes.Python = selectedVersion
						case "rust":
							m.config.Runtimes.Rust = selectedVersion
						case "bun":
							m.config.Runtimes.Bun = selectedVersion
						case "deno":
							m.config.Runtimes.Deno = selectedVersion
						case "go":
							m.config.Runtimes.Go = selectedVersion
						}
						m.statusMessage = fmt.Sprintf("Updated %s runtime to version %s", m.selectedRuntimeType, selectedVersion)
					}
					
					config.SaveConfig(m.config)
					m.showVersionList = false
					m.versionCursor = 0
					m.selectedService = nil
					m.selectedRuntimeType = ""
				}
			}
			return m, nil
		}

		switch msg.String() {
		case "esc":
			// Go back to services view
			if m.currentView == addServiceView {
				m.currentView = servicesView
				m.cursor = 0
				m.activePanel = mainPanel
			}

		case "ctrl+c", "q":
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
			} else if m.activePanel == mainPanel {
				if m.currentView == logsView {
					// Scroll logs up
					maxScroll := len(m.logs) - 10
					if maxScroll < 0 {
						maxScroll = 0
					}
					if m.logScrollOffset < maxScroll {
						m.logScrollOffset++
					}
				} else if m.currentView == addServiceView {
					// Navigate service types (14 types)
					if m.cursor > 0 {
						m.cursor--
					}
				} else {
					// Generic navigation for other views
					if m.cursor > 0 {
						m.cursor--
					}
				}
			} else if m.activePanel == detailPanel {
				// Scroll detail panel up
				if m.detailScrollOffset > 0 {
					m.detailScrollOffset--
				}
			}

		case "down", "j":
			if m.activePanel == sidebarPanel && m.sidebarCursor < len(m.sidebarItems)-1 {
				m.sidebarCursor++
			} else if m.activePanel == mainPanel {
				if m.currentView == logsView {
					// Scroll logs down
					if m.logScrollOffset > 0 {
						m.logScrollOffset--
					}
				} else if m.currentView == addServiceView {
					// Navigate service types (14 types)
					if m.cursor < 13 {
						m.cursor++
					}
				} else if m.currentView == servicesView {
					// Navigate services
					if m.cursor < len(m.config.Services)-1 {
						m.cursor++
					}
				} else if m.currentView == projectsView {
					// Navigate projects
					if m.cursor < len(m.config.Projects)-1 {
						m.cursor++
					}
				} else if m.currentView == databasesView {
					// Navigate databases (assuming we have a databases list)
					// For now, just allow cursor movement
					if m.cursor < 10 { // placeholder max
						m.cursor++
					}
				} else if m.currentView == runtimesView {
					// Navigate runtimes (7 runtimes: PHP, Node, Python, Rust, Bun, Deno, Go)
					if m.cursor < 6 {
						m.cursor++
					}
				}
			} else if m.activePanel == detailPanel {
				// Scroll detail panel down
				m.detailScrollOffset++
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
				case 4: // Logs
					m.currentView = logsView
					m.activePanel = mainPanel
				case 5: // Background Tasks
					m.currentView = backgroundTasksView
					m.activePanel = mainPanel
				case 6: // Settings
					m.currentView = settingsView
					m.activePanel = mainPanel
				case 7: // Refresh
					return m, m.checkStatus()
				case 8: // Quit
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
			if m.activePanel == mainPanel {
				if m.currentView == servicesView && m.cursor < len(m.config.Services) {
					m.selectedService = &m.config.Services[m.cursor]
					return m, m.fetchVersions(m.selectedService.Type)
				} else if m.currentView == runtimesView {
					// Handle runtime version selection
					runtimeTypes := []string{"php", "node", "python", "rust", "bun", "deno", "go"}
					if m.cursor < len(runtimeTypes) {
						return m, m.fetchVersions(runtimeTypes[m.cursor])
					}
				}
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
			// Add new service or project
			if m.activePanel == mainPanel {
				if m.currentView == servicesView {
					m.currentView = addServiceView
				} else if m.currentView == projectsView {
					// Open new project modal
					m.showProjectCreate = true
					m.projectTypeCursor = 0
				}
			}
		}
	}

	return m, nil
}

func (m *model) addLog(level, service, message string) {
	entry := logEntry{
		timestamp: time.Now().Format("15:04:05"),
		level:     level,
		service:   service,
		message:   message,
	}
	m.logs = append(m.logs, entry)
	
	// Keep only last 1000 logs
	if len(m.logs) > 1000 {
		m.logs = m.logs[len(m.logs)-1000:]
	}
}

func (m *model) addBackgroundTask(name, status, message string) {
	task := backgroundTask{
		id:        fmt.Sprintf("task-%d", len(m.backgroundTasks)),
		name:      name,
		status:    status,
		startTime: time.Now().Format("15:04:05"),
		message:   message,
	}
	m.backgroundTasks = append(m.backgroundTasks, task)
}

func (m *model) updateBackgroundTask(id, status, message string) {
	for i := range m.backgroundTasks {
		if m.backgroundTasks[i].id == id {
			m.backgroundTasks[i].status = status
			m.backgroundTasks[i].message = message
			break
		}
	}
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
				taskID := fmt.Sprintf("start-%s", s.Name)
				m.addBackgroundTask(fmt.Sprintf("Starting %s", s.Name), "running", "Pulling image...")
				m.addLog("info", s.Name, "Starting service...")
				
				if err := m.docker.StartService(ctx, s); err != nil {
					m.addLog("error", s.Name, fmt.Sprintf("Failed to start: %v", err))
					m.updateBackgroundTask(taskID, "failed", err.Error())
				} else {
					m.addLog("success", s.Name, "Started successfully")
					m.updateBackgroundTask(taskID, "completed", "Service started")
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
				taskID := fmt.Sprintf("stop-%s", s.Name)
				m.addBackgroundTask(fmt.Sprintf("Stopping %s", s.Name), "running", "Stopping container...")
				m.addLog("info", s.Name, "Stopping service...")
				
				if err := m.docker.StopService(ctx, s.Name); err != nil {
					m.addLog("error", s.Name, fmt.Sprintf("Failed to stop: %v", err))
					m.updateBackgroundTask(taskID, "failed", err.Error())
				} else {
					m.addLog("success", s.Name, "Stopped successfully")
					m.updateBackgroundTask(taskID, "completed", "Service stopped")
				}
			}(service)
		}
	}
}

func (m model) View() string {
	if !m.ready {
		// Loading screen with background
		loadingText := lipgloss.NewStyle().
			Foreground(primaryColor).
			Bold(true).
			Render("✨ Initializing Lumine...")
		
		loadingBox := lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(primaryColor).
			Padding(2, 4).
			Background(bgColor).
			Render(loadingText)
		
		return lipgloss.Place(
			m.width,
			m.height,
			lipgloss.Center,
			lipgloss.Center,
			loadingBox,
			lipgloss.WithWhitespaceChars("░"),
			lipgloss.WithWhitespaceForeground(lipgloss.Color("#45475a")),
		)
	}

	// Show version selector overlay if active
	if m.showVersionList {
		return m.renderVersionSelector()
	}

	// Show project creation modal if active
	if m.showProjectCreate {
		return m.renderProjectCreateModal()
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

	// Title bar with better formatting
	titleLeft := lipgloss.NewStyle().
		Bold(true).
		Foreground(bgColor).
		Background(primaryColor).
		Padding(0, 2).
		Render("✨ LUMINE")
	
	titleRight := lipgloss.NewStyle().
		Foreground(bgColor).
		Background(primaryColor).
		Padding(0, 2).
		Render("Docker Development Manager ⚡")
	
	titlePadding := m.width - lipgloss.Width(titleLeft) - lipgloss.Width(titleRight) - 4
	if titlePadding < 0 {
		titlePadding = 0
	}
	
	titleBar := lipgloss.JoinHorizontal(
		lipgloss.Top,
		titleLeft,
		lipgloss.NewStyle().
			Background(primaryColor).
			Width(titlePadding).
			Render(""),
		titleRight,
	)
	
	s.WriteString(titleBar + "\n")

	// Calculate dimensions
	sidebarWidth := 22
	contentHeight := m.height - 4 // Title (1) + status (1) + help (1) + spacing (1)
	
	// Main layout: Fixed Sidebar | Main Content | Detail Panel
	sidebar := m.renderSidebarFixed(sidebarWidth, contentHeight)
	mainContent := m.renderMainContent()
	detailPanel := m.renderDetailPanelDynamic()

	content := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebar,
		mainContent,
		detailPanel,
	)
	s.WriteString(content + "\n")

	// Status bar
	statusBar := m.renderStatusBar()
	s.WriteString(statusBar + "\n")

	// Help
	help := m.renderHelp()
	s.WriteString(help)

	// Wrap everything with background
	fullContent := s.String()
	
	// Add background to the entire view
	styledContent := lipgloss.NewStyle().
		Background(bgColor).
		Width(m.width).
		Height(m.height).
		Render(fullContent)

	return styledContent
}

func (m *model) performCleanup() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		
		switch m.cleanupCursor {
		case 0: // Remove Container
			if m.selectedService != nil {
				m.addLog("info", m.selectedService.Name, "Removing container...")
				if err := m.docker.RemoveContainer(ctx, m.selectedService.Name); err != nil {
					m.addLog("error", m.selectedService.Name, fmt.Sprintf("Failed to remove: %v", err))
					m.statusMessage = fmt.Sprintf("Failed to remove %s", m.selectedService.Name)
				} else {
					m.addLog("success", m.selectedService.Name, "Container removed")
					m.statusMessage = fmt.Sprintf("Removed %s successfully", m.selectedService.Name)
				}
			}
			
		case 1: // Remove with Volume
			if m.selectedService != nil {
				m.addLog("info", m.selectedService.Name, "Removing container and volume...")
				// Stop and remove container
				m.docker.StopService(ctx, m.selectedService.Name)
				if err := m.docker.RemoveContainer(ctx, m.selectedService.Name); err != nil {
					m.addLog("error", m.selectedService.Name, fmt.Sprintf("Failed to remove: %v", err))
				}
				
				// Remove volume
				volumeName := fmt.Sprintf("lumine_%s_data", m.selectedService.Name)
				if err := m.docker.RemoveVolume(ctx, volumeName); err != nil {
					m.addLog("warning", m.selectedService.Name, "Volume not found or already removed")
				} else {
					m.addLog("success", m.selectedService.Name, "Container and volume removed")
					m.statusMessage = fmt.Sprintf("Removed %s with volume", m.selectedService.Name)
				}
			}
			
		case 2: // Remove All Containers
			m.addLog("info", "system", "Removing all containers...")
			if err := m.docker.RemoveAllContainers(ctx, true); err != nil {
				m.addLog("error", "system", fmt.Sprintf("Failed to remove containers: %v", err))
				m.statusMessage = "Failed to remove all containers"
			} else {
				m.addLog("success", "system", "All containers removed")
				m.statusMessage = "All containers removed"
			}
			
		case 3: // Nuclear Cleanup
			m.addLog("warning", "system", "Starting nuclear cleanup...")
			opts := docker.CleanupOptions{
				RemoveContainers: true,
				RemoveVolumes:    true,
				RemoveNetworks:   true,
				Force:            true,
			}
			
			if err := m.docker.Cleanup(ctx, opts); err != nil {
				m.addLog("error", "system", fmt.Sprintf("Cleanup failed: %v", err))
				m.statusMessage = "Cleanup failed"
			} else {
				m.addLog("success", "system", "Nuclear cleanup completed - all resources removed")
				m.statusMessage = "Complete cleanup successful"
			}
		}
		
		m.showCleanupDialog = false
		m.cleanupCursor = 0
		m.selectedService = nil
		
		return m.checkStatus()()
	}
}
