package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var validPaths = map[string]struct{}{
	"provider":                    {},
	"providers.lmstudio.base_url": {},
	"providers.lmstudio.model":    {},
	"providers.ollama.base_url":   {},
	"providers.ollama.model":      {},
}

var setCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := strings.ToLower(args[0])
		value := args[1]

		// Validate key
		if _, ok := validPaths[key]; !ok {
			return fmt.Errorf("invalid configuration key: %s\nValid keys: %s",
				key, strings.Join(getValidKeys(), ", "))
		}

		// Validate provider if setting provider
		if key == "provider" && !isValidProvider(value) {
			return fmt.Errorf("invalid provider: %s\nValid providers: lmstudio, ollama", value)
		}

		// Set the value
		viper.Set(key, value)

		// Save the configuration
		if err := viper.WriteConfig(); err != nil {
			return fmt.Errorf("failed to save configuration: %v", err)
		}

		fmt.Printf("Successfully set %s to %s\n", key, value)
		return nil
	},
}

func isValidProvider(provider string) bool {
	return provider == "lmstudio" || provider == "ollama"
}

func getValidKeys() []string {
	keys := make([]string, 0, len(validPaths))
	for k := range validPaths {
		keys = append(keys, k)
	}
	return keys
}
