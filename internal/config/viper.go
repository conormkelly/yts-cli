// internal/config/viper.go
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

func InitializeViper() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	ytsConfigDir := filepath.Join(configDir, "yts")
	if err := os.MkdirAll(ytsConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(ytsConfigDir)

	// Set defaults
	viper.SetDefault("llm_base_url", defaultLLMURL)
	viper.SetDefault("model", defaultModel)
	viper.SetDefault("output_format", "markdown")
	viper.SetDefault("summary_type", "medium")
	viper.SetDefault("max_retries", 3)
	viper.SetDefault("timeout_seconds", 30)

	// Bind environment variables
	viper.BindEnv("llm_base_url", "YTS_LLM_URL")
	viper.BindEnv("model", "YTS_MODEL")

	// Read config file
	if err := viper.ReadInConfig(); err != nil {
		// Create default config if not exists
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			if err := viper.SafeWriteConfig(); err != nil {
				return fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() *Config {
	return &Config{
		LLMBaseURL:   viper.GetString("llm_base_url"),
		Model:        viper.GetString("model"),
		OutputFormat: viper.GetString("output_format"),
		SummaryType:  viper.GetString("summary_type"),
		MaxRetries:   viper.GetInt("max_retries"),
		Timeout:      viper.GetInt("timeout_seconds"),
	}
}
