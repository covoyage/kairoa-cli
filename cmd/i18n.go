package cmd

import (
	"fmt"

	"github.com/covoyage/kairoa-cli/internal/i18n"
	"github.com/spf13/cobra"
)

var i18nCmd = &cobra.Command{
	Use:   "lang",
	Short: "Language settings",
	Long:  `Manage language settings for the CLI.`,
}

var i18nSetCmd = &cobra.Command{
	Use:   "set [lang]",
	Short: "Set language",
	Long:  `Set the language for the CLI. Supported languages: en, zh`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		lang := args[0]
		switch lang {
		case "en", "english":
			i18n.SetLocale(i18n.English)
			fmt.Println("Language set to English")
		case "zh", "chinese", "中文":
			i18n.SetLocale(i18n.Chinese)
			fmt.Println("语言已设置为中文")
		default:
			return fmt.Errorf("unsupported language: %s (supported: en, zh)", lang)
		}
		return nil
	},
}

var i18nGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get current language",
	Long:  `Display the current language setting.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		locale := i18n.GetLocale()
		switch locale {
		case i18n.Chinese:
			fmt.Println("Current language: Chinese (zh)")
		default:
			fmt.Println("Current language: English (en)")
		}
		return nil
	},
}

var i18nListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available languages",
	Long:  `List all available languages.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available languages:")
		fmt.Println("  en - English")
		fmt.Println("  zh - 中文")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(i18nCmd)
	i18nCmd.AddCommand(i18nSetCmd)
	i18nCmd.AddCommand(i18nGetCmd)
	i18nCmd.AddCommand(i18nListCmd)
}
