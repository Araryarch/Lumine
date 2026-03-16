package gui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Araryarch/Lumine/pkg/commands"
	"github.com/Araryarch/Lumine/pkg/config"
	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/gui/types"
	"github.com/Araryarch/Lumine/pkg/i18n"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/tasks"
	throttle "github.com/boz/go-throttle"
	"github.com/jesseduffield/gocui"
	lcUtils "github.com/jesseduffield/lazycore/pkg/utils"
	"github.com/sasha-s/go-deadlock"
	"github.com/sirupsen/logrus"
)

// Gui wraps the gocui Gui object which handles rendering and events
type Gui struct {
	g             *gocui.Gui
	Log           *logrus.Entry
	OSCommand     *commands.OSCommand
	State         guiState
	Config        *config.AppConfig
	Tr            *i18n.TranslationSet
	statusManager *statusManager
	taskManager   *tasks.TaskManager
	ErrorChan     chan error
	Views         Views
	Orchestrator  *lumine.Orchestrator
	DockerCommand *DockerCommand

	// if we've suspended the gui (e.g. because we've switched to a subprocess)
	// we typically want to pause some things that are running like background
	// file refreshes
	PauseBackgroundThreads bool

	Mutexes

	Panels Panels
}

type Panels struct {
	LumineDocker    *panels.SideListPanel[*DockerControl]
	LumineServers   *panels.SideListPanel[*lumine.Service]
	LumineLanguages *panels.SideListPanel[*lumine.Service]
	LumineFiles     *panels.SideListPanel[*lumine.Service]
	LumineProjects  *panels.SideListPanel[*lumine.Project]
	LumineDatabases *panels.SideListPanel[*lumine.Service]
	Menu            *panels.SideListPanel[*types.MenuItem]
}

type DockerCommand struct {
	InDockerComposeProject bool
}

type Mutexes struct {
	SubprocessMutex deadlock.Mutex
	ViewStackMutex  deadlock.Mutex
}

type mainPanelState struct {
	// ObjectKey tells us what context we are in. For example, if we are looking at the logs of a particular service in the services panel this key might be 'services-<service id>-logs'. The key is made so that if something changes which might require us to re-run the logs command or run a different command, the key will be different, and we'll then know to do whatever is required. Object key probably isn't the best name for this but Context is already used to refer to tabs. Maybe I should just call them tabs.
	ObjectKey string
}

type panelStates struct {
	Main *mainPanelState
}

type guiState struct {
	// the names of views in the current focus stack (last item is the current view)
	ViewStack        []string
	Platform         commands.Platform
	Panels           *panelStates
	SubProcessOutput string

	// if true, we show containers with an 'exited' status in the containers panel
	ShowExitedContainers bool

	ScreenMode WindowMaximisation

	// Maintains the state of manual filtering i.e. typing in a substring
	// to filter on in the current panel.
	Filter filterState
}

type filterState struct {
	// If true then we're either currently inside the filter view
	// or we've committed the filter and we're back in the list view
	active bool
	// The panel that we're filtering.
	panel panels.ISideListPanel
	// The string that we're filtering on
	needle string
}

// screen sizing determines how much space your selected window takes up (window
// as in panel, not your terminal's window). Sometimes you want a bit more space
// to see the contents of a panel, and this keeps track of how much maximisation
// you've set
type WindowMaximisation int

const (
	SCREEN_NORMAL WindowMaximisation = iota
	SCREEN_HALF
	SCREEN_FULL
)

func getScreenMode(config *config.AppConfig) WindowMaximisation {
	switch config.UserConfig.Gui.ScreenMode {
	case "normal":
		return SCREEN_NORMAL
	case "half":
		return SCREEN_HALF
	case "fullscreen":
		return SCREEN_FULL
	default:
		return SCREEN_NORMAL
	}
}

// NewGui builds a new gui handler
func NewGui(log *logrus.Entry, oSCommand *commands.OSCommand, tr *i18n.TranslationSet, config *config.AppConfig, errorChan chan error) (*Gui, error) {
	initialState := guiState{
		Platform: *oSCommand.Platform,
		Panels: &panelStates{
			Main: &mainPanelState{
				ObjectKey: "",
			},
		},
		ViewStack: []string{},

		ShowExitedContainers: true,
		ScreenMode:           getScreenMode(config),
	}

	// Initialize Lumine orchestrator
	orchestrator, err := lumine.NewOrchestrator()
	if err != nil {
		log.Warnf("Failed to initialize Lumine orchestrator: %v", err)
		// Continue without Lumine features
	}

	gui := &Gui{
		Log:           log,
		OSCommand:     oSCommand,
		State:         initialState,
		Config:        config,
		Tr:            tr,
		statusManager: &statusManager{},
		taskManager:   tasks.NewTaskManager(log, tr),
		ErrorChan:     errorChan,
		Orchestrator:  orchestrator,
		DockerCommand: &DockerCommand{
			InDockerComposeProject: false,
		},
	}

	deadlock.Opts.Disable = !gui.Config.Debug
	deadlock.Opts.DeadlockTimeout = 10 * time.Second

	return gui, nil
}

