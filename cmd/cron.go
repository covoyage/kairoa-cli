package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/spf13/cobra"
)

var cronCmd = &cobra.Command{
	Use:   "cron",
	Short: "Cron expression utilities",
	Long:  `Parse and explain cron expressions, show next execution times.`,
}

var cronParseCmd = &cobra.Command{
	Use:   "parse [expression]",
	Short: "Parse and explain cron expression",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expression := args[0]

		// Try to parse
		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(expression)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}

		// Print expression parts
		parts := strings.Fields(expression)
		if len(parts) == 5 {
			fmt.Println("Cron Expression Parts:")
			fmt.Println()
			fmt.Printf("  %-15s %s\n", "Minute:", parts[0])
			fmt.Printf("  %-15s %s\n", "Hour:", parts[1])
			fmt.Printf("  %-15s %s\n", "Day of Month:", parts[2])
			fmt.Printf("  %-15s %s\n", "Month:", parts[3])
			fmt.Printf("  %-15s %s\n", "Day of Week:", parts[4])
		}

		// Show next execution times
		count, _ := cmd.Flags().GetInt("count")
		fmt.Printf("\nNext %d execution times:\n\n", count)

		now := time.Now()
		for i := 0; i < count; i++ {
			next := schedule.Next(now)
			fmt.Printf("  %d. %s\n", i+1, next.Format("2006-01-02 15:04:05"))
			now = next
		}

		return nil
	},
}

var cronNextCmd = &cobra.Command{
	Use:   "next [expression]",
	Short: "Show next execution times",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expression := args[0]
		count, _ := cmd.Flags().GetInt("count")

		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		schedule, err := parser.Parse(expression)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}

		now := time.Now()
		for i := 0; i < count; i++ {
			next := schedule.Next(now)
			fmt.Println(next.Format("2006-01-02 15:04:05"))
			now = next
		}

		return nil
	},
}

var cronValidateCmd = &cobra.Command{
	Use:   "validate [expression]",
	Short: "Validate cron expression",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		expression := args[0]

		parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow)
		_, err := parser.Parse(expression)
		if err != nil {
			fmt.Printf("Invalid: %s\n", err)
			return nil
		}

		fmt.Println("Valid cron expression")
		return nil
	},
}

var cronExamplesCmd = &cobra.Command{
	Use:   "examples",
	Short: "Show common cron expression examples",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Common Cron Expression Examples:")
		fmt.Println()
		fmt.Println("  Expression          Description")
		fmt.Println("  ------------------  ----------------------------------------")
		fmt.Println("  * * * * *           Every minute")
		fmt.Println("  0 * * * *           Every hour")
		fmt.Println("  0 0 * * *           Every day at midnight")
		fmt.Println("  0 12 * * *          Every day at noon")
		fmt.Println("  0 0 * * 0           Every Sunday")
		fmt.Println("  0 0 1 * *           Every month")
		fmt.Println("  0 0 1 1 *           Every year")
		fmt.Println("  */5 * * * *         Every 5 minutes")
		fmt.Println("  0 9-17 * * 1-5      Every hour from 9-5 on weekdays")
		fmt.Println("  0 0 * * 1-5         Every weekday at midnight")
		fmt.Println()
		fmt.Println("Special strings:")
		fmt.Println("  @yearly, @annually  Run once a year (0 0 1 1 *)")
		fmt.Println("  @monthly            Run once a month (0 0 1 * *)")
		fmt.Println("  @weekly             Run once a week (0 0 * * 0)")
		fmt.Println("  @daily, @midnight   Run once a day (0 0 * * *)")
		fmt.Println("  @hourly             Run once an hour (0 * * * *)")
		return nil
	},
}

func getFieldDescription(value string, min, max int) string {
	switch value {
	case "*":
		return fmt.Sprintf("Every %d-%d", min, max)
	case "*/2":
		return "Every 2"
	case "*/5":
		return "Every 5"
	case "*/10":
		return "Every 10"
	case "*/15":
		return "Every 15"
	case "*/30":
		return "Every 30"
	default:
		if strings.Contains(value, ",") {
			return "Multiple values"
		}
		if strings.Contains(value, "-") {
			return "Range"
		}
		if strings.Contains(value, "/") {
			return "Step"
		}
		return "Specific value"
	}
}

func init() {
	rootCmd.AddCommand(cronCmd)
	cronCmd.AddCommand(cronParseCmd)
	cronCmd.AddCommand(cronNextCmd)
	cronCmd.AddCommand(cronValidateCmd)
	cronCmd.AddCommand(cronExamplesCmd)

	cronParseCmd.Flags().IntP("count", "n", 5, "Number of next execution times to show")
	cronNextCmd.Flags().IntP("count", "n", 10, "Number of execution times to show")
}
