package cmd

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

var uuidCmd = &cobra.Command{
	Use:   "uuid",
	Short: "Generate UUIDs",
	Long:  `Generate UUIDs (Universally Unique Identifiers) with various options.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		noHyphens, _ := cmd.Flags().GetBool("no-hyphens")
		version, _ := cmd.Flags().GetString("version")

		for i := 0; i < count; i++ {
			var id uuid.UUID
			switch version {
			case "v1":
				id = uuid.Must(uuid.NewUUID())
			case "v4":
				id = uuid.New()
			default:
				id = uuid.New()
			}

			result := id.String()
			if noHyphens {
				result = strings.ReplaceAll(result, "-", "")
			}
			fmt.Println(result)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(uuidCmd)
	uuidCmd.Flags().IntP("count", "c", 1, "Number of UUIDs to generate")
	uuidCmd.Flags().BoolP("no-hyphens", "n", false, "Remove hyphens from output")
	uuidCmd.Flags().StringP("version", "v", "v4", "UUID version (v1, v4)")
}
