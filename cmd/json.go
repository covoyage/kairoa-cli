package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var jsonCmd = &cobra.Command{
	Use:   "json",
	Short: "JSON formatting utilities",
	Long:  `Format, minify, and validate JSON data.`,
}

var jsonFormatCmd = &cobra.Command{
	Use:   "format [file]",
	Short: "Format JSON with indentation",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		indent, _ := cmd.Flags().GetString("indent")

		var input []byte
		var err error

		if len(args) > 0 {
			input, err = os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
		} else {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		var data interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		indentBytes := []byte(indent)
		output, err := json.MarshalIndent(data, "", string(indentBytes))
		if err != nil {
			return fmt.Errorf("failed to format JSON: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var jsonMinifyCmd = &cobra.Command{
	Use:   "minify [file]",
	Short: "Minify JSON by removing whitespace",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		var err error

		if len(args) > 0 {
			input, err = os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
		} else {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		var data interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		output, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("failed to minify JSON: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var jsonValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate JSON syntax",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		var err error

		if len(args) > 0 {
			input, err = os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
		} else {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		var data interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		fmt.Println("Valid JSON")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(jsonCmd)
	jsonCmd.AddCommand(jsonFormatCmd)
	jsonCmd.AddCommand(jsonMinifyCmd)
	jsonCmd.AddCommand(jsonValidateCmd)

	jsonFormatCmd.Flags().StringP("indent", "i", "  ", "Indentation string (default: two spaces)")
}
