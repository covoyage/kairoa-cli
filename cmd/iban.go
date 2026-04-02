package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var ibanCmd = &cobra.Command{
	Use:   "iban",
	Short: "IBAN validator and formatter",
	Long:  `Validate and format IBAN (International Bank Account Number).`,
}

var countryLengths = map[string]int{
	"AL": 28, "AD": 24, "AT": 20, "AZ": 28, "BE": 16, "BH": 22, "BA": 20, "BR": 29,
	"BG": 22, "CR": 22, "HR": 21, "CY": 28, "CZ": 24, "DK": 18, "DO": 28, "EE": 20,
	"FO": 18, "FI": 18, "FR": 27, "GE": 22, "DE": 22, "GI": 23, "GR": 27, "GL": 18,
	"GT": 28, "HU": 28, "IS": 26, "IE": 22, "IL": 23, "IT": 27, "JO": 30, "KZ": 20,
	"KW": 30, "LV": 21, "LB": 28, "LI": 21, "LT": 20, "LU": 20, "MK": 19, "MT": 31,
	"MR": 27, "MU": 30, "MD": 24, "MC": 27, "ME": 22, "NL": 18, "NO": 15, "PK": 24,
	"PS": 29, "PL": 28, "PT": 25, "QA": 29, "RO": 24, "SM": 27, "SA": 24, "RS": 22,
	"SK": 24, "SI": 19, "ES": 24, "SE": 24, "CH": 21, "TN": 24, "TR": 26, "UA": 29,
	"AE": 23, "GB": 22, "VG": 24, "XK": 20,
}

var ibanValidateCmd = &cobra.Command{
	Use:   "validate [iban]",
	Short: "Validate IBAN",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		iban := normalizeIBAN(args[0])

		result := validateIBAN(iban)

		fmt.Printf("IBAN: %s\n", formatIBAN(iban))
		fmt.Printf("Normalized: %s\n", iban)
		fmt.Printf("Country Code: %s\n", iban[:2])
		fmt.Printf("Check Digits: %s\n", iban[2:4])
		fmt.Printf("BBAN: %s\n", iban[4:])

		if expectedLen, ok := countryLengths[iban[:2]]; ok {
			fmt.Printf("Expected Length: %d\n", expectedLen)
			fmt.Printf("Actual Length: %d\n", len(iban))
		}

		if result {
			fmt.Println("\n✓ Valid IBAN")
		} else {
			fmt.Println("\n✗ Invalid IBAN")
		}

		return nil
	},
}

var ibanFormatCmd = &cobra.Command{
	Use:   "format [iban]",
	Short: "Format IBAN",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		iban := normalizeIBAN(args[0])
		fmt.Println(formatIBAN(iban))
		return nil
	},
}

func normalizeIBAN(iban string) string {
	// Remove spaces and convert to uppercase
	iban = strings.ToUpper(iban)
	iban = regexp.MustCompile(`[^A-Z0-9]`).ReplaceAllString(iban, "")
	return iban
}

func formatIBAN(iban string) string {
	// Add spaces every 4 characters
	var parts []string
	for i := 0; i < len(iban); i += 4 {
		end := i + 4
		if end > len(iban) {
			end = len(iban)
		}
		parts = append(parts, iban[i:end])
	}
	return strings.Join(parts, " ")
}

func validateIBAN(iban string) bool {
	if len(iban) < 4 {
		return false
	}

	// Check country code length
	countryCode := iban[:2]
	if expectedLen, ok := countryLengths[countryCode]; ok {
		if len(iban) != expectedLen {
			return false
		}
	}

	// Move first 4 chars to end
	rearranged := iban[4:] + iban[:4]

	// Convert letters to numbers (A=10, B=11, ...)
	var numeric strings.Builder
	for _, char := range rearranged {
		if char >= 'A' && char <= 'Z' {
			numeric.WriteString(fmt.Sprintf("%d", int(char-'A'+10)))
		} else {
			numeric.WriteRune(char)
		}
	}

	// Calculate mod 97
	return mod97(numeric.String()) == 1
}

func mod97(s string) int {
	// Calculate mod 97 for large numbers
	remainder := 0
	for _, char := range s {
		digit := int(char - '0')
		remainder = (remainder*10 + digit) % 97
	}
	return remainder
}

func init() {
	rootCmd.AddCommand(ibanCmd)
	ibanCmd.AddCommand(ibanValidateCmd)
	ibanCmd.AddCommand(ibanFormatCmd)
}
