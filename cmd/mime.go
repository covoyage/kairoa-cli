package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

type MimeInfo struct {
	MIME        string
	Extensions  []string
	Description string
}

var mimeDB = map[string]MimeInfo{
	// Text
	"text/plain":        {MIME: "text/plain", Extensions: []string{"txt"}, Description: "Plain text"},
	"text/html":         {MIME: "text/html", Extensions: []string{"html", "htm"}, Description: "HTML document"},
	"text/css":          {MIME: "text/css", Extensions: []string{"css"}, Description: "CSS stylesheet"},
	"text/javascript":   {MIME: "text/javascript", Extensions: []string{"js", "mjs"}, Description: "JavaScript"},
	"application/json":  {MIME: "application/json", Extensions: []string{"json"}, Description: "JSON data"},
	"text/xml":          {MIME: "text/xml", Extensions: []string{"xml"}, Description: "XML document"},
	"text/csv":          {MIME: "text/csv", Extensions: []string{"csv"}, Description: "CSV data"},
	"text/markdown":     {MIME: "text/markdown", Extensions: []string{"md", "markdown"}, Description: "Markdown"},
	"text/yaml":         {MIME: "text/yaml", Extensions: []string{"yaml", "yml"}, Description: "YAML"},

	// Images
	"image/jpeg":        {MIME: "image/jpeg", Extensions: []string{"jpg", "jpeg"}, Description: "JPEG image"},
	"image/png":         {MIME: "image/png", Extensions: []string{"png"}, Description: "PNG image"},
	"image/gif":         {MIME: "image/gif", Extensions: []string{"gif"}, Description: "GIF image"},
	"image/webp":        {MIME: "image/webp", Extensions: []string{"webp"}, Description: "WebP image"},
	"image/svg+xml":     {MIME: "image/svg+xml", Extensions: []string{"svg"}, Description: "SVG image"},
	"image/bmp":         {MIME: "image/bmp", Extensions: []string{"bmp"}, Description: "BMP image"},
	"image/tiff":        {MIME: "image/tiff", Extensions: []string{"tiff", "tif"}, Description: "TIFF image"},
	"image/x-icon":      {MIME: "image/x-icon", Extensions: []string{"ico"}, Description: "Icon"},

	// Audio
	"audio/mpeg":        {MIME: "audio/mpeg", Extensions: []string{"mp3"}, Description: "MP3 audio"},
	"audio/wav":         {MIME: "audio/wav", Extensions: []string{"wav"}, Description: "WAV audio"},
	"audio/ogg":         {MIME: "audio/ogg", Extensions: []string{"ogg"}, Description: "OGG audio"},
	"audio/aac":         {MIME: "audio/aac", Extensions: []string{"aac"}, Description: "AAC audio"},
	"audio/flac":        {MIME: "audio/flac", Extensions: []string{"flac"}, Description: "FLAC audio"},

	// Video
	"video/mp4":         {MIME: "video/mp4", Extensions: []string{"mp4"}, Description: "MP4 video"},
	"video/webm":        {MIME: "video/webm", Extensions: []string{"webm"}, Description: "WebM video"},
	"video/ogg":         {MIME: "video/ogg", Extensions: []string{"ogv"}, Description: "OGG video"},
	"video/x-msvideo":   {MIME: "video/x-msvideo", Extensions: []string{"avi"}, Description: "AVI video"},
	"video/quicktime":   {MIME: "video/quicktime", Extensions: []string{"mov"}, Description: "QuickTime video"},

	// Application
	"application/pdf":           {MIME: "application/pdf", Extensions: []string{"pdf"}, Description: "PDF document"},
	"application/zip":           {MIME: "application/zip", Extensions: []string{"zip"}, Description: "ZIP archive"},
	"application/gzip":          {MIME: "application/gzip", Extensions: []string{"gz"}, Description: "GZIP archive"},
	"application/x-tar":         {MIME: "application/x-tar", Extensions: []string{"tar"}, Description: "TAR archive"},
	"application/x-7z-compressed": {MIME: "application/x-7z-compressed", Extensions: []string{"7z"}, Description: "7-Zip archive"},
	"application/x-rar-compressed": {MIME: "application/x-rar-compressed", Extensions: []string{"rar"}, Description: "RAR archive"},
	"application/octet-stream":  {MIME: "application/octet-stream", Extensions: []string{"bin"}, Description: "Binary data"},
	"application/msword":        {MIME: "application/msword", Extensions: []string{"doc"}, Description: "Word document"},
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": {MIME: "application/vnd.openxmlformats-officedocument.wordprocessingml.document", Extensions: []string{"docx"}, Description: "Word document (OpenXML)"},
	"application/vnd.ms-excel":  {MIME: "application/vnd.ms-excel", Extensions: []string{"xls"}, Description: "Excel spreadsheet"},
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": {MIME: "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", Extensions: []string{"xlsx"}, Description: "Excel spreadsheet (OpenXML)"},
}

