package cmd

import (
	"encoding/base64"
	"fmt"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var imageCmd = &cobra.Command{
	Use:   "image",
	Short: "Image processing utilities",
	Long:  `Process and convert images (resize, rotate, convert format, base64).`,
}

var imageInfoCmd = &cobra.Command{
	Use:   "info [file]",
	Short: "Show image information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		img, format, err := loadImage(file)
		if err != nil {
			return err
		}

		bounds := img.Bounds()
		width := bounds.Dx()
		height := bounds.Dy()

		// Get file size
		info, err := os.Stat(file)
		if err != nil {
			return err
		}

		fmt.Printf("File: %s\n", file)
		fmt.Printf("Format: %s\n", strings.ToUpper(format))
		fmt.Printf("Dimensions: %dx%d\n", width, height)
		fmt.Printf("Size: %d bytes (%.2f KB)\n", info.Size(), float64(info.Size())/1024)

		return nil
	},
}

var imageConvertCmd = &cobra.Command{
	Use:   "convert [input] [output]",
	Short: "Convert image format",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]
		output := args[1]

		img, _, err := loadImage(input)
		if err != nil {
			return err
		}

		return saveImage(output, img)
	},
}

var imageBase64Cmd = &cobra.Command{
	Use:   "base64 [file]",
	Short: "Convert image to base64",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		file := args[0]

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		encoded := base64.StdEncoding.EncodeToString(data)
		
		// Get MIME type
		ext := strings.ToLower(filepath.Ext(file))
		var mimeType string
		switch ext {
		case ".png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		case ".gif":
			mimeType = "image/gif"
		case ".webp":
			mimeType = "image/webp"
		default:
			mimeType = "image/png"
		}

		fmt.Printf("data:%s;base64,%s\n", mimeType, encoded)
		return nil
	},
}

var imageBase64DecodeCmd = &cobra.Command{
	Use:   "base64-decode [base64] [output]",
	Short: "Decode base64 to image",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		encoded := args[0]
		output := args[1]

		// Remove data URI prefix if present
		if idx := strings.Index(encoded, ","); idx != -1 {
			encoded = encoded[idx+1:]
		}

		data, err := base64.StdEncoding.DecodeString(encoded)
		if err != nil {
			return fmt.Errorf("invalid base64: %w", err)
		}

		return os.WriteFile(output, data, 0644)
	},
}

func loadImage(filename string) (image.Image, string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	return image.Decode(file)
}

func saveImage(filename string, img image.Image) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".png":
		return png.Encode(file, img)
	case ".jpg", ".jpeg":
		return jpeg.Encode(file, img, &jpeg.Options{Quality: 85})
	case ".gif":
		return gif.Encode(file, img, nil)
	default:
		return fmt.Errorf("unsupported format: %s", ext)
	}
}

func init() {
	rootCmd.AddCommand(imageCmd)
	imageCmd.AddCommand(imageInfoCmd)
	imageCmd.AddCommand(imageConvertCmd)
	imageCmd.AddCommand(imageBase64Cmd)
	imageCmd.AddCommand(imageBase64DecodeCmd)
}
