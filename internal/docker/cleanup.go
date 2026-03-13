package docker

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/volume"
)

// CleanupOptions defines cleanup behavior
type CleanupOptions struct {
	RemoveContainers bool
	RemoveVolumes    bool
	RemoveNetworks   bool
	Force            bool
}

// ListLumineContainers lists all containers with "lumine-" prefix
func (m *Manager) ListLumineContainers(ctx context.Context) ([]types.Container, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "lumine-")
	
	return m.client.ContainerList(ctx, types.ContainerListOptions{
		All:     true,
		Filters: filterArgs,
	})
}

// StopAllContainers stops all Lumine containers
func (m *Manager) StopAllContainers(ctx context.Context) error {
	containers, err := m.ListLumineContainers(ctx)
	if err != nil {
		return err
	}

	for _, container := range containers {
		if container.State == "running" {
			if err := m.client.ContainerStop(ctx, container.ID, container.StopOptions{}); err != nil {
				return fmt.Errorf("failed to stop container %s: %w", container.Names[0], err)
			}
		}
	}

	return nil
}

// RemoveContainer removes a specific container
func (m *Manager) RemoveContainer(ctx context.Context, containerName string) error {
	// Add lumine- prefix if not present
	if !strings.HasPrefix(containerName, "lumine-") {
		containerName = "lumine-" + containerName
	}

	return m.client.ContainerRemove(ctx, containerName, types.ContainerRemoveOptions{
		Force:         true,
		RemoveVolumes: false,
	})
}

// RemoveAllContainers removes all Lumine containers
func (m *Manager) RemoveAllContainers(ctx context.Context, force bool) error {
	containers, err := m.ListLumineContainers(ctx)
	if err != nil {
		return err
	}

	for _, container := range containers {
		if err := m.client.ContainerRemove(ctx, container.ID, types.ContainerRemoveOptions{
			Force:         force,
			RemoveVolumes: false,
		}); err != nil {
			return fmt.Errorf("failed to remove container %s: %w", container.Names[0], err)
		}
	}

	return nil
}

// ListLumineVolumes lists all volumes with "lumine_" prefix
func (m *Manager) ListLumineVolumes(ctx context.Context) ([]*volume.Volume, error) {
	filterArgs := filters.NewArgs()
	filterArgs.Add("name", "lumine_")
	
	volumeList, err := m.client.VolumeList(ctx, volume.ListOptions{Filters: filterArgs})
	if err != nil {
		return nil, err
	}

	var lumineVolumes []*volume.Volume
	for _, vol := range volumeList.Volumes {
		if strings.HasPrefix(vol.Name, "lumine_") {
			lumineVolumes = append(lumineVolumes, vol)
		}
	}

	return lumineVolumes, nil
}

// RemoveVolume removes a specific volume
func (m *Manager) RemoveVolume(ctx context.Context, volumeName string) error {
	return m.client.VolumeRemove(ctx, volumeName, true)
}

// RemoveAllVolumes removes all Lumine volumes
func (m *Manager) RemoveAllVolumes(ctx context.Context) error {
	volumes, err := m.ListLumineVolumes(ctx)
	if err != nil {
		return err
	}

	for _, volume := range volumes {
		if err := m.client.VolumeRemove(ctx, volume.Name, true); err != nil {
			return fmt.Errorf("failed to remove volume %s: %w", volume.Name, err)
		}
	}

	return nil
}

// RemoveNetwork removes the Lumine network
func (m *Manager) RemoveNetwork(ctx context.Context) error {
	return m.client.NetworkRemove(ctx, "lumine")
}

// Cleanup performs cleanup based on options
func (m *Manager) Cleanup(ctx context.Context, opts CleanupOptions) error {
	if opts.RemoveContainers {
		if !opts.Force {
			if err := m.StopAllContainers(ctx); err != nil {
				return fmt.Errorf("failed to stop containers: %w", err)
			}
		}
		
		if err := m.RemoveAllContainers(ctx, opts.Force); err != nil {
			return fmt.Errorf("failed to remove containers: %w", err)
		}
	}

	if opts.RemoveVolumes {
		if err := m.RemoveAllVolumes(ctx); err != nil {
			return fmt.Errorf("failed to remove volumes: %w", err)
		}
	}

	if opts.RemoveNetworks {
		if err := m.RemoveNetwork(ctx); err != nil {
			// Ignore error if network doesn't exist
			if !strings.Contains(err.Error(), "not found") {
				return fmt.Errorf("failed to remove network: %w", err)
			}
		}
	}

	return nil
}

// GetContainerStats returns statistics for a container
func (m *Manager) GetContainerStats(ctx context.Context, containerName string) (types.ContainerStats, error) {
	stats, err := m.client.ContainerStats(ctx, containerName, false)
	if err != nil {
		return types.ContainerStats{}, err
	}
	defer stats.Body.Close()

	return stats, nil
}

// GetContainerLogs returns logs for a container
func (m *Manager) GetContainerLogs(ctx context.Context, containerName string, tail string) (string, error) {
	options := types.ContainerLogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}

	logs, err := m.client.ContainerLogs(ctx, containerName, options)
	if err != nil {
		return "", err
	}
	defer logs.Close()

	buf := new(strings.Builder)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
