package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	// LLM Configuration
	LLMBaseURL string `json:"llm_base_url"`
	Model      string `json:"model"`

	// Output Configuration
	OutputFormat string `json:"output_format"`
	SummaryType  string `json:"summary_type"` // "short", "medium", "long"

	// API Configuration
	MaxRetries int `json:"max_retries"`
	Timeout    int `json:"timeout_seconds"`
}

const (
	defaultConfigFile = "config.json"
	defaultLLMURL     = "http://localhost:1234"
	defaultModel      = "llama-3.2-3b-instruct"
)

var DefaultConfig = Config{
	LLMBaseURL:   defaultLLMURL,
	Model:        defaultModel,
	OutputFormat: "markdown",
	SummaryType:  "medium",
	MaxRetries:   3,
	Timeout:      30,
}

// Load reads the configuration file from the user's config directory
func Load() (*Config, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config directory: %w", err)
	}

	ytsConfigDir := filepath.Join(configDir, "yts")
	if err := os.MkdirAll(ytsConfigDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	configPath := filepath.Join(ytsConfigDir, defaultConfigFile)

	// If config doesn't exist, create it with defaults
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return createDefaultConfig(configPath)
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

// Save writes the current configuration to disk
func (c *Config) Save() error {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	configPath := filepath.Join(configDir, "yts", defaultConfigFile)

	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func createDefaultConfig(path string) (*Config, error) {
	config := DefaultConfig

	data, err := json.MarshalIndent(config, "", "    ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write default config: %w", err)
	}

	return &config, nil
}
