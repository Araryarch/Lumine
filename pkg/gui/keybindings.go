package gui

import (
	"fmt"

	"github.com/jesseduffield/gocui"
)

func (gc *GuiController) openProject(g *gocui.Gui, v *gocui.View) error {
	if gc.CurrentView == ViewProjects && len(gc.ProjectList) > 0 && gc.SelectedIdx < len(gc.ProjectList) {
		project := gc.ProjectList[gc.SelectedIdx]
		gc.showMessage(fmt.Sprintf("Open: %s", project.URL), "info")
	}
	return nil
}

func (gc *GuiController) showHelp(g *gocui.Gui, v *gocui.View) error {
	gc.showMessage("j/k: navigate | Enter: select | b: back | q: quit | ?: help", "info")
	return nil
}
