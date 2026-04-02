package cmd

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/spf13/cobra"
)

var rsaCmd = &cobra.Command{
	Use:   "rsa",
	Short: "RSA key pair generator",
	Long:  `Generate RSA public/private key pairs.`,
}

var rsaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate RSA key pair",
	RunE: func(cmd *cobra.Command, args []string) error {
		bits, _ := cmd.Flags().GetInt("bits")
		format, _ := cmd.Flags().GetString("format")

		// Generate key pair
		privateKey, err := rsa.GenerateKey(rand.Reader, bits)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		// Encode private key
		var privateKeyPEM []byte
		if format == "pkcs8" {
			privateKeyBytes, _ := x509.MarshalPKCS8PrivateKey(privateKey)
			privateKeyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "PRIVATE KEY",
				Bytes: privateKeyBytes,
			})
		} else {
			privateKeyPEM = pem.EncodeToMemory(&pem.Block{
				Type:  "RSA PRIVATE KEY",
				Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
			})
		}

		// Encode public key
		publicKeyBytes, _ := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		publicKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "PUBLIC KEY",
			Bytes: publicKeyBytes,
		})

		fmt.Println("=== Private Key ===")
		fmt.Println(string(privateKeyPEM))
		fmt.Println("=== Public Key ===")
		fmt.Println(string(publicKeyPEM))

		return nil
	},
}

var rsaInfoCmd = &cobra.Command{
	Use:   "info [key]",
	Short: "Show RSA key information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		keyPEM := args[0]

		block, _ := pem.Decode([]byte(keyPEM))
		if block == nil {
			return fmt.Errorf("failed to decode PEM")
		}

		switch block.Type {
		case "RSA PRIVATE KEY":
			key, err := x509.ParsePKCS1PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse key: %w", err)
			}
			fmt.Printf("Type: RSA Private Key\n")
			fmt.Printf("Bits: %d\n", key.N.BitLen())
			fmt.Printf("Public Exponent: %d\n", key.E)

		case "PRIVATE KEY":
			key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse key: %w", err)
			}
			if rsaKey, ok := key.(*rsa.PrivateKey); ok {
				fmt.Printf("Type: RSA Private Key (PKCS#8)\n")
				fmt.Printf("Bits: %d\n", rsaKey.N.BitLen())
				fmt.Printf("Public Exponent: %d\n", rsaKey.E)
			}

		case "PUBLIC KEY":
			key, err := x509.ParsePKIXPublicKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse key: %w", err)
			}
			if rsaKey, ok := key.(*rsa.PublicKey); ok {
				fmt.Printf("Type: RSA Public Key\n")
				fmt.Printf("Bits: %d\n", rsaKey.N.BitLen())
				fmt.Printf("Exponent: %d\n", rsaKey.E)
			}

		case "RSA PUBLIC KEY":
			key, err := x509.ParsePKCS1PublicKey(block.Bytes)
			if err != nil {
				return fmt.Errorf("failed to parse key: %w", err)
			}
			fmt.Printf("Type: RSA Public Key (PKCS#1)\n")
			fmt.Printf("Bits: %d\n", key.N.BitLen())
			fmt.Printf("Exponent: %d\n", key.E)

		default:
			return fmt.Errorf("unknown key type: %s", block.Type)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(rsaCmd)
	rsaCmd.AddCommand(rsaGenerateCmd)
	rsaCmd.AddCommand(rsaInfoCmd)

	rsaGenerateCmd.Flags().IntP("bits", "b", 2048, "Key size in bits (1024, 2048, 3072, 4096)")
	rsaGenerateCmd.Flags().StringP("format", "f", "pkcs1", "Private key format (pkcs1, pkcs8)")
}
