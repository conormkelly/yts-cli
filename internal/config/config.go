package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/conormkelly/yts-cli/internal/constants"
	"github.com/spf13/viper"
)

// Current config value
const currentConfigVersion = "1.0.1"

// Config holds all configuration values
type Config struct {
	Version     string           `mapstructure:"version"`
	Provider    string           `mapstructure:"provider"` // Current provider (lmstudio, ollama)
	Providers   ProvidersConfig  `mapstructure:"providers"`
	Summaries   SummaryConfig    `mapstructure:"summaries"`
	Transcripts TranscriptConfig `mapstructure:"transcripts"`
	Queries     QueryConfig      `mapstructure:"queries"`
}

// ProvidersConfig holds settings for each provider
type ProvidersConfig struct {
	LMStudio LMStudioConfig `mapstructure:"lmstudio"`
	Ollama   OllamaConfig   `mapstructure:"ollama"`
	Claude   ClaudeConfig   `mapstructure:"claude"`
	OpenAI   OpenAIConfig   `mapstructure:"openai"`
}

type LMStudioConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type OllamaConfig struct {
	BaseURL string `mapstructure:"base_url"`
	Model   string `mapstructure:"model"`
}

type ClaudeConfig struct {
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	TimeoutSecs int     `mapstructure:"timeout_seconds"`
	MaxRetries  int     `mapstructure:"max_retries"`
}

type OpenAIConfig struct {
	Model       string  `mapstructure:"model"`
	Temperature float64 `mapstructure:"temperature"`
	MaxTokens   int     `mapstructure:"max_tokens"`
	TimeoutSecs int     `mapstructure:"timeout_seconds"`
	MaxRetries  int     `mapstructure:"max_retries"`
	OrgID       string  `mapstructure:"organization_id"`
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

// QueryConfig holds the query template
type QueryConfig struct {
	SystemPrompt string `mapstructure:"system_prompt"`
}

const (
	// Global
	defaultProvider = "lmstudio"
	configFileName  = "config"
	configFileType  = "json"
	configDirName   = "yts"

	defaultLMStudioURL   = "http://localhost:1234"
	defaultLMStudioModel = "llama-3.2-3b-instruct"

	defaultOllamaURL   = "http://localhost:11434"
	defaultOllamaModel = "llama3.2"

	defaultClaudeModel          = "claude-3-5-sonnet-20241022"
	defaultClaudeTemperature    = 0.3  // Lower for more focused summaries
	defaultClaudeMaxTokens      = 8192 // Generous limit for detailed analysis
	defaultClaudeTimeoutSeconds = 120  // Long timeout for big transcripts
	defaultClaudeMaxRetries     = 3

	defaultOpenAIModel          = "gpt-4o" // Latest model for best summaries
	defaultOpenAITemperature    = 0.3      // Lower for more focused summaries
	defaultOpenAIMaxTokens      = 8192     // Conservative limit for most transcripts
	defaultOpenAITimeoutSeconds = 120      // Long timeout for big transcripts
	defaultOpenAIMaxRetries     = 3
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
			setDefaults()
			viper.Set("version", currentConfigVersion)

			configPath := filepath.Join(ytsConfigDir, configFileName+"."+configFileType)
			if err := viper.SafeWriteConfigAs(configPath); err != nil {
				return fmt.Errorf("failed to write default config: %w", err)
			}
		} else {
			return fmt.Errorf("failed to read config: %w", err)
		}
	} else {
		// Config exists, check for migrations
		var cfg Config
		if err := viper.Unmarshal(&cfg); err != nil {
			return fmt.Errorf("failed to parse config: %w", err)
		}

		// Apply migrations if needed
		if cfg.Version != currentConfigVersion {
			applyMigrations(&cfg)

			// Update viper with migrated config
			for k, v := range structToMap(cfg) {
				viper.Set(k, v)
			}

			// Save the updated config
			if err := viper.WriteConfig(); err != nil {
				return fmt.Errorf("failed to save migrated config: %w", err)
			}

			fmt.Fprintf(os.Stderr, "Configuration updated to version %s\n", currentConfigVersion)
		}
	}

	return nil
}

