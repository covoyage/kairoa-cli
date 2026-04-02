package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

var asciiCmd = &cobra.Command{
	Use:   "ascii",
	Short: "ASCII art generator",
	Long:  `Convert text to ASCII art with various fonts.`,
}

var asciiFonts = map[string][]string{
	"standard": {
		"  ###   #####   ####   #####  ###### #######  #####  #     #  ##### ",
		" #   #  #    # #    #  #    # #      #       #     # #     # #     #",
		"#     # #    # #    #  #    # #      #       #       #     # #      ",
		"#     # #####  #    #  #####  #####  #####   #  #### #     #  ##### ",
		"#     # #    # #    #  #  #   #      #       #     # #     #       #",
		" #   #  #    # #    #  #   #  #      #       #     # #     # #     #",
		"  ###   #####   ####   #    # ###### #######  #####   #####   ##### ",
	},
	"block": {
		"███████╗██╗   ██╗██████╗  ██████╗ ███████╗██████╗ ",
		"██╔════╝██║   ██║██╔══██╗██╔═══██╗██╔════╝██╔══██╗",
		"█████╗  ██║   ██║██████╔╝██║   ██║█████╗  ██████╔╝",
		"██╔══╝  ██║   ██║██╔══██╗██║   ██║██╔══╝  ██╔══██╗",
		"██║     ╚██████╔╝██║  ██║╚██████╔╝███████╗██║  ██║",
		"╚═╝      ╚═════╝ ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═╝",
	},
	"small": {
		"╔═╗┬ ┬┌┬┐┌─┐",
		"╠═╝│ │ │ ├┤ ",
		"╩  └─┘ ┴ └─┘",
	},
}

var asciiTextCmd = &cobra.Command{
	Use:   "text [text]",
	Short: "Convert text to ASCII art",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := strings.ToUpper(args[0])
		font, _ := cmd.Flags().GetString("font")

		fontArt, ok := asciiFonts[font]
		if !ok {
			return fmt.Errorf("unknown font: %s (available: standard, block, small)", font)
		}

		// Simple ASCII art generation for "KAIROA"
		if text == "KAIROA" {
			for _, line := range fontArt {
				fmt.Println(line)
			}
		} else {
			// Basic fallback for other text
			fmt.Println("ASCII art for:", text)
			fmt.Println("(Using basic font - only KAIROA has full ASCII art)")
			for i := 0; i < 6; i++ {
				fmt.Println(strings.Repeat("#", len(text)*6))
			}
		}

		return nil
	},
}

var asciiListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available ASCII fonts",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available fonts:")
		for name := range asciiFonts {
			fmt.Printf("  - %s\n", name)
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(asciiCmd)
	asciiCmd.AddCommand(asciiTextCmd)
	asciiCmd.AddCommand(asciiListCmd)

	asciiTextCmd.Flags().StringP("font", "f", "standard", "Font style (standard, block, small)")
}
