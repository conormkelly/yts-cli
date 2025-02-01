package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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
}

type CompletionResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewClient(baseURL string) *Client {
	return &Client{baseURL: baseURL}
}

func (c *Client) Summarize(transcript string) (string, error) {
	prompt := `You are a skilled video summarizer. Please provide a clear, well-structured summary of the following video transcript. Focus on:
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
	}

	jsonData, err := json.Marshal(req)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %w", err)
	}

	resp, err := http.Post(
		c.baseURL+"/v1/chat/completions",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("error making request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errorResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&errorResponse); err != nil {
			return "", fmt.Errorf("LLM API error (status %d)", resp.StatusCode)
		}
		return "", fmt.Errorf("LLM API error: %v", errorResponse)
	}

	var result CompletionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %w", err)
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no completion choices returned")
	}

	return result.Choices[0].Message.Content, nil
}
