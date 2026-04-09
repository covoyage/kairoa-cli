package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var httpCmd = &cobra.Command{
	Use:   "http",
	Short: "HTTP client utilities",
	Long:  `Send HTTP requests and inspect responses.`,
}

var httpRequestCmd = &cobra.Command{
	Use:     "request [url]",
	Short:   "Send HTTP request",
	Aliases: []string{"req"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		method, _ := cmd.Flags().GetString("method")
		body, _ := cmd.Flags().GetString("body")
		headers, _ := cmd.Flags().GetStringSlice("header")
		showHeaders, _ := cmd.Flags().GetBool("show-headers")

		// Ensure URL has scheme
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		// Create request
		var bodyReader io.Reader
		if body != "" {
			bodyReader = strings.NewReader(body)
		}

		req, err := http.NewRequest(method, url, bodyReader)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		// Add headers
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				req.Header.Set(strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]))
			}
		}

		// Send request
		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		start := time.Now()
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		duration := time.Since(start)

		// Print response info
		fmt.Printf("%s %s\n", color.GreenString("HTTP"), resp.Status)
		fmt.Printf("Time: %v\n", duration)

		if showHeaders {
			fmt.Println("\nHeaders:")
			for name, values := range resp.Header {
				for _, value := range values {
					fmt.Printf("  %s: %s\n", name, value)
				}
			}
		}

		// Read body
		respBody, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if len(respBody) > 0 {
			fmt.Println("\nBody:")

			// Try to format JSON
			var jsonData interface{}
			if err := json.Unmarshal(respBody, &jsonData); err == nil {
				if formatted, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
					fmt.Println(string(formatted))
				} else {
					fmt.Println(string(respBody))
				}
			} else {
				fmt.Println(string(respBody))
			}
		}

		return nil
	},
}

var httpGetCmd = &cobra.Command{
	Use:   "get [url]",
	Short: "Send HTTP GET request",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		showHeaders, _ := cmd.Flags().GetBool("show-headers")

		// Ensure URL has scheme
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		start := time.Now()
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()
		duration := time.Since(start)

		fmt.Printf("%s %s\n", color.GreenString("HTTP"), resp.Status)
		fmt.Printf("Time: %v\n", duration)

		if showHeaders {
			fmt.Println("\nHeaders:")
			for name, values := range resp.Header {
				for _, value := range values {
					fmt.Printf("  %s: %s\n", name, value)
				}
			}
		}

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if len(body) > 0 {
			fmt.Println("\nBody:")

			// Try to format JSON
			var jsonData interface{}
			if err := json.Unmarshal(body, &jsonData); err == nil {
				if formatted, err := json.MarshalIndent(jsonData, "", "  "); err == nil {
					fmt.Println(string(formatted))
				} else {
					fmt.Println(string(body))
				}
			} else {
				fmt.Println(string(body))
			}
		}

		return nil
	},
}

var httpPostCmd = &cobra.Command{
	Use:   "post [url] [data]",
	Short: "Send HTTP POST request",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		data := args[1]

		// Ensure URL has scheme
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Post(url, "application/json", strings.NewReader(data))
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		fmt.Printf("%s %s\n", color.GreenString("HTTP"), resp.Status)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response body: %w", err)
		}

		if len(body) > 0 {
			fmt.Println("\nBody:")
			fmt.Println(string(body))
		}

		return nil
	},
}

var httpHeadersCmd = &cobra.Command{
	Use:   "headers [url]",
	Short: "Show HTTP headers only",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		// Ensure URL has scheme
		if !strings.HasPrefix(url, "http://") && !strings.HasPrefix(url, "https://") {
			url = "https://" + url
		}

		client := &http.Client{
			Timeout: 30 * time.Second,
		}

		resp, err := client.Head(url)
		if err != nil {
			// Try GET if HEAD fails
			resp, err = client.Get(url)
			if err != nil {
				return fmt.Errorf("request failed: %w", err)
			}
			defer resp.Body.Close()
			io.Copy(io.Discard, resp.Body)
		} else {
			defer resp.Body.Close()
		}

		fmt.Printf("%s %s\n\n", color.GreenString("HTTP"), resp.Status)

		fmt.Println("Request Headers:")
		fmt.Println("  (Not available for HEAD request)")

		fmt.Println("\nResponse Headers:")
		for name, values := range resp.Header {
			for _, value := range values {
				fmt.Printf("  %s: %s\n", name, value)
			}
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpCmd)
	httpCmd.AddCommand(httpRequestCmd)
	httpCmd.AddCommand(httpGetCmd)
	httpCmd.AddCommand(httpPostCmd)
	httpCmd.AddCommand(httpHeadersCmd)

	httpRequestCmd.Flags().StringP("method", "X", "GET", "HTTP method")
	httpRequestCmd.Flags().StringP("body", "d", "", "Request body")
	httpRequestCmd.Flags().StringSliceP("header", "H", []string{}, "Request headers (format: 'Key: Value')")
	httpRequestCmd.Flags().BoolP("show-headers", "i", false, "Show response headers")

	httpGetCmd.Flags().BoolP("show-headers", "i", false, "Show response headers")
}
