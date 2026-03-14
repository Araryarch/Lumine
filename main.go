package main

import (
	"context"
	"fmt"
	"log"
	"os"

	appProject "github.com/Araryarch/lumine/internal/application/project"
	appService "github.com/Araryarch/lumine/internal/application/service"
	"github.com/Araryarch/lumine/internal/infrastructure/config"
	"github.com/Araryarch/lumine/internal/infrastructure/docker"
	"github.com/Araryarch/lumine/internal/infrastructure/repository"
	"github.com/Araryarch/lumine/internal/presentation/tui"
	"github.com/jesseduffield/gocui"
)

func main() {
	if !docker.IsRunning() {
		fmt.Println("❌ Docker is not running")
		fmt.Println("Please start Docker and try again")
		os.Exit(1)
	}

	fmt.Println("✓ Docker is running")

	ctx := context.Background()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Error loading config: %v\n", err)
	}

	dockerClient, err := docker.NewClient(ctx)
	if err != nil {
		log.Fatalf("Error creating Docker client: %v\n", err)
	}
	defer dockerClient.Close()

	serviceRepo := repository.NewServiceRepository(dockerClient, cfg)
	projectRepo := repository.NewProjectRepository()

	serviceSvc := appService.NewService(serviceRepo)
	projectSvc := appProject.NewService(projectRepo)

	g := gocui.NewGui()
	if err := g.Init(); err != nil {
		log.Fatalf("Error initializing GUI: %v\n", err)
	}
	defer g.Close()

	// Enable mouse support
	g.Mouse = true
	g.Cursor = true

	controller := tui.NewController(g, cfg, serviceSvc, projectSvc)

	g.SetLayout(controller.Layout)

	if err := controller.SetupKeybindings(); err != nil {
		log.Fatalf("Error setting up keybindings: %v\n", err)
	}

	// Start auto-refresh
	controller.StartAutoRefresh(g)

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Fatalf("Error in main loop: %v\n", err)
	}

	fmt.Println("\n✓ Goodbye!")
}
