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

type Client struct {
	baseURL string
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

type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

type StreamResponse struct {
	Choices []struct {
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
	} `json:"choices"`
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) SummarizeStream(transcript string, callback func(string)) error {
	prompt := `
	You are a skilled video summarizer. Please provide a clear, well-structured summary of the following video transcript. Focus on:
	- Main topics and key points
	- Important details and insights
	- Clear structure and readability

	Transcript:
	`

	req := CompletionRequest{
		Model: "llama-3.2-3b-instruct",
		Messages: []Message{
			{Role: "system", Content: prompt},
			{Role: "user", Content: transcript},
		},
		Stream: true,
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		c.baseURL+"/v1/chat/completions",
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
			return fmt.Errorf("LLM API error (status %d)", resp.StatusCode)
		}
		return fmt.Errorf("LLM API error: %v", errorResponse)
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

// Regular non-streaming method kept for reference
func (c *Client) Summarize(transcript string) (string, error) {
	var summary strings.Builder
	err := c.SummarizeStream(transcript, func(chunk string) {
		summary.WriteString(chunk)
	})
	if err != nil {
		return "", err
	}
	return summary.String(), nil
}
