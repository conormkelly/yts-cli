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

type LMStudioProvider struct {
	baseURL string
	model   string
}

func NewLMStudioProvider(baseURL string, model string) *LMStudioProvider {
	return &LMStudioProvider{baseURL: baseURL, model: model}
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type CompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   bool      `json:"stream"`
}

type StreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func (p *LMStudioProvider) Stream(systemPrompt string, transcript string, callback func(string)) error {
	req := CompletionRequest{
		Model: p.model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: transcript},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		p.baseURL+"/v1/chat/completions",
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
			return fmt.Errorf("LM Studio API error (status %d)", resp.StatusCode)
		}
		return fmt.Errorf("LM Studio API error: %v", errorResponse)
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
		if line == "event: error" {
			return fmt.Errorf("LM Studio error - check the Developer tab > Developer Logs for details")
		}

		// Remove "data: " prefix if present
		line = strings.TrimPrefix(line, "data: ")

		var streamResp StreamResponse
		if err := json.Unmarshal([]byte(line), &streamResp); err != nil {
			return fmt.Errorf("error parsing stream response: %w", err)
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				callback(content)
			}
		}
	}

	return nil
}
