package docker

import (
	"context"
	"fmt"
	"io"
	"os/exec"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type ContainerInfo struct {
	Name   string
	Status string
	State  string
	Ports  []string
	Image  string
}

func IsRunning() bool {
	cmd := exec.Command("docker", "ps")
	if err := cmd.Run(); err != nil {
		return false
	}
	return true
}

func GetClient() (*client.Client, error) {
	return client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
}

func ListContainers(ctx context.Context, cli *client.Client) ([]ContainerInfo, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return nil, err
	}

	var result []ContainerInfo
	for _, c := range containers {
		name := c.Names[0]
		if len(name) > 0 && name[0] == '/' {
			name = name[1:]
		}

		ports := []string{}
		for _, port := range c.Ports {
			if port.PublicPort > 0 {
				ports = append(ports, fmt.Sprintf("%d:%d", port.PublicPort, port.PrivatePort))
			}
		}

		result = append(result, ContainerInfo{
			Name:   name,
			Status: c.Status,
			State:  c.State,
			Ports:  ports,
			Image:  c.Image,
		})
	}

	return result, nil
}

func GetContainerStatus(ctx context.Context, cli *client.Client, containerName string) (string, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{All: true})
	if err != nil {
		return "unknown", err
	}

	for _, c := range containers {
		for _, name := range c.Names {
			if name == "/"+containerName || name == containerName {
				return c.State, nil
			}
		}
	}

	return "not found", nil
}

func StartContainer(ctx context.Context, cli *client.Client, containerName string) error {
	return cli.ContainerStart(ctx, containerName, container.StartOptions{})
}

func StopContainer(ctx context.Context, cli *client.Client, containerName string) error {
	timeout := 10
	return cli.ContainerStop(ctx, containerName, container.StopOptions{Timeout: &timeout})
}

func RestartContainer(ctx context.Context, cli *client.Client, containerName string) error {
	timeout := 10
	return cli.ContainerRestart(ctx, containerName, container.StopOptions{Timeout: &timeout})
}

func RemoveContainer(ctx context.Context, cli *client.Client, containerName string) error {
	return cli.ContainerRemove(ctx, containerName, container.RemoveOptions{Force: true})
}

func GetContainerLogs(ctx context.Context, cli *client.Client, containerName string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
	}

	logs, err := cli.ContainerLogs(ctx, containerName, options)
	if err != nil {
		return "", err
	}
	defer logs.Close()

	logBytes, err := io.ReadAll(logs)
	if err != nil {
		return "", err
	}

	return string(logBytes), nil
}

func ExecCommand(ctx context.Context, cli *client.Client, containerName string, cmd []string) error {
	execConfig := types.ExecConfig{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execID, err := cli.ContainerExecCreate(ctx, containerName, execConfig)
	if err != nil {
		return err
	}

	return cli.ContainerExecStart(ctx, execID.ID, types.ExecStartCheck{})
}

func GetDockerVersion(ctx context.Context, cli *client.Client) (string, error) {
	version, err := cli.ServerVersion(ctx)
	if err != nil {
		return "", err
	}
	return version.Version, nil
}
