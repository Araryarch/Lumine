package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"lumine/internal/config"
	"lumine/internal/docker"
	"lumine/internal/ui"
)

func main() {
	// Validate Docker installation
	fmt.Println("🔍 Checking Docker installation...")
	validation := docker.ValidateDocker()

	if !validation.DockerInstalled {
		fmt.Println("❌ Docker is not installed!\n")
		fmt.Println(docker.InstallDocker())
		os.Exit(1)
	}

	if !validation.DockerRunning {
		fmt.Println("⚠️  Docker is installed but not running.")
		fmt.Println("🚀 Attempting to start Docker...")

		if err := docker.StartDocker(); err != nil {
			fmt.Printf("❌ Failed to start Docker: %v\n", err)
			fmt.Println("\nPlease start Docker manually and run 'lumine' again.")
			os.Exit(1)
		}

		// Wait for Docker to start
		fmt.Println("⏳ Waiting for Docker to start...")
		for i := 0; i < 30; i++ {
			time.Sleep(time.Second)
			validation = docker.ValidateDocker()
			if validation.DockerRunning {
				break
			}
		}

		if !validation.DockerRunning {
			fmt.Println("❌ Docker failed to start. Please start it manually.")
			os.Exit(1)
		}
	}

	fmt.Printf("✅ Docker is running (version %s)\n", validation.DockerVersion)

	if validation.ComposeInstalled {
		fmt.Printf("✅ Docker Compose is available (version %s)\n", validation.ComposeVersion)
	} else {
		fmt.Println("⚠️  Docker Compose not found (optional)")
	}

	// Initialize config
	fmt.Println("📝 Initializing configuration...")
	if err := config.InitConfig(); err != nil {
		fmt.Printf("❌ Error initializing config: %v\n", err)
		os.Exit(1)
	}

	// Ensure Docker network
	dockerMgr, err := docker.NewManager()
	if err != nil {
		fmt.Printf("❌ Error connecting to Docker: %v\n", err)
		os.Exit(1)
	}

	ctx := context.Background()
	if err := dockerMgr.EnsureDockerNetwork(ctx); err != nil {
		fmt.Printf("⚠️  Warning: Could not create Docker network: %v\n", err)
	}

	fmt.Println("✅ Configuration ready")
	fmt.Println("\n🌟 Starting Lumine...\n")

	time.Sleep(time.Millisecond * 500)

	// Start TUI
	p := tea.NewProgram(ui.NewModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
