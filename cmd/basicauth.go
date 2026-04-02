package cmd

import (
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var basicauthCmd = &cobra.Command{
	Use:   "basicauth",
	Short: "Basic Auth generator",
	Long:  `Generate Basic Authentication headers and credentials.`,
}

var basicauthGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Basic Auth header",
	RunE: func(cmd *cobra.Command, args []string) error {
		username, _ := cmd.Flags().GetString("username")
		password, _ := cmd.Flags().GetString("password")

		if username == "" && password == "" {
			return fmt.Errorf("username or password required")
		}

		// Generate Basic Auth string
		credentials := fmt.Sprintf("%s:%s", username, password)
		encoded := base64.StdEncoding.EncodeToString([]byte(credentials))

		fmt.Printf("Credentials: %s\n", credentials)
		fmt.Printf("Base64: %s\n", encoded)
		fmt.Printf("Authorization Header: Basic %s\n", encoded)
		fmt.Printf("cURL: curl -u '%s:%s' https://example.com\n", username, password)
		fmt.Printf("cURL (header): curl -H 'Authorization: Basic %s' https://example.com\n", encoded)

		return nil
	},
}

var basicauthDecodeCmd = &cobra.Command{
	Use:   "decode [base64]",
	Short: "Decode Basic Auth string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		encoded := args[0]

		// Remove "Basic " prefix if present
		if len(encoded) > 6 && encoded[:6] == "Basic " {
			encoded = encoded[6:]
		}

		decoded, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return fmt.Errorf("invalid base64: %w", err)
		}

		parts := string(decoded)
		for i, c := range parts {
			if c == ':' {
				username := parts[:i]
				password := parts[i+1:]
				fmt.Printf("Username: %s\n", username)
				fmt.Printf("Password: %s\n", password)
				return nil
			}
		}

		return fmt.Errorf("invalid credentials format")
	},
}

func init() {
	rootCmd.AddCommand(basicauthCmd)
	basicauthCmd.AddCommand(basicauthGenerateCmd)
	basicauthCmd.AddCommand(basicauthDecodeCmd)

	basicauthGenerateCmd.Flags().StringP("username", "u", "", "Username")
	basicauthGenerateCmd.Flags().StringP("password", "p", "", "Password")
}
