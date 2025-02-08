package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/conormkelly/yts-cli/internal/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Display current configuration",
	Long: `Display the current configuration settings for YTS CLI.
This includes both default values and any user overrides from config file or environment variables.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %v", err)
		}

		// Get config file location
		configDir, err := os.UserConfigDir()
		if err != nil {
			return fmt.Errorf("failed to get config directory: %v", err)
		}
		configPath := filepath.Join(configDir, "yts", "config.json")

		// Print header
		header := "YouTube Transcript Summarizer (YTS) Configuration"
		fmt.Println(header)
		fmt.Println(strings.Repeat("-", len(header)))

		// Display config file location
		fmt.Printf("Config file: %s\n", configPath)
		if !fileExists(configPath) {
			fmt.Println("Status: Using default configuration (no config file found)")
		} else {
			fmt.Println("Status: Configuration file loaded")
		}
		fmt.Println()

		// Display provider settings
		fmt.Println("Provider Settings:")
		fmt.Printf("  Active Provider: %s\n", cfg.Provider)
		if isEnvOverride("YTS_PROVIDER") {
			fmt.Println("    └─ (set by YTS_PROVIDER environment variable)")
		}

		// LM Studio settings
		fmt.Println("\n  LM Studio:")
		fmt.Printf("    Base URL: %s\n", cfg.Providers.LMStudio.BaseURL)
		if isEnvOverride("YTS_LMSTUDIO_URL") {
			fmt.Println("      └─ (set by YTS_LMSTUDIO_URL environment variable)")
		}
		fmt.Printf("    Model: %s\n", cfg.Providers.LMStudio.Model)
		if isEnvOverride("YTS_LMSTUDIO_MODEL") {
			fmt.Println("      └─ (set by YTS_LMSTUDIO_MODEL environment variable)")
		}

		// Ollama settings
		fmt.Println("\n  Ollama:")
		fmt.Printf("    Base URL: %s\n", cfg.Providers.Ollama.BaseURL)
		if isEnvOverride("YTS_OLLAMA_URL") {
			fmt.Println("      └─ (set by YTS_OLLAMA_URL environment variable)")
		}
		fmt.Printf("    Model: %s\n", cfg.Providers.Ollama.Model)
		if isEnvOverride("YTS_OLLAMA_MODEL") {
			fmt.Println("      └─ (set by YTS_OLLAMA_MODEL environment variable)")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(configCmd)
}

// fileExists checks if a file exists and is not a directory
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// isEnvOverride checks if an environment variable is set
func isEnvOverride(envVar string) bool {
	return viper.GetString(strings.ToLower(strings.TrimPrefix(envVar, "YTS_"))) != viper.GetViper().GetString(strings.ToLower(strings.TrimPrefix(envVar, "YTS_")))
}
