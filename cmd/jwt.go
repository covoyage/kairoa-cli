package cmd

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var jwtCmd = &cobra.Command{
	Use:   "jwt [token]",
	Short: "Decode JWT tokens",
	Long:  `Decode and display the contents of JWT (JSON Web Token) tokens.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		token := args[0]

		// Remove "Bearer " prefix if present
		token = strings.TrimPrefix(token, "Bearer ")
		token = strings.TrimPrefix(token, "bearer ")

		parts := strings.Split(token, ".")
		if len(parts) != 3 {
			return fmt.Errorf("invalid JWT format: expected 3 parts, got %d", len(parts))
		}

		showHeader, _ := cmd.Flags().GetBool("header")
		showPayload, _ := cmd.Flags().GetBool("payload")
		showSignature, _ := cmd.Flags().GetBool("signature")

		// If no specific parts requested, show all
		if !showHeader && !showPayload && !showSignature {
			showHeader = true
			showPayload = true
			showSignature = true
		}

		if showHeader {
			header, err := decodeJWTPart(parts[0])
			if err != nil {
				return fmt.Errorf("failed to decode header: %w", err)
			}
			fmt.Println("Header:")
			fmt.Println(header)
		}

		if showPayload {
			payload, err := decodeJWTPart(parts[1])
			if err != nil {
				return fmt.Errorf("failed to decode payload: %w", err)
			}
			if showHeader {
				fmt.Println()
			}
			fmt.Println("Payload:")
			fmt.Println(payload)
		}

		if showSignature {
			if showHeader || showPayload {
				fmt.Println()
			}
			fmt.Println("Signature:")
			fmt.Println(parts[2])
		}

		return nil
	},
}

func decodeJWTPart(part string) (string, error) {
	// Add padding if needed
	padding := 4 - len(part)%4
	if padding != 4 {
		part += strings.Repeat("=", padding)
	}

	// Replace URL-safe characters
	part = strings.ReplaceAll(part, "-", "+")
	part = strings.ReplaceAll(part, "_", "/")

	decoded, err := base64.StdEncoding.DecodeString(part)
	if err != nil {
		return "", err
	}

	var data interface{}
	if err := json.Unmarshal(decoded, &data); err != nil {
		return string(decoded), nil
	}

	formatted, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return string(decoded), nil
	}

	return string(formatted), nil
}

func init() {
	rootCmd.AddCommand(jwtCmd)
	jwtCmd.Flags().BoolP("header", "H", false, "Show header only")
	jwtCmd.Flags().BoolP("payload", "p", false, "Show payload only")
	jwtCmd.Flags().BoolP("signature", "s", false, "Show signature only")
}
