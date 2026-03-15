package gui

import (
	"context"
	"fmt"
	"time"

	"github.com/Araryarch/Lumine/pkg/gui/panels"
	"github.com/Araryarch/Lumine/pkg/gui/types"
	"github.com/Araryarch/Lumine/pkg/lumine"
	"github.com/Araryarch/Lumine/pkg/tasks"
	"github.com/Araryarch/Lumine/pkg/utils"
	"github.com/fatih/color"
	"github.com/jesseduffield/gocui"
)

func (gui *Gui) getLumineDatabasesPanel() *panels.SideListPanel[*lumine.Database] {
	return &panels.SideListPanel[*lumine.Database]{
		ContextState: &panels.ContextState[*lumine.Database]{
			GetMainTabs: func() []panels.MainTab[*lumine.Database] {
				return []panels.MainTab[*lumine.Database]{
					{
						Key:    "info",
						Title:  "Database Info",
						Render: gui.renderLumineDatabaseInfo,
					},
					{
						Key:    "logs",
						Title:  "Query Logs",
						Render: gui.renderLumineDatabaseLogs,
					},
					{
						Key:    "slow",
						Title:  "Slow Queries",
						Render: gui.renderLumineDatabaseSlowQueries,
					},
				}
			},
			GetItemContextCacheKey: func(db *lumine.Database) string {
				return "lumine-database-" + db.Name
			},
		},
		ListPanel: panels.ListPanel[*lumine.Database]{
			List: panels.NewFilteredList[*lumine.Database](),
			View: gui.Views.LumineDatabases,
		},
		NoItemsMessage: "No databases",
		Gui:            gui.intoInterface(),
		Sort: func(a *lumine.Database, b *lumine.Database) bool {
			return a.Name < b.Name
		},
		GetTableCells: func(db *lumine.Database) []string {
			activeIndicator := ""
			if db.Name == gui.Orchestrator.DatabaseManager.GetActiveConnection().Database {
				activeIndicator = utils.ColoredString("●", color.FgGreen)
			}

			return []string{
				activeIndicator,
				utils.ColoredString(db.Name, color.FgCyan),
				string(db.Type),
				db.Size,
			}
		},
	}
}

func (gui *Gui) renderLumineDatabaseInfo(db *lumine.Database) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		output := ""
		output += utils.WithPadding("Name: ", 15) + utils.ColoredString(db.Name, color.FgCyan) + "\n"
		output += utils.WithPadding("Type: ", 15) + string(db.Type) + "\n"
		output += utils.WithPadding("Size: ", 15) + db.Size + "\n"
		output += utils.WithPadding("Created: ", 15) + db.CreatedAt.Format("2006-01-02 15:04:05") + "\n"

		activeConn := gui.Orchestrator.DatabaseManager.GetActiveConnection()
		if db.Name == activeConn.Database {
			output += "\n" + utils.ColoredString("● Active Connection", color.FgGreen) + "\n"
			output += utils.WithPadding("Host: ", 15) + activeConn.Host + "\n"
			output += utils.WithPadding("Port: ", 15) + fmt.Sprintf("%d", activeConn.Port) + "\n"
			output += utils.WithPadding("User: ", 15) + activeConn.Username + "\n"
		}

		return output
	})
}

func (gui *Gui) renderLumineDatabaseLogs(db *lumine.Database) tasks.TaskFunc {
	return gui.NewTickerTask(TickerTaskOpts{
		Func: func(ctx context.Context, notifyStopped chan struct{}) {
			logs, err := gui.Orchestrator.DatabaseManager.GetQueryLogs(db.Name, 50)
			if err != nil {
				gui.RenderStringMain(fmt.Sprintf("Error fetching logs: %v", err))
				return
			}

			output := utils.ColoredString("Recent Queries:\n\n", color.FgYellow)
			for _, log := range logs {
				output += fmt.Sprintf("[%s] %s (%.2fms)\n",
					log.Timestamp.Format("15:04:05"),
					log.Query,
					log.Duration.Seconds()*1000,
				)
			}

			gui.reRenderStringMain(output)
		},
		Duration:   time.Second * 2,
		Before:     func(ctx context.Context) { gui.clearMainView() },
		Wrap:       true,
		Autoscroll: true,
	})
}

