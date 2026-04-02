package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var ipCmd = &cobra.Command{
	Use:   "ip",
	Short: "IP lookup utilities",
	Long:  `Look up IP address information including geolocation.`,
}

var ipLookupCmd = &cobra.Command{
	Use:   "lookup [ip|domain]",
	Short: "Look up IP information",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		var query string
		if len(args) > 0 {
			query = args[0]
		} else {
			// Get own IP
			resp, err := http.Get("https://api.ipify.org?format=json")
			if err != nil {
				return fmt.Errorf("failed to get IP: %w", err)
			}
			defer resp.Body.Close()

			var result struct {
				IP string `json:"ip"`
			}
			if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
				return fmt.Errorf("failed to parse response: %w", err)
			}
			query = result.IP
		}

		// Resolve IP if domain
		ip := net.ParseIP(query)
		if ip == nil {
			ips, err := net.LookupIP(query)
			if err != nil {
				return fmt.Errorf("failed to resolve: %w", err)
			}
			if len(ips) > 0 {
				ip = ips[0]
			}
		}

		fmt.Printf("IP Address: %s\n\n", color.GreenString(ip.String()))

		// Get geolocation info
		url := fmt.Sprintf("http://ip-api.com/json/%s", ip.String())
		client := &http.Client{Timeout: 10 * time.Second}
		resp, err := client.Get(url)
		if err != nil {
			return fmt.Errorf("failed to get geolocation: %w", err)
		}
		defer resp.Body.Close()

		var geo struct {
			Status      string  `json:"status"`
			Country     string  `json:"country"`
			CountryCode string  `json:"countryCode"`
			Region      string  `json:"region"`
			RegionName  string  `json:"regionName"`
			City        string  `json:"city"`
			Zip         string  `json:"zip"`
			Lat         float64 `json:"lat"`
			Lon         float64 `json:"lon"`
			Timezone    string  `json:"timezone"`
			ISP         string  `json:"isp"`
			Org         string  `json:"org"`
			AS          string  `json:"as"`
			Query       string  `json:"query"`
		}

		if err := json.NewDecoder(resp.Body).Decode(&geo); err != nil {
			return fmt.Errorf("failed to parse geolocation: %w", err)
		}

		if geo.Status != "success" {
			return fmt.Errorf("geolocation lookup failed")
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.Append([]string{"Country", fmt.Sprintf("%s (%s)", geo.Country, geo.CountryCode)})
		table.Append([]string{"Region", geo.RegionName})
		table.Append([]string{"City", geo.City})
		table.Append([]string{"ZIP", geo.Zip})
		table.Append([]string{"Coordinates", fmt.Sprintf("%.4f, %.4f", geo.Lat, geo.Lon)})
		table.Append([]string{"Timezone", geo.Timezone})
		table.Append([]string{"ISP", geo.ISP})
		table.Append([]string{"Organization", geo.Org})
		table.Append([]string{"AS", geo.AS})
		table.Render()

		return nil
	},
}

var ipLocalCmd = &cobra.Command{
	Use:   "local",
	Short: "Show local IP addresses",
	RunE: func(cmd *cobra.Command, args []string) error {
		addrs, err := net.InterfaceAddrs()
		if err != nil {
			return fmt.Errorf("failed to get addresses: %w", err)
		}

		fmt.Println("Local IP Addresses:")
		for _, addr := range addrs {
			if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					fmt.Printf("  IPv4: %s\n", ipnet.IP.String())
				} else {
					fmt.Printf("  IPv6: %s\n", ipnet.IP.String())
				}
			}
		}

		return nil
	},
}

var ipPublicCmd = &cobra.Command{
	Use:   "public",
	Short: "Show public IP address",
	RunE: func(cmd *cobra.Command, args []string) error {
		resp, err := http.Get("https://api.ipify.org?format=json")
		if err != nil {
			return fmt.Errorf("failed to get IP: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read response: %w", err)
		}

		var result struct {
			IP string `json:"ip"`
		}
		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("failed to parse response: %w", err)
		}

		fmt.Printf("Public IP: %s\n", color.GreenString(result.IP))
		return nil
	},
}

func init() {
	rootCmd.AddCommand(ipCmd)
	ipCmd.AddCommand(ipLookupCmd)
	ipCmd.AddCommand(ipLocalCmd)
	ipCmd.AddCommand(ipPublicCmd)
}
