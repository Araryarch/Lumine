package gui

import (
	"fmt"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/jesseduffield/lazydocker/pkg/utils"
)

func (gui *Gui) refreshNotifications() error {
	if gui.Orchestrator == nil {
		return nil
	}
	
	// Notifications are automatically managed by the NotificationManager
	// This function just triggers a re-render if needed
	return nil
}

func (gui *Gui) renderNotifications() string {
	if gui.Orchestrator == nil {
		return ""
	}
	
	notifications := gui.Orchestrator.NotificationMgr.GetActive()
	if len(notifications) == 0 {
		return ""
	}
	
	// Show only the most recent notification
	notification := notifications[len(notifications)-1]
	
	// Calculate time since notification
	elapsed := time.Since(notification.Timestamp)
	if elapsed > 5*time.Second {
		// Auto-dismiss after 5 seconds
		return ""
	}
	
	var icon string
	var colorAttr color.Attribute
	
	switch notification.Type {
	case "success":
		icon = "✓"
		colorAttr = color.FgGreen
	case "error":
		icon = "✗"
		colorAttr = color.FgRed
	case "warning":
		icon = "⚠"
		colorAttr = color.FgYellow
	case "info":
		icon = "ℹ"
		colorAttr = color.FgCyan
	default:
		icon = "•"
		colorAttr = color.FgWhite
	}
	
	// Format: [icon] message
	message := fmt.Sprintf("%s %s", icon, notification.Message)
	return utils.ColoredString(message, colorAttr)
}

// Helper function to display notification in status bar or overlay
func (gui *Gui) displayNotification(message string, notifType string) {
	if gui.Orchestrator == nil {
		return
	}
	
	switch notifType {
	case "success":
		gui.Orchestrator.NotificationMgr.ShowSuccess(message)
	case "error":
		gui.Orchestrator.NotificationMgr.ShowError(message)
	case "warning":
		gui.Orchestrator.NotificationMgr.ShowWarning(message)
	case "info":
		gui.Orchestrator.NotificationMgr.ShowInfo(message)
	}
}

// Update the information view to include notifications
func (gui *Gui) getInformationContentWithNotifications() string {
	baseInfo := gui.getInformationContent()
	
	notification := gui.renderNotifications()
	if notification == "" {
		return baseInfo
	}
	
	// Add notification to the right side of information bar
	padding := strings.Repeat(" ", 5)
	return baseInfo + padding + notification
}
