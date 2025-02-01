// cmd/root.go
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/yourusername/yts-cli/internal/config"
	"github.com/yourusername/yts-cli/internal/llm"
	"github.com/yourusername/yts-cli/internal/transcript"
)

var (
	cfgFile     string
	summaryType string
	outputFile  string
)

var rootCmd = &cobra.Command{
	Use:   "yts [youtube-url]",
	Short: "Summarize YouTube video transcripts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		videoURL := args[0]

		// Initialize transcript fetcher
		fetcher, err := transcript.NewFetcher()
		if err != nil {
			return fmt.Errorf("failed to initialize transcript fetcher: %v", err)
		}
		defer fetcher.Cleanup()

		// Initialize LLM client using config
		llmClient := llm.NewClient(viper.GetString("llm_base_url"))

		// Fetch transcript
		fmt.Println("Fetching transcript...")
		text, err := fetcher.Fetch(videoURL)
		if err != nil {
			return fmt.Errorf("failed to fetch transcript: %v", err)
		}

		// Generate summary using streaming
		fmt.Println("\nGenerating summary...")
		err = llmClient.SummarizeStream(text, func(chunk string) {
			fmt.Print(chunk)
		})
		if err != nil {
			return fmt.Errorf("failed to generate summary: %v", err)
		}

		// Handle output file if specified
		if outputFile != "" {
			// Get the full summary first
			summary, err := llmClient.Summarize(text)
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

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/yts/config.json)")
	rootCmd.Flags().StringVarP(&summaryType, "summary", "s", "medium", "summary type (short, medium, long)")
	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "", "output file path")

	// Bind flags to viper
	viper.BindPFlag("summary_type", rootCmd.Flags().Lookup("summary"))
	viper.BindPFlag("output_file", rootCmd.Flags().Lookup("output"))
}

func initConfig() {
	if err := config.InitializeViper(); err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing config:", err)
		os.Exit(1)
	}
}
