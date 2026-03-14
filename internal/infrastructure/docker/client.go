package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

type Client struct {
	cli *client.Client
	ctx context.Context
}

func NewClient(ctx context.Context) (*Client, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}

	return &Client{
		cli: cli,
		ctx: ctx,
	}, nil
}

func (c *Client) Close() error {
	return c.cli.Close()
}

func (c *Client) ListContainers() ([]types.Container, error) {
	return c.cli.ContainerList(c.ctx, container.ListOptions{All: true})
}

func (c *Client) GetContainerStatus(name string) (string, error) {
	containers, err := c.ListContainers()
	if err != nil {
		return "", err
	}

	for _, cont := range containers {
		for _, n := range cont.Names {
			if n == "/"+name || n == name {
				return cont.State, nil
			}
		}
	}

	return "not found", nil
}

func (c *Client) StartContainer(name string) error {
	return c.cli.ContainerStart(c.ctx, name, container.StartOptions{})
}

func (c *Client) StopContainer(name string) error {
	timeout := 10
	return c.cli.ContainerStop(c.ctx, name, container.StopOptions{Timeout: &timeout})
}

func (c *Client) RestartContainer(name string) error {
	timeout := 10
	return c.cli.ContainerRestart(c.ctx, name, container.StopOptions{Timeout: &timeout})
}

func (c *Client) GetLogs(name string, tail string) (string, error) {
	options := container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Tail:       tail,
		Details:    true,
		Timestamps: true,
	}

	logs, err := c.cli.ContainerLogs(c.ctx, name, options)
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

func (c *Client) GetVersion() (string, error) {
	version, err := c.cli.ServerVersion(c.ctx)
	if err != nil {
		return "", err
	}
	return version.Version, nil
}

func (c *Client) ExecCommand(name string, cmd []string) error {
	execConfig := container.ExecOptions{
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          cmd,
	}

	execID, err := c.cli.ContainerExecCreate(c.ctx, name, execConfig)
	if err != nil {
		return err
	}

	return c.cli.ContainerExecStart(c.ctx, execID.ID, container.ExecStartOptions{
		Detach: false,
	})
}

func (c *Client) GetContainerInfo(name string) (*types.Container, error) {
	containers, err := c.ListContainers()
	if err != nil {
		return nil, err
	}

	for _, cont := range containers {
		for _, n := range cont.Names {
			if n == "/"+name || n == name {
				return &cont, nil
			}
		}
	}

	return nil, nil
}

func IsRunning() bool {
	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return false
	}
	defer cli.Close()

	_, err = cli.Ping(context.Background())
	return err == nil
}
