package cmd

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/argon2"
	"golang.org/x/term"
)

var vaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Password vault manager",
	Long:  `Securely store and manage passwords with encryption.`,
}

var vaultFile string

type VaultEntry struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Username  string    `json:"username"`
	Password  string    `json:"password"`
	URL       string    `json:"url,omitempty"`
	Notes     string    `json:"notes,omitempty"`
	Category  string    `json:"category"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Vault struct {
	Version int          `json:"version"`
	Entries []VaultEntry `json:"entries"`
}

var vaultInitCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize password vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(vaultFile); err == nil {
			return fmt.Errorf("vault already exists at %s", vaultFile)
		}

		flagPw, _ := cmd.Flags().GetString("password")
		password, err := readPassword(flagPw)
		if err != nil {
			return err
		}

		vault := Vault{
			Version: 1,
			Entries: []VaultEntry{},
		}

		return saveVault(vaultFile, password, vault)
	},
}

var vaultListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all entries in vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		flagPw, _ := cmd.Flags().GetString("password")
		password, err := readPassword(flagPw)
		if err != nil {
			return err
		}

		vault, err := loadVault(vaultFile, password)
		if err != nil {
			return err
		}

		category, _ := cmd.Flags().GetString("category")

		fmt.Printf("%-36s %-20s %-20s %s\n", "ID", "Title", "Username", "Category")
		fmt.Println(strings.Repeat("-", 100))

		for _, entry := range vault.Entries {
			if category != "" && entry.Category != category {
				continue
			}
			title := entry.Title
			if len(title) > 20 {
				title = title[:17] + "..."
			}
			username := entry.Username
			if len(username) > 20 {
				username = username[:17] + "..."
			}
			fmt.Printf("%-36s %-20s %-20s %s\n", entry.ID, title, username, entry.Category)
		}

		return nil
	},
}

var vaultAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add entry to vault",
	RunE: func(cmd *cobra.Command, args []string) error {
		flagPw, _ := cmd.Flags().GetString("password")
		password, err := readPassword(flagPw)
		if err != nil {
			return err
		}

		vault, err := loadVault(vaultFile, password)
		if err != nil {
			return err
		}

		title, _ := cmd.Flags().GetString("title")
		username, _ := cmd.Flags().GetString("username")
		entryPassword, _ := cmd.Flags().GetString("entry-password")
		url, _ := cmd.Flags().GetString("url")
		notes, _ := cmd.Flags().GetString("notes")
		category, _ := cmd.Flags().GetString("category")

		if title == "" {
			return fmt.Errorf("title is required")
		}

		entry := VaultEntry{
			ID:        generateID(),
			Title:     title,
			Username:  username,
			Password:  entryPassword,
			URL:       url,
			Notes:     notes,
			Category:  category,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}

		vault.Entries = append(vault.Entries, entry)

		if err := saveVault(vaultFile, password, vault); err != nil {
			return err
		}

		fmt.Printf("Added entry: %s\n", entry.ID)
		return nil
	},
}

var vaultGetCmd = &cobra.Command{
	Use:   "get [id]",
	Short: "Get entry from vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		flagPw, _ := cmd.Flags().GetString("password")
		password, err := readPassword(flagPw)
		if err != nil {
			return err
		}

		vault, err := loadVault(vaultFile, password)
		if err != nil {
			return err
		}

		for _, entry := range vault.Entries {
			if entry.ID == id {
				fmt.Printf("ID: %s\n", entry.ID)
				fmt.Printf("Title: %s\n", entry.Title)
				fmt.Printf("Username: %s\n", entry.Username)
				fmt.Printf("Password: %s\n", entry.Password)
				if entry.URL != "" {
					fmt.Printf("URL: %s\n", entry.URL)
				}
				if entry.Notes != "" {
					fmt.Printf("Notes: %s\n", entry.Notes)
				}
				fmt.Printf("Category: %s\n", entry.Category)
				fmt.Printf("Created: %s\n", entry.CreatedAt.Format("2006-01-02 15:04:05"))
				return nil
			}
		}

		return fmt.Errorf("entry not found: %s", id)
	},
}

var vaultRemoveCmd = &cobra.Command{
	Use:   "remove [id]",
	Short: "Remove entry from vault",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		id := args[0]

		flagPw, _ := cmd.Flags().GetString("password")
		password, err := readPassword(flagPw)
		if err != nil {
			return err
		}

		vault, err := loadVault(vaultFile, password)
		if err != nil {
			return err
		}

		for i, entry := range vault.Entries {
			if entry.ID == id {
				vault.Entries = append(vault.Entries[:i], vault.Entries[i+1:]...)
				if err := saveVault(vaultFile, password, vault); err != nil {
					return err
				}
				fmt.Printf("Removed entry: %s\n", id)
				return nil
			}
		}

		return fmt.Errorf("entry not found: %s", id)
	},
}

func loadVault(filename, password string) (Vault, error) {
	var vault Vault

	data, err := os.ReadFile(filename)
	if err != nil {
		return vault, err
	}

	decrypted, err := decrypt(data, password)
	if err != nil {
		return vault, fmt.Errorf("failed to decrypt vault (wrong password?)")
	}

	err = json.Unmarshal(decrypted, &vault)
	return vault, err
}

func saveVault(filename, password string, vault Vault) error {
	// Ensure the parent directory exists
	if err := os.MkdirAll(filepath.Dir(filename), 0700); err != nil {
		return fmt.Errorf("failed to create vault directory: %w", err)
	}

	data, err := json.Marshal(vault)
	if err != nil {
		return err
	}

	encrypted, err := encrypt(data, password)
	if err != nil {
		return err
	}

	return os.WriteFile(filename, encrypted, 0600)
}

// saltSize is the length of the random salt prepended to every ciphertext.
const saltSize = 16

// deriveKey stretches password into a 32-byte AES key using Argon2id.
func deriveKey(password string, salt []byte) []byte {
	return argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, 32)
}

func encrypt(plaintext []byte, password string) ([]byte, error) {
	salt := make([]byte, saltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return nil, fmt.Errorf("failed to generate salt: %w", err)
	}

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, fmt.Errorf("failed to generate nonce: %w", err)
	}

	// Layout: salt | nonce | ciphertext+tag
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return append(salt, ciphertext...), nil
}

func decrypt(ciphertext []byte, password string) ([]byte, error) {
	if len(ciphertext) < saltSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	salt := ciphertext[:saltSize]
	ciphertext = ciphertext[saltSize:]

	key := deriveKey(password, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	return gcm.Open(nil, nonce, ciphertext, nil)
}

func generateID() string {
	b := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, b); err != nil {
		panic(fmt.Sprintf("failed to generate random ID: %v", err))
	}
	return base64.URLEncoding.EncodeToString(b)
}

// readPassword reads a password from the terminal without echoing it.
// It falls back to the provided flag value if one is set.
func readPassword(flagValue string) (string, error) {
	if flagValue != "" {
		return flagValue, nil
	}
	fmt.Fprint(os.Stderr, "Master password: ")
	pw, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Fprintln(os.Stderr)
	if err != nil {
		return "", fmt.Errorf("failed to read password: %w", err)
	}
	if len(pw) == 0 {
		return "", fmt.Errorf("password cannot be empty")
	}
	return string(pw), nil
}

func init() {
	rootCmd.AddCommand(vaultCmd)
	vaultCmd.AddCommand(vaultInitCmd)
	vaultCmd.AddCommand(vaultListCmd)
	vaultCmd.AddCommand(vaultAddCmd)
	vaultCmd.AddCommand(vaultGetCmd)
	vaultCmd.AddCommand(vaultRemoveCmd)

	defaultVaultPath := filepath.Join(os.Getenv("HOME"), ".kairoa", "vault.dat")
	vaultCmd.PersistentFlags().StringVar(&vaultFile, "vault-file", defaultVaultPath, "Vault file path")
	vaultCmd.PersistentFlags().StringP("password", "p", "", "Master password")

	vaultInitCmd.Flags().StringP("password", "p", "", "Master password")
	
	vaultListCmd.Flags().StringP("category", "c", "", "Filter by category")
	
	vaultAddCmd.Flags().StringP("title", "t", "", "Entry title")
	vaultAddCmd.Flags().StringP("username", "u", "", "Username")
	vaultAddCmd.Flags().StringP("entry-password", "", "", "Password")
	vaultAddCmd.Flags().String("url", "", "URL")
	vaultAddCmd.Flags().String("notes", "", "Notes")
	vaultAddCmd.Flags().String("category", "general", "Category")
}