func (gui *Gui) renderGlobalOptions() error {
	return gui.renderOptionsMap(map[string]string{
		"PgUp/PgDn": gui.Tr.Scroll,
		"← → ↑ ↓":   gui.Tr.Navigate,
		"q":         gui.Tr.Quit,
		"m":         gui.Tr.Menu,
		"w":         "Setup Wizard",
		"1-6":       "Switch Panel",
	})
}

func (gui *Gui) goEvery(interval time.Duration, function func() error) {
	_ = function() // time.Tick doesn't run immediately so we'll do that here // TODO: maybe change
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			if !gui.PauseBackgroundThreads {
				_ = function()
			}
		}
	}()
}

// Run setup the gui with keybindings and start the mainloop
func (gui *Gui) Run() error {
	// closing our task manager which in turn closes the current task if there is any, so we aren't leaving processes lying around after closing lazydocker
	defer gui.taskManager.Close()

	// Close Lumine orchestrator on exit
	if gui.Orchestrator != nil {
		defer gui.Orchestrator.Close()
	}

	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:       gocui.OutputTrue,
		RuneReplacements: map[rune]string{},
	})
	if err != nil {
		return err
	}
	defer g.Close()

	// forgive the double-negative, this is because of my yaml `omitempty` woes
	if !gui.Config.UserConfig.Gui.IgnoreMouseEvents {
		g.Mouse = true
	}

	gui.g = g // TODO: always use gui.g rather than passing g around everywhere

	// if the deadlock package wants to report a deadlock, we first need to
	// close the gui so that we can actually read what it prints.
	deadlock.Opts.LogBuf = lcUtils.NewOnceWriter(os.Stderr, func() {
		gui.g.Close()
	})

	if err := gui.SetColorScheme(); err != nil {
		return err
	}

	throttledRefresh := throttle.ThrottleFunc(time.Millisecond*50, true, gui.refresh)
	defer throttledRefresh.Stop()

	go func() {
		for err := range gui.ErrorChan {
			if err == nil {
				continue
			}
			if strings.Contains(err.Error(), "No such container") {
				// this happens all the time when e.g. restarting containers so we won't worry about it
				gui.Log.Warn(err)
				continue
			}
			_ = gui.createErrorPanel(err.Error())
		}
	}()

	g.SetManager(gocui.ManagerFunc(gui.layout), gocui.ManagerFunc(gui.getFocusLayout()))

	if err := gui.createAllViews(); err != nil {
		return err
	}
	if err := gui.setInitialViewContent(); err != nil {
		return err
	}

	// TODO: see if we can avoid the circular dependency
	gui.setPanels()

	if err = gui.keybindings(g); err != nil {
		return err
	}

	if gui.g.CurrentView() == nil {
		viewName := gui.initiallyFocusedViewName()
		view, err := gui.g.View(viewName)
		if err != nil {
			return err
		}

		if err := gui.switchFocus(view); err != nil {
			return err
		}
	}
	
	// Show first setup wizard on first run
	if gui.Orchestrator != nil && gui.Orchestrator.ConfigManager.IsFirstRun() {
		go func() {
			// Wait a bit for GUI to fully initialize
			time.Sleep(500 * time.Millisecond)
			gui.g.Update(func(g *gocui.Gui) error {
				return gui.ShowFirstSetupWizard()
			})
		}()
	}

	ctx, finish := context.WithCancel(context.Background())
	defer finish()

	go gui.listenForEvents(ctx, throttledRefresh.Trigger)

	go func() {
		throttledRefresh.Trigger()

		gui.goEvery(time.Millisecond*30, gui.reRenderMain)

		// Lumine refresh cycles
		if gui.Orchestrator != nil {
			gui.goEvery(time.Millisecond*2000, gui.refreshDockerControl)
			gui.goEvery(time.Millisecond*2000, gui.refreshLumineServers)
			gui.goEvery(time.Millisecond*2000, gui.refreshLumineLanguages)
			gui.goEvery(time.Millisecond*2000, gui.refreshLumineFiles)
			gui.goEvery(time.Millisecond*5000, gui.refreshLumineProjects)
			gui.goEvery(time.Millisecond*3000, gui.refreshLumineDatabases)
			gui.goEvery(time.Millisecond*1000, gui.refreshNotifications)
		}
	}()

	err = g.MainLoop()
	if err == gocui.ErrQuit {
		return nil
	}
	return err
}

func (gui *Gui) setPanels() {
	gui.Panels = Panels{
		Menu: gui.getMenuPanel(),
	}

	// Initialize Lumine panels
	if gui.Orchestrator != nil {
		gui.Panels.LumineDocker = gui.getLumineDockerPanel()
		gui.Panels.LumineServers = gui.getLumineServersPanel()
		gui.Panels.LumineLanguages = gui.getLumineLanguagesPanel()
		gui.Panels.LumineFiles = gui.getLumineFilesPanel()
		gui.Panels.LumineProjects = gui.getLumineProjectsPanel()
		gui.Panels.LumineDatabases = gui.getLumineDatabasesPanel()
	}
}