// applyMigrations updates an existing config to the current version without losing user settings
func applyMigrations(cfg *Config) {
	// If there's no version or it's a new config, just set the current version
	if cfg.Version == "" {
		cfg.Version = currentConfigVersion
		return
	}

	// Handle specific version migrations
	switch cfg.Version {
	case "1.0.0":
		// Migrate from 1.0.0 to 1.1.0
		if cfg.Queries.SystemPrompt == "" {
			cfg.Queries = QueryConfig{
				SystemPrompt: constants.QueryPrompt,
			}
		}
		cfg.Version = "1.1.0"
		// fallthrough

		// Add future version migrations here with fallthrough to apply sequentially
		// case "1.1.0":
		//    // Migrate from 1.1.0 to 1.2.0
		//    cfg.Version = "1.2.0"
		//    fallthrough
	}
}

// Helper function to convert struct to map for viper
func structToMap(obj interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	val := reflect.ValueOf(obj)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		tag := typ.Field(i).Tag.Get("mapstructure")
		if tag != "" && tag != "-" {
			if field.Kind() == reflect.Struct {
				// Handle nested structs recursively
				nested := structToMap(field.Interface())
				for k, v := range nested {
					result[tag+"."+k] = v
				}
			} else {
				result[tag] = field.Interface()
			}
		}
	}

	return result
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

	viper.SetDefault("providers.claude.model", defaultClaudeModel)
	viper.SetDefault("providers.claude.temperature", defaultClaudeTemperature)
	viper.SetDefault("providers.claude.max_tokens", defaultClaudeMaxTokens)
	viper.SetDefault("providers.claude.timeout_seconds", defaultClaudeTimeoutSeconds)
	viper.SetDefault("providers.claude.max_retries", defaultClaudeMaxRetries)

	viper.SetDefault("providers.openai.model", defaultOpenAIModel)
	viper.SetDefault("providers.openai.temperature", defaultOpenAITemperature)
	viper.SetDefault("providers.openai.max_tokens", defaultOpenAIMaxTokens)
	viper.SetDefault("providers.openai.timeout_seconds", defaultOpenAITimeoutSeconds)
	viper.SetDefault("providers.openai.max_retries", defaultOpenAIMaxRetries)
	viper.SetDefault("providers.openai.organization_id", "") // Empty default for optional org ID

	// Set summary template defaults
	viper.SetDefault("summaries.short.system_prompt", constants.ShortSummaryPrompt)
	viper.SetDefault("summaries.long.system_prompt", constants.LongSummaryPrompt)
	viper.SetDefault("transcripts.system_prompt", constants.TranscriptPrompt)
	viper.SetDefault("queries.system_prompt", constants.QueryPrompt)
}

func bindEnvVars() {
	viper.BindEnv("provider", "YTS_PROVIDER")

	// LM Studio env vars
	viper.BindEnv("providers.lmstudio.base_url", "YTS_LMSTUDIO_URL")
	viper.BindEnv("providers.lmstudio.model", "YTS_LMSTUDIO_MODEL")

	// Ollama env vars
	viper.BindEnv("providers.ollama.base_url", "YTS_OLLAMA_URL")
	viper.BindEnv("providers.ollama.model", "YTS_OLLAMA_MODEL")

	// Claude env vars
	viper.BindEnv("providers.claude.model", "YTS_CLAUDE_MODEL")
	viper.BindEnv("providers.claude.temperature", "YTS_CLAUDE_TEMPERATURE")
	viper.BindEnv("providers.claude.max_tokens", "YTS_CLAUDE_MAX_TOKENS")
	viper.BindEnv("providers.claude.timeout", "YTS_CLAUDE_TIMEOUT")

	// OpenAI env vars
	viper.BindEnv("providers.openai.model", "YTS_OPENAI_MODEL")
	viper.BindEnv("providers.openai.temperature", "YTS_OPENAI_TEMPERATURE")
	viper.BindEnv("providers.openai.max_tokens", "YTS_OPENAI_MAX_TOKENS")
	viper.BindEnv("providers.openai.timeout", "YTS_OPENAI_TIMEOUT")
	viper.BindEnv("providers.openai.organization_id", "YTS_OPENAI_ORG_ID")
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
