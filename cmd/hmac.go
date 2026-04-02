package cmd

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/spf13/cobra"
)

var hmacCmd = &cobra.Command{
	Use:   "hmac",
	Short: "Calculate HMAC",
	Long:  `Calculate HMAC (Hash-based Message Authentication Code) for text or files.`,
}

var hmacTextCmd = &cobra.Command{
	Use:   "text [string]",
	Short: "Calculate HMAC of a text string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]
		key, _ := cmd.Flags().GetString("key")
		algorithm, _ := cmd.Flags().GetString("algorithm")

		if key == "" {
			return fmt.Errorf("key is required")
		}

		h := getHMACHasher(algorithm, []byte(key))
		if h == nil {
			return fmt.Errorf("unknown algorithm: %s", algorithm)
		}

		h.Write([]byte(text))
		result := hex.EncodeToString(h.Sum(nil))
		fmt.Println(result)
		return nil
	},
}

var hmacFileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: "Calculate HMAC of a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		key, _ := cmd.Flags().GetString("key")
		algorithm, _ := cmd.Flags().GetString("algorithm")

		if key == "" {
			return fmt.Errorf("key is required")
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer file.Close()

		h := getHMACHasher(algorithm, []byte(key))
		if h == nil {
			return fmt.Errorf("unknown algorithm: %s", algorithm)
		}

		if _, err := io.Copy(h, file); err != nil {
			return fmt.Errorf("failed to calculate HMAC: %w", err)
		}

		result := hex.EncodeToString(h.Sum(nil))
		fmt.Println(result)
		return nil
	},
}

func getHMACHasher(algo string, key []byte) hash.Hash {
	switch algo {
	case "md5":
		return hmac.New(md5.New, key)
	case "sha1":
		return hmac.New(sha1.New, key)
	case "sha256":
		return hmac.New(sha256.New, key)
	case "sha512":
		return hmac.New(sha512.New, key)
	default:
		return nil
	}
}

func init() {
	rootCmd.AddCommand(hmacCmd)
	hmacCmd.AddCommand(hmacTextCmd)
	hmacCmd.AddCommand(hmacFileCmd)

	hmacTextCmd.Flags().StringP("key", "k", "", "Secret key (required)")
	hmacTextCmd.Flags().StringP("algorithm", "a", "sha256", "Hash algorithm (md5, sha1, sha256, sha512)")
	hmacFileCmd.Flags().StringP("key", "k", "", "Secret key (required)")
	hmacFileCmd.Flags().StringP("algorithm", "a", "sha256", "Hash algorithm (md5, sha1, sha256, sha512)")
}
