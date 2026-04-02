package cmd

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

var dataCmd = &cobra.Command{
	Use:   "data",
	Short: "Data format conversion",
	Long:  `Convert between CSV and JSON formats.`,
}

var dataCsv2JsonCmd = &cobra.Command{
	Use:   "csv2json [file]",
	Short: "Convert CSV to JSON",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input io.Reader
		if len(args) > 0 {
			file, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("failed to open file: %w", err)
			}
			defer file.Close()
			input = file
		} else {
			input = os.Stdin
		}

		delimiter, _ := cmd.Flags().GetString("delimiter")
		var sep rune
		switch delimiter {
		case "comma", ",":
			sep = ','
		case "semicolon", ";":
			sep = ';'
		case "tab", "\t":
			sep = '\t'
		case "pipe", "|":
			sep = '|'
		default:
			sep = ','
		}

		reader := csv.NewReader(input)
		reader.Comma = sep

		records, err := reader.ReadAll()
		if err != nil {
			return fmt.Errorf("failed to parse CSV: %w", err)
		}

		if len(records) == 0 {
			return fmt.Errorf("empty CSV")
		}

		headers := records[0]
		var result []map[string]interface{}

		for _, record := range records[1:] {
			row := make(map[string]interface{})
			for i, value := range record {
				if i < len(headers) {
					// Try to convert to number
					if num, err := strconv.ParseFloat(value, 64); err == nil {
						row[headers[i]] = num
					} else if value == "true" || value == "TRUE" {
						row[headers[i]] = true
					} else if value == "false" || value == "FALSE" {
						row[headers[i]] = false
					} else {
						row[headers[i]] = value
					}
				}
			}
			result = append(result, row)
		}

		output, err := json.MarshalIndent(result, "", "  ")
		if err != nil {
			return fmt.Errorf("failed to convert to JSON: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

var dataJson2CsvCmd = &cobra.Command{
	Use:   "json2csv [file]",
	Short: "Convert JSON to CSV",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var input []byte
		var err error

		if len(args) > 0 {
			input, err = os.ReadFile(args[0])
			if err != nil {
				return fmt.Errorf("failed to read file: %w", err)
			}
		} else {
			input, err = io.ReadAll(os.Stdin)
			if err != nil {
				return fmt.Errorf("failed to read stdin: %w", err)
			}
		}

		delimiter, _ := cmd.Flags().GetString("delimiter")
		var sep rune
		switch delimiter {
		case "comma", ",":
			sep = ','
		case "semicolon", ";":
			sep = ';'
		case "tab", "\t":
			sep = '\t'
		case "pipe", "|":
			sep = '|'
		default:
			sep = ','
		}

		var data []map[string]interface{}
		if err := json.Unmarshal(input, &data); err != nil {
			return fmt.Errorf("invalid JSON: %w", err)
		}

		if len(data) == 0 {
			return fmt.Errorf("empty JSON array")
		}

		// Collect all unique keys
		keys := make(map[string]bool)
		for _, row := range data {
			for key := range row {
				keys[key] = true
			}
		}

		// Get ordered keys
		var headers []string
		for key := range keys {
			headers = append(headers, key)
		}

		writer := csv.NewWriter(os.Stdout)
		writer.Comma = sep

		// Write headers
		if err := writer.Write(headers); err != nil {
			return err
		}

		// Write data
		for _, row := range data {
			record := make([]string, len(headers))
			for i, key := range headers {
				if val, ok := row[key]; ok {
					record[i] = fmt.Sprintf("%v", val)
				}
			}
			if err := writer.Write(record); err != nil {
				return err
			}
		}

		writer.Flush()
		return writer.Error()
	},
}

var dataSizeCmd = &cobra.Command{
	Use:   "size [value]",
	Short: "Convert data size units",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.ToUpper(strings.TrimSpace(args[0]))

		// Parse value and unit
		var value float64
		var unit string

		for i := 0; i < len(input); i++ {
			if (input[i] < '0' || input[i] > '9') && input[i] != '.' {
				value, _ = strconv.ParseFloat(input[:i], 64)
				unit = input[i:]
				break
			}
		}

		if unit == "" {
			value, _ = strconv.ParseFloat(input, 64)
			unit = "B"
		}

		// Convert to bytes
		units := map[string]float64{
			"B":  1,
			"KB": 1024,
			"MB": 1024 * 1024,
			"GB": 1024 * 1024 * 1024,
			"TB": 1024 * 1024 * 1024 * 1024,
			"PB": 1024 * 1024 * 1024 * 1024 * 1024,
		}

		multiplier, ok := units[unit]
		if !ok {
			return fmt.Errorf("unknown unit: %s", unit)
		}

		bytes := value * multiplier

		fmt.Printf("%.2f %s =\n\n", value, unit)

		// Display in all units
		for _, u := range []string{"B", "KB", "MB", "GB", "TB", "PB"} {
			converted := bytes / units[u]
			fmt.Printf("  %10.4f %s\n", converted, u)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dataCmd)
	dataCmd.AddCommand(dataCsv2JsonCmd)
	dataCmd.AddCommand(dataJson2CsvCmd)
	dataCmd.AddCommand(dataSizeCmd)

	dataCsv2JsonCmd.Flags().StringP("delimiter", "d", "comma", "CSV delimiter (comma, semicolon, tab, pipe)")
	dataJson2CsvCmd.Flags().StringP("delimiter", "d", "comma", "CSV delimiter (comma, semicolon, tab, pipe)")
}
