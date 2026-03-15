package gui

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/tasks"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

type DockerControl struct {
	Name   string
	Status string
}

func (gui *Gui) getLumineDockerPanel() *panels.SideListPanel[*DockerControl] {
	return &panels.SideListPanel[*DockerControl]{
		ContextState: &panels.ContextState[*DockerControl]{
			GetMainTabs: func() []panels.MainTab[*DockerControl] {
				return []panels.MainTab[*DockerControl]{
					{
						Key:    "status",
						Title:  "Docker Status",
						Render: gui.renderDockerStatus,
					},
					{
						Key:    "stats",
						Title:  "Resource Usage",
						Render: gui.renderDockerStats,
					},
					{
						Key:    "containers",
						Title:  "Containers",
						Render: gui.renderDockerContainers,
					},
					{
						Key:    "logs",
						Title:  "Docker Logs",
						Render: gui.renderDockerLogs,
					},
					{
						Key:    "info",
						Title:  "System Info",
						Render: gui.renderDockerInfo,
					},
				}
			},
			GetItemContextCacheKey: func(dc *DockerControl) string {
				return "docker-control-" + dc.Status
			},
		},
		ListPanel: panels.ListPanel[*DockerControl]{
			List: panels.NewFilteredList[*DockerControl](),
			View: gui.Views.LumineDocker,
		},
		NoItemsMessage: "Docker Control",
		Gui:            gui.intoInterface(),
		GetTableCells: func(dc *DockerControl) []string {
			statusColor := color.FgRed
			if dc.Status == "running" {
				statusColor = color.FgGreen
			}

			return []string{
				utils.ColoredString(dc.Name, color.FgCyan),
				utils.ColoredString(dc.Status, statusColor),
			}
		},
	}
}

func (gui *Gui) renderDockerStatus(dc *DockerControl) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Docker: ", 15) + gui.getColoredStatus(dc.Status) + "\n\n"
		
		if dc.Status == "running" {
			// Get Docker version
			version, _ := gui.OSCommand.RunCommandWithOutput("docker version --format '{{.Server.Version}}'")
			output += utils.WithPadding("Version: ", 15) + strings.TrimSpace(version) + "\n"
			
			// Get running containers count
			containersCount, _ := gui.OSCommand.RunCommandWithOutput("docker ps -q | wc -l")
			output += utils.WithPadding("Containers: ", 15) + strings.TrimSpace(containersCount) + " running\n"
			
			// Get images count
			imagesCount, _ := gui.OSCommand.RunCommandWithOutput("docker images -q | wc -l")
			output += utils.WithPadding("Images: ", 15) + strings.TrimSpace(imagesCount) + " total\n"
			
			// Get volumes count
			volumesCount, _ := gui.OSCommand.RunCommandWithOutput("docker volume ls -q | wc -l")
			output += utils.WithPadding("Volumes: ", 15) + strings.TrimSpace(volumesCount) + " total\n"
			
			output += "\n" + utils.ColoredString("Press 'S' to stop Docker\n", color.FgYellow)
		} else {
			output += utils.ColoredString("Press 's' to start Docker\n", color.FgGreen)
		}
		
		return output
	})
}

func (gui *Gui) renderDockerStats(dc *DockerControl) tasks.TaskFunc {
	return gui.NewTickerTask(TickerTaskOpts{
		Func: func(ctx context.Context, notifyStopped chan struct{}) {
			if dc.Status != "running" {
				gui.RenderStringMain("Docker is not running")
				return
			}

			// Get Docker stats
			stats, err := gui.OSCommand.RunCommandWithOutput("docker stats --no-stream --format 'table {{.Name}}\t{{.CPUPerc}}\t{{.MemUsage}}\t{{.NetIO}}\t{{.BlockIO}}'")
			if err != nil {
				gui.RenderStringMain(fmt.Sprintf("Error fetching stats: %v", err))
				return
			}

			output := utils.ColoredString("Docker Resource Usage:\n\n", color.FgYellow)
			output += stats

			// Get system-wide Docker info
			systemInfo, _ := gui.OSCommand.RunCommandWithOutput("docker system df")
			output += "\n\n" + utils.ColoredString("Disk Usage:\n", color.FgCyan)
			output += systemInfo

			gui.reRenderStringMain(output)
		},
		Duration:   time.Second * 3,
		Before:     func(ctx context.Context) { gui.clearMainView() },
		Wrap:       false,
		Autoscroll: false,
	})
}

