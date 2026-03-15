package i18n

import (
	"github.com/sirupsen/logrus"
)

// TranslationSet contains all translations
type TranslationSet struct {
	// Lumine specific
	LumineServices  string
	LumineProjects  string
	LumineDatabases string

	// Common
	Confirm                                    string
	Cancel                                     string
	Close                                      string
	Quit                                       string
	Return                                     string
	Navigate                                   string
	Scroll                                     string
	Menu                                       string
	Filter                                     string
	FilterPrompt                               string
	ErrorOccurred                              string
	ErrorTitle                                 string
	ConnectionFailed                           string
	ConfirmQuit                                string
	NotEnoughSpace                             string
	Donate                                     string
	CannotKillChildError                       string
	Execute                                    string
	Yes                                        string
	No                                         string
	NoViewMachingNewLineFocusedSwitchStatement string

	// Actions
	Remove   string
	Start    string
	Stop     string
	Restart  string
	ViewLogs string

	// Screen modes
	LcNextScreenMode string
	LcPrevScreenMode string
	LcFilter         string

	// Context
	PreviousContext  string
	NextContext      string
	FocusMain        string
	ViewBulkCommands string

	// Custom commands
	CustomCommandTitle string

	// Menu
	MenuTitle string

	// Panels
	ProjectTitle    string
	ServicesTitle   string
	ContainersTitle string
	ImagesTitle     string
	VolumesTitle    string
	StacksTitle     string

	// Containers
	StandaloneContainersTitle string

	// Networks
	NetworksTitle string
}

// NewTranslationSetFromConfig creates a new translation set
func NewTranslationSetFromConfig(log *logrus.Entry, language string) (*TranslationSet, error) {
	return &TranslationSet{
		// Lumine
		LumineServices:  "Lumine Services",
		LumineProjects:  "Lumine Projects",
		LumineDatabases: "Lumine Databases",

		// Common
		Confirm:              "Confirm",
		Cancel:               "Cancel",
		Close:                "Close",
		Quit:                 "Quit",
		Return:               "Return",
		Navigate:             "Navigate",
		Scroll:               "Scroll",
		Menu:                 "Menu",
		Filter:               "Filter",
		FilterPrompt:         "Type to filter",
		ErrorOccurred:        "An error occurred",
		ErrorTitle:           "Error",
		ConnectionFailed:     "Connection failed",
		ConfirmQuit:          "Are you sure you want to quit?",
		NotEnoughSpace:       "Not enough space",
		Donate:               "Donate",
		CannotKillChildError: "Cannot kill child process",
		Execute:              "Execute",
		Yes:                  "Yes",
		No:                   "No",
		NoViewMachingNewLineFocusedSwitchStatement: "No view matching new line focused switch statement",

		// Actions
		Remove:   "Remove",
		Start:    "Start",
		Stop:     "Stop",
		Restart:  "Restart",
		ViewLogs: "View Logs",

		// Screen modes
		LcNextScreenMode: "Next Screen Mode",
		LcPrevScreenMode: "Previous Screen Mode",
		LcFilter:         "Filter",

		// Context
		PreviousContext:  "Previous Context",
		NextContext:      "Next Context",
		FocusMain:        "Focus Main",
		ViewBulkCommands: "View Bulk Commands",

		// Custom commands
		CustomCommandTitle: "Custom Command",

		// Menu
		MenuTitle: "Menu",

		// Panels
		ProjectTitle:    "Projects",
		ServicesTitle:   "Services",
		ContainersTitle: "Containers",
		ImagesTitle:     "Images",
		VolumesTitle:    "Volumes",
		StacksTitle:     "Stacks",

		// Containers
		StandaloneContainersTitle: "Containers",

		// Networks
		NetworksTitle: "Networks",
	}, nil
}
