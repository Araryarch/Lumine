package app

import (
	"io"
	"strings"

	"github.com/Araryarch/Lumine/pkg/commands"
	"github.com/Araryarch/Lumine/pkg/config"
	"github.com/Araryarch/Lumine/pkg/gui"
	"github.com/Araryarch/Lumine/pkg/i18n"
	"github.com/Araryarch/Lumine/pkg/log"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/sirupsen/logrus"
)

// App struct
type App struct {
	closers []io.Closer

	Config    *config.AppConfig
	Log       *logrus.Entry
	OSCommand *commands.OSCommand
	Gui       *gui.Gui
	Tr        *i18n.TranslationSet
	ErrorChan chan error
}

// NewApp bootstrap a new application
func NewApp(config *config.AppConfig) (*App, error) {
	app := &App{
		closers:   []io.Closer{},
		Config:    config,
		ErrorChan: make(chan error),
	}
	var err error
	app.Log = log.NewLogger(config, "23432119147a4367abf7c0de2aa99a2d")
	app.Tr, err = i18n.NewTranslationSetFromConfig(app.Log, config.UserConfig.Gui.Language)
	if err != nil {
		return app, err
	}
	app.OSCommand = commands.NewOSCommand(app.Log, config)

	app.Gui, err = gui.NewGui(app.Log, app.OSCommand, app.Tr, config, app.ErrorChan)
	if err != nil {
		return app, err
	}
	return app, nil
}

func (app *App) Run() error {
	return app.Gui.Run()
}

func (app *App) Close() error {
	return utils.CloseMany(app.closers)
}

type errorMapping struct {
	originalError string
	newError      string
}

// KnownError takes an error and tells us whether it's an error that we know about
func (app *App) KnownError(err error) (string, bool) {
	errorMessage := err.Error()

	mappings := []errorMapping{
		{
			originalError: "connection refused",
			newError:      "Cannot connect to services. Please ensure services are running.",
		},
	}

	for _, mapping := range mappings {
		if strings.Contains(errorMessage, mapping.originalError) {
			return mapping.newError, true
		}
	}

	return "", false
}
