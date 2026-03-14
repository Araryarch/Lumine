package app

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/jesseduffield/gocui"
	"github.com/Araryarch/lumine/pkg/config"
	"github.com/Araryarch/lumine/pkg/gui"
)

type App struct {
	Config        *config.Config
	Version       string
	DockerClient  *client.Client
	Gui           *gocui.Gui
	GuiController *gui.GuiController
	Context       context.Context
}

func NewApp(cfg *config.Config, version string) (*App, error) {
	ctx := context.Background()

	// Create Docker client
	dockerClient, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, fmt.Errorf("failed to create docker client: %w", err)
	}

	// Create GUI
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return nil, fmt.Errorf("failed to create gui: %w", err)
	}

	app := &App{
		Config:       cfg,
		Version:      version,
		DockerClient: dockerClient,
		Gui:          g,
		Context:      ctx,
	}

	// Create GUI controller
	app.GuiController = gui.NewGuiController(app.Gui, app.Config, app.DockerClient, app.Context)

	return app, nil
}

func (a *App) Run() error {
	defer a.Gui.Close()

	a.Gui.SetManagerFunc(a.GuiController.Layout)

	if err := a.GuiController.SetupKeybindings(); err != nil {
		return err
	}

	if err := a.Gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}