var mimeCmd = &cobra.Command{
	Use:   "mime",
	Short: "MIME type lookup",
	Long:  `Look up MIME types by file extension or vice versa.`,
}

var mimeLookupCmd = &cobra.Command{
	Use:   "lookup [file-or-mime]",
	Short: "Look up MIME type",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := strings.ToLower(args[0])

		// Check if it's a file extension
		if strings.HasPrefix(input, ".") {
			input = input[1:]
		}

		// Try to find by extension
		for _, info := range mimeDB {
			for _, ext := range info.Extensions {
				if ext == input {
					printMimeInfo(info)
					return nil
				}
			}
		}

		// Try to find by MIME type
		if info, ok := mimeDB[input]; ok {
			printMimeInfo(info)
			return nil
		}

		return fmt.Errorf("MIME type not found for: %s", args[0])
	},
}

var mimeFromFileCmd = &cobra.Command{
	Use:   "file [filename]",
	Short: "Get MIME type from filename",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ext := strings.ToLower(filepath.Ext(args[0]))
		if ext == "" {
			return fmt.Errorf("no extension found in filename: %s", args[0])
		}
		ext = ext[1:] // Remove leading dot

		for _, info := range mimeDB {
			for _, e := range info.Extensions {
				if e == ext {
					fmt.Println(info.MIME)
					return nil
				}
			}
		}

		return fmt.Errorf("unknown file extension: %s", ext)
	},
}

var mimeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all MIME types",
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")

		fmt.Println("MIME Types:")
		fmt.Println()

		for _, info := range mimeDB {
			if category != "" && !strings.HasPrefix(info.MIME, category) {
				continue
			}
			fmt.Printf("%-50s %s\n", info.MIME, strings.Join(info.Extensions, ", "))
		}

		return nil
	},
}

var mimeExtCmd = &cobra.Command{
	Use:   "ext [extension]",
	Short: "Get MIME type for extension",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ext := strings.ToLower(args[0])
		if strings.HasPrefix(ext, ".") {
			ext = ext[1:]
		}

		for _, info := range mimeDB {
			for _, e := range info.Extensions {
				if e == ext {
					fmt.Println(info.MIME)
					return nil
				}
			}
		}

		return fmt.Errorf("unknown extension: %s", ext)
	},
}

func printMimeInfo(info MimeInfo) {
	fmt.Printf("MIME Type: %s\n", info.MIME)
	fmt.Printf("Extensions: %s\n", strings.Join(info.Extensions, ", "))
	fmt.Printf("Description: %s\n", info.Description)
}

func init() {
	rootCmd.AddCommand(mimeCmd)
	mimeCmd.AddCommand(mimeLookupCmd)
	mimeCmd.AddCommand(mimeFromFileCmd)
	mimeCmd.AddCommand(mimeListCmd)
	mimeCmd.AddCommand(mimeExtCmd)

	mimeListCmd.Flags().StringP("category", "c", "", "Filter by category (text/, image/, audio/, video/, application/)")
}
