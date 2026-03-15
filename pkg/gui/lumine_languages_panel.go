package gui

import (
	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/tasks"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getLumineLanguagesPanel() *panels.SideListPanel[*lumine.Tool] {
	return &panels.SideListPanel[*lumine.Tool]{
		ContextState: &panels.ContextState[*lumine.Tool]{
			GetMainTabs: func() []panels.MainTab[*lumine.Tool] {
				return []panels.MainTab[*lumine.Tool]{
					{
						Key:    "info",
						Title:  "Tool Info",
						Render: gui.renderLumineToolInfo,
					},
					{
						Key:    "usage",
						Title:  "Usage",
						Render: gui.renderLumineToolUsage,
					},
				}
			},
			GetItemContextCacheKey: func(tool *lumine.Tool) string {
				return "lumine-tool-" + tool.Name + "-" + tool.Version
			},
		},
		ListPanel: panels.ListPanel[*lumine.Tool]{
			List: panels.NewFilteredList[*lumine.Tool](),
			View: gui.Views.LumineLanguages,
		},
		NoItemsMessage: "No development tools",
		Gui:            gui.intoInterface(),
		Sort: func(a *lumine.Tool, b *lumine.Tool) bool {
			if a.Installed && !b.Installed {
				return true
			}
			if !a.Installed && b.Installed {
				return false
			}
			return a.Name < b.Name
		},
		GetTableCells: func(tool *lumine.Tool) []string {
			statusText := "not installed"
			statusColor := color.FgRed
			if tool.Installed {
				statusText = "installed"
				statusColor = color.FgGreen
			}

			return []string{
				utils.ColoredString(tool.DisplayName, color.FgCyan),
				tool.Version,
				utils.ColoredString(statusText, statusColor),
			}
		},
	}
}

func (gui *Gui) renderLumineToolInfo(tool *lumine.Tool) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(tool.DisplayName, color.FgCyan) + "\n"
		output += utils.WithPadding("Type: ", 15) + tool.Type + "\n"
		output += utils.WithPadding("Version: ", 15) + tool.Version + "\n"
		
		statusText := "Not Installed"
		statusColor := color.FgRed
		if tool.Installed {
			statusText = "Installed"
			statusColor = color.FgGreen
		}
		output += utils.WithPadding("Status: ", 15) + utils.ColoredString(statusText, statusColor) + "\n"
		
		if tool.Description != "" {
			output += "\n" + tool.Description + "\n"
		}

		return output
	})
}

func (gui *Gui) renderLumineToolUsage(tool *lumine.Tool) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		if !tool.Installed {
			return "Tool is not installed.\n\nInstall instructions:\n" + gui.getToolInstallInstructions(tool)
		}
		
		return gui.getToolUsageInfo(tool)
	})
}

func (gui *Gui) getToolInstallInstructions(tool *lumine.Tool) string {
	instructions := map[string]string{
		"php":      "Install PHP via your package manager:\n  Ubuntu/Debian: sudo apt install php\n  Fedora: sudo dnf install php\n  Arch: sudo pacman -S php",
		"composer": "Install Composer:\n  curl -sS https://getcomposer.org/installer | php\n  sudo mv composer.phar /usr/local/bin/composer",
		"node":     "Install Node.js via nvm or package manager:\n  nvm: nvm install 20\n  Ubuntu: sudo apt install nodejs\n  Or download from: https://nodejs.org",
		"npm":      "npm comes with Node.js installation",
		"yarn":     "Install Yarn:\n  npm install -g yarn",
		"pnpm":     "Install pnpm:\n  npm install -g pnpm",
		"bun":      "Install Bun:\n  curl -fsSL https://bun.sh/install | bash",
		"deno":     "Install Deno:\n  curl -fsSL https://deno.land/install.sh | sh",
		"python":   "Install Python via your package manager:\n  Ubuntu/Debian: sudo apt install python3\n  Fedora: sudo dnf install python3",
		"pip":      "pip comes with Python installation",
		"poetry":   "Install Poetry:\n  curl -sSL https://install.python-poetry.org | python3 -",
		"pipenv":   "Install Pipenv:\n  pip3 install --user pipenv",
	}
	
	if inst, ok := instructions[tool.Name]; ok {
		return inst
	}
	return "No installation instructions available"
}

func (gui *Gui) getToolUsageInfo(tool *lumine.Tool) string {
	usage := map[string]string{
		"php":      "PHP CLI:\n  php -v          # Check version\n  php script.php  # Run script\n  php -S localhost:8000  # Built-in server",
		"composer": "Composer:\n  composer install      # Install dependencies\n  composer require pkg  # Add package\n  composer update       # Update packages",
		"node":     "Node.js:\n  node -v         # Check version\n  node script.js  # Run script",
		"npm":      "npm:\n  npm install     # Install dependencies\n  npm install pkg # Add package\n  npm run dev     # Run script",
		"yarn":     "Yarn:\n  yarn install    # Install dependencies\n  yarn add pkg    # Add package\n  yarn dev        # Run script",
		"pnpm":     "pnpm:\n  pnpm install    # Install dependencies\n  pnpm add pkg    # Add package\n  pnpm dev        # Run script",
		"bun":      "Bun:\n  bun install     # Install dependencies\n  bun add pkg     # Add package\n  bun run dev     # Run script",
		"deno":     "Deno:\n  deno run script.ts  # Run script\n  deno task dev       # Run task",
		"python":   "Python:\n  python3 -V      # Check version\n  python3 script.py  # Run script",
		"pip":      "pip:\n  pip3 install pkg    # Install package\n  pip3 install -r requirements.txt",
		"poetry":   "Poetry:\n  poetry install  # Install dependencies\n  poetry add pkg  # Add package\n  poetry run python script.py",
		"pipenv":   "Pipenv:\n  pipenv install  # Install dependencies\n  pipenv install pkg  # Add package\n  pipenv run python script.py",
	}
	
	if info, ok := usage[tool.Name]; ok {
		return info
	}
	return "No usage information available"
}

// Keybinding handlers
func (gui *Gui) handleLumineLanguageRefresh(g *gocui.Gui, v *gocui.View) error {
	tool, err := gui.Panels.LumineLanguages.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.WithWaitingStatus("Checking installation...", func() error {
		if err := gui.Orchestrator.ToolManager.RefreshToolStatus(tool.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return gui.refreshLumineLanguages()
	})
}

func (gui *Gui) refreshLumineLanguages() error {
	if gui.Orchestrator == nil || gui.Panels.LumineLanguages == nil {
		return nil
	}

	tools := gui.Orchestrator.ToolManager.GetAllTools()
	gui.Panels.LumineLanguages.SetItems(tools)
	return gui.Panels.LumineLanguages.RerenderList()
}
