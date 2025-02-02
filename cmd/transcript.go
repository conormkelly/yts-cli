package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/conormkelly/yts-cli/internal/config"
	"github.com/conormkelly/yts-cli/internal/llm"
	"github.com/conormkelly/yts-cli/internal/transcript"
	"github.com/spf13/cobra"
)

var transcriptCmd = &cobra.Command{
	Use:   "transcript [youtube-url]",
	Short: "Format the raw transcript with proper capitalization and punctuation",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		videoURL := args[0]

		// Get configuration
		cfg, err := config.GetConfig()
		if err != nil {
			return fmt.Errorf("failed to get config: %v", err)
		}

		// Initialize transcript fetcher
		fetcher, err := transcript.NewFetcher()
		if err != nil {
			return fmt.Errorf("failed to initialize transcript fetcher: %v", err)
		}
		defer fetcher.Cleanup()

		// Initialize LLM client
		llmClient := llm.NewClient(cfg.LLMBaseURL)

		// Fetch transcript
		rawTranscript, err := fetcher.Fetch(videoURL)
		if err != nil {
			return fmt.Errorf("failed to fetch transcript: %v", err)
		}

		// Format transcript using streaming
		var formattedTranscript strings.Builder

		err = llmClient.Stream(
			cfg.Transcripts.SystemPrompt,
			cfg.Model,
			rawTranscript,
			func(chunk string) {
				fmt.Print(chunk)
				formattedTranscript.WriteString(chunk)
			},
		)
		if err != nil {
			return fmt.Errorf("failed to format transcript: %v", err)
		}
		// Add a newline at the end of the stream
		formattedTranscript.WriteString("\n")
		fmt.Println()

		// Handle output file if specified
		if outputFile != "" {
			if err := os.WriteFile(outputFile, []byte(formattedTranscript.String()), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %v", err)
			}
			fmt.Printf("\nFormatted transcript saved to %s\n", outputFile)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(transcriptCmd)
	transcriptCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")
}
