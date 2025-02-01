package main

import (
	"fmt"
	"log"

	"github.com/yourusername/yts-cli/internal/transcript"
)

func main() {
	fetcher, err := transcript.NewFetcher()
	if err != nil {
		log.Fatalf("Failed to initialize transcript fetcher: %v", err)
	}
	defer fetcher.Cleanup()

	videoURL := "https://www.youtube.com/watch?v=Fh2Q-XW5LyU"

	fmt.Println("Fetching transcript...")
	text, err := fetcher.Fetch(videoURL)
	if err != nil {
		log.Fatalf("Failed to fetch transcript: %v", err)
	}

	fmt.Println("\nTranscript:")
	fmt.Println(text)
}
