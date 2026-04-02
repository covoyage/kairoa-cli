package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var romanCmd = &cobra.Command{
	Use:   "roman",
	Short: "Roman numeral converter",
	Long:  `Convert between Roman numerals and Arabic numbers.`,
}

var romanToArabicMap = map[byte]int{
	'I': 1,
	'V': 5,
	'X': 10,
	'L': 50,
	'C': 100,
	'D': 500,
	'M': 1000,
}

var arabicToRomanMap = []struct {
	value int
	symbol string
}{
	{1000, "M"}, {900, "CM"}, {500, "D"}, {400, "CD"},
	{100, "C"}, {90, "XC"}, {50, "L"}, {40, "XL"},
	{10, "X"}, {9, "IX"}, {5, "V"}, {4, "IV"}, {1, "I"},
}

var romanToCmd = &cobra.Command{
	Use:   "to-arabic [roman]",
	Short: "Convert Roman numeral to Arabic number",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		roman := strings.ToUpper(args[0])
		result := romanToArabic(roman)
		fmt.Printf("%s = %d\n", roman, result)
		return nil
	},
}

var romanFromCmd = &cobra.Command{
	Use:   "from-arabic [number]",
	Short: "Convert Arabic number to Roman numeral",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var num int
		fmt.Sscanf(args[0], "%d", &num)
		result := arabicToRoman(num)
		fmt.Printf("%d = %s\n", num, result)
		return nil
	},
}

var romanConvertCmd = &cobra.Command{
	Use:   "convert [input]",
	Short: "Auto-detect and convert",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		// Try to parse as number first
		var num int
		if _, err := fmt.Sscanf(input, "%d", &num); err == nil && num > 0 {
			fmt.Printf("%d = %s\n", num, arabicToRoman(num))
		} else {
			// Treat as Roman numeral
			roman := strings.ToUpper(input)
			result := romanToArabic(roman)
			fmt.Printf("%s = %d\n", roman, result)
		}

		return nil
	},
}

func romanToArabic(roman string) int {
	result := 0
	prevValue := 0

	for i := len(roman) - 1; i >= 0; i-- {
		value := romanToArabicMap[roman[i]]
		if value < prevValue {
			result -= value
		} else {
			result += value
		}
		prevValue = value
	}

	return result
}

func arabicToRoman(num int) string {
	if num <= 0 || num > 3999 {
		return "Invalid (must be 1-3999)"
	}

	var result strings.Builder
	for _, pair := range arabicToRomanMap {
		for num >= pair.value {
			result.WriteString(pair.symbol)
			num -= pair.value
		}
	}

	return result.String()
}

func init() {
	rootCmd.AddCommand(romanCmd)
	romanCmd.AddCommand(romanToCmd)
	romanCmd.AddCommand(romanFromCmd)
	romanCmd.AddCommand(romanConvertCmd)
}
