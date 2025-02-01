// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/conormkelly/yts-cli/internal/config"
	"github.com/conormkelly/yts-cli/internal/llm"
	"github.com/conormkelly/yts-cli/internal/transcript"
	"github.com/spf13/cobra"
)

var (
	summaryType string
	outputFile  string
)

var rootCmd = &cobra.Command{
	Use:   "yts [youtube-url]",
	Short: "Summarize YouTube video transcripts",
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

		// Initialize LLM client using config
		llmClient := llm.NewClient(cfg.LLMBaseURL)

		// Get appropriate system prompt based on summary type
		systemPrompt := config.GetSystemPrompt(summaryType)
		if systemPrompt == "" {
			return fmt.Errorf("failed to get system prompt for summary type: %s", summaryType)
		}

		// Fetch transcript
		fmt.Println("Fetching transcript...")
		transcript, err := fetcher.Fetch(videoURL)
		if err != nil {
			return fmt.Errorf("failed to fetch transcript: %v", err)
		}

		// Generate summary using streaming
		fmt.Println("\nGenerating summary...")
		err = llmClient.SummarizeStream(systemPrompt, cfg.Model, transcript, func(chunk string) {
			fmt.Print(chunk)
		})
		if err != nil {
			return fmt.Errorf("failed to generate summary: %v", err)
		}

		// Handle output file if specified
		if outputFile != "" {
			summary, err := llmClient.Summarize(systemPrompt, cfg.Model, transcript)
			if err != nil {
				return fmt.Errorf("failed to generate summary for file: %v", err)
			}
			if err := os.WriteFile(outputFile, []byte(summary), 0644); err != nil {
				return fmt.Errorf("failed to write output file: %v", err)
			}
			fmt.Printf("\nSummary saved to %s\n", outputFile)
		}

		return nil
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.Flags().StringVarP(&summaryType, "summary", "s", "medium", "summary type (short, medium, long)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")
}

func initConfig() {
	if err := config.Initialize(); err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing config:", err)
		os.Exit(1)
	}
}
