package cmd

import (
	"crypto/hmac"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base32"
	"encoding/binary"
	"fmt"
	"hash"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var otpCmd = &cobra.Command{
	Use:   "otp",
	Short: "OTP (One-Time Password) generator",
	Long:  `Generate TOTP and HOTP codes like Google Authenticator.`,
}

var otpTotpCmd = &cobra.Command{
	Use:   "totp [secret]",
	Short: "Generate TOTP code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := args[0]
		digits, _ := cmd.Flags().GetInt("digits")
		period, _ := cmd.Flags().GetInt("period")
		algorithm, _ := cmd.Flags().GetString("algorithm")

		code, err := generateTOTP(secret, digits, period, algorithm)
		if err != nil {
			return err
		}

		fmt.Printf("TOTP Code: %s\n", code)

		// Show remaining time
		remaining := period - int(time.Now().Unix()%int64(period))
		fmt.Printf("Expires in: %d seconds\n", remaining)

		return nil
	},
}

var otpHotpCmd = &cobra.Command{
	Use:   "hotp [secret] [counter]",
	Short: "Generate HOTP code",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := args[0]
		var counter uint64
		fmt.Sscanf(args[1], "%d", &counter)

		digits, _ := cmd.Flags().GetInt("digits")
		algorithm, _ := cmd.Flags().GetString("algorithm")

		code, err := generateHOTP(secret, counter, digits, algorithm)
		if err != nil {
			return err
		}

		fmt.Printf("HOTP Code: %s\n", code)
		return nil
	},
}

var otpVerifyCmd = &cobra.Command{
	Use:   "verify [secret] [code]",
	Short: "Verify TOTP code",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		secret := args[0]
		code := args[1]
		digits, _ := cmd.Flags().GetInt("digits")
		period, _ := cmd.Flags().GetInt("period")
		algorithm, _ := cmd.Flags().GetString("algorithm")
		window, _ := cmd.Flags().GetInt("window")

		// Check current and adjacent time windows
		now := time.Now().Unix()
		for i := -window; i <= window; i++ {
			timestamp := now + int64(i*period)
			expectedCode, err := generateTOTPAt(secret, digits, period, algorithm, timestamp)
			if err != nil {
				continue
			}
			if expectedCode == code {
				fmt.Println("✓ Valid code")
				return nil
			}
		}

		fmt.Println("✗ Invalid code")
		return nil
	},
}

func generateTOTP(secret string, digits, period int, algorithm string) (string, error) {
	timestamp := time.Now().Unix()
	return generateTOTPAt(secret, digits, period, algorithm, timestamp)
}

func generateTOTPAt(secret string, digits, period int, algorithm string, timestamp int64) (string, error) {
	counter := uint64(timestamp / int64(period))
	return generateHOTP(secret, counter, digits, algorithm)
}

func generateHOTP(secret string, counter uint64, digits int, algorithm string) (string, error) {
	// Normalize secret
	secret = strings.ToUpper(secret)
	secret = strings.ReplaceAll(secret, " ", "")

	// Decode base32 secret
	key, err := base32.StdEncoding.DecodeString(secret)
	if err != nil {
		return "", fmt.Errorf("invalid secret: %w", err)
	}

	// Create HMAC
	var hasher func() hash.Hash
	switch strings.ToUpper(algorithm) {
	case "SHA256":
		hasher = sha256.New
	case "SHA512":
		hasher = sha512.New
	default:
		hasher = sha1.New
	}

	h := hmac.New(hasher, key)
	
	// Write counter as 8 bytes big-endian
	counterBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(counterBytes, counter)
	h.Write(counterBytes)
	
	hash := h.Sum(nil)

	// Dynamic truncation
	offset := hash[len(hash)-1] & 0x0f
	code := (int(hash[offset])&0x7f)<<24 |
		(int(hash[offset+1])&0xff)<<16 |
		(int(hash[offset+2])&0xff)<<8 |
		(int(hash[offset+3]) & 0xff)

	// Modulo to get desired digits
	code = code % pow10(digits)

	return fmt.Sprintf(fmt.Sprintf("%%0%dd", digits), code), nil
}

func pow10(n int) int {
	result := 1
	for i := 0; i < n; i++ {
		result *= 10
	}
	return result
}

func init() {
	rootCmd.AddCommand(otpCmd)
	otpCmd.AddCommand(otpTotpCmd)
	otpCmd.AddCommand(otpHotpCmd)
	otpCmd.AddCommand(otpVerifyCmd)

	otpTotpCmd.Flags().IntP("digits", "d", 6, "Number of digits")
	otpTotpCmd.Flags().IntP("period", "p", 30, "Time period in seconds")
	otpTotpCmd.Flags().StringP("algorithm", "a", "SHA1", "Hash algorithm (SHA1, SHA256, SHA512)")

	otpHotpCmd.Flags().IntP("digits", "d", 6, "Number of digits")
	otpHotpCmd.Flags().StringP("algorithm", "a", "SHA1", "Hash algorithm (SHA1, SHA256, SHA512)")

	otpVerifyCmd.Flags().IntP("digits", "d", 6, "Number of digits")
	otpVerifyCmd.Flags().IntP("period", "p", 30, "Time period in seconds")
	otpVerifyCmd.Flags().StringP("algorithm", "a", "SHA1", "Hash algorithm (SHA1, SHA256, SHA512)")
	otpVerifyCmd.Flags().IntP("window", "w", 1, "Time window tolerance")
}
