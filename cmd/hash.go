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
	"github.com/zeebo/blake3"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/blake2s"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

var hashCmd = &cobra.Command{
	Use:   "hash",
	Short: i18n.T("hash.short"),
	Long:  i18n.T("hash.long"),
}

// applyHashLocale refreshes all translatable strings in the hash command tree.
func applyHashLocale() {
	hashCmd.Short = i18n.T("hash.short")
	hashCmd.Long = i18n.T("hash.long")
	hashTextCmd.Short = i18n.T("hash.text")
	hashFileCmd.Short = i18n.T("hash.file")
	if f := hashTextCmd.Flags().Lookup("algorithm"); f != nil {
		f.Usage = i18n.T("hash.algorithm")
	}
	if f := hashFileCmd.Flags().Lookup("algorithm"); f != nil {
		f.Usage = i18n.T("hash.algorithm")
	}
}

var hashTextCmd = &cobra.Command{
	Use:   "text [string]",
	Short: i18n.T("hash.text"),
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]
		algorithms, _ := cmd.Flags().GetStringSlice("algorithm")

		if len(algorithms) == 0 {
			algorithms = []string{"md5", "sha1", "sha256", "sha384", "sha512", "sha3-256", "sha3-512", "blake2b-256", "blake2b-512", "blake2s-256", "blake3", "ripemd160"}
		}

		for _, algo := range algorithms {
			h := getHasher(algo)
			if h == nil {
				fmt.Fprintln(os.Stderr, i18n.T("error.unknown", algo))
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
			algorithms = []string{"md5", "sha1", "sha256", "sha384", "sha512", "sha3-256", "sha3-512", "blake2b-256", "blake2b-512", "blake2s-256", "blake3", "ripemd160"}
		}

		file, err := os.Open(filePath)
		if err != nil {
			return fmt.Errorf("%s: %w", i18n.T("error.readFile", filePath), err)
		}
		defer file.Close()

		for _, algo := range algorithms {
			h := getHasher(algo)
			if h == nil {
				fmt.Fprintln(os.Stderr, i18n.T("error.unknown", algo))
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
	case "sha3-224":
		return sha3.New224()
	case "sha3-256":
		return sha3.New256()
	case "sha3-384":
		return sha3.New384()
	case "sha3-512":
		return sha3.New512()
	case "blake2b-256":
		h, _ := blake2b.New256(nil)
		return h
	case "blake2b-384":
		h, _ := blake2b.New384(nil)
		return h
	case "blake2b-512":
		h, _ := blake2b.New512(nil)
		return h
	case "blake2s-256":
		h, _ := blake2s.New256(nil)
		return h
	case "blake3":
		return blake3.New()
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

	RegisterLocaleApplier(applyHashLocale)
}
