package cmd

import (
	"fmt"
	"io"
	"net/url"
	"os"

	"github.com/spf13/cobra"
)

var urlCmd = &cobra.Command{
	Use:   "url",
	Short: "URL encoding/decoding",
	Long:  `Encode or decode URL strings.`,
}

var urlEncodeCmd = &cobra.Command{
	Use:   "encode [text]",
	Short: "Encode text to URL format",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		} else {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			input = string(data)
		}

		encoded := url.QueryEscape(input)
		fmt.Println(encoded)
		return nil
	},
}

var urlDecodeCmd = &cobra.Command{
	Use:   "decode [text]",
	Short: "Decode URL-encoded text",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			input = args[0]
		} else {
			data, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			input = string(data)
		}

		decoded, err := url.QueryUnescape(input)
		if err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}
		fmt.Println(decoded)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(urlCmd)
	urlCmd.AddCommand(urlEncodeCmd)
	urlCmd.AddCommand(urlDecodeCmd)
}
