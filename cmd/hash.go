package cmd

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"hash"
	"io"
	"os"

	"github.com/covoyage/kairoa-cli/internal/i18n"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ripemd160"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: i18n.T("hash.short"),
	Long:  i18n.T("hash.long"),
}

var hashTextCmd = &cobra.Command{
	Use:   "text [string]",
	Short: i18n.T("hash.text"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]
		algorithms, _ := cmd.Flags().GetStringSlice("algorithm")

		if len(algorithms) == 0 {
			algorithms = []string{"md5", "sha1", "sha256", "sha384", "sha512", "ripemd160"}
		}

		for _, algo := range algorithms {
			h := getHasher(algo)
			if h == nil {
				fmt.Fprintf(os.Stderr, i18n.T("error.unknown")+": %s\n", algo)
				continue
			}
			h.Write([]byte(text))
			result := hex.EncodeToString(h.Sum(nil))
			fmt.Printf("%s: %s\n", algo, result)
		}
		return nil
	},
}

var hashFileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: i18n.T("hash.file"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		algorithms, _ := cmd.Flags().GetStringSlice("algorithm")

		if len(algorithms) == 0 {
			algorithms = []string{"md5", "sha1", "sha256", "sha384", "sha512", "ripemd160"}
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T("error.readFile", filePath), err)
		}
		defer file.Close()

		for _, algo := range algorithms {
			h := getHasher(algo)
			if h == nil {
				fmt.Fprintf(os.Stderr, i18n.T("error.unknown")+": %s\n", algo)
				continue
			}

			file.Seek(0, 0)
			if _, err := io.Copy(h, file); err != nil {
				return fmt.Errorf("%s: %w", i18n.T("error.readFile", filePath), err)
			}
			result := hex.EncodeToString(h.Sum(nil))
			fmt.Printf("%s: %s\n", algo, result)
		}

		return nil
	},
}

func getHasher(algo string) hash.Hash {
	switch algo {
	case "md5":
		return md5.New()
	case "sha1":
		return sha1.New()
	case "sha256":
		return sha256.New()
	case "sha384":
		return sha512.New384()
	case "sha512":
		return sha512.New()
	case "ripemd160":
		return ripemd160.New()
	default:
		return nil
	}
}

func init() {
	rootCmd.AddCommand(hashCmd)
	hashCmd.AddCommand(hashTextCmd)
	hashCmd.AddCommand(hashFileCmd)

	hashTextCmd.Flags().StringSliceP("algorithm", "a", []string{}, i18n.T("hash.algorithm"))
	hashFileCmd.Flags().StringSliceP("algorithm", "a", []string{}, i18n.T("hash.algorithm"))
}
