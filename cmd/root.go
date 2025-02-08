// cmd/root.go
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
	"github.com/spf13/viper"
)

var (
	longSummary bool
	provider    string // ollama, LM Studio etc
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
		fetcher := transcript.NewTranscriptFetcher()
		if err != nil {
			return fmt.Errorf("failed to initialize transcript fetcher: %v", err)
		}

		// Initialize LLM client using config
		llmClient, err := llm.NewProvider(cfg)
		if err != nil {
			return fmt.Errorf("failed to initialize llm client: %v", err)
		}

		// Get appropriate system prompt based on resolved summary type
		resolvedSummaryType := getSummaryType()
		systemPrompt := config.GetSystemPrompt(resolvedSummaryType)
		if systemPrompt == "" {
			return fmt.Errorf("failed to get system prompt for summary type: %s", resolvedSummaryType)
		}

		// Fetch transcript
		title, transcript, err := fetcher.Fetch(videoURL)
		if err != nil {
			return fmt.Errorf("failed to fetch transcript: %v", err)
		}

		fmt.Printf("\nTitle: %s\n\n", title)

		var transcriptText strings.Builder
		for i := range transcript {
			transcriptText.WriteString(transcript[i].Text + "\n")
		}

		// Generate summary using streaming
		var summary strings.Builder
		err = llmClient.Stream(systemPrompt, transcriptText.String(), func(chunk string) {
			fmt.Print(chunk)
			summary.WriteString(chunk)
		})
		if err != nil {
			return fmt.Errorf("failed to generate summary: %v", err)
		}
		// Add newline
		summary.WriteString("\n")
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

			if err := os.WriteFile(outputFile, []byte(title+"\n\n"+summary.String()), 0644); err != nil {
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

	rootCmd.Flags().StringVarP(&provider, "provider", "p", "", "LLM provider (lmstudio, ollama)")
	rootCmd.Flags().BoolVarP(&longSummary, "long", "l", false, "Generate a detailed summary")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")

	// Bind provider flag to viper
	viper.BindPFlag("provider", rootCmd.Flags().Lookup("provider"))
}

func initConfig() {
	if err := config.Initialize(); err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing config:", err)
		os.Exit(1)
	}
}

// getSummaryType returns the appropriate summary type based on flags
func getSummaryType() string {
	if longSummary {
		return "long"
	}
	return "short" // default
}
