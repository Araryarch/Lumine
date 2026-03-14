package services

import (
	"context"
	"fmt"

	"github.com/docker/docker/client"
	"github.com/Araryarch/lumine/pkg/config"
	"github.com/Araryarch/lumine/pkg/docker"
)

type ServiceStatus struct {
	Name    string
	Status  string
	Running bool
	Port    int
	Image   string
}

func GetServicesStatus(ctx context.Context, cli *client.Client, cfg *config.Config) ([]ServiceStatus, error) {
	containers, err := docker.ListContainers(ctx, cli)
	if err != nil {
		return nil, err
	}

	var statuses []ServiceStatus
	for name, service := range cfg.Services {
		if !service.Enabled {
			continue
		}

		containerName := "lumine-" + name
		status := ServiceStatus{
			Name:    name,
			Status:  "stopped",
			Running: false,
			Port:    service.Port,
			Image:   service.Image,
		}

		// Find container in list
		for _, container := range containers {
			if container.Name == containerName {
				status.Status = container.State
				status.Running = container.State == "running"
				break
			}
		}

		statuses = append(statuses, status)
	}

	return statuses, nil
}

func StartAllServices(ctx context.Context, cli *client.Client, cfg *config.Config) error {
	for name, service := range cfg.Services {
		if !service.Enabled {
			continue
		}

		containerName := "lumine-" + name
		if err := docker.StartContainer(ctx, cli, containerName); err != nil {
			return fmt.Errorf("failed to start %s: %w", name, err)
		}
	}
	return nil
}

func StopAllServices(ctx context.Context, cli *client.Client, cfg *config.Config) error {
	for name, service := range cfg.Services {
		if !service.Enabled {
			continue
		}

		containerName := "lumine-" + name
		if err := docker.StopContainer(ctx, cli, containerName); err != nil {
			return fmt.Errorf("failed to stop %s: %w", name, err)
		}
	}
	return nil
}

func RestartAllServices(ctx context.Context, cli *client.Client, cfg *config.Config) error {
	for name, service := range cfg.Services {
		if !service.Enabled {
			continue
		}

		containerName := "lumine-" + name
		if err := docker.RestartContainer(ctx, cli, containerName); err != nil {
			return fmt.Errorf("failed to restart %s: %w", name, err)
		}
	}
	return nil
}
