package main

import (
	"fmt"
	"log"

	"github.com/yourusername/yts-cli/internal/llm"
	"github.com/yourusername/yts-cli/internal/transcript"
)

func main() {
	// Initialize transcript fetcher
	fetcher, err := transcript.NewFetcher()
	if err != nil {
		log.Fatalf("Failed to initialize transcript fetcher: %v", err)
	}
	defer fetcher.Cleanup()

	// Initialize LLM client
	llmClient := llm.NewClient("http://localhost:1234")

	videoURL := "https://www.youtube.com/watch?v=Fh2Q-XW5LyU"

	// Fetch transcript
	fmt.Println("Fetching transcript...")
	text, err := fetcher.Fetch(videoURL)
	if err != nil {
		log.Fatalf("Failed to fetch transcript: %v", err)
	}

	// Generate summary
	fmt.Println("\nGenerating summary...")
	summary, err := llmClient.Summarize(text)
	if err != nil {
		log.Fatalf("Failed to generate summary: %v", err)
	}

	fmt.Printf("\nSummary:\n%s\n", summary)
}
