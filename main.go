package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
	appService "github.com/Araryarch/lumine/internal/application/service"
	appProject "github.com/Araryarch/lumine/internal/application/project"
	"github.com/Araryarch/lumine/internal/infrastructure/config"
	"github.com/Araryarch/lumine/internal/infrastructure/docker"
	"github.com/Araryarch/lumine/internal/infrastructure/repository"
	"github.com/Araryarch/lumine/internal/presentation/tui"
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

	g, err := gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		log.Fatalf("Error creating GUI: %v\n", err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorBlack
	g.SelBgColor = gocui.ColorGreen

	controller := tui.NewController(g, cfg, serviceSvc, projectSvc)

	g.SetManagerFunc(controller.Layout)

	if err := controller.SetupKeybindings(); err != nil {
		log.Fatalf("Error setting up keybindings: %v\n", err)
	}

	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Fatalf("Error in main loop: %v\n", err)
	}

	fmt.Println("\n✓ Goodbye!")
}
