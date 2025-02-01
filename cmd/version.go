package cmd

import (
	"fmt"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

// Build information. Populated at build-time using -ldflags:
var (
	version   = "dev"
	commit    = "none"
	buildDate = "unknown"
	goVersion = runtime.Version()
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Display version and build information",
	Long: `Display detailed version information about the YTS CLI, including:
- Semantic version number
- Git commit hash
- Build date
- Go version used for compilation
- Supported platforms`,
	Run: func(cmd *cobra.Command, args []string) {
		// Format the version information
		versionInfo := []string{
			fmt.Sprintf("Version:    %s", version),
			fmt.Sprintf("Commit:     %s", commit),
			fmt.Sprintf("Built:      %s", buildDate),
			fmt.Sprintf("Go version: %s", goVersion),
			fmt.Sprintf("OS/Arch:    %s/%s", runtime.GOOS, runtime.GOARCH),
			fmt.Sprintf("Supported platforms: %s", getSupportedPlatforms()),
		}

		// Print version information
		fmt.Println("YouTube Transcript Summarizer (YTS) CLI")
		fmt.Println(strings.Repeat("-", 50))
		fmt.Println(strings.Join(versionInfo, "\n"))
	},
}

// getSupportedPlatforms returns a string of supported OS/Arch combinations
func getSupportedPlatforms() string {
	platforms := []string{
		"linux/amd64",
		"linux/arm64",
		"darwin/amd64",
		"darwin/arm64",
		"windows/amd64",
	}
	return strings.Join(platforms, ", ")
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
