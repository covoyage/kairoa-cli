package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

var diffCmd = &cobra.Command{
	Use:   "diff [file1] [file2]",
	Short: "Compare text differences",
	Long:  `Compare two texts or files and show the differences.`,
}

var diffTextCmd = &cobra.Command{
	Use:   "text [text1] [text2]",
	Short: "Compare two text strings",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		text1, text2 := args[0], args[1]
		return showDiff(text1, text2, "text1", "text2")
	},
}

var diffFileCmd = &cobra.Command{
	Use:   "file [file1] [file2]",
	Short: "Compare two files",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		file1, file2 := args[0], args[1]

		content1, err := os.ReadFile(file1)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file1, err)
		}

		content2, err := os.ReadFile(file2)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file2, err)
		}

		return showDiff(string(content1), string(content2), file1, file2)
	},
}

var diffStdinCmd = &cobra.Command{
	Use:   "-",
	Short: "Compare stdin with a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		// Read stdin
		var stdinContent strings.Builder
		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			stdinContent.WriteString(scanner.Text())
			stdinContent.WriteString("\n")
		}

		fileContent, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("failed to read %s: %w", file, err)
		}

		return showDiff(stdinContent.String(), string(fileContent), "stdin", file)
	},
}

func showDiff(text1, text2, name1, name2 string) error {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)

	fmt.Printf("--- %s\n", name1)
	fmt.Printf("+++ %s\n\n", name2)

	lineNum1, lineNum2 := 1, 1

	for _, diff := range diffs {
		lines := strings.Split(diff.Text, "\n")
		// Remove last empty element if text ends with newline
		if len(lines) > 0 && lines[len(lines)-1] == "" {
			lines = lines[:len(lines)-1]
		}

		switch diff.Type {
		case diffmatchpatch.DiffDelete:
			for _, line := range lines {
				fmt.Printf("%s %s\n", color.RedString("-%d", lineNum1), color.RedString(line))
				lineNum1++
			}
		case diffmatchpatch.DiffInsert:
			for _, line := range lines {
				fmt.Printf("%s %s\n", color.GreenString("+%d", lineNum2), color.GreenString(line))
				lineNum2++
			}
		case diffmatchpatch.DiffEqual:
			for range lines {
				fmt.Printf(" %d\n", lineNum1)
				lineNum1++
				lineNum2++
			}
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(diffCmd)
	diffCmd.AddCommand(diffTextCmd)
	diffCmd.AddCommand(diffFileCmd)
}