func (gui *Gui) renderLumineDatabaseSlowQueries(db *lumine.Database) tasks.TaskFunc {
	return gui.NewSimpleRenderStringTask(func() string {
		logs, err := gui.Orchestrator.DatabaseManager.GetSlowQueries(db.Name, 100*time.Millisecond, 20)
		if err != nil {
			return fmt.Sprintf("Error fetching slow queries: %v", err)
		}

		if len(logs) == 0 {
			return "No slow queries found (threshold: 100ms)"
		}

		output := utils.ColoredString("Slow Queries (>100ms):\n\n", color.FgRed)
		for _, log := range logs {
			output += fmt.Sprintf("[%s] %.2fms\n%s\n\n",
				log.Timestamp.Format("15:04:05"),
				log.Duration.Seconds()*1000,
				log.Query,
			)
		}

		return output
	})
}

// Keybinding handlers for Lumine databases
func (gui *Gui) handleLumineDatabaseCreate(g *gocui.Gui, v *gocui.View) error {
	return gui.createPromptPanel("Database Name", func(g *gocui.Gui, v *gocui.View) error {
		dbName := gui.trimmedContent(v)
		if dbName == "" {
			return gui.createErrorPanel("Database name cannot be empty")
		}

		return gui.WithWaitingStatus("Creating database...", func() error {
			if err := gui.Orchestrator.CreateDatabase(dbName); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineDatabases()
		})
	})
}

func (gui *Gui) handleLumineDatabaseDrop(g *gocui.Gui, v *gocui.View) error {
	db, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.createConfirmationPanel("Confirm", fmt.Sprintf("Drop database '%s'? This cannot be undone.", db.Name), func(g *gocui.Gui, v *gocui.View) error {
		return gui.WithWaitingStatus("Dropping database...", func() error {
			if err := gui.Orchestrator.DropDatabase(db.Name); err != nil {
				return gui.createErrorPanel(err.Error())
			}
			return gui.refreshLumineDatabases()
		})
	}, nil)
}

func (gui *Gui) handleLumineDatabaseBackup(g *gocui.Gui, v *gocui.View) error {
	db, err := gui.Panels.LumineDatabases.GetSelectedItem()
	if err != nil {
		return nil
	}

	return gui.WithWaitingStatus("Backing up database...", func() error {
		if err := gui.Orchestrator.BackupDatabase(db.Name); err != nil {
			return gui.createErrorPanel(err.Error())
		}
		return nil
	})
}

func (gui *Gui) handleLumineDatabaseSwitch(g *gocui.Gui, v *gocui.View) error {
	profiles := gui.Orchestrator.DatabaseManager.ListProfiles()

	menuItems := make([]*types.MenuItem, len(profiles))
	for i, profile := range profiles {
		p := profile // capture for closure
		menuItems[i] = &types.MenuItem{
			LabelColumns: []string{p.Name, string(p.Type)},
			OnPress: func() error {
				return gui.WithWaitingStatus("Switching connection...", func() error {
					if err := gui.Orchestrator.DatabaseManager.SwitchConnection(p.Name); err != nil {
						return gui.createErrorPanel(err.Error())
					}
					gui.Orchestrator.NotificationMgr.ShowSuccess(fmt.Sprintf("Switched to %s", p.Name))
					return gui.refreshLumineDatabases()
				})
			},
		}
	}

	return gui.Menu(CreateMenuOptions{
		Title: "Switch Database Connection",
		Items: menuItems,
	})
}

func (gui *Gui) refreshLumineDatabases() error {
	if gui.Orchestrator == nil || gui.Panels.LumineDatabases == nil {
		return nil
	}

	databases, err := gui.Orchestrator.DatabaseManager.ListDatabases()
	if err != nil {
		gui.Log.Error(err)
		return nil
	}

	gui.Panels.LumineDatabases.SetItems(databases)

	return gui.Panels.LumineDatabases.RerenderList()
}
