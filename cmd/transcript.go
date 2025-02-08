package cmd

import (
	"fmt"
	"os"
	"path/filepath"
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
		fetcher := transcript.NewTranscriptFetcher()
		if err != nil {
			return fmt.Errorf("failed to initialize transcript fetcher: %v", err)
		}

		// Initialize LLM client
		llmClient, err := llm.NewProvider(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize provider: %v", err)
		}

		// Fetch transcript
		title, rawTranscript, err := fetcher.Fetch(videoURL)
		if err != nil {
			return fmt.Errorf("failed to fetch transcript: %v", err)
		}

		fmt.Printf("\nTitle: %s\n\n", title)

		// TODO: make this a flag
		includeTimestamps := false

		var transcriptText strings.Builder
		for i := range rawTranscript {
			if includeTimestamps {
				transcriptText.WriteString(fmt.Sprintf("[%.1fs]: %s\n", rawTranscript[i].Start, rawTranscript[i].Text))
			} else {
				transcriptText.WriteString(rawTranscript[i].Text + "\n")
			}
		}

		// Format transcript using streaming
		var formattedTranscript strings.Builder

		err = llmClient.Stream(
			cfg.Transcripts.SystemPrompt,
			transcriptText.String(),
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
			// Handle home directory expansion
			if outputFile[:2] == "~/" {
				homeDir, err := os.UserHomeDir()
				if err != nil {
					return fmt.Errorf("failed to get home directory: %v", err)
				}
				outputFile = filepath.Join(homeDir, outputFile[2:])
			}

			// Ensure directory exists
			outputDir := filepath.Dir(outputFile)
			if err := os.MkdirAll(outputDir, 0755); err != nil {
				return fmt.Errorf("failed to create output directory: %v", err)
			}

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
