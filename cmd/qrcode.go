package cmd

import (
	"fmt"

	"github.com/skip2/go-qrcode"
	"github.com/spf13/cobra"
)

var qrCmd = &cobra.Command{
	Use:   "qr [text]",
	Short: "Generate QR code",
	Long:  `Generate QR code as ASCII art or save to file.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]
		output, _ := cmd.Flags().GetString("output")
		size, _ := cmd.Flags().GetInt("size")
		ascii, _ := cmd.Flags().GetBool("ascii")

		if ascii || output == "" {
			// Generate ASCII QR code
			q, err := qrcode.New(text, qrcode.Medium)
			if err != nil {
				return fmt.Errorf("failed to generate QR code: %w", err)
			}
			fmt.Println(q.ToSmallString(false))
			return nil
		}

		// Save to file
		err := qrcode.WriteFile(text, qrcode.Medium, size, output)
		if err != nil {
			return fmt.Errorf("failed to save QR code: %w", err)
		}

		fmt.Printf("QR code saved to: %s\n", output)
		return nil
	},
}

var qrTerminalCmd = &cobra.Command{
	Use:   "terminal [text]",
	Short: "Generate QR code for terminal display",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		text := args[0]
		inverted, _ := cmd.Flags().GetBool("inverted")

		q, err := qrcode.New(text, qrcode.Medium)
		if err != nil {
			return fmt.Errorf("failed to generate QR code: %w", err)
		}

		fmt.Println(q.ToSmallString(inverted))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(qrCmd)
	rootCmd.AddCommand(qrTerminalCmd)

	qrCmd.Flags().StringP("output", "o", "", "Output file path (PNG)")
	qrCmd.Flags().IntP("size", "s", 256, "Image size in pixels")
	qrCmd.Flags().BoolP("ascii", "a", false, "Output as ASCII art")

	qrTerminalCmd.Flags().BoolP("inverted", "i", false, "Invert colors")
}
