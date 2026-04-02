package cmd

import (
	"fmt"
	"regexp"
	"unicode"

	"github.com/spf13/cobra"
)

var passwordstrengthCmd = &cobra.Command{
	Use:   "password-strength [password]",
	Short: "Password strength checker",
	Long:  `Analyze password strength and provide recommendations.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		password := args[0]

		result := checkPasswordStrength(password)

		fmt.Printf("Password Length: %d\n", len(password))
		fmt.Printf("Strength Score: %d/100\n", result.Score)
		fmt.Printf("Strength Level: %s\n", result.Level)
		fmt.Println()
		fmt.Println("Criteria:")
		fmt.Printf("  ✓ Lowercase letters: %v\n", result.HasLower)
		fmt.Printf("  ✓ Uppercase letters: %v\n", result.HasUpper)
		fmt.Printf("  ✓ Numbers: %v\n", result.HasNumber)
		fmt.Printf("  ✓ Special characters: %v\n", result.HasSpecial)
		fmt.Printf("  ✓ Minimum length (8+): %v\n", result.HasMinLength)
		fmt.Println()

		if len(result.Suggestions) > 0 {
			fmt.Println("Suggestions:")
			for _, s := range result.Suggestions {
				fmt.Printf("  • %s\n", s)
			}
		}

		return nil
	},
}

type PasswordStrength struct {
	Score        int
	Level        string
	HasLower     bool
	HasUpper     bool
	HasNumber    bool
	HasSpecial   bool
	HasMinLength bool
	Suggestions  []string
}

func checkPasswordStrength(password string) PasswordStrength {
	result := PasswordStrength{
		Suggestions: []string{},
	}

	// Check character types
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			result.HasLower = true
		case unicode.IsUpper(char):
			result.HasUpper = true
		case unicode.IsDigit(char):
			result.HasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			result.HasSpecial = true
		}
	}

	// Check length
	result.HasMinLength = len(password) >= 8

	// Calculate score
	score := 0
	if result.HasLower {
		score += 10
	}
	if result.HasUpper {
		score += 10
	}
	if result.HasNumber {
		score += 10
	}
	if result.HasSpecial {
		score += 15
	}
	if result.HasMinLength {
		score += 10
	}

	// Length bonus
	if len(password) >= 12 {
		score += 15
	}
	if len(password) >= 16 {
		score += 20
	}

	// Variety bonus
	variety := 0
	if result.HasLower {
		variety++
	}
	if result.HasUpper {
		variety++
	}
	if result.HasNumber {
		variety++
	}
	if result.HasSpecial {
		variety++
	}
	score += variety * 5

	// Cap score at 100
	if score > 100 {
		score = 100
	}
	result.Score = score

	// Determine level
	switch {
	case score < 30:
		result.Level = "Very Weak"
	case score < 50:
		result.Level = "Weak"
	case score < 70:
		result.Level = "Fair"
	case score < 85:
		result.Level = "Strong"
	default:
		result.Level = "Very Strong"
	}

	// Generate suggestions
	if !result.HasLower {
		result.Suggestions = append(result.Suggestions, "Add lowercase letters")
	}
	if !result.HasUpper {
		result.Suggestions = append(result.Suggestions, "Add uppercase letters")
	}
	if !result.HasNumber {
		result.Suggestions = append(result.Suggestions, "Add numbers")
	}
	if !result.HasSpecial {
		result.Suggestions = append(result.Suggestions, "Add special characters (!@#$%^&*)")
	}
	if !result.HasMinLength {
		result.Suggestions = append(result.Suggestions, "Use at least 8 characters")
	}
	if len(password) < 12 {
		result.Suggestions = append(result.Suggestions, "Consider using 12+ characters for better security")
	}

	// Check for common patterns
	if matched, _ := regexp.MatchString(`(?i)(password|123456|qwerty|abc123)`, password); matched {
		result.Suggestions = append(result.Suggestions, "Avoid common words and patterns")
	}
	if matched, _ := regexp.MatchString(`(.)\1{2,}`, password); matched {
		result.Suggestions = append(result.Suggestions, "Avoid repeating characters")
	}

	return result
}

func init() {
	rootCmd.AddCommand(passwordstrengthCmd)
}
