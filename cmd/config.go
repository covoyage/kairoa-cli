package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/beevik/etree"
	"github.com/pelletier/go-toml/v2"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Config file format conversion",
	Long:  `Convert between JSON, YAML, TOML, and XML formats.`,
}

var configConvertCmd = &cobra.Command{
	Use:   "convert [file]",
	Short: "Convert config file to another format",
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

		from, _ := cmd.Flags().GetString("from")
		to, _ := cmd.Flags().GetString("to")

		if from == "" {
			// Auto-detect format
			from = detectFormat(input)
		}

		// Parse input
		var data interface{}
		switch from {
		case "json":
			if err := json.Unmarshal(input, &data); err != nil {
				return fmt.Errorf("failed to parse JSON: %w", err)
			}
		case "yaml", "yml":
			if err := yaml.Unmarshal(input, &data); err != nil {
				return fmt.Errorf("failed to parse YAML: %w", err)
			}
		case "toml":
			if err := toml.Unmarshal(input, &data); err != nil {
				return fmt.Errorf("failed to parse TOML: %w", err)
			}
		case "xml":
			doc := etree.NewDocument()
			if err := doc.ReadFromBytes(input); err != nil {
				return fmt.Errorf("failed to parse XML: %w", err)
			}
			data = xmlToMap(doc.Root())
		default:
			return fmt.Errorf("unsupported input format: %s", from)
		}

		// Convert to output format
		var output []byte
		switch to {
		case "json":
			output, err = json.MarshalIndent(data, "", "  ")
		case "yaml", "yml":
			output, err = yaml.Marshal(data)
		case "toml":
			output, err = toml.Marshal(data)
		default:
			return fmt.Errorf("unsupported output format: %s", to)
		}

		if err != nil {
			return fmt.Errorf("failed to convert: %w", err)
		}

		fmt.Println(string(output))
		return nil
	},
}

func detectFormat(data []byte) string {
	trimmed := strings.TrimSpace(string(data))

	if strings.HasPrefix(trimmed, "<") {
		return "xml"
	}
	if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
		return "json"
	}
	if strings.Contains(trimmed, "=") && !strings.Contains(trimmed, ":") {
		return "toml"
	}
	return "yaml"
}

func xmlToMap(element *etree.Element) map[string]interface{} {
	result := make(map[string]interface{})

	// Add attributes
	for _, attr := range element.Attr {
		result["@"+attr.Key] = attr.Value
	}

	// Add child elements
	childMap := make(map[string][]interface{})
	for _, child := range element.ChildElements() {
		childData := xmlToMap(child)
		childMap[child.Tag] = append(childMap[child.Tag], childData)
	}

	for tag, children := range childMap {
		if len(children) == 1 {
			result[tag] = children[0]
		} else {
			result[tag] = children
		}
	}

	// Add text content
	if text := strings.TrimSpace(element.Text()); text != "" {
		if len(result) == 0 {
			return map[string]interface{}{"#text": text}
		}
		result["#text"] = text
	}

	return result
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configConvertCmd)

	configConvertCmd.Flags().StringP("from", "f", "", "Input format (json, yaml, toml, xml)")
	configConvertCmd.Flags().StringP("to", "t", "json", "Output format (json, yaml, toml)")
}
