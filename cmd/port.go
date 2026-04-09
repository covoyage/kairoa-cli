package cmd

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var portCmd = &cobra.Command{
	Use:   "port",
	Short: "Port scanner",
	Long:  `Scan ports on a host to check if they are open.`,
}

var portScanCmd = &cobra.Command{
	Use:   "scan [host]",
	Short: "Scan ports on a host",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		host := args[0]
		startPort, _ := cmd.Flags().GetInt("start")
		endPort, _ := cmd.Flags().GetInt("end")
		timeout, _ := cmd.Flags().GetInt("timeout")
		concurrency, _ := cmd.Flags().GetInt("concurrency")
		showClosed, _ := cmd.Flags().GetBool("show-closed")

		if startPort < 1 || startPort > 65535 {
			return fmt.Errorf("start port must be between 1 and 65535")
		}
		if endPort < 1 || endPort > 65535 {
			return fmt.Errorf("end port must be between 1 and 65535")
		}
		if endPort < startPort {
			return fmt.Errorf("end port (%d) must be >= start port (%d)", endPort, startPort)
		}

		fmt.Printf("Scanning %s (ports %d-%d)...\n\n", host, startPort, endPort)
		startTime := time.Now()

		var openPorts []struct {
			port    int
			latency time.Duration
		}
		var mu sync.Mutex
		var wg sync.WaitGroup

		semaphore := make(chan struct{}, concurrency)

		for port := startPort; port <= endPort; port++ {
			wg.Add(1)
			semaphore <- struct{}{}

			go func(p int) {
				defer wg.Done()
				defer func() { <-semaphore }()

				address := fmt.Sprintf("%s:%d", host, p)
				start := time.Now()
				conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Millisecond)
				latency := time.Since(start)

				if err == nil {
					conn.Close()
					mu.Lock()
					openPorts = append(openPorts, struct {
						port    int
						latency time.Duration
					}{port: p, latency: latency})
					mu.Unlock()

					if showClosed {
						fmt.Printf("\rPort %d: %s (%.2f ms)\n", p, color.GreenString("OPEN"), float64(latency.Microseconds())/1000)
					}
				} else if showClosed {
					fmt.Printf("\rPort %d: %s\n", p, color.RedString("CLOSED"))
				}
			}(port)
		}

		wg.Wait()
		duration := time.Since(startTime)

		fmt.Printf("\n\nScan completed in %.2f seconds\n", duration.Seconds())
		fmt.Printf("Scanned %d ports, found %d open\n\n", endPort-startPort+1, len(openPorts))

		if len(openPorts) > 0 {
			table := tablewriter.NewWriter(os.Stdout)
			table.SetHeader([]string{"Port", "Status", "Latency", "Service"})

			for _, p := range openPorts {
				service := getServiceName(p.port)
				table.Append([]string{
					strconv.Itoa(p.port),
					color.GreenString("OPEN"),
					fmt.Sprintf("%.2f ms", float64(p.latency.Microseconds())/1000),
					service,
				})
			}
			table.Render()
		}

		return nil
	},
}

var portCheckCmd = &cobra.Command{
	Use:   "check [host:port]",
	Short: "Check if a specific port is open",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		address := args[0]
		timeout, _ := cmd.Flags().GetInt("timeout")

		if !strings.Contains(address, ":") {
			return fmt.Errorf("address must be in format host:port")
		}

		parts := strings.Split(address, ":")
		port, err := strconv.Atoi(parts[len(parts)-1])
		if err != nil {
			return fmt.Errorf("invalid port: %w", err)
		}

		fmt.Printf("Checking %s... ", address)
		conn, err := net.DialTimeout("tcp", address, time.Duration(timeout)*time.Millisecond)
		if err != nil {
			fmt.Println(color.RedString("CLOSED"))
			return nil
		}
		conn.Close()
		fmt.Println(color.GreenString("OPEN (%s)", getServiceName(port)))
		return nil
	},
}

func getServiceName(port int) string {
	services := map[int]string{
		20:    "FTP-Data",
		21:    "FTP",
		22:    "SSH",
		23:    "Telnet",
		25:    "SMTP",
		53:    "DNS",
		80:    "HTTP",
		110:   "POP3",
		143:   "IMAP",
		443:   "HTTPS",
		3306:  "MySQL",
		3389:  "RDP",
		5432:  "PostgreSQL",
		6379:  "Redis",
		8080:  "HTTP-Proxy",
		8443:  "HTTPS-Alt",
		9200:  "Elasticsearch",
		27017: "MongoDB",
	}
	if name, ok := services[port]; ok {
		return name
	}
	return "Unknown"
}

func init() {
	rootCmd.AddCommand(portCmd)
	portCmd.AddCommand(portScanCmd)
	portCmd.AddCommand(portCheckCmd)

	portScanCmd.Flags().IntP("start", "s", 1, "Start port")
	portScanCmd.Flags().IntP("end", "e", 1024, "End port")
	portScanCmd.Flags().IntP("timeout", "t", 700, "Timeout in milliseconds")
	portScanCmd.Flags().IntP("concurrency", "c", 200, "Number of concurrent connections")
	portScanCmd.Flags().Bool("show-closed", false, "Show closed ports (slower)")

	portCheckCmd.Flags().IntP("timeout", "t", 3000, "Timeout in milliseconds")
}
