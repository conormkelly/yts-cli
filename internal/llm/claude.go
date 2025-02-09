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
	claudeAPIURL   = "https://api.anthropic.com/v1/messages"
	keyringService = "yts-cli"
	retryBaseDelay = 1 * time.Second
)

type ClaudeProvider struct {
	model       string
	apiKey      string
	maxTokens   int
	maxRetries  int
	temperature float64
	client      *http.Client
}

type ClaudeRequest struct {
	Model       string          `json:"model"`
	Messages    []ClaudeMessage `json:"messages"`
	System      string          `json:"system,omitempty"`
	Stream      bool            `json:"stream"`
	MaxTokens   int             `json:"max_tokens,omitempty"`
	Temperature float64         `json:"temperature,omitempty"`
}

type ClaudeMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ClaudeStreamEvent struct {
	Type string `json:"type"`

	// For message_start
	Message *struct {
		ID      string `json:"id"`
		Role    string `json:"role"`
		Content []any  `json:"content"`
	} `json:"message,omitempty"`

	// For content_block_delta
	Index int `json:"index,omitempty"`
	Delta *struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"delta,omitempty"`
}

func NewClaudeProvider(cfg *config.Config) (*ClaudeProvider, error) {
	// Get API key from keyring
	apiKey, err := keyring.Get(keyringService, "claude")
	if err != nil {
		return nil, fmt.Errorf("failed to get Claude API key from keyring: %w", err)
	}

	return &ClaudeProvider{
		model:       cfg.Providers.Claude.Model,
		apiKey:      apiKey,
		temperature: cfg.Providers.Claude.Temperature,
		maxTokens:   cfg.Providers.Claude.MaxTokens,
		maxRetries:  cfg.Providers.Claude.MaxRetries,
		client: &http.Client{
			Timeout: time.Duration(cfg.Providers.Claude.TimeoutSecs) * time.Second,
		},
	}, nil
}

func (p *ClaudeProvider) Stream(systemPrompt string, transcript string, callback func(string)) error {
	req := ClaudeRequest{
		Model: p.model,
		Messages: []ClaudeMessage{
			{Role: "user", Content: transcript},
		},
		System:      systemPrompt,
		Stream:      true,
		MaxTokens:   p.maxTokens,
		Temperature: p.temperature,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	request, err := http.NewRequest("POST", claudeAPIURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set headers
	request.Header.Set("x-api-key", p.apiKey)
	request.Header.Set("anthropic-version", "2023-06-01")
	request.Header.Set("content-type", "application/json")

	resp, err := p.client.Do(request)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(body))
	}

	// Read the stream line by line
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
		if line == "" {
			continue
		}

		// Handle event: lines
		if strings.HasPrefix(line, "event:") {
			eventType := strings.TrimSpace(strings.TrimPrefix(line, "event:"))
			if eventType == "error" {
				// Get the next line which should contain the error details
				errLine, err := reader.ReadString('\n')
				if err != nil {
					return fmt.Errorf("error reading error details: %w", err)
				}
				// Try to parse error details if present
				if strings.HasPrefix(errLine, "data:") {
					errData := strings.TrimSpace(strings.TrimPrefix(errLine, "data:"))
					return fmt.Errorf("Claude streaming error: %s", errData)
				}
				return fmt.Errorf("Claude streaming error (no details available)")
			}
			continue
		}

		// Only process data: lines which contain the actual JSON
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		// Remove "data: " prefix and parse the JSON
		line = strings.TrimPrefix(line, "data: ")

		var event ClaudeStreamEvent
		if err := json.Unmarshal([]byte(line), &event); err != nil {
			return fmt.Errorf("error parsing stream event: %w", err)
		}

		// Handle different event types
		switch event.Type {
		case "content_block_delta":
			if event.Delta != nil && event.Delta.Type == "text_delta" {
				callback(event.Delta.Text)
			}
		}
	}

	return nil
}
