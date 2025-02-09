// internal/llm/provider.go
package llm

import (
	"fmt"

	"github.com/conormkelly/yts-cli/internal/config"
)

type Provider interface {
	Stream(systemPrompt string, transcript string, callback func(string)) error
}

func NewProvider(cfg *config.Config) (Provider, error) {
	switch cfg.Provider {
	case "lmstudio":
		return NewLMStudioProvider(cfg.Providers.LMStudio.BaseURL, cfg.Providers.LMStudio.Model), nil
	case "ollama":
		return NewOllamaProvider(cfg.Providers.Ollama.BaseURL, cfg.Providers.Ollama.Model), nil
	case "claude":
		return NewClaudeProvider(cfg)
	case "openai":
		return NewOpenAIProvider(cfg)
	default:
		return nil, fmt.Errorf("unsupported provider type: %s", cfg.Provider)
	}
}
