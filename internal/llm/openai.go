package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/conormkelly/yts-cli/internal/config"
	"github.com/zalando/go-keyring"
)

const (
	openaiAPIURL         = "https://api.openai.com/v1/chat/completions"
	openaiRetryBaseDelay = 1 * time.Second
)

type OpenAIProvider struct {
	model       string
	apiKey      string
	orgID       string
	maxTokens   int
	maxRetries  int
	temperature float64
	client      *http.Client
}

type OpenAIRequest struct {
	Model       string          `json:"model"`
	Messages    []OpenAIMessage `json:"messages"`
	Stream      bool            `json:"stream"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type OpenAIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type OpenAIStreamResponse struct {
	ID      string `json:"id"`
	Object  string `json:"object"`
	Created int64  `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index        int    `json:"index"`
		Delta        Delta  `json:"delta"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
}

type Delta struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

func NewOpenAIProvider(cfg *config.Config) (*OpenAIProvider, error) {
	// Get API key from keyring
	apiKey, err := keyring.Get(config.KeyringService, "openai")
	if err != nil {
		return nil, fmt.Errorf("failed to get OpenAI API key from keyring: %w", err)
	}

	return &OpenAIProvider{
		model:       cfg.Providers.OpenAI.Model,
		apiKey:      apiKey,
		orgID:       cfg.Providers.OpenAI.OrgID,
		temperature: cfg.Providers.OpenAI.Temperature,
		maxTokens:   cfg.Providers.OpenAI.MaxTokens,
		maxRetries:  cfg.Providers.OpenAI.MaxRetries,
		client: &http.Client{
			Timeout: time.Duration(cfg.Providers.OpenAI.TimeoutSecs) * time.Second,
		},
	}, nil
}

func (p *OpenAIProvider) Stream(systemPrompt string, transcript string, callback func(string)) error {
	req := OpenAIRequest{
		Model: p.model,
		Messages: []OpenAIMessage{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: transcript},
		},
		Stream:      true,
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	request, err := http.NewRequest("POST", openaiAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	request.Header.Set("Authorization", "Bearer "+p.apiKey)
	request.Header.Set("Content-Type", "application/json")
	if p.orgID != "" {
		request.Header.Set("OpenAI-Organization", p.orgID)
	}

	resp, err := p.client.Do(request)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(body))
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return fmt.Errorf("error reading stream: %w", err)
		}

		line = strings.TrimSpace(line)
		if line == "" || line == "data: [DONE]" {
			continue
		}

		// Remove "data: " prefix
		line = strings.TrimPrefix(line, "data: ")

		var streamResp OpenAIStreamResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			return fmt.Errorf("error parsing stream response: %w", err)
		}

		// Process choices
		for _, choice := range streamResp.Choices {
			if choice.Delta.Content != "" {
				callback(choice.Delta.Content)
			}
		}
	}

	return nil
}

// Helper method to handle rate limits and retries
func (p *OpenAIProvider) handleRetry(attempt int, err error) (bool, error) {
	if attempt >= p.maxRetries {
		return false, fmt.Errorf("max retries exceeded: %w", err)
	}

	// Check if error is rate limit related
	if strings.Contains(err.Error(), "rate limit") {
		delay := time.Duration(1<<uint(attempt)) * openaiRetryBaseDelay
		time.Sleep(delay)
		return true, nil
	}

	return false, err
}
