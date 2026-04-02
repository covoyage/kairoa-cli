package cmd

import (
	"crypto/sha256"
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var certCmd = &cobra.Command{
	Use:   "cert",
	Short: "Certificate utilities",
	Long:  `View and inspect SSL/TLS certificates.`,
}

var certViewCmd = &cobra.Command{
	Use:   "view [host:port]",
	Short: "View certificate from remote host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		if !strings.Contains(address, ":") {
			address = address + ":443"
		}

		insecure, _ := cmd.Flags().GetBool("insecure")

		// Connect and get certificate
		conf := &tls.Config{
			InsecureSkipVerify: insecure,
		}

		conn, err := tls.Dial("tcp", address, conf)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer conn.Close()

		certs := conn.ConnectionState().PeerCertificates
		if len(certs) == 0 {
			return fmt.Errorf("no certificates found")
		}

		fmt.Printf("Certificate chain for %s:\n\n", address)

		for i, cert := range certs {
			fmt.Printf("=== Certificate %d ===\n", i+1)
			printCertificate(cert)
			if i < len(certs)-1 {
				fmt.Println()
			}
		}

		return nil
	},
}

var certFileCmd = &cobra.Command{
	Use:   "file [path]",
	Short: "View certificate from file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]

		data, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read file: %w", err)
		}

		// Try to parse as PEM
		block, _ := pem.Decode(data)
		if block == nil {
			return fmt.Errorf("failed to decode PEM block")
		}

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return fmt.Errorf("failed to parse certificate: %w", err)
		}

		printCertificate(cert)
		return nil
	},
}

var certVerifyCmd = &cobra.Command{
	Use:   "verify [host:port]",
	Short: "Verify certificate chain",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		if !strings.Contains(address, ":") {
			address = address + ":443"
		}

		// Connect and get certificate
		conf := &tls.Config{}

		conn, err := tls.Dial("tcp", address, conf)
		if err != nil {
			return fmt.Errorf("failed to connect: %w", err)
		}
		defer conn.Close()

		state := conn.ConnectionState()
		certs := state.PeerCertificates

		fmt.Printf("Verifying certificate chain for %s:\n\n", address)

		// Build cert pool from system roots
		roots, err := x509.SystemCertPool()
		if err != nil {
			roots = x509.NewCertPool()
		}

		// Verify each certificate
		for i, cert := range certs {
			opts := x509.VerifyOptions{
				Roots:         roots,
				CurrentTime:   time.Now(),
				DNSName:       state.ServerName,
				Intermediates: x509.NewCertPool(),
			}

			// Add remaining certs as intermediates
			for j := i + 1; j < len(certs); j++ {
				opts.Intermediates.AddCert(certs[j])
			}

			chains, err := cert.Verify(opts)
			if err != nil {
				fmt.Printf("Certificate %d: %s\n", i+1, color.RedString("FAILED - %v", err))
			} else {
				fmt.Printf("Certificate %d: %s (chain length: %d)\n", i+1, color.GreenString("VERIFIED"), len(chains[0]))
			}
		}

		return nil
	},
}

func printCertificate(cert *x509.Certificate) {
	table := tablewriter.NewWriter(os.Stdout)

	// Subject
	table.Append([]string{"Subject", cert.Subject.String()})

	// Issuer
	table.Append([]string{"Issuer", cert.Issuer.String()})

	// Serial Number
	table.Append([]string{"Serial Number", cert.SerialNumber.String()})

	// Validity
	notBefore := cert.NotBefore.Format("2006-01-02 15:04:05 MST")
	notAfter := cert.NotAfter.Format("2006-01-02 15:04:05 MST")

	validityColor := color.GreenString
	if time.Now().After(cert.NotAfter) {
		validityColor = color.RedString
	} else if time.Now().Before(cert.NotBefore) {
		validityColor = color.YellowString
	}

	table.Append([]string{"Not Before", notBefore})
	table.Append([]string{"Not After", validityColor(notAfter)})

	// Key Algorithm
	table.Append([]string{"Key Algorithm", cert.PublicKeyAlgorithm.String()})

	// Signature Algorithm
	table.Append([]string{"Signature Algorithm", cert.SignatureAlgorithm.String()})

	// DNS Names
	if len(cert.DNSNames) > 0 {
		table.Append([]string{"DNS Names", strings.Join(cert.DNSNames, ", ")})
	}

	// IP Addresses
	if len(cert.IPAddresses) > 0 {
		var ips []string
		for _, ip := range cert.IPAddresses {
			ips = append(ips, ip.String())
		}
		table.Append([]string{"IP Addresses", strings.Join(ips, ", ")})
	}

	// Fingerprint
	fingerprint := sha256.Sum256(cert.Raw)
	table.Append([]string{"SHA-256 Fingerprint", hex.EncodeToString(fingerprint[:])})

	table.Render()
}

func init() {
	rootCmd.AddCommand(certCmd)
	certCmd.AddCommand(certViewCmd)
	certCmd.AddCommand(certFileCmd)
	certCmd.AddCommand(certVerifyCmd)

	certViewCmd.Flags().BoolP("insecure", "k", false, "Skip certificate verification")
}
