package docker

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"lumine/internal/config"
	"lumine/internal/port"
)

type Manager struct {
	client *client.Client
}

func NewManager() (*Manager, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Manager{client: cli}, nil
}

func (m *Manager) StartService(ctx context.Context, service config.Service) error {
	// Check if port is available, find alternative if not
	portMgr := port.NewManager()
	actualPort, err := portMgr.GetAlternativePort(service.Port)
	if err != nil {
		return fmt.Errorf("failed to find available port: %w", err)
	}

	// Log if port was changed
	if actualPort != service.Port {
		fmt.Printf("⚠️  Port %d is in use, using alternative port %d for %s\n", 
			service.Port, actualPort, service.Name)
	}

	// Pull image
	imageName := fmt.Sprintf("%s:%s", service.Type, service.Version)
	reader, err := m.client.ImagePull(ctx, imageName, types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)
	reader.Close()

	// Create container with actual port
	containerConfig := &container.Config{
		Image: imageName,
		Env:   envMapToSlice(service.Env),
	}

	hostConfig := &container.HostConfig{
		PortBindings: nat.PortMap{
			nat.Port(fmt.Sprintf("%d/tcp", service.Port)): []nat.PortBinding{
				{HostIP: "0.0.0.0", HostPort: fmt.Sprintf("%d", actualPort)},
			},
		},
	}

	resp, err := m.client.ContainerCreate(ctx, containerConfig, hostConfig, nil, nil, fmt.Sprintf("lumine-%s", service.Name))
	if err != nil {
		return err
	}

	return m.client.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{})
}

func (m *Manager) StopService(ctx context.Context, serviceName string) error {
	containerName := fmt.Sprintf("lumine-%s", serviceName)
	timeout := 10
	return m.client.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout})
}

func (m *Manager) ListContainers(ctx context.Context) ([]types.Container, error) {
	return m.client.ContainerList(ctx, types.ContainerListOptions{All: true})
}

func envMapToSlice(envMap map[string]string) []string {
	var envSlice []string
	for k, v := range envMap {
		envSlice = append(envSlice, fmt.Sprintf("%s=%s", k, v))
	}
	return envSlice
}
