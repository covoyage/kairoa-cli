package cmd

import (
	"fmt"

	"github.com/covoyage/kairoa-cli/internal/i18n"
	"github.com/spf13/cobra"
)

var (
	buildVersion = "dev"
	buildCommit  = "none"
	buildDate    = "unknown"
)

// SetVersionInfo is called from main with values injected via -ldflags.
func SetVersionInfo(version, commit, date string) {
	buildVersion = version
	buildCommit = commit
	buildDate = date
}

var rootCmd = &cobra.Command{
	Use:   "kairoa",
	Short: "Kairoa CLI - A developer toolbox",
	Long: `Kairoa CLI is a command-line version of the Kairoa developer toolbox.
It provides various utilities for developers including:
  - Hash calculation (MD5, SHA1, SHA256, etc.)
  - UUID generation
  - Base64 encoding/decoding
  - JSON formatting
  - URL encoding/decoding
  - JWT decoding
  - And more...`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		lang, _ := cmd.Flags().GetString("lang")
		if lang != "" {
			switch lang {
			case "zh", "zh-CN", "zh-TW", "zh-HK":
				i18n.SetLocale(i18n.Chinese)
			case "en":
				i18n.SetLocale(i18n.English)
			}
		}
	},
}

func Execute() error {
	i18n.Init()
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().StringP("lang", "l", "", "Language (en, zh)")
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("kairoa %s (commit: %s, built: %s)\n", buildVersion, buildCommit, buildDate)
	},
}
