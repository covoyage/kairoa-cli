package cmd

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/spf13/cobra"
)

var baseCmd = &cobra.Command{
	Use:   "base",
	Short: "Base conversion utilities",
	Long:  `Convert numbers between different bases (binary, octal, decimal, hexadecimal).`,
}

var baseConvertCmd = &cobra.Command{
	Use:   "convert [number]",
	Short: "Convert number to all bases",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.TrimSpace(args[0])
		fromBase, _ := cmd.Flags().GetInt("from")

		// Auto-detect base if not specified
		if fromBase == 0 {
			fromBase = detectBase(input)
		}

		// Clean input
		input = cleanInput(input, fromBase)

		// Parse number
		num := new(big.Int)
		_, ok := num.SetString(input, fromBase)
		if !ok {
			return fmt.Errorf("invalid number for base %d: %s", fromBase, input)
		}

		// Convert to all bases
		fmt.Println("Binary (Base 2):  ", num.Text(2))
		fmt.Println("Octal (Base 8):   ", num.Text(8))
		fmt.Println("Decimal (Base 10):", num.Text(10))
		fmt.Println("Hex (Base 16):    ", strings.ToUpper(num.Text(16)))
		return nil
	},
}

var baseToCmd = &cobra.Command{
	Use:   "to [number]",
	Short: "Convert number to specific base",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.TrimSpace(args[0])
		fromBase, _ := cmd.Flags().GetInt("from")
		toBase, _ := cmd.Flags().GetInt("base")

		if toBase < 2 || toBase > 36 {
			return fmt.Errorf("target base must be between 2 and 36")
		}

		// Auto-detect base if not specified
		if fromBase == 0 {
			fromBase = detectBase(input)
		}

		// Clean input
		input = cleanInput(input, fromBase)

		// Parse number
		num := new(big.Int)
		_, ok := num.SetString(input, fromBase)
		if !ok {
			return fmt.Errorf("invalid number for base %d: %s", fromBase, input)
		}

		fmt.Println(num.Text(toBase))
		return nil
	},
}

func detectBase(input string) int {
	input = strings.ToLower(input)

	// Check for prefixes
	if strings.HasPrefix(input, "0b") {
		return 2
	}
	if strings.HasPrefix(input, "0o") {
		return 8
	}
	if strings.HasPrefix(input, "0x") {
		return 16
	}

	// Check for non-decimal characters
	if strings.ContainsAny(input, "abcdef") {
		return 16
	}
	if strings.ContainsAny(input, "89") {
		return 10
	}
	if strings.ContainsAny(input, "234567") {
		return 10
	}

	return 10
}

func cleanInput(input string, base int) string {
	input = strings.ToLower(input)

	// Remove prefixes
	switch base {
	case 2:
		input = strings.TrimPrefix(input, "0b")
	case 8:
		input = strings.TrimPrefix(input, "0o")
	case 16:
		input = strings.TrimPrefix(input, "0x")
	}

	return strings.TrimSpace(input)
}

func init() {
	rootCmd.AddCommand(baseCmd)
	baseCmd.AddCommand(baseConvertCmd)
	baseCmd.AddCommand(baseToCmd)

	baseConvertCmd.Flags().IntP("from", "f", 0, "Input base (auto-detect if not specified)")
	baseToCmd.Flags().IntP("from", "f", 0, "Input base (auto-detect if not specified)")
	baseToCmd.Flags().IntP("base", "b", 10, "Target base (2-36)")
}
