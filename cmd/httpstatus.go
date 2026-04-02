package cmd

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

type StatusCodeInfo struct {
	Code        int
	Name        string
	Description string
	Category    string
}

var httpStatusDB = map[int]StatusCodeInfo{
	// 1xx Informational
	100: {100, "Continue", "The server has received the request headers", "informational"},
	101: {101, "Switching Protocols", "Server is switching to different protocol", "informational"},
	102: {102, "Processing", "Server is processing the request", "informational"},
	103: {103, "Early Hints", "Preliminary hints for the response", "informational"},

	// 2xx Success
	200: {200, "OK", "Request succeeded", "success"},
	201: {201, "Created", "Resource created successfully", "success"},
	202: {202, "Accepted", "Request accepted for processing", "success"},
	204: {204, "No Content", "Request succeeded with no content to return", "success"},
	206: {206, "Partial Content", "Partial content returned", "success"},

	// 3xx Redirection
	300: {300, "Multiple Choices", "Multiple options for the resource", "redirection"},
	301: {301, "Moved Permanently", "Resource moved permanently", "redirection"},
	302: {302, "Found", "Resource found at different URI", "redirection"},
	304: {304, "Not Modified", "Resource not modified since last request", "redirection"},
	307: {307, "Temporary Redirect", "Temporary redirect to different URI", "redirection"},
	308: {308, "Permanent Redirect", "Permanent redirect to different URI", "redirection"},

	// 4xx Client Error
	400: {400, "Bad Request", "Request syntax error", "client-error"},
	401: {401, "Unauthorized", "Authentication required", "client-error"},
	403: {403, "Forbidden", "Server refusing action", "client-error"},
	404: {404, "Not Found", "Resource not found", "client-error"},
	405: {405, "Method Not Allowed", "HTTP method not allowed", "client-error"},
	408: {408, "Request Timeout", "Request timeout", "client-error"},
	409: {409, "Conflict", "Request conflicts with current state", "client-error"},
	410: {410, "Gone", "Resource no longer available", "client-error"},
	429: {429, "Too Many Requests", "Rate limit exceeded", "client-error"},

	// 5xx Server Error
	500: {500, "Internal Server Error", "Server encountered an error", "server-error"},
	501: {501, "Not Implemented", "Server does not support functionality", "server-error"},
	502: {502, "Bad Gateway", "Invalid response from upstream", "server-error"},
	503: {503, "Service Unavailable", "Server temporarily unavailable", "server-error"},
	504: {504, "Gateway Timeout", "Upstream server timeout", "server-error"},
}

var httpstatusCmd = &cobra.Command{
	Use:   "httpstatus",
	Short: "HTTP status code reference",
	Long:  `Look up HTTP status codes and their meanings.`,
}

var httpstatusLookupCmd = &cobra.Command{
	Use:   "lookup [code]",
	Short: "Look up HTTP status code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		code, err := strconv.Atoi(args[0])
		if err != nil {
			return fmt.Errorf("invalid status code: %s", args[0])
		}

		info, ok := httpStatusDB[code]
		if !ok {
			return fmt.Errorf("unknown status code: %d", code)
		}

		fmt.Printf("Code: %d\n", info.Code)
		fmt.Printf("Name: %s\n", info.Name)
		fmt.Printf("Description: %s\n", info.Description)
		fmt.Printf("Category: %s\n", info.Category)

		return nil
	},
}

var httpstatusListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all HTTP status codes",
	RunE: func(cmd *cobra.Command, args []string) error {
		category, _ := cmd.Flags().GetString("category")

		fmt.Println("HTTP Status Codes:")
		fmt.Println()

		categories := map[string]string{
			"informational": "1xx Informational",
			"success":       "2xx Success",
			"redirection":   "3xx Redirection",
			"client-error":  "4xx Client Error",
			"server-error":  "5xx Server Error",
		}

		for catKey, catName := range categories {
			if category != "" && category != catKey {
				continue
			}

			fmt.Printf("%s:\n", catName)
			for _, info := range httpStatusDB {
				if info.Category == catKey {
					fmt.Printf("  %d %s - %s\n", info.Code, info.Name, info.Description)
				}
			}
			fmt.Println()
		}

		return nil
	},
}

var httpstatusSearchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Search HTTP status codes",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := strings.ToLower(args[0])

		fmt.Printf("Search results for '%s':\n\n", args[0])

		found := false
		for _, info := range httpStatusDB {
			if strings.Contains(strings.ToLower(info.Name), query) ||
				strings.Contains(strings.ToLower(info.Description), query) ||
				strconv.Itoa(info.Code) == query {
				fmt.Printf("%d %s - %s\n", info.Code, info.Name, info.Description)
				found = true
			}
		}

		if !found {
			fmt.Println("No results found.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(httpstatusCmd)
	httpstatusCmd.AddCommand(httpstatusLookupCmd)
	httpstatusCmd.AddCommand(httpstatusListCmd)
	httpstatusCmd.AddCommand(httpstatusSearchCmd)

	httpstatusListCmd.Flags().StringP("category", "c", "", "Filter by category (informational, success, redirection, client-error, server-error)")
}
