package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

// Binding - a keybinding mapping a key and modifier to a handler
type Binding struct {
	ViewName    string
	Handler     func(*gocui.Gui, *gocui.View) error
	Key         interface{}
	Modifier    gocui.Modifier
	Description string
}

// GetKey returns the key as a string
func (b *Binding) GetKey() string {
	key := 0

	switch b.Key.(type) {
	case rune:
		key = int(b.Key.(rune))
	case gocui.Key:
		key = int(b.Key.(gocui.Key))
	}

	// special keys
	switch key {
	case 27:
		return "esc"
	case 13:
		return "enter"
	case 32:
		return "space"
	case 65514:
		return "►"
	case 65515:
		return "◄"
	case 65517:
		return "▲"
	case 65516:
		return "▼"
	case 65508:
		return "PgUp"
	case 65507:
		return "PgDn"
	}

	return fmt.Sprintf("%c", key)
}

// GetInitialKeybindings returns all keybindings
func (gui *Gui) GetInitialKeybindings() []*Binding {
	bindings := []*Binding{
		// Global bindings
		{
			ViewName: "",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.escape),
		},
		{
			ViewName: "",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlC,
			Modifier: gocui.ModNone,
			Handler:  gui.quit,
		},
		{
			ViewName: "",
			Key:      gocui.KeyPgup,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.scrollUpMain),
		},
		{
			ViewName: "",
			Key:      gocui.KeyPgdn,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.scrollDownMain),
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlU,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.scrollUpMain),
		},
		{
			ViewName: "",
			Key:      gocui.KeyCtrlD,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.scrollDownMain),
		},
		{
			ViewName: "",
			Key:      gocui.KeyEnd,
			Modifier: gocui.ModNone,
			Handler:  gui.autoScrollMain,
		},
		{
			ViewName: "",
			Key:      gocui.KeyHome,
			Modifier: gocui.ModNone,
			Handler:  gui.jumpToTopMain,
		},
		{
			ViewName: "",
			Key:      'x',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		},
		{
			ViewName: "",
			Key:      '?',
			Modifier: gocui.ModNone,
			Handler:  gui.handleCreateOptionsMenu,
		},
		// Menu bindings
		{
			ViewName: "menu",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuClose),
		},
		{
			ViewName: "menu",
			Key:      'q',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuClose),
		},
		{
			ViewName: "menu",
			Key:      ' ',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuPress),
		},
		{
			ViewName: "menu",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuPress),
		},
		{
			ViewName: "menu",
			Key:      'y',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.handleMenuPress),
		},
		// Main panel bindings
		{
			ViewName:    "main",
			Key:         gocui.KeyEsc,
			Modifier:    gocui.ModNone,
			Handler:     gui.handleExitMain,
			Description: "Return",
		},
		{
			ViewName: "main",
			Key:      gocui.KeyArrowLeft,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollLeftMain,
		},
		{
			ViewName: "main",
			Key:      gocui.KeyArrowRight,
			Modifier: gocui.ModNone,
			Handler:  gui.scrollRightMain,
		},
		{
			ViewName: "main",
			Key:      'h',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollLeftMain,
		},
		{
			ViewName: "main",
			Key:      'l',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollRightMain,
		},
		// Filter bindings
		{
			ViewName: "filter",
			Key:      gocui.KeyEnter,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.commitFilter),
		},
		{
			ViewName: "filter",
			Key:      gocui.KeyEsc,
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.escapeFilterPrompt),
		},
		// Global scroll bindings
		{
			ViewName: "",
			Key:      'J',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.scrollDownMain),
		},
		{
			ViewName: "",
			Key:      'K',
			Modifier: gocui.ModNone,
			Handler:  wrappedHandler(gui.scrollUpMain),
		},
		{
			ViewName: "",
			Key:      'H',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollLeftMain,
		},
		{
			ViewName: "",
			Key:      'L',
			Modifier: gocui.ModNone,
			Handler:  gui.scrollRightMain,
		},
		{
			ViewName:    "",
			Key:         '+',
			Handler:     wrappedHandler(gui.nextScreenMode),
			Description: "Next Screen Mode",
		},
		{
			ViewName:    "",
			Key:         '_',
			Handler:     wrappedHandler(gui.prevScreenMode),
			Description: "Prev Screen Mode",
		},
	}

	// Lumine Docker Control panel bindings
	lumineDockerBindings := []*Binding{
		{
			ViewName:    "lumineDocker",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDockerStart,
			Description: "Start Docker",
		},
		{
			ViewName:    "lumineDocker",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDockerStop,
			Description: "Stop Docker",
		},
		{
			ViewName:    "lumineDocker",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleDockerRestart,
			Description: "Restart Docker",
		},
	}

	// Lumine Servers panel bindings
	lumineServersBindings := []*Binding{
		{
			ViewName:    "lumineServers",
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineServerAdd,
			Description: "New Server",
		},
		{
			ViewName:    "lumineServers",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineServerEdit,
			Description: "Edit Settings",
		},
		{
			ViewName:    "lumineServers",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineServerStart,
			Description: "Start Server",
		},
		{
			ViewName:    "lumineServers",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineServerStop,
			Description: "Stop Server",
		},
		{
			ViewName:    "lumineServers",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineServerRestart,
			Description: "Restart Server",
		},
		{
			ViewName:    "lumineServers",
			Key:         'x',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineServerRemove,
			Description: "Remove Server",
		},
		{
			ViewName:    "lumineServers",
			Key:         'c',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineSettings,
			Description: "Settings",
		},
	}

	// Lumine Languages panel bindings
	lumineLanguagesBindings := []*Binding{
		{
			ViewName:    "lumineLanguages",
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineLanguageAdd,
			Description: "New Runtime",
		},
		{
			ViewName:    "lumineLanguages",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineLanguageEdit,
			Description: "Edit Settings",
		},
		{
			ViewName:    "lumineLanguages",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineLanguageStart,
			Description: "Start Runtime",
		},
		{
			ViewName:    "lumineLanguages",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineLanguageStop,
			Description: "Stop Runtime",
		},
		{
			ViewName:    "lumineLanguages",
			Key:         'v',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineLanguageVersionSwitch,
			Description: "Switch Version",
		},
	}


	// Lumine Projects panel bindings
	lumineProjectsBindings := []*Binding{
		{
			ViewName:    "lumineProjects",
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineProjectCreate,
			Description: "New Project",
		},
		{
			ViewName:    "lumineProjects",
			Key:         'E',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineProjectEdit,
			Description: "Edit Settings",
		},
		{
			ViewName:    "lumineProjects",
			Key:         'd',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineProjectDelete,
			Description: "Delete Project",
		},
		{
			ViewName:    "lumineProjects",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineProjectExpose,
			Description: "Expose (Tunnel)",
		},
		{
			ViewName:    "lumineProjects",
			Key:         'o',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineProjectOpen,
			Description: "Open in Browser",
		},
		{
			ViewName:    "lumineProjects",
			Key:         't',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineProjectTerminal,
			Description: "Open Terminal",
		},
	}

	// Lumine Databases panel bindings
	lumineDatabasesBindings := []*Binding{
		{
			ViewName:    "lumineDatabases",
			Key:         'n',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineDatabaseServiceAdd,
			Description: "New Database Service",
		},
		{
			ViewName:    "lumineDatabases",
			Key:         'e',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineDatabaseServiceEdit,
			Description: "Edit Settings",
		},
		{
			ViewName:    "lumineDatabases",
			Key:         's',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineDatabaseServiceStart,
			Description: "Start Service",
		},
		{
			ViewName:    "lumineDatabases",
			Key:         'S',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineDatabaseServiceStop,
			Description: "Stop Service",
		},
		{
			ViewName:    "lumineDatabases",
			Key:         'r',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineDatabaseServiceRestart,
			Description: "Restart Service",
		},
		{
			ViewName:    "lumineDatabases",
			Key:         'x',
			Modifier:    gocui.ModNone,
			Handler:     gui.handleLumineDatabaseServiceRemove,
			Description: "Remove Service",
		},
	}

	bindings = append(bindings, lumineDockerBindings...)
	bindings = append(bindings, lumineServersBindings...)
	bindings = append(bindings, lumineLanguagesBindings...)
	bindings = append(bindings, lumineProjectsBindings...)
	bindings = append(bindings, lumineDatabasesBindings...)

	// Panel navigation bindings
	bindings = append(bindings, []*Binding{
		{Handler: gui.handleGoTo(gui.Panels.LumineDocker.View), Key: '1', Description: "Focus Docker"},
		{Handler: gui.handleGoTo(gui.Panels.LumineServers.View), Key: '2', Description: "Focus Servers"},
		{Handler: gui.handleGoTo(gui.Panels.LumineLanguages.View), Key: '3', Description: "Focus Languages"},
		{Handler: gui.handleGoTo(gui.Panels.LumineProjects.View), Key: '4', Description: "Focus Projects"},
		{Handler: gui.handleGoTo(gui.Panels.LumineDatabases.View), Key: '5', Description: "Focus Databases"},
	}...)

	// Add up/down/click bindings for all panels
	for _, panel := range gui.allSidePanels() {
		bindings = append(bindings, []*Binding{
			{ViewName: panel.GetView().Name(), Key: gocui.KeyArrowLeft, Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: panel.GetView().Name(), Key: gocui.KeyArrowRight, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: panel.GetView().Name(), Key: 'h', Modifier: gocui.ModNone, Handler: gui.previousView},
			{ViewName: panel.GetView().Name(), Key: 'l', Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: panel.GetView().Name(), Key: gocui.KeyTab, Modifier: gocui.ModNone, Handler: gui.nextView},
			{ViewName: panel.GetView().Name(), Key: gocui.KeyBacktab, Modifier: gocui.ModNone, Handler: gui.previousView},
		}...)
	}

	setUpDownClickBindings := func(viewName string, onUp func() error, onDown func() error, onClick func() error) {
		bindings = append(bindings, []*Binding{
			{ViewName: viewName, Key: 'k', Modifier: gocui.ModNone, Handler: wrappedHandler(onUp)},
			{ViewName: viewName, Key: gocui.KeyArrowUp, Modifier: gocui.ModNone, Handler: wrappedHandler(onUp)},
			{ViewName: viewName, Key: gocui.MouseWheelUp, Modifier: gocui.ModNone, Handler: wrappedHandler(onUp)},
			{ViewName: viewName, Key: 'j', Modifier: gocui.ModNone, Handler: wrappedHandler(onDown)},
			{ViewName: viewName, Key: gocui.KeyArrowDown, Modifier: gocui.ModNone, Handler: wrappedHandler(onDown)},
			{ViewName: viewName, Key: gocui.MouseWheelDown, Modifier: gocui.ModNone, Handler: wrappedHandler(onDown)},
			{ViewName: viewName, Key: gocui.MouseLeft, Modifier: gocui.ModNone, Handler: wrappedHandler(onClick)},
		}...)
	}

	for _, panel := range gui.allListPanels() {
		setUpDownClickBindings(panel.GetView().Name(), panel.HandlePrevLine, panel.HandleNextLine, panel.HandleClick)
	}

	setUpDownClickBindings("main", gui.scrollUpMain, gui.scrollDownMain, gui.handleMainClick)

	for _, panel := range gui.allSidePanels() {
		bindings = append(bindings,
			&Binding{
				ViewName:    panel.GetView().Name(),
				Key:         gocui.KeyEnter,
				Modifier:    gocui.ModNone,
				Handler:     gui.handleEnterMain,
				Description: "Focus Main",
			},
			&Binding{
				ViewName:    panel.GetView().Name(),
				Key:         '[',
				Modifier:    gocui.ModNone,
				Handler:     wrappedHandler(panel.HandlePrevMainTab),
				Description: "Previous Context",
			},
			&Binding{
				ViewName:    panel.GetView().Name(),
				Key:         ']',
				Modifier:    gocui.ModNone,
				Handler:     wrappedHandler(panel.HandleNextMainTab),
				Description: "Next Context",
			},
		)
	}

	for _, panel := range gui.allListPanels() {
		if !panel.IsFilterDisabled() {
			bindings = append(bindings, &Binding{
				ViewName:    panel.GetView().Name(),
				Key:         '/',
				Modifier:    gocui.ModNone,
				Handler:     wrappedHandler(gui.handleOpenFilter),
				Description: "Filter",
			})
		}
	}

	return bindings
}

func (gui *Gui) keybindings(g *gocui.Gui) error {
	bindings := gui.GetInitialKeybindings()

	for _, binding := range bindings {
		if err := g.SetKeybinding(binding.ViewName, binding.Key, binding.Modifier, binding.Handler); err != nil {
			return err
		}
	}

	if err := g.SetTabClickBinding("main", gui.onMainTabClick); err != nil {
		return err
	}

	return nil
}

func wrappedHandler(f func() error) func(*gocui.Gui, *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		return f()
	}
}
