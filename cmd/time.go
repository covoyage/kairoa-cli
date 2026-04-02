package cmd

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var timeCmd = &cobra.Command{
	Use:   "time",
	Short: "Time utilities",
	Long:  `Convert between timestamps and dates, get current time, etc.`,
}

var timeNowCmd = &cobra.Command{
	Use:   "now",
	Short: "Get current time in various formats",
	RunE: func(cmd *cobra.Command, args []string) error {
		unix, _ := cmd.Flags().GetBool("unix")
		unixMs, _ := cmd.Flags().GetBool("unix-ms")
		iso, _ := cmd.Flags().GetBool("iso")
		rfc, _ := cmd.Flags().GetBool("rfc")

		now := time.Now()

		if !unix && !unixMs && !iso && !rfc {
			// Default: show all formats
			fmt.Printf("Unix (seconds): %d\n", now.Unix())
			fmt.Printf("Unix (ms):      %d\n", now.UnixMilli())
			fmt.Printf("ISO 8601:       %s\n", now.Format(time.RFC3339))
			fmt.Printf("RFC 1123:       %s\n", now.Format(time.RFC1123))
			fmt.Printf("Local:          %s\n", now.Format("2006-01-02 15:04:05"))
		} else {
			if unix {
				fmt.Println(now.Unix())
			}
			if unixMs {
				fmt.Println(now.UnixMilli())
			}
			if iso {
				fmt.Println(now.Format(time.RFC3339))
			}
			if rfc {
				fmt.Println(now.Format(time.RFC1123))
			}
		}
		return nil
	},
}

var timeConvertCmd = &cobra.Command{
	Use:   "convert [timestamp]",
	Short: "Convert timestamp to date or vice versa",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.TrimSpace(args[0])
		tz, _ := cmd.Flags().GetString("timezone")

		loc, err := time.LoadLocation(tz)
		if err != nil {
			return fmt.Errorf("invalid timezone: %w", err)
		}

		// Try to parse as Unix timestamp
		if ts, err := strconv.ParseInt(input, 10, 64); err == nil {
			var t time.Time
			if ts > 1e12 {
				// Milliseconds
				t = time.UnixMilli(ts).In(loc)
			} else {
				// Seconds
				t = time.Unix(ts, 0).In(loc)
			}
			fmt.Printf("Local:    %s\n", t.Format("2006-01-02 15:04:05"))
			fmt.Printf("UTC:      %s\n", t.UTC().Format("2006-01-02 15:04:05"))
			fmt.Printf("ISO 8601: %s\n", t.Format(time.RFC3339))
			return nil
		}

		// Try to parse as date string
		formats := []string{
			time.RFC3339,
			"2006-01-02 15:04:05",
			"2006-01-02",
			"2006/01/02 15:04:05",
			"2006/01/02",
		}

		for _, format := range formats {
			if t, err := time.ParseInLocation(format, input, loc); err == nil {
				fmt.Printf("Unix (seconds): %d\n", t.Unix())
				fmt.Printf("Unix (ms):      %d\n", t.UnixMilli())
				fmt.Printf("ISO 8601:       %s\n", t.Format(time.RFC3339))
				return nil
			}
		}

		return fmt.Errorf("unable to parse input: %s", input)
	},
}

func init() {
	rootCmd.AddCommand(timeCmd)
	timeCmd.AddCommand(timeNowCmd)
	timeCmd.AddCommand(timeConvertCmd)

	timeNowCmd.Flags().BoolP("unix", "u", false, "Show Unix timestamp (seconds)")
	timeNowCmd.Flags().BoolP("unix-ms", "m", false, "Show Unix timestamp (milliseconds)")
	timeNowCmd.Flags().BoolP("iso", "i", false, "Show ISO 8601 format")
	timeNowCmd.Flags().BoolP("rfc", "r", false, "Show RFC 1123 format")

	timeConvertCmd.Flags().StringP("timezone", "z", "Local", "Timezone (e.g., UTC, America/New_York, Asia/Shanghai)")
}