func (gui *Gui) renderDockerContainers(dc *DockerControl) tasks.TaskFunc {
	return gui.NewTickerTask(TickerTaskOpts{
		Func: func(ctx context.Context, notifyStopped chan struct{}) {
			if dc.Status != "running" {
				gui.RenderStringMain("Docker is not running")
				return
			}

			// Get all containers (running and stopped)
			containers, err := gui.OSCommand.RunCommandWithOutput("docker ps -a --format 'table {{.Names}}\t{{.Status}}\t{{.Ports}}\t{{.Image}}'")
			if err != nil {
				gui.RenderStringMain(fmt.Sprintf("Error fetching containers: %v", err))
				return
			}

			output := utils.ColoredString("All Containers:\n\n", color.FgYellow)
			output += containers

			gui.reRenderStringMain(output)
		},
		Duration:   time.Second * 2,
		Before:     func(ctx context.Context) { gui.clearMainView() },
		Wrap:       false,
		Autoscroll: false,
	})
}

func (gui *Gui) renderDockerLogs(dc *DockerControl) tasks.TaskFunc {
	return gui.NewTickerTask(TickerTaskOpts{
		Func: func(ctx context.Context, notifyStopped chan struct{}) {
			if dc.Status != "running" {
				gui.RenderStringMain("Docker is not running")
				return
			}

			// Get Docker daemon logs (last 50 lines)
			logs, err := gui.OSCommand.RunCommandWithOutput("journalctl -u docker.service -n 50 --no-pager")
			if err != nil {
				// Fallback for systems without systemd
				logs, err = gui.OSCommand.RunCommandWithOutput("tail -n 50 /var/log/docker.log 2>/dev/null || echo 'Docker logs not available'")
				if err != nil {
					gui.RenderStringMain("Docker logs not available")
					return
				}
			}

			output := utils.ColoredString("Docker Daemon Logs:\n\n", color.FgYellow)
			output += logs

			gui.reRenderStringMain(output)
		},
		Duration:   time.Second * 3,
		Before:     func(ctx context.Context) { gui.clearMainView() },
		Wrap:       true,
		Autoscroll: true,
	})
}

func (gui *Gui) renderDockerInfo(dc *DockerControl) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		info, err := gui.OSCommand.RunCommandWithOutput("docker info --format '{{.ServerVersion}}\n{{.OperatingSystem}}\n{{.Architecture}}'")
		if err != nil {
			return "Docker not available or not running"
		}
		
		return "Docker System Info:\n\n" + info
	})
}

// Keybinding handlers
func (gui *Gui) handleDockerStart(g *gocui.Gui, v *gocui.View) error {
	return gui.WithWaitingStatus("Starting Docker...", func() error {
		_, err := gui.OSCommand.RunCommandWithOutput("sudo systemctl start docker")
		if err != nil {
			return gui.createErrorPanel(fmt.Sprintf("Failed to start Docker: %v", err))
		}
		gui.Orchestrator.NotificationMgr.ShowSuccess("Docker started successfully")
		return gui.refreshDockerControl()
	})
}

func (gui *Gui) handleDockerStop(g *gocui.Gui, v *gocui.View) error {
	return gui.createConfirmationPanel("Confirm", "Stop Docker daemon?", func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Stopping Docker...", func() error {
			_, err := gui.OSCommand.RunCommandWithOutput("sudo systemctl stop docker")
			if err != nil {
				return gui.createErrorPanel(fmt.Sprintf("Failed to stop Docker: %v", err))
			}
			gui.Orchestrator.NotificationMgr.ShowSuccess("Docker stopped")
			return gui.refreshDockerControl()
		})
	}, nil)
}

func (gui *Gui) handleDockerRestart(g *gocui.Gui, v *gocui.View) error {
	return gui.WithWaitingStatus("Restarting Docker...", func() error {
		_, err := gui.OSCommand.RunCommandWithOutput("sudo systemctl restart docker")
		if err != nil {
			return gui.createErrorPanel(fmt.Sprintf("Failed to restart Docker: %v", err))
		}
		gui.Orchestrator.NotificationMgr.ShowSuccess("Docker restarted")
		return gui.refreshDockerControl()
	})
}

func (gui *Gui) refreshDockerControl() error {
	if gui.Panels.LumineDocker == nil {
		return nil
	}

	// Check Docker status
	_, err := gui.OSCommand.RunCommandWithOutput("docker info")
	status := "stopped"
	if err == nil {
		status = "running"
	}

	dockerControl := []*DockerControl{
		{
			Name:   "Docker Daemon",
			Status: status,
		},
	}

	gui.Panels.LumineDocker.SetItems(dockerControl)
	return gui.Panels.LumineDocker.RerenderList()
}
