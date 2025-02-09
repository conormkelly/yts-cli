package cmd

import (
	"fmt"
	"strings"

	"github.com/conormkelly/yts-cli/internal/config"
	"github.com/spf13/cobra"
)

var apikeyCmd = &cobra.Command{
	Use:   "apikey",
	Short: "Manage API keys for LLM providers",
	Long: `Manage API keys for language model providers securely.
Keys are stored in your system's secure keyring.`,
}

var setKeyCmd = &cobra.Command{
	Use:   "set [provider] [api-key]",
	Short: "Set API key for a provider",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := strings.ToLower(args[0])
		apiKey := args[1]

		// Validate provider
		if !isValidAPIKeyProvider(provider) {
			return fmt.Errorf("invalid provider: %s. Valid providers: claude, openai", provider)
		}

		if provider == "claude" && !strings.HasPrefix(apiKey, "sk-") {
			return fmt.Errorf("Claude API keys should start with 'sk-'")
		}

		keyManager := config.NewAPIKeyManager()
		if err := keyManager.SetAPIKey(provider, apiKey); err != nil {
			return fmt.Errorf("failed to set API key: %w", err)
		}

		return nil
	},
}

var deleteKeyCmd = &cobra.Command{
	Use:   "delete [provider]",
	Short: "Delete API key for a provider",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		provider := strings.ToLower(args[0])

		// Validate provider
		if !isValidAPIKeyProvider(provider) {
			return fmt.Errorf("invalid provider: %s. Valid providers: claude, openai", provider)
		}

		keyManager := config.NewAPIKeyManager()
		if err := keyManager.DeleteAPIKey(provider); err != nil {
			return fmt.Errorf("failed to delete API key: %w", err)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(apikeyCmd)
	apikeyCmd.AddCommand(setKeyCmd)
	apikeyCmd.AddCommand(deleteKeyCmd)
}

func isValidAPIKeyProvider(provider string) bool {
	validProviders := map[string]bool{
		"claude": true,
		"openai": true,
	}
	return validProviders[provider]
}
