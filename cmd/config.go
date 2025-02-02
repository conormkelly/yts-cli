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

		// Display configuration values
		fmt.Println("Current Settings:")
		fmt.Printf("  LLM Base URL: %s\n", cfg.LLMBaseURL)
		if isEnvOverride("YTS_LLM_URL") {
			fmt.Println("    └─ (set by YTS_LLM_URL environment variable)")
		}

		fmt.Printf("  Model: %s\n", cfg.Model)
		if isEnvOverride("YTS_MODEL") {
			fmt.Println("    └─ (set by YTS_MODEL environment variable)")
		}

		fmt.Printf("  Default Summary Type: %s\n", cfg.SummaryType)

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
