package cmd

import (
	"encoding/base64"
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var base64Cmd = &cobra.Command{
	Use:   "base64",
	Short: "Base64 encode/decode",
	Long:  `Encode or decode text using Base64 encoding.`,
}

var base64EncodeCmd = &cobra.Command{
	Use:   "encode [text]",
	Short: "Encode text to Base64",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		if len(args) > 0 {
			input = []byte(args[0])
		} else {
			var err error
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		encoded := base64.StdEncoding.EncodeToString(input)
		fmt.Println(encoded)
		return nil
	},
}

var base64DecodeCmd = &cobra.Command{
	Use:   "decode [text]",
	Short: "Decode Base64 to text",
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

		decoded, err := base64.StdEncoding.DecodeString(input)
		if err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}
		fmt.Print(string(decoded))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(base64Cmd)
	base64Cmd.AddCommand(base64EncodeCmd)
	base64Cmd.AddCommand(base64DecodeCmd)
}
