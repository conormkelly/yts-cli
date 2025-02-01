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

	// Generate summary using streaming
	fmt.Println("\nGenerating summary...")
	err = llmClient.SummarizeStream(text, func(chunk string) {
		fmt.Print(chunk)
	})
	if err != nil {
		log.Fatalf("Failed to generate summary: %v", err)
	}
	fmt.Println() // Add newline at the end
}
