package cmd

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/spf13/cobra"
)

var ecdsaCmd = &cobra.Command{
	Use:   "ecdsa",
	Short: "ECDSA key pair generator and signer",
	Long:  `Generate ECDSA (Elliptic Curve Digital Signature Algorithm) key pairs and sign/verify messages.`,
}

var ecdsaGenerateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate ECDSA key pair",
	RunE: func(cmd *cobra.Command, args []string) error {
		curveName, _ := cmd.Flags().GetString("curve")

		var curve elliptic.Curve
		switch curveName {
		case "P-224":
			curve = elliptic.P224()
		case "P-256":
			curve = elliptic.P256()
		case "P-384":
			curve = elliptic.P384()
		case "P-521":
			curve = elliptic.P521()
		default:
			return fmt.Errorf("unsupported curve: %s (use P-224, P-256, P-384, or P-521)", curveName)
		}

		privateKey, err := ecdsa.GenerateKey(curve, rand.Reader)
		if err != nil {
			return fmt.Errorf("failed to generate key: %w", err)
		}

		// Encode private key
		privateKeyBytes, err := x509.MarshalECPrivateKey(privateKey)
		if err != nil {
			return fmt.Errorf("failed to marshal private key: %w", err)
		}
		privateKeyPEM := pem.EncodeToMemory(&pem.Block{
			Type:  "EC PRIVATE KEY",
			Bytes: privateKeyBytes,
		})

		// Encode public key
		publicKeyBytes, err := x509.MarshalPKIXPublicKey(&privateKey.PublicKey)
		if err != nil {
			return fmt.Errorf("failed to marshal public key: %w", err)
		}
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

var ecdsaSignCmd = &cobra.Command{
	Use:   "sign [message]",
	Short: "Sign a message with ECDSA private key",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]
		keyPEM, _ := cmd.Flags().GetString("key")

		if keyPEM == "" {
			return fmt.Errorf("private key is required (--key)")
		}

		block, _ := pem.Decode([]byte(keyPEM))
		if block == nil {
			return fmt.Errorf("failed to decode PEM")
		}

		privateKey, err := x509.ParseECPrivateKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse private key: %w", err)
		}

		signature, err := ecdsa.SignASN1(rand.Reader, privateKey, []byte(message))
		if err != nil {
			return fmt.Errorf("failed to sign: %w", err)
		}

		fmt.Println(base64.StdEncoding.EncodeToString(signature))
		return nil
	},
}

var ecdsaVerifyCmd = &cobra.Command{
	Use:   "verify [message] [signature]",
	Short: "Verify a signature with ECDSA public key",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := args[0]
		signatureB64 := args[1]
		keyPEM, _ := cmd.Flags().GetString("key")

		if keyPEM == "" {
			return fmt.Errorf("public key is required (--key)")
		}

		signature, err := base64.StdEncoding.DecodeString(signatureB64)
		if err != nil {
			return fmt.Errorf("failed to decode signature: %w", err)
		}

		block, _ := pem.Decode([]byte(keyPEM))
		if block == nil {
			return fmt.Errorf("failed to decode PEM")
		}

		publicKeyInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse public key: %w", err)
		}

		publicKey, ok := publicKeyInterface.(*ecdsa.PublicKey)
		if !ok {
			return fmt.Errorf("not an ECDSA public key")
		}

		valid := ecdsa.VerifyASN1(publicKey, []byte(message), signature)
		if valid {
			fmt.Println("Signature is valid")
		} else {
			fmt.Println("Signature is invalid")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(ecdsaCmd)
	ecdsaCmd.AddCommand(ecdsaGenerateCmd)
	ecdsaCmd.AddCommand(ecdsaSignCmd)
	ecdsaCmd.AddCommand(ecdsaVerifyCmd)

	ecdsaGenerateCmd.Flags().StringP("curve", "c", "P-256", "Elliptic curve (P-224, P-256, P-384, P-521)")
	ecdsaSignCmd.Flags().StringP("key", "k", "", "Private key PEM")
	ecdsaVerifyCmd.Flags().StringP("key", "k", "", "Public key PEM")
}
