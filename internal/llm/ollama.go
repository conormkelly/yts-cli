package llm

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type OllamaProvider struct {
	baseURL string
	model   string
}

func NewOllamaProvider(baseURL, model string) *OllamaProvider {
	return &OllamaProvider{
		baseURL: baseURL,
		model:   model,
	}
}

type OllamaRequest struct {
	Model   string                 `json:"model"`
	Prompt  string                 `json:"prompt"`
	Stream  bool                   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

func (p *OllamaProvider) Stream(systemPrompt string, transcript string, callback func(string)) error {
	// Combine system prompt and transcript into a single prompt
	// Using a clear delimiter to separate system context from user input
	req := OllamaRequest{
		Model: p.model, // Use provider's configured model
		Prompt: fmt.Sprintf("### System Instructions:\n%s\n\n### Content to Process:\n%s",
			systemPrompt, transcript),
		Stream: true,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		p.baseURL+"/api/generate",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return fmt.Errorf("Ollama API error (status %d)", resp.StatusCode)
		}
		return fmt.Errorf("Ollama API error: %v", errorResponse)
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
		if line == "" {
			continue
		}

		var streamResp OllamaResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			return fmt.Errorf("error parsing stream response: %w", err)
		}

		if streamResp.Response != "" {
			callback(streamResp.Response)
		}

		if streamResp.Done {
			break
		}
	}

	return nil
}
