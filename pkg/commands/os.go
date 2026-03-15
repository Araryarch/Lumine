package commands

import (
	"os"
	"os/exec"

	"github.com/jesseduffield/lazydocker/pkg/config"
	"github.com/sirupsen/logrus"
)

// Platform represents the OS platform
type Platform struct {
	OS              string
	Shell           string
	ShellArg        string
	OpenCommand     string
	OpenLinkCommand string
}

// OSCommand is a command runner for OS commands
type OSCommand struct {
	Log      *logrus.Entry
	Platform *Platform
	Config   *config.AppConfig
}

// NewOSCommand creates a new OS command
func NewOSCommand(log *logrus.Entry, config *config.AppConfig) *OSCommand {
	platform := &Platform{
		OS:              "linux",
		Shell:           "bash",
		ShellArg:        "-c",
		OpenCommand:     "xdg-open",
		OpenLinkCommand: "xdg-open",
	}

	return &OSCommand{
		Log:      log,
		Platform: platform,
		Config:   config,
	}
}

// RunCommandWithOutput runs a command and returns output
func (c *OSCommand) RunCommandWithOutput(command string) (string, error) {
	cmd := exec.Command(c.Platform.Shell, c.Platform.ShellArg, command)
	output, err := cmd.CombinedOutput()
	return string(output), err
}

// ExecutableFromString creates an executable command
func (c *OSCommand) ExecutableFromString(commandStr string) *exec.Cmd {
	return exec.Command(c.Platform.Shell, c.Platform.ShellArg, commandStr)
}

// OpenLink opens a URL in browser
func (c *OSCommand) OpenLink(link string) error {
	cmd := exec.Command(c.Platform.OpenLinkCommand, link)
	return cmd.Start()
}

// OpenFile opens a file with default application
func (c *OSCommand) OpenFile(filename string) error {
	cmd := exec.Command(c.Platform.OpenCommand, filename)
	return cmd.Start()
}

// EditFile opens a file in editor
func (c *OSCommand) EditFile(filename string) (*exec.Cmd, error) {
	editor := "vim" // default editor
	if envEditor := os.Getenv("EDITOR"); envEditor != "" {
		editor = envEditor
	}
	return exec.Command(editor, filename), nil
}

// RunCustomCommand runs a custom command
func (c *OSCommand) RunCustomCommand(command string) *exec.Cmd {
	return c.ExecutableFromString(command)
}
