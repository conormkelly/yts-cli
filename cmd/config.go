package cmd

import (
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage YTS configuration",
	Long: `Manage YTS CLI configuration settings.
Available subcommands:
  - view:  Display current configuration
  - set:   Set a configuration value
  - edit:  Edit configuration in your default text editor`,
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Add subcommands
	configCmd.AddCommand(viewCmd)
	configCmd.AddCommand(setCmd)
	configCmd.AddCommand(editCmd)
}
