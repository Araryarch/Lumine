package docker

import (
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

type ValidationResult struct {
	DockerInstalled bool
	DockerRunning   bool
	DockerVersion   string
	ComposeInstalled bool
	ComposeVersion  string
	Error           error
}

// ValidateDocker checks if Docker is installed and running
func ValidateDocker() *ValidationResult {
	result := &ValidationResult{}

	// Check if docker command exists
	dockerPath, err := exec.LookPath("docker")
	if err != nil {
		result.Error = fmt.Errorf("Docker is not installed")
		return result
	}
	result.DockerInstalled = true

	// Check Docker version
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	versionCmd := exec.CommandContext(ctx, dockerPath, "version", "--format", "{{.Server.Version}}")
	output, err := versionCmd.Output()
	if err != nil {
		result.Error = fmt.Errorf("Docker is installed but not running")
		return result
	}
	result.DockerVersion = strings.TrimSpace(string(output))
	result.DockerRunning = true

	// Check Docker Compose
	composeCmd := exec.CommandContext(ctx, dockerPath, "compose", "version", "--short")
	output, err = composeCmd.Output()
	if err == nil {
		result.ComposeInstalled = true
		result.ComposeVersion = strings.TrimSpace(string(output))
	}

	return result
}

// InstallDocker provides instructions to install Docker
func InstallDocker() string {
	switch runtime.GOOS {
	case "darwin":
		return `Docker is not installed. Please install Docker Desktop for Mac:

1. Download from: https://www.docker.com/products/docker-desktop
2. Install Docker Desktop
3. Start Docker Desktop
4. Run 'lumine' again

Or install via Homebrew:
  brew install --cask docker`

	case "windows":
		return `Docker is not installed. Please install Docker Desktop for Windows:

1. Download from: https://www.docker.com/products/docker-desktop
2. Install Docker Desktop
3. Start Docker Desktop
4. Run 'lumine' again

Or install via Chocolatey:
  choco install docker-desktop`

	case "linux":
		return `Docker is not installed. Install Docker using your package manager:

Ubuntu/Debian:
  curl -fsSL https://get.docker.com -o get-docker.sh
  sudo sh get-docker.sh
  sudo usermod -aG docker $USER
  newgrp docker

Fedora:
  sudo dnf install docker
  sudo systemctl start docker
  sudo usermod -aG docker $USER

Arch Linux:
  sudo pacman -S docker
  sudo systemctl start docker
  sudo usermod -aG docker $USER

After installation, run 'lumine' again.`

	default:
		return "Please install Docker from https://docs.docker.com/get-docker/"
	}
}

// StartDocker attempts to start Docker daemon
func StartDocker() error {
	switch runtime.GOOS {
	case "darwin":
		// Start Docker Desktop on macOS
		return exec.Command("open", "-a", "Docker").Run()

	case "windows":
		// Start Docker Desktop on Windows
		return exec.Command("cmd", "/C", "start", "Docker Desktop").Run()

	case "linux":
		// Start Docker daemon on Linux
		return exec.Command("sudo", "systemctl", "start", "docker").Run()

	default:
		return fmt.Errorf("unsupported operating system")
	}
}

// EnsureDockerNetwork creates the lumine network if it doesn't exist
func (m *Manager) EnsureDockerNetwork(ctx context.Context) error {
	// Check if network exists
	cmd := exec.CommandContext(ctx, "docker", "network", "ls", "--filter", "name=lumine", "--format", "{{.Name}}")
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	if strings.TrimSpace(string(output)) == "lumine" {
		return nil // Network already exists
	}

	// Create network
	createCmd := exec.CommandContext(ctx, "docker", "network", "create", "lumine")
	return createCmd.Run()
}
