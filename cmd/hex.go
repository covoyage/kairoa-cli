package cmd

import (
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var hexCmd = &cobra.Command{
	Use:   "hex",
	Short: "Hex encoding/decoding",
	Long:  `Encode or decode text using hexadecimal encoding.`,
}

var hexEncodeCmd = &cobra.Command{
	Use:   "encode [text]",
	Short: "Encode text to hex",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		upper, _ := cmd.Flags().GetBool("upper")
		var input []byte
		var err error

		if len(args) > 0 {
			input = []byte(args[0])
		} else {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		encoded := hex.EncodeToString(input)
		if upper {
			encoded = strings.ToUpper(encoded)
		}
		fmt.Println(encoded)
		return nil
	},
}

var hexDecodeCmd = &cobra.Command{
	Use:   "decode [text]",
	Short: "Decode hex to text",
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

		decoded, err := hex.DecodeString(input)
		if err != nil {
			return fmt.Errorf("failed to decode: %w", err)
		}
		fmt.Print(string(decoded))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(hexCmd)
	hexCmd.AddCommand(hexEncodeCmd)
	hexCmd.AddCommand(hexDecodeCmd)

	hexEncodeCmd.Flags().BoolP("upper", "u", false, "Use uppercase letters")
}
