package cmd

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var regexCmd = &cobra.Command{
	Use:   "regex",
	Short: "Regular expression utilities",
	Long:  `Test and match regular expressions.`,
}

var regexTestCmd = &cobra.Command{
	Use:   "test [pattern] [text]",
	Short: "Test regex pattern against text",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		text := args[1]

		flags := ""
		if f, _ := cmd.Flags().GetBool("ignore-case"); f {
			flags += "i"
		}
		if f, _ := cmd.Flags().GetBool("multiline"); f {
			flags += "m"
		}
		if f, _ := cmd.Flags().GetBool("dotall"); f {
			flags += "s"
		}

		re, err := regexp.Compile("(?" + flags + ")" + pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		matches := re.FindAllStringIndex(text, -1)
		if len(matches) == 0 {
			fmt.Println("No matches found")
			return nil
		}

		fmt.Printf("Found %d match(es):\n\n", len(matches))

		for i, match := range matches {
			start, end := match[0], match[1]
			matchedText := text[start:end]

			// Show context
			contextStart := max(0, start-20)
			contextEnd := min(len(text), end+20)

			before := text[contextStart:start]
			after := text[end:contextEnd]

			fmt.Printf("Match %d: %s\n", i+1, color.GreenString(matchedText))
			fmt.Printf("  Position: %d-%d\n", start, end)
			fmt.Printf("  Context: ...%s%s%s...\n", before, color.GreenString(matchedText), after)
			fmt.Println()
		}

		return nil
	},
}

var regexMatchCmd = &cobra.Command{
	Use:   "match [pattern]",
	Short: "Match regex pattern against stdin or file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		file, _ := cmd.Flags().GetString("file")

		var text string
		if file != "" {
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			text = string(content)
		} else {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			text = string(content)
		}

		flags := ""
		if f, _ := cmd.Flags().GetBool("ignore-case"); f {
			flags += "i"
		}
		if f, _ := cmd.Flags().GetBool("multiline"); f {
			flags += "m"
		}

		re, err := regexp.Compile("(?" + flags + ")" + pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		lines := strings.Split(text, "\n")
		for i, line := range lines {
			if re.MatchString(line) {
				// Highlight matches
				highlighted := re.ReplaceAllStringFunc(line, func(s string) string {
					return color.GreenString(s)
				})
				fmt.Printf("%d: %s\n", i+1, highlighted)
			}
		}

		return nil
	},
}

var regexReplaceCmd = &cobra.Command{
	Use:   "replace [pattern] [replacement]",
	Short: "Replace matches with replacement string",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		pattern := args[0]
		replacement := args[1]
		file, _ := cmd.Flags().GetString("file")

		var text string
		if file != "" {
			content, err := os.ReadFile(file)
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			text = string(content)
		} else {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			text = string(content)
		}

		re, err := regexp.Compile(pattern)
		if err != nil {
			return fmt.Errorf("invalid regex: %w", err)
		}

		result := re.ReplaceAllString(text, replacement)
		fmt.Print(result)
		return nil
	},
}

var regexExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show common regex patterns",
	RunE: func(cmd *cobra.Command, args []string) error {
		examples := []struct {
			Name    string
			Pattern string
			Desc    string
		}{
			{"Email", `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`, "Match email addresses"},
			{"URL", `https?://[\w\-]+(\.[\w\-]+)+([\w\-\.,@?^=%&:/~\+#]*[\w\-\@?^=%&/~\+#])?`, "Match URLs"},
			{"IPv4", `\b(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\b`, "Match IPv4 addresses"},
			{"Phone (US)", `\b\d{3}[-.]?\d{3}[-.]?\d{4}\b`, "Match US phone numbers"},
			{"Date (YYYY-MM-DD)", `\d{4}-\d{2}-\d{2}`, "Match dates in YYYY-MM-DD format"},
			{"Credit Card", `\b\d{4}[\s-]?\d{4}[\s-]?\d{4}[\s-]?\d{4}\b`, "Match credit card numbers"},
			{"Hex Color", `#[0-9A-Fa-f]{6}`, "Match hex color codes"},
			{"Numbers Only", `^\d+$`, "Match only numeric strings"},
		}

		fmt.Println("Common Regex Patterns:")
		fmt.Println()
		for _, ex := range examples {
			fmt.Printf("%s:\n", color.CyanString(ex.Name))
			fmt.Printf("  Pattern: %s\n", ex.Pattern)
			fmt.Printf("  Description: %s\n\n", ex.Desc)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(regexCmd)
	regexCmd.AddCommand(regexTestCmd)
	regexCmd.AddCommand(regexMatchCmd)
	regexCmd.AddCommand(regexReplaceCmd)
	regexCmd.AddCommand(regexExamplesCmd)

	regexTestCmd.Flags().BoolP("ignore-case", "i", false, "Case insensitive matching")
	regexTestCmd.Flags().BoolP("multiline", "m", false, "Multiline mode")
	regexTestCmd.Flags().BoolP("dotall", "s", false, "Dot matches newlines")

	regexMatchCmd.Flags().StringP("file", "f", "", "Read from file instead of stdin")
	regexMatchCmd.Flags().BoolP("ignore-case", "i", false, "Case insensitive matching")
	regexMatchCmd.Flags().BoolP("multiline", "m", false, "Multiline mode")

	regexReplaceCmd.Flags().StringP("file", "f", "", "Read from file instead of stdin")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
