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
)

var (
	shortSummary bool
	longSummary  bool
	summaryType  string // maintain for backward compatibility
	outputFile   string
)

var rootCmd = &cobra.Command{
	Use:   "yts [youtube-url]",
	Short: "Summarize YouTube video transcripts",
	Args:  cobra.ExactArgs(1),
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFlags(cmd, args)
	},
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
		llmClient := llm.NewClient(cfg.LLMBaseURL)

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
		err = llmClient.Stream(systemPrompt, cfg.Model, transcriptText.String(), func(chunk string) {
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

	// Add new boolean flags for summary types
	rootCmd.Flags().BoolVarP(&shortSummary, "short", "s", false, "Generate a short summary")
	rootCmd.Flags().BoolVarP(&longSummary, "long", "l", false, "Generate a detailed summary")

	// Keep original summary flag for backward compatibility but mark as deprecated
	rootCmd.Flags().StringVarP(&summaryType, "summary", "", "", "Summary type (short, medium, long) [deprecated]")
	rootCmd.Flags().MarkDeprecated("summary", "use --short or --long flags instead")

	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")
}

func initConfig() {
	if err := config.Initialize(); err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing config:", err)
		os.Exit(1)
	}
}

// validateFlags ensures that summary type flags are used correctly
func validateFlags(_ *cobra.Command, _ []string) error {
	// Count active summary flags
	activeFlagCount := 0
	if shortSummary {
		activeFlagCount++
	}
	if longSummary {
		activeFlagCount++
	}
	if summaryType != "" {
		activeFlagCount++
	}

	// Handle multiple flags
	if activeFlagCount > 1 {
		return fmt.Errorf("only one summary type flag can be specified (--short, --long, or --summary)")
	}

	// Handle legacy --summary flag
	if summaryType != "" {
		switch summaryType {
		case "short":
			shortSummary = true
		case "long":
			longSummary = true
		case "medium":
			// default behavior
		default:
			return fmt.Errorf("invalid summary type: %s (must be short, medium, or long)", summaryType)
		}
	}

	return nil
}

// getSummaryType returns the appropriate summary type based on flags
func getSummaryType() string {
	if shortSummary {
		return "short"
	}
	if longSummary {
		return "long"
	}
	return "medium" // default
}
