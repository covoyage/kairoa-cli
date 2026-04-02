package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var textCmd = &cobra.Command{
	Use:   "text",
	Short: "Text processing utilities",
	Long:  `Various text processing tools including statistics and case conversion.`,
}

var textStatsCmd = &cobra.Command{
	Use:   "stats [file]",
	Short: "Show text statistics",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var text string

		if len(args) > 0 {
			content, err := os.ReadFile(args[0])
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

		stats := calculateTextStats(text)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Metric", "Count"})

		table.Append([]string{"Characters (with spaces)", fmt.Sprintf("%d", stats.charsWithSpaces)})
		table.Append([]string{"Characters (without spaces)", fmt.Sprintf("%d", stats.charsWithoutSpaces)})
		table.Append([]string{"Words", fmt.Sprintf("%d", stats.words)})
		table.Append([]string{"Lines", fmt.Sprintf("%d", stats.lines)})
		table.Append([]string{"Paragraphs", fmt.Sprintf("%d", stats.paragraphs)})
		table.Append([]string{"Chinese Characters", fmt.Sprintf("%d", stats.chineseChars)})
		table.Append([]string{"English Characters", fmt.Sprintf("%d", stats.englishChars)})
		table.Append([]string{"Numbers", fmt.Sprintf("%d", stats.numbers)})
		table.Append([]string{"Punctuation", fmt.Sprintf("%d", stats.punctuation)})
		table.Append([]string{"Bytes", fmt.Sprintf("%d", stats.bytes)})

		table.Render()
		return nil
	},
}

var textCaseCmd = &cobra.Command{
	Use:   "case [type]",
	Short: "Convert text case",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		caseType := strings.ToLower(args[0])

		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}
		text := string(content)

		switch caseType {
		case "upper", "u":
			fmt.Println(strings.ToUpper(text))
		case "lower", "l":
			fmt.Println(strings.ToLower(text))
		case "title", "t":
			fmt.Println(strings.Title(text))
		case "camel", "c":
			fmt.Println(toCamelCase(text))
		case "snake", "s":
			fmt.Println(toSnakeCase(text))
		case "kebab", "k":
			fmt.Println(toKebabCase(text))
		default:
			return fmt.Errorf("unknown case type: %s (use: upper, lower, title, camel, snake, kebab)", caseType)
		}

		return nil
	},
}

var textReplaceCmd = &cobra.Command{
	Use:   "replace [old] [new]",
	Short: "Replace text",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		oldStr := args[0]
		newStr := args[1]

		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		result := strings.ReplaceAll(string(content), oldStr, newStr)
		fmt.Print(result)
		return nil
	},
}

var textLinesCmd = &cobra.Command{
	Use:   "lines",
	Short: "Process lines (sort, unique, reverse)",
	RunE: func(cmd *cobra.Command, args []string) error {
		sortLines, _ := cmd.Flags().GetBool("sort")
		unique, _ := cmd.Flags().GetBool("unique")
		reverse, _ := cmd.Flags().GetBool("reverse")

		content, err := io.ReadAll(os.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read stdin: %w", err)
		}

		lines := strings.Split(string(content), "\n")

		if unique {
			seen := make(map[string]bool)
			var uniqueLines []string
			for _, line := range lines {
				if !seen[line] {
					seen[line] = true
					uniqueLines = append(uniqueLines, line)
				}
			}
			lines = uniqueLines
		}

		if sortLines {
			// Simple bubble sort
			for i := 0; i < len(lines); i++ {
				for j := i + 1; j < len(lines); j++ {
					if lines[i] > lines[j] {
						lines[i], lines[j] = lines[j], lines[i]
					}
				}
			}
		}

		if reverse {
			for i, j := 0, len(lines)-1; i < j; i, j = i+1, j-1 {
				lines[i], lines[j] = lines[j], lines[i]
			}
		}

		fmt.Println(strings.Join(lines, "\n"))
		return nil
	},
}

type textStats struct {
	charsWithSpaces    int
	charsWithoutSpaces int
	words              int
	lines              int
	paragraphs         int
	chineseChars       int
	englishChars       int
	numbers            int
	punctuation        int
	bytes              int
}

func calculateTextStats(text string) textStats {
	var stats textStats

	stats.charsWithSpaces = utf8.RuneCountInString(text)
	stats.bytes = len(text)

	// Characters without spaces
	stats.charsWithoutSpaces = utf8.RuneCountInString(strings.ReplaceAll(text, " ", ""))

	// Lines
	if text == "" {
		stats.lines = 0
	} else {
		stats.lines = strings.Count(text, "\n") + 1
	}

	// Paragraphs
	if text == "" {
		stats.paragraphs = 0
	} else {
		paragraphs := strings.Split(text, "\n\n")
		stats.paragraphs = 0
		for _, p := range paragraphs {
			if strings.TrimSpace(p) != "" {
				stats.paragraphs++
			}
		}
	}

	// Word count
	if text != "" {
		words := strings.Fields(text)
		stats.words = len(words)
	}

	// Character types
	for _, r := range text {
		switch {
		case r >= '\u4e00' && r <= '\u9fff':
			stats.chineseChars++
		case unicode.IsLetter(r) && r < 128:
			stats.englishChars++
		case unicode.IsNumber(r):
			stats.numbers++
		case unicode.IsPunct(r) || unicode.IsSymbol(r):
			stats.punctuation++
		}
	}

	return stats
}

func toCamelCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_' || r == '-'
	})

	for i := 1; i < len(words); i++ {
		words[i] = strings.Title(words[i])
	}

	return strings.Join(words, "")
}

func toSnakeCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '-'
	})

	for i := 0; i < len(words); i++ {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "_")
}

func toKebabCase(s string) string {
	words := strings.FieldsFunc(s, func(r rune) bool {
		return r == ' ' || r == '_'
	})

	for i := 0; i < len(words); i++ {
		words[i] = strings.ToLower(words[i])
	}

	return strings.Join(words, "-")
}

func init() {
	rootCmd.AddCommand(textCmd)
	textCmd.AddCommand(textStatsCmd)
	textCmd.AddCommand(textCaseCmd)
	textCmd.AddCommand(textReplaceCmd)
	textCmd.AddCommand(textLinesCmd)

	textLinesCmd.Flags().BoolP("sort", "s", false, "Sort lines")
	textLinesCmd.Flags().BoolP("unique", "u", false, "Remove duplicate lines")
	textLinesCmd.Flags().BoolP("reverse", "r", false, "Reverse line order")
}
