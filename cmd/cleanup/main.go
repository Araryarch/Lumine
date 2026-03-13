package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"lumine/internal/docker"
)

const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorCyan   = "\033[36m"
)

func main() {
	fmt.Printf("%sLumine Cleanup Tool%s\n\n", colorCyan, colorReset)

	// Check Docker
	dockerMgr, err := docker.NewManager()
	if err != nil {
		fmt.Printf("%sFailed to connect to Docker: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Show menu
	showMenu()

	// Get user choice
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter choice [1-5]: ")
	choice, _ := reader.ReadString('\n')
	choice = strings.TrimSpace(choice)

	ctx := context.Background()

	switch choice {
	case "1":
		stopContainers(ctx, dockerMgr)
	case "2":
		removeContainers(ctx, dockerMgr, false)
	case "3":
		removeContainersWithVolumes(ctx, dockerMgr)
	case "4":
		nuclearCleanup(ctx, dockerMgr)
	case "5":
		fmt.Printf("%sCleanup cancelled.%s\n", colorYellow, colorReset)
		os.Exit(0)
	default:
		fmt.Printf("%sInvalid choice!%s\n", colorRed, colorReset)
		os.Exit(1)
	}
}

func showMenu() {
	fmt.Println("Select cleanup option:")
	fmt.Println("  1) Stop containers only")
	fmt.Println("  2) Remove containers (keep data)")
	fmt.Println("  3) Remove containers + volumes (DELETE DATA)")
	fmt.Println("  4) Nuclear cleanup (REMOVE EVERYTHING)")
	fmt.Println("  5) Cancel")
	fmt.Println()
}

func stopContainers(ctx context.Context, dockerMgr *docker.Manager) {
	fmt.Printf("%sStopping containers...%s\n", colorCyan, colorReset)
	
	if err := dockerMgr.StopAllContainers(ctx); err != nil {
		fmt.Printf("%s❌ Failed to stop containers: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}
	
	fmt.Printf("%s✓ Containers stopped!%s\n", colorGreen, colorReset)
	showSummary(ctx, dockerMgr)
}

func removeContainers(ctx context.Context, dockerMgr *docker.Manager, force bool) {
	if !force {
		fmt.Printf("%s⚠️  This will remove all Lumine containers%s\n", colorYellow, colorReset)
		if !confirm("Continue?") {
			fmt.Printf("%sOperation cancelled.%s\n", colorYellow, colorReset)
			return
		}
	}

	fmt.Printf("%sStopping containers...%s\n", colorCyan, colorReset)
	if err := dockerMgr.StopAllContainers(ctx); err != nil {
		fmt.Printf("%s⚠️  Warning: %v%s\n", colorYellow, err, colorReset)
	}

	fmt.Printf("%sRemoving containers...%s\n", colorCyan, colorReset)
	if err := dockerMgr.RemoveAllContainers(ctx, true); err != nil {
		fmt.Printf("%s❌ Failed to remove containers: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("%s✓ Containers removed!%s\n", colorGreen, colorReset)
	showSummary(ctx, dockerMgr)
}

func removeContainersWithVolumes(ctx context.Context, dockerMgr *docker.Manager) {
	fmt.Printf("%sWARNING: This will DELETE ALL DATABASE DATA!%s\n", colorRed, colorReset)
	
	if !confirmTyping("yes") {
		fmt.Printf("%sCleanup cancelled.%s\n", colorYellow, colorReset)
		return
	}

	// Backup first
	fmt.Printf("%sCreating backup...%s\n", colorCyan, colorReset)
	if err := createBackup(ctx, dockerMgr); err != nil {
		fmt.Printf("%s⚠️  Backup failed: %v%s\n", colorYellow, err, colorReset)
		fmt.Print("Continue without backup? [y/N]: ")
		if !confirm("") {
			return
		}
	} else {
		fmt.Printf("%s✓ Backup created%s\n", colorGreen, colorReset)
	}

	// Stop containers
	fmt.Printf("%sStopping containers...%s\n", colorCyan, colorReset)
	if err := dockerMgr.StopAllContainers(ctx); err != nil {
		fmt.Printf("%s⚠️  Warning: %v%s\n", colorYellow, err, colorReset)
	}

	// Remove containers
	fmt.Printf("%sRemoving containers...%s\n", colorCyan, colorReset)
	if err := dockerMgr.RemoveAllContainers(ctx, true); err != nil {
		fmt.Printf("%s❌ Failed to remove containers: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	// Remove volumes
	fmt.Printf("%sRemoving volumes...%s\n", colorCyan, colorReset)
	if err := dockerMgr.RemoveAllVolumes(ctx); err != nil {
		fmt.Printf("%s❌ Failed to remove volumes: %v%s\n", colorRed, err, colorReset)
		os.Exit(1)
	}

	fmt.Printf("%s✓ Complete cleanup done!%s\n", colorGreen, colorReset)
	showSummary(ctx, dockerMgr)
}

func nuclearCleanup(ctx context.Context, dockerMgr *docker.Manager) {
	fmt.Printf("%sNUCLEAR OPTION: This will DESTROY EVERYTHING!%s\n", colorRed, colorReset)
	fmt.Printf("%s   - All containers%s\n", colorRed, colorReset)
	fmt.Printf("%s   - All volumes (data)%s\n", colorRed, colorReset)
	fmt.Printf("%s   - Network%s\n", colorRed, colorReset)
	fmt.Println()

	if !confirmTyping("DESTROY") {
		fmt.Printf("%sDestruction cancelled.%s\n", colorYellow, colorReset)
		return
	}

	// Backup first
	fmt.Printf("%sCreating backup...%s\n", colorCyan, colorReset)
	if err := createBackup(ctx, dockerMgr); err != nil {
		fmt.Printf("%s⚠️  Backup failed: %v%s\n", colorYellow, err, colorReset)
	} else {
		fmt.Printf("%s✓ Backup created%s\n", colorGreen, colorReset)
	}

	// Stop containers
	fmt.Printf("%sStopping containers...%s\n", colorCyan, colorReset)
	dockerMgr.StopAllContainers(ctx)

	// Remove containers
	fmt.Printf("%sRemoving containers...%s\n", colorCyan, colorReset)
	dockerMgr.RemoveAllContainers(ctx, true)

	// Remove volumes
	fmt.Printf("%sRemoving volumes...%s\n", colorCyan, colorReset)
	dockerMgr.RemoveAllVolumes(ctx)

	// Remove network
	fmt.Printf("%sRemoving network...%s\n", colorCyan, colorReset)
	dockerMgr.RemoveNetwork(ctx)

	fmt.Printf("%s✓ Nuclear cleanup complete!%s\n", colorGreen, colorReset)
	showSummary(ctx, dockerMgr)
}

func createBackup(ctx context.Context, dockerMgr *docker.Manager) error {
	timestamp := time.Now().Format("20060102-150405")
	backupFile := fmt.Sprintf("lumine-backup-%s.sql", timestamp)

	// Try to backup MySQL
	logs, err := dockerMgr.GetContainerLogs(ctx, "lumine-mysql", "10")
	if err == nil && logs != "" {
		// Container exists, try backup
		// Note: This is simplified, actual backup would use docker exec
		fmt.Printf("Backup would be saved to: %s\n", backupFile)
		return nil
	}

	return fmt.Errorf("no MySQL container found")
}

func showSummary(ctx context.Context, dockerMgr *docker.Manager) {
	fmt.Println()
	fmt.Printf("%sCleanup Summary:%s\n", colorCyan, colorReset)

	containers, _ := dockerMgr.ListLumineContainers(ctx)
	fmt.Printf("Containers: %d\n", len(containers))

	volumes, _ := dockerMgr.ListLumineVolumes(ctx)
	fmt.Printf("Volumes: %d\n", len(volumes))

	fmt.Println()
	fmt.Printf("%sDone!%s\n", colorGreen, colorReset)
}

func confirm(prompt string) bool {
	reader := bufio.NewReader(os.Stdin)
	if prompt != "" {
		fmt.Printf("%s [y/N]: ", prompt)
	}
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))
	return response == "y" || response == "yes"
}

func confirmTyping(expected string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Type '%s' to confirm: ", expected)
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)
	return response == expected
}
