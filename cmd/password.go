package cmd

import (
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/spf13/cobra"
)

var passwordCmd = &cobra.Command{
	Use:   "password",
	Short: "Password generator",
	Long:  `Generate secure random passwords.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		length, _ := cmd.Flags().GetInt("length")
		count, _ := cmd.Flags().GetInt("count")
		noUpper, _ := cmd.Flags().GetBool("no-upper")
		noLower, _ := cmd.Flags().GetBool("no-lower")
		noNumbers, _ := cmd.Flags().GetBool("no-numbers")
		noSpecial, _ := cmd.Flags().GetBool("no-special")

		// Build character set
		var chars string
		if !noLower {
			chars += "abcdefghijklmnopqrstuvwxyz"
		}
		if !noUpper {
			chars += "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
		}
		if !noNumbers {
			chars += "0123456789"
		}
		if !noSpecial {
			chars += "!@#$%^&*()_+-=[]{}|;:,.<>?"
		}

		if chars == "" {
			return fmt.Errorf("at least one character type must be enabled")
		}

		for i := 0; i < count; i++ {
			password, err := generatePassword(chars, length)
			if err != nil {
				return err
			}
			fmt.Println(password)
		}
		return nil
	},
}

func generatePassword(chars string, length int) (string, error) {
	result := make([]byte, length)
	max := big.NewInt(int64(len(chars)))

	for i := 0; i < length; i++ {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", err
		}
		result[i] = chars[n.Int64()]
	}

	return string(result), nil
}

func init() {
	rootCmd.AddCommand(passwordCmd)
	passwordCmd.Flags().IntP("length", "n", 16, "Password length")
	passwordCmd.Flags().IntP("count", "c", 1, "Number of passwords to generate")
	passwordCmd.Flags().Bool("no-upper", false, "Exclude uppercase letters")
	passwordCmd.Flags().Bool("no-lower", false, "Exclude lowercase letters")
	passwordCmd.Flags().Bool("no-numbers", false, "Exclude numbers")
	passwordCmd.Flags().Bool("no-special", false, "Exclude special characters")
}
