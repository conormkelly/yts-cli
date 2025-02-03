package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration values
type Config struct {
	LLMBaseURL  string           `mapstructure:"llm_base_url"`
	Model       string           `mapstructure:"model"`
	SummaryType string           `mapstructure:"summary_type"`
	Summaries   SummaryConfig    `mapstructure:"summaries"`
	Transcripts TranscriptConfig `mapstructure:"transcripts"`
}

// SummaryConfig holds the different summary templates
type SummaryConfig struct {
	Short  SummaryTemplate `mapstructure:"short"`
	Medium SummaryTemplate `mapstructure:"medium"`
	Long   SummaryTemplate `mapstructure:"long"`
}

type SummaryTemplate struct {
	SystemPrompt string `mapstructure:"system_prompt"`
}

type TranscriptConfig struct {
	SystemPrompt string `mapstructure:"system_prompt"`
}

const (
	defaultLLMURL  = "http://localhost:1234"
	defaultModel   = "llama-3.2-3b-instruct"
	configFileName = "config"
	configFileType = "json"
	configDirName  = "yts"
)

// Initialize sets up Viper with our configuration
func Initialize() error {
	// Get user config directory
	configDir, err := os.UserConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config directory: %w", err)
	}

	// Ensure our config directory exists
	ytsConfigDir := filepath.Join(configDir, configDirName)
	if err := os.MkdirAll(ytsConfigDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Set up Viper
	viper.SetConfigName(configFileName)
	viper.SetConfigType(configFileType)
	viper.AddConfigPath(ytsConfigDir)

	// Set defaults
	setDefaults()

	// Bind environment variables
	bindEnvVars()

	// Try to read config file
	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file doesn't exist, create it with defaults
			configPath := filepath.Join(ytsConfigDir, configFileName+"."+configFileType)
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	}

	return nil
}

// GetConfig returns the current configuration
func GetConfig() (*Config, error) {
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	return &config, nil
}

func setDefaults() {
	viper.SetDefault("llm_base_url", defaultLLMURL)
	viper.SetDefault("model", defaultModel)
	viper.SetDefault("output_format", "markdown")
	viper.SetDefault("summary_type", "medium")

	// Set summary template defaults
	viper.SetDefault("summaries.short.system_prompt", `Create a concise summary of the following video transcript. Focus on:
- Key points only (3-5 bullet points)
- Main conclusion or takeaway
- Keep it brief and to the point`)

	viper.SetDefault("summaries.medium.system_prompt", `Provide a comprehensive summary of the following video transcript. Include:
- Main topics and key points
- Important details and insights
- Supporting examples or evidence
- Organize with clear headings`)

	viper.SetDefault("summaries.long.system_prompt", `Create a detailed analysis of the following video transcript. Include:
- Thorough coverage of all major topics
- Detailed examples and supporting information
- Analysis of key concepts and their relationships
- Clear structure with sections and subsections
- Any relevant technical details or specifications`)

	viper.SetDefault("transcripts.system_prompt", `Format the following raw transcript text.
- Add appropriate capitalization and punctuation
- Keep all original words exactly as they appear
- Never add any additional commentary
- Do not correct spelling or grammar
- Add paragraph breaks where appropriate
- Do not otherwise modify the content in any way`)
}

func bindEnvVars() {
	viper.BindEnv("llm_base_url", "YTS_LLM_URL")
	viper.BindEnv("model", "YTS_MODEL")
}

// GetSystemPrompt returns the appropriate system prompt based on summary type
func GetSystemPrompt(summaryType string) string {
	cfg, err := GetConfig()
	if err != nil {
		return "" // Handle error appropriately in your application
	}

	switch summaryType {
	case "short":
		return cfg.Summaries.Short.SystemPrompt
	case "long":
		return cfg.Summaries.Long.SystemPrompt
	default:
		return cfg.Summaries.Medium.SystemPrompt
	}
}
