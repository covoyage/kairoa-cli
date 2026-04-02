package cmd

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/spf13/cobra"
)

var dnsCmd = &cobra.Command{
	Use:   "dns [domain]",
	Short: "DNS lookup utilities",
	Long:  `Perform DNS lookups for domains.`,
}

var dnsLookupCmd = &cobra.Command{
	Use:   "lookup [domain]",
	Short: "Look up DNS records",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordTypes, _ := cmd.Flags().GetStringSlice("type")

		if len(recordTypes) == 0 {
			recordTypes = []string{"A", "AAAA", "MX", "TXT", "NS", "CNAME"}
		}

		for _, recordType := range recordTypes {
			fmt.Printf("\n=== %s Records ===\n", recordType)
			if err := lookupDNSRecord(domain, recordType); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
		return nil
	},
}

func lookupDNSRecord(domain, recordType string) error {
	switch strings.ToUpper(recordType) {
	case "A":
		ips, err := net.LookupIP(domain)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 != nil {
				fmt.Printf("  %s\n", ipv4.String())
			}
		}

	case "AAAA":
		ips, err := net.LookupIP(domain)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		for _, ip := range ips {
			if ipv4 := ip.To4(); ipv4 == nil {
				fmt.Printf("  %s\n", ip.String())
			}
		}

	case "MX":
		mxs, err := net.LookupMX(domain)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		if len(mxs) == 0 {
			fmt.Println("  No MX records found")
		}
		for _, mx := range mxs {
			fmt.Printf("  %d %s\n", mx.Pref, mx.Host)
		}

	case "TXT":
		txts, err := net.LookupTXT(domain)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		if len(txts) == 0 {
			fmt.Println("  No TXT records found")
		}
		for _, txt := range txts {
			fmt.Printf("  %s\n", txt)
		}

	case "NS":
		nss, err := net.LookupNS(domain)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		if len(nss) == 0 {
			fmt.Println("  No NS records found")
		}
		for _, ns := range nss {
			fmt.Printf("  %s\n", ns.Host)
		}

	case "CNAME":
		cname, err := net.LookupCNAME(domain)
		if err != nil {
			return fmt.Errorf("lookup failed: %w", err)
		}
		fmt.Printf("  %s\n", cname)

	case "SOA":
		// SOA record lookup is not directly supported in standard library
		fmt.Println("  SOA record lookup not supported in this implementation")

	default:
		return fmt.Errorf("unsupported record type: %s", recordType)
	}
	return nil
}

var dnsCloudflareCmd = &cobra.Command{
	Use:   "cloudflare [domain]",
	Short: "Query DNS using Cloudflare API",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		domain := args[0]
		recordType, _ := cmd.Flags().GetString("type")

		url := fmt.Sprintf("https://cloudflare-dns.com/dns-query?name=%s&type=%s",
			domain, recordType)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return err
		}
		req.Header.Set("Accept", "application/dns-json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return fmt.Errorf("request failed: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("read failed: %w", err)
		}

		var result struct {
			Status   int `json:"Status"`
			Answer   []struct {
				Name string `json:"name"`
				Type int    `json:"type"`
				TTL  int    `json:"TTL"`
				Data string `json:"data"`
			} `json:"Answer"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return fmt.Errorf("parse failed: %w", err)
		}

		fmt.Printf("Domain: %s\n", domain)
		fmt.Printf("Record Type: %s\n", recordType)
		fmt.Printf("Status: %d\n\n", result.Status)

		if len(result.Answer) == 0 {
			fmt.Println("No records found")
			return nil
		}

		fmt.Println("Answers:")
		for _, ans := range result.Answer {
			fmt.Printf("  [%d] %s (TTL: %d)\n", ans.Type, ans.Data, ans.TTL)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(dnsCmd)
	dnsCmd.AddCommand(dnsLookupCmd)
	dnsCmd.AddCommand(dnsCloudflareCmd)

	dnsLookupCmd.Flags().StringSliceP("type", "t", []string{}, "Record types to query (A, AAAA, MX, TXT, NS, CNAME, SOA)")
	dnsCloudflareCmd.Flags().StringP("type", "t", "A", "Record type to query")
}
