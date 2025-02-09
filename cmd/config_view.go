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

		keyManager := config.NewAPIKeyManager()
		isClaudeAPIKeySet := keyManager.HasAPIKey("claude")

		// Print header with nice formatting
		fmt.Printf("YTS Configuration\n%s\n\n", strings.Repeat("=", 17))

		// Config file info
		fmt.Printf("Config File: %s\n", configPath)
		if !fileExists(configPath) {
			fmt.Println("Status: Using default configuration (no config file found)")
		} else {
			fmt.Println("Status: Loaded")
		}
		fmt.Println()

		// Provider settings
		fmt.Printf("Active Provider: %s\n", cfg.Provider)

		// Show available providers
		fmt.Println("\nProvider Settings")
		fmt.Println("│")
		fmt.Println("├── LM Studio")
		fmt.Printf("│   ├── Base URL: %s\n", cfg.Providers.LMStudio.BaseURL)
		fmt.Printf("│   └── Model: %s\n", cfg.Providers.LMStudio.Model)
		fmt.Println("├── Ollama")
		fmt.Printf("│   ├── Base URL: %s\n", cfg.Providers.Ollama.BaseURL)
		fmt.Printf("│   └── Model: %s\n", cfg.Providers.Ollama.Model)
		fmt.Println("└── Claude")
		fmt.Printf("    ├── Model: %s\n", cfg.Providers.Claude.Model)
		fmt.Printf("    ├── Temperature: %.1f\n", cfg.Providers.Claude.Temperature)
		fmt.Printf("    ├── Max Tokens: %d\n", cfg.Providers.Claude.MaxTokens)
		fmt.Printf("    ├── Timeout: %d seconds\n", cfg.Providers.Claude.TimeoutSecs)
		fmt.Printf("    ├── Max Retries: %d\n", cfg.Providers.Claude.MaxRetries)
		fmt.Printf("    └── API Key Set: %v\n", isClaudeAPIKeySet)

		return nil
	},
}
