package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var editCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration file in your default text editor",
	RunE: func(cmd *cobra.Command, args []string) error {
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %v", err)
		}

		configPath := filepath.Join(configDir, "yts", "config.json")

		// Ensure config file exists
		if !fileExists(configPath) {
			// Initialize with defaults if it doesn't exist
			if err := initializeConfig(configPath); err != nil {
				return fmt.Errorf("failed to initialize config: %v", err)
			}
		}

		editor := getEditor()
		if editor == "" {
			return fmt.Errorf("no suitable text editor found")
		}

		// Launch editor
		editorCmd := exec.Command(editor, configPath)
		editorCmd.Stdin = os.Stdin
		editorCmd.Stdout = os.Stdout
		editorCmd.Stderr = os.Stderr

		return editorCmd.Run()
	},
}

func getEditor() string {
	// Check environment variables first
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor
	}

	// Platform-specific defaults
	switch runtime.GOOS {
	case "windows":
		return "notepad"
	default:
		// Try common editors on Unix-like systems
		editors := []string{"nano", "vim", "vi"}
		for _, editor := range editors {
			if _, err := exec.LookPath(editor); err == nil {
				return editor
			}
		}
	}

	return ""
}

func initializeConfig(path string) error {
	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Write default config if it doesn't exist
	if err := viper.SafeWriteConfigAs(path); err != nil {
		return fmt.Errorf("failed to write default config: %v", err)
	}

	return nil
}

// Helper function used by multiple commands
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
