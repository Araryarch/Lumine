package gui

import (
	"fmt"

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
			output += utils.ColoredString("Press 'S' to stop Docker\n", color.FgYellow)
		} else {
			output += utils.ColoredString("Press 's' to start Docker\n", color.FgGreen)
		}
		
		return output
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
