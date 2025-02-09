package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/conormkelly/yts-cli/internal/config"
	"github.com/spf13/cobra"
)

var viewCmd = &cobra.Command{
	Use:   "view",
	Short: "View current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %v", err)
		}

		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %v", err)
		}
		configPath := filepath.Join(configDir, "yts", "config.json")

		// Print header with nice formatting
		fmt.Printf("YTS Configuration\n%s\n\n", strings.Repeat("=", 17))

		// Config file info
		fmt.Printf("Config File: %s\n", configPath)
		if !fileExists(configPath) {
			fmt.Println("Status: Using default configuration (no config file found)")
		} else {
			fmt.Println("Status: Configuration file loaded")
		}
		fmt.Println()

		// Provider settings
		fmt.Println("Active Provider Settings")
		fmt.Printf("├── Provider: %s\n", cfg.Provider)
		if cfg.Provider == "lmstudio" {
			fmt.Printf("├── Base URL: %s\n", cfg.Providers.LMStudio.BaseURL)
			fmt.Printf("└── Model: %s\n", cfg.Providers.LMStudio.Model)
		} else {
			fmt.Printf("├── Base URL: %s\n", cfg.Providers.Ollama.BaseURL)
			fmt.Printf("└── Model: %s\n", cfg.Providers.Ollama.Model)
		}
		fmt.Println()

		// Show available providers
		fmt.Println("Available Providers")
		fmt.Println("├── LM Studio")
		fmt.Printf("│   ├── Base URL: %s\n", cfg.Providers.LMStudio.BaseURL)
		fmt.Printf("│   └── Model: %s\n", cfg.Providers.LMStudio.Model)
		fmt.Println("└── Ollama")
		fmt.Printf("    ├── Base URL: %s\n", cfg.Providers.Ollama.BaseURL)
		fmt.Printf("    └── Model: %s\n", cfg.Providers.Ollama.Model)

		return nil
	},
}
