package cmd

import (
	"github.com/covoyage/kairoa-cli/internal/i18n"
	"github.com/spf13/cobra"
)

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
		// Handle language flag
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
	// Initialize i18n with detected locale
	i18n.Init()
	
	return rootCmd.Execute()
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	
	// Add global flags
	rootCmd.PersistentFlags().StringP("lang", "l", "", "Language (en, zh)")
}
