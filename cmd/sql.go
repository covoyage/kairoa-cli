package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var sqlCmd = &cobra.Command{
	Use:   "sql",
	Short: "SQL utilities",
	Long:  `Format and validate SQL queries.`,
}

var sqlFormatCmd = &cobra.Command{
	Use:   "format [file]",
	Short: "Format SQL query",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			content, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			input = string(content)
		} else {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			input = string(content)
		}

		indent, _ := cmd.Flags().GetString("indent")
		uppercase, _ := cmd.Flags().GetBool("uppercase")

		formatted := formatSQL(input, indent, uppercase)
		fmt.Println(formatted)
		return nil
	},
}

var sqlMinifyCmd = &cobra.Command{
	Use:   "minify [file]",
	Short: "Minify SQL query",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input string
		if len(args) > 0 {
			content, err := os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
			input = string(content)
		} else {
			content, err := io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
			input = string(content)
		}

		minified := minifySQL(input)
		fmt.Println(minified)
		return nil
	},
}

// Simple SQL formatter
func formatSQL(sql, indent string, uppercase bool) string {
	keywords := []string{
		"SELECT", "FROM", "WHERE", "AND", "OR", "INSERT", "INTO", "VALUES",
		"UPDATE", "SET", "DELETE", "JOIN", "LEFT", "RIGHT", "INNER", "OUTER",
		"ON", "GROUP", "BY", "ORDER", "HAVING", "LIMIT", "OFFSET", "UNION",
		"ALL", "DISTINCT", "AS", "CREATE", "TABLE", "DROP", "ALTER", "INDEX",
		"PRIMARY", "KEY", "FOREIGN", "REFERENCES", "NOT", "NULL", "DEFAULT",
		"AUTO_INCREMENT", "UNIQUE", "CHECK", "CONSTRAINT", "IF", "EXISTS",
		"CASE", "WHEN", "THEN", "ELSE", "END", "CAST", "CONVERT", "LIKE",
		"IN", "BETWEEN", "IS", "EXISTS", "COUNT", "SUM", "AVG", "MIN", "MAX",
	}

	// Normalize whitespace
	sql = strings.TrimSpace(sql)
	sql = strings.ReplaceAll(sql, "\t", " ")
	for strings.Contains(sql, "  ") {
		sql = strings.ReplaceAll(sql, "  ", " ")
	}

	// Uppercase/lowercase keywords
	for _, kw := range keywords {
		lowerKw := strings.ToLower(kw)
		if uppercase {
			sql = replaceWord(sql, lowerKw, kw)
		} else {
			sql = replaceWord(sql, kw, lowerKw)
		}
	}

	// Add newlines before major keywords
	majorKeywords := []string{"SELECT", "FROM", "WHERE", "GROUP BY", "ORDER BY", "HAVING", "LIMIT", "UNION"}
	for _, kw := range majorKeywords {
		lowerKw := strings.ToLower(kw)
		sql = strings.ReplaceAll(sql, " "+kw+" ", "\n"+kw+" ")
		sql = strings.ReplaceAll(sql, " "+lowerKw+" ", "\n"+lowerKw+" ")
	}

	// Add newlines for JOINs
	joinKeywords := []string{"JOIN", "LEFT JOIN", "RIGHT JOIN", "INNER JOIN", "OUTER JOIN", "CROSS JOIN"}
	for _, kw := range joinKeywords {
		lowerKw := strings.ToLower(kw)
		sql = strings.ReplaceAll(sql, " "+kw+" ", "\n"+kw+" ")
		sql = strings.ReplaceAll(sql, " "+lowerKw+" ", "\n"+lowerKw+" ")
	}

	// Indent lines
	lines := strings.Split(sql, "\n")
	var result strings.Builder
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}

		// Add proper indentation
		if strings.HasPrefix(strings.ToUpper(trimmed), "FROM") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "WHERE") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "GROUP") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "ORDER") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "HAVING") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "LIMIT") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "UNION") ||
			strings.Contains(strings.ToUpper(trimmed), "JOIN") {
			result.WriteString(trimmed)
		} else if strings.HasPrefix(strings.ToUpper(trimmed), "SELECT") {
			result.WriteString(trimmed)
		} else if strings.HasPrefix(strings.ToUpper(trimmed), "AND") ||
			strings.HasPrefix(strings.ToUpper(trimmed), "OR") {
			result.WriteString(indent + trimmed)
		} else {
			result.WriteString(indent + trimmed)
		}

		if i < len(lines)-1 {
			result.WriteString("\n")
		}
	}

	return result.String()
}

func minifySQL(sql string) string {
	// Remove all extra whitespace
	sql = strings.TrimSpace(sql)
	sql = strings.ReplaceAll(sql, "\n", " ")
	sql = strings.ReplaceAll(sql, "\t", " ")
	for strings.Contains(sql, "  ") {
		sql = strings.ReplaceAll(sql, "  ", " ")
	}
	return sql
}

func replaceWord(s, old, new string) string {
	// Simple word replacement - replace standalone words
	result := s
	oldUpper := strings.ToUpper(old)
	oldLower := strings.ToLower(old)

	// Replace if surrounded by non-word characters or boundaries
	result = strings.ReplaceAll(result, " "+oldUpper+" ", " "+new+" ")
	result = strings.ReplaceAll(result, " "+oldLower+" ", " "+new+" ")
	result = strings.ReplaceAll(result, "("+oldUpper+" ", "("+new+" ")
	result = strings.ReplaceAll(result, "("+oldLower+" ", "("+new+" ")
	result = strings.ReplaceAll(result, " "+oldUpper+")", " "+new+")")
	result = strings.ReplaceAll(result, " "+oldLower+")", " "+new+")")

	return result
}

func init() {
	rootCmd.AddCommand(sqlCmd)
	sqlCmd.AddCommand(sqlFormatCmd)
	sqlCmd.AddCommand(sqlMinifyCmd)

	sqlFormatCmd.Flags().StringP("indent", "i", "  ", "Indentation string")
	sqlFormatCmd.Flags().BoolP("uppercase", "u", false, "Uppercase keywords")
}
