package cmd

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"fmt"

	"github.com/spf13/cobra"
)

var ed25519Cmd = &cobra.Command{
	Use:   "ed25519",
	Short: "Ed25519 key pair generator and signer",
	Long:  `Generate Ed25519 key pairs and sign/verify messages. Ed25519 is a modern, fast, and secure digital signature algorithm.`,
}

var ed25519GenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate Ed25519 key pair",
	RunE: func(cmd *cobra.Command, args []string) error {
		publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		fmt.Println("=== Private Key (Base64) ===")
		fmt.Println(base64.StdEncoding.EncodeToString(privateKey))
		fmt.Println("=== Public Key (Base64) ===")
		fmt.Println(base64.StdEncoding.EncodeToString(publicKey))

		return nil
	},
}

var ed25519SignCmd = &cobra.Command{
	Use:   "sign [message]",
	Short: "Sign a message with Ed25519 private key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]
		keyB64, _ := cmd.Flags().GetString("key")

		if keyB64 == "" {
			return fmt.Errorf("private key is required (--key)")
		}

		privateKeyBytes, err := base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			return fmt.Errorf("failed to decode private key: %w", err)
		}

		if len(privateKeyBytes) != ed25519.PrivateKeySize {
			return fmt.Errorf("invalid private key size: expected %d, got %d", ed25519.PrivateKeySize, len(privateKeyBytes))
		}

		privateKey := ed25519.PrivateKey(privateKeyBytes)
		signature := ed25519.Sign(privateKey, []byte(message))

		fmt.Println(base64.StdEncoding.EncodeToString(signature))
		return nil
	},
}

var ed25519VerifyCmd = &cobra.Command{
	Use:   "verify [message] [signature]",
	Short: "Verify a signature with Ed25519 public key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]
		signatureB64 := args[1]
		keyB64, _ := cmd.Flags().GetString("key")

		if keyB64 == "" {
			return fmt.Errorf("public key is required (--key)")
		}

		publicKeyBytes, err := base64.StdEncoding.DecodeString(keyB64)
		if err != nil {
			return fmt.Errorf("failed to decode public key: %w", err)
		}

		if len(publicKeyBytes) != ed25519.PublicKeySize {
			return fmt.Errorf("invalid public key size: expected %d, got %d", ed25519.PublicKeySize, len(publicKeyBytes))
		}

		signature, err := base64.StdEncoding.DecodeString(signatureB64)
		if err != nil {
			return fmt.Errorf("failed to decode signature: %w", err)
		}

		publicKey := ed25519.PublicKey(publicKeyBytes)
		valid := ed25519.Verify(publicKey, []byte(message), signature)

		if valid {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature is invalid")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ed25519Cmd)
	ed25519Cmd.AddCommand(ed25519GenerateCmd)
	ed25519Cmd.AddCommand(ed25519SignCmd)
	ed25519Cmd.AddCommand(ed25519VerifyCmd)

	ed25519SignCmd.Flags().StringP("key", "k", "", "Private key (Base64)")
	ed25519VerifyCmd.Flags().StringP("key", "k", "", "Public key (Base64)")
}
