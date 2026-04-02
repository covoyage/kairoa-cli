package cmd

import (
	"crypto/tls"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var tlsCmd = &cobra.Command{
	Use:   "tls",
	Short: "TLS/SSL checker",
	Long:  `Check TLS configuration and certificate details.`,
}

var tlsCheckCmd = &cobra.Command{
	Use:   "check [host:port]",
	Short: "Check TLS configuration",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		if !strings.Contains(address, ":") {
			address = address + ":443"
		}

		// Try different TLS versions
		versions := map[string]uint16{
			"TLS 1.0": tls.VersionTLS10,
			"TLS 1.1": tls.VersionTLS11,
			"TLS 1.2": tls.VersionTLS12,
			"TLS 1.3": tls.VersionTLS13,
		}

		fmt.Printf("Checking TLS configuration for %s...\n\n", address)

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Version", "Supported", "Cipher Suites"})

		for name, version := range versions {
			config := &tls.Config{
				MinVersion: version,
				MaxVersion: version,
			}

			conn, err := tls.Dial("tcp", address, config)
			if err != nil {
				table.Append([]string{name, color.RedString("No"), "-"})
				continue
			}
			conn.Close()

			// Get cipher suites
			config.CipherSuites = nil
			conn, err = tls.Dial("tcp", address, config)
			if err != nil {
				table.Append([]string{name, color.GreenString("Yes"), "-"})
				continue
			}
			
			state := conn.ConnectionState()
			cipherSuite := tls.CipherSuiteName(state.CipherSuite)
			conn.Close()

			table.Append([]string{name, color.GreenString("Yes"), cipherSuite})
		}

		table.Render()

		// Check certificate
		fmt.Println("\nCertificate Info:")
		config := &tls.Config{}
		conn, err := tls.Dial("tcp", address, config)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) > 0 {
			cert := certs[0]
			
			certTable := tablewriter.NewWriter(os.Stdout)
			certTable.SetAutoWrapText(false)

			certTable.Append([]string{"Subject", cert.Subject.String()})
			certTable.Append([]string{"Issuer", cert.Issuer.String()})
			certTable.Append([]string{"Valid From", cert.NotBefore.Format("2006-01-02")})
			
			validColor := color.GreenString
			if time.Now().After(cert.NotAfter) {
				validColor = color.RedString
			}
			certTable.Append([]string{"Valid Until", validColor(cert.NotAfter.Format("2006-01-02"))})
			
			daysUntil := int(time.Until(cert.NotAfter).Hours() / 24)
			certTable.Append([]string{"Days Until Expiry", fmt.Sprintf("%d", daysUntil)})
			
			certTable.Render()
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(tlsCmd)
	tlsCmd.AddCommand(tlsCheckCmd)
}
