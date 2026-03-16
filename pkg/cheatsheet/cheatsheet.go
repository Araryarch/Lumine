package cheatsheet

import (
	"fmt"
	"os"
	"path/filepath"
)

// Generate generates cheatsheet documentation
func Generate() {
	fmt.Println("Generating cheatsheets...")
	// TODO: Implement cheatsheet generation
}

// Check checks cheatsheet validity
func Check() {
	fmt.Println("Checking cheatsheets...")
	// TODO: Implement cheatsheet checking
}

// GetKeybindingsDir returns the keybindings directory
func GetKeybindingsDir() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".lumine", "keybindings")
}