func (gui *Gui) refresh() {
	// Refresh Lumine panels
	if gui.Orchestrator != nil {
		go func() {
			if err := gui.refreshDockerControl(); err != nil {
				gui.Log.Error(err)
			}
		}()
		go func() {
			if err := gui.refreshLumineServers(); err != nil {
				gui.Log.Error(err)
			}
		}()
		go func() {
			if err := gui.refreshLumineLanguages(); err != nil {
				gui.Log.Error(err)
			}
		}()
		go func() {
			if err := gui.refreshLumineFiles(); err != nil {
				gui.Log.Error(err)
			}
		}()
		go func() {
			if err := gui.refreshLumineProjects(); err != nil {
				gui.Log.Error(err)
			}
		}()
		go func() {
			if err := gui.refreshLumineDatabases(); err != nil {
				gui.Log.Error(err)
			}
		}()
	}
}

func (gui *Gui) listenForEvents(ctx context.Context, refresh func()) {
	// Lumine doesn't need Docker event listening
	// Just refresh periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			refresh()
		}
	}
}

// checkForContextChange runs the currently focused panel's 'select' function, simulating the current item having just been selected. This will then trigger a check to see if anything's changed (e.g. a service has a new container) and if so, the appropriate code will run. For example, if you're reading logs from a service and all of a sudden its container changes, this will trigger the 'select' function, which will work out that the context is not different because of the new container, and then it will re-attempt to get the logs, this time for the correct container. This 'context' is stored in the main panel's ObjectKey. I'm using the term 'context' here more broadly than just the different tabs you can view in a panel.
func (gui *Gui) checkForContextChange() error {
	return gui.newLineFocused(gui.g.CurrentView())
}

func (gui *Gui) reRenderMain() error {
	mainView := gui.Views.Main
	if mainView == nil {
		return nil
	}
	if mainView.IsTainted() {
		gui.g.Update(func(g *gocui.Gui) error {
			return nil
		})
	}
	return nil
}

func (gui *Gui) quit(g *gocui.Gui, v *gocui.View) error {
	if gui.Config.UserConfig.ConfirmOnQuit {
		return gui.createConfirmationPanel("", gui.Tr.ConfirmQuit, func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}, nil)
	}
	return gocui.ErrQuit
}

// this handler is executed when we press escape when there is only one view
// on the stack.
func (gui *Gui) escape() error {
	if gui.State.Filter.active {
		return gui.clearFilter()
	}

	return nil
}

func (gui *Gui) handleAppInfo(g *gocui.Gui, v *gocui.View) error {
	if !gui.g.Mouse {
		return nil
	}

	cx, _ := v.Cursor()
	if cx > len("Lumine") {
		return nil
	}
	
	// Show app info or open project page
	return gui.createConfirmationPanel("Lumine - Local Development Environment Manager", 
		"A Docker-based development environment manager\nVersion: "+gui.Config.Version, 
		func(g *gocui.Gui, v *gocui.View) error {
			return nil
		}, nil)
}

func (gui *Gui) editFile(filename string) error {
	cmd, err := gui.OSCommand.EditFile(filename)
	if err != nil {
		return gui.createErrorPanel(err.Error())
	}

	return gui.runSubprocess(cmd)
}

func (gui *Gui) openFile(filename string) error {
	if err := gui.OSCommand.OpenFile(filename); err != nil {
		return gui.createErrorPanel(err.Error())
	}
	return nil
}

func (gui *Gui) handleCustomCommand(g *gocui.Gui, v *gocui.View) error {
	return gui.createPromptPanel(gui.Tr.CustomCommandTitle, func(g *gocui.Gui, v *gocui.View) error {
		command := gui.trimmedContent(v)
		return gui.runSubprocess(gui.OSCommand.RunCustomCommand(command))
	})
}

func (gui *Gui) ShouldRefresh(key string) bool {
	if gui.State.Panels.Main.ObjectKey == key {
		return false
	}

	gui.State.Panels.Main.ObjectKey = key
	return true
}

func (gui *Gui) FilterString(view *gocui.View) string {
	return ""
}

func (gui *Gui) initiallyFocusedViewName() string {
	return "lumineDocker"
}

func (gui *Gui) IgnoreStrings() []string {
	return gui.Config.UserConfig.Ignore
}

func (gui *Gui) Update(f func() error) {
	gui.g.Update(func(*gocui.Gui) error { return f() })
}

func (gui *Gui) promptToReturn() {
	fmt.Print("\nPress return to continue...")
	fmt.Scanln()
}

// this is used by our cheatsheet code to generate keybindings. We need some views
// and panels to exist for us to know what keybindings there are, so we invoke
// gocui in headless mode and create them.
func (gui *Gui) SetupFakeGui() {
	g, err := gocui.NewGui(gocui.NewGuiOpts{
		OutputMode:       gocui.OutputTrue,
		RuneReplacements: map[rune]string{},
		Headless:         true,
	})
	if err != nil {
		panic(err)
	}
	gui.g = g
	defer g.Close()
	if err := gui.createAllViews(); err != nil {
		panic(err)
	}

	gui.setPanels()
}
