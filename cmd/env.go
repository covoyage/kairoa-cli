package cmd

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/spf13/cobra"
)

var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Environment variable manager",
	Long:  `Manage environment variables from .env files.`,
}

var envLoadCmd = &cobra.Command{
	Use:   "load [file]",
	Short: "Load and display .env file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		
		vars, err := loadEnvFile(file)
		if err != nil {
			return err
		}

		format, _ := cmd.Flags().GetString("format")
		
		if format == "export" {
			for key, value := range vars {
				fmt.Printf("export %s=\"%s\"\n", key, value)
			}
		} else {
			// Sort keys
			var keys []string
			for key := range vars {
				keys = append(keys, key)
			}
			sort.Strings(keys)

			for _, key := range keys {
				value := vars[key]
				if strings.Contains(value, " ") {
					fmt.Printf("%s=\"%s\"\n", key, value)
				} else {
					fmt.Printf("%s=%s\n", key, value)
				}
			}
		}

		return nil
	},
}

var envGetCmd = &cobra.Command{
	Use:   "get [key]",
	Short: "Get environment variable",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		file, _ := cmd.Flags().GetString("file")
		
		var value string
		if file != "" {
			vars, err := loadEnvFile(file)
			if err != nil {
				return err
			}
			value = vars[key]
		} else {
			value = os.Getenv(key)
		}

		if value == "" {
			return fmt.Errorf("environment variable not found: %s", key)
		}

		fmt.Println(value)
		return nil
	},
}

var envSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set environment variable in .env file",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]
		value := args[1]
		file, _ := cmd.Flags().GetString("file")
		
		if file == "" {
			file = ".env"
		}

		vars, _ := loadEnvFile(file)
		vars[key] = value

		return saveEnvFile(file, vars)
	},
}

var envListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all environment variables",
	RunE: func(cmd *cobra.Command, args []string) error {
		file, _ := cmd.Flags().GetString("file")
		
		var vars map[string]string
		var err error
		
		if file != "" {
			vars, err = loadEnvFile(file)
			if err != nil {
				return err
			}
		} else {
			vars = make(map[string]string)
			for _, e := range os.Environ() {
				parts := strings.SplitN(e, "=", 2)
				if len(parts) == 2 {
					vars[parts[0]] = parts[1]
				}
			}
		}

		// Sort keys
		var keys []string
		for key := range vars {
			keys = append(keys, key)
		}
		sort.Strings(keys)

		for _, key := range keys {
			value := vars[key]
			fmt.Printf("%s=%s\n", key, value)
		}

		return nil
	},
}

var envValidateCmd = &cobra.Command{
	Use:   "validate [file]",
	Short: "Validate .env file format",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]
		
		vars, err := loadEnvFile(file)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		fmt.Printf("✓ Valid .env file with %d variables\n", len(vars))
		
		// Check for duplicates
		seen := make(map[string]bool)
		f, _ := os.Open(file)
		defer f.Close()
		
		scanner := bufio.NewScanner(f)
		lineNum := 0
		hasDuplicates := false
		
		for scanner.Scan() {
			lineNum++
			line := strings.TrimSpace(scanner.Text())
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				if seen[key] {
					fmt.Printf("⚠ Duplicate key at line %d: %s\n", lineNum, key)
					hasDuplicates = true
				}
				seen[key] = true
			}
		}

		if !hasDuplicates {
			fmt.Println("✓ No duplicate keys found")
		}

		return nil
	},
}

func loadEnvFile(filename string) (map[string]string, error) {
	vars := make(map[string]string)

	file, err := os.Open(filename)
	if err != nil {
		return vars, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Parse key=value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])
			
			// Remove quotes if present
			if (strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")) ||
			   (strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")) {
				value = value[1 : len(value)-1]
			}
			
			vars[key] = value
		}
	}

	return vars, scanner.Err()
}

func saveEnvFile(filename string, vars map[string]string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Sort keys
	var keys []string
	for key := range vars {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		value := vars[key]
		if strings.Contains(value, " ") || strings.Contains(value, "$") {
			fmt.Fprintf(file, "%s=\"%s\"\n", key, value)
		} else {
			fmt.Fprintf(file, "%s=%s\n", key, value)
		}
	}

	return nil
}

func init() {
	rootCmd.AddCommand(envCmd)
	envCmd.AddCommand(envLoadCmd)
	envCmd.AddCommand(envGetCmd)
	envCmd.AddCommand(envSetCmd)
	envCmd.AddCommand(envListCmd)
	envCmd.AddCommand(envValidateCmd)

	envLoadCmd.Flags().StringP("format", "f", "default", "Output format (default, export)")
	envGetCmd.Flags().StringP("file", "f", "", "Env file path")
	envSetCmd.Flags().StringP("file", "f", ".env", "Env file path")
	envListCmd.Flags().StringP("file", "f", "", "Env file path (if not set, uses system env)")
}
