package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

// Config holds all configuration values
type Config struct {
	Provider    string           `mapstructure:"provider"`  // Current provider (lmstudio, ollama)
	Providers   ProvidersConfig  `mapstructure:"providers"` // Provider-specific configs
	Summaries   SummaryConfig    `mapstructure:"summaries"`
	Transcripts TranscriptConfig `mapstructure:"transcripts"`
}

// ProvidersConfig holds settings for each provider
type ProvidersConfig struct {
	LMStudio LMStudioConfig `mapstructure:"lmstudio"`
	Ollama   OllamaConfig   `mapstructure:"ollama"`
}

type LMStudioConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type OllamaConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

// SummaryConfig holds the different summary templates
type SummaryConfig struct {
	Short SummaryTemplate `mapstructure:"short"`
	Long  SummaryTemplate `mapstructure:"long"`
}

type SummaryTemplate struct {
	SystemPrompt string `mapstructure:"system_prompt"`
}

type TranscriptConfig struct {
	SystemPrompt string `mapstructure:"system_prompt"`
}

const (
	defaultProvider      = "lmstudio"
	defaultLMStudioURL   = "http://localhost:1234"
	defaultLMStudioModel = "llama-3.2-3b-instruct"
	defaultOllamaURL     = "http://localhost:11434"
	defaultOllamaModel   = "llama3.2"
	configFileName       = "config"
	configFileType       = "json"
	configDirName        = "yts"
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
	// Global defaults
	viper.SetDefault("provider", defaultProvider)

	// Provider-specific defaults
	viper.SetDefault("providers.lmstudio.base_url", defaultLMStudioURL)
	viper.SetDefault("providers.lmstudio.model", defaultLMStudioModel)
	viper.SetDefault("providers.ollama.base_url", defaultOllamaURL)
	viper.SetDefault("providers.ollama.model", defaultOllamaModel)

	// Set summary template defaults
	viper.SetDefault("summaries.short.system_prompt", `Create a concise summary of the following transcript. Focus on:
- Core message in 1-2 sentences
- 3-5 key points that support or develop the core message
- If applicable, note any specific calls to action or main conclusions
Keep total length under 150 words.`)

	// Long summary
	viper.SetDefault("summaries.long.system_prompt", `Create a detailed analysis of the following transcript that preserves the original context and depth while making it accessible. Structure as follows:

1. Executive Summary (3-4 sentences)
2. Context and Background
   - Identify the apparent purpose/context
   - Note any assumed knowledge or prerequisites
3. Main Content Analysis
   - Break down major themes and arguments
   - Highlight key terminology and concepts
   - Connect related ideas and show progression
4. Evidence and Support
   - Note specific examples, data, or case studies
   - Identify methodologies or frameworks used
5. Implications and Conclusions
   - Summarize main takeaways
   - Note potential applications or next steps

Preserve technical accuracy while ensuring readability. Include relevant quotes when they significantly support key points.`)

	viper.SetDefault("transcripts.system_prompt", `Format the following raw transcript text.
- Add appropriate capitalization and punctuation
- Keep all original words exactly as they appear
- Never add any additional commentary
- Do not correct spelling or grammar
- Add paragraph breaks where appropriate
- Do not otherwise modify the content in any way`)
}

func bindEnvVars() {
	viper.BindEnv("provider", "YTS_PROVIDER")

	// LM Studio env vars
	viper.BindEnv("providers.lmstudio.base_url", "YTS_LMSTUDIO_URL")
	viper.BindEnv("providers.lmstudio.model", "YTS_LMSTUDIO_MODEL")

	// Ollama env vars
	viper.BindEnv("providers.ollama.base_url", "YTS_OLLAMA_URL")
	viper.BindEnv("providers.ollama.model", "YTS_OLLAMA_MODEL")
}

// GetSystemPrompt returns the appropriate system prompt based on summary type
func GetSystemPrompt(summaryType string) string {
	cfg, err := GetConfig()
	if err != nil {
		return ""
	}

	switch summaryType {
	case "long":
		return cfg.Summaries.Long.SystemPrompt
	default:
		return cfg.Summaries.Short.SystemPrompt
	}
}

// GetActiveProvider returns the config for the currently selected provider
func (c *Config) GetActiveProvider() (baseURL string, model string, err error) {
	switch c.Provider {
	case "lmstudio":
		return c.Providers.LMStudio.BaseURL, c.Providers.LMStudio.Model, nil
	case "ollama":
		return c.Providers.Ollama.BaseURL, c.Providers.Ollama.Model, nil
	default:
		return "", "", fmt.Errorf("unsupported provider: %s", c.Provider)
	}
}
