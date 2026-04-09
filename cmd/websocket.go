package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fatih/color"
	"github.com/gorilla/websocket"
	"github.com/spf13/cobra"
)

var wsCmd = &cobra.Command{
	Use:   "ws [url]",
	Short: "WebSocket client",
	Long:  `Connect to WebSocket server and send/receive messages.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]
		headers, _ := cmd.Flags().GetStringSlice("header")

		// Ensure URL has ws:// or wss:// scheme
		if len(url) < 3 || (url[:3] != "ws:" && url[:4] != "wss:") {
			if len(url) >= 5 && url[:5] == "http:" {
				url = "ws:" + url[5:]
			} else if len(url) >= 6 && url[:6] == "https:" {
				url = "wss:" + url[6:]
			} else {
				url = "wss://" + url
			}
		}

		// Create dialer
		dialer := websocket.DefaultDialer
		dialer.HandshakeTimeout = 10 * time.Second

		// Parse custom headers
		headerMap := make(map[string][]string)
		for _, h := range headers {
			parts := splitHeader(h)
			if len(parts) == 2 {
				headerMap[parts[0]] = append(headerMap[parts[0]], parts[1])
			}
		}

		// Connect
		fmt.Printf("Connecting to %s...\n", url)
		conn, _, err := dialer.Dial(url, headerMap)
		if err != nil {
			return fmt.Errorf("connection failed: %w", err)
		}
		defer conn.Close()

		fmt.Println(color.GreenString("Connected!"))
		fmt.Println("Type messages and press Enter to send. Press Ctrl+C to exit.")
		fmt.Println()

		// Context for cancellation
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle interrupt signal
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Start reader goroutine — blocks on ReadMessage, no busy-loop.
		go func() {
			defer cancel()
			for {
				_, message, err := conn.ReadMessage()
				if err != nil {
					if ctx.Err() == nil && websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
						fmt.Printf("\n%s: %v\n", color.RedString("Error"), err)
					}
					return
				}
				fmt.Printf("\r%s %s\n", color.CyanString("<--"), string(message))
				fmt.Print("> ")
			}
		}()

		// Read from stdin and send messages
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("> ")

		for {
			select {
			case <-sigChan:
				fmt.Println("\nDisconnecting...")
				conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
				return nil
			default:
				line, err := reader.ReadString('\n')
				if err != nil {
					return err
				}

				line = line[:len(line)-1] // Remove newline
				if line == "" {
					fmt.Print("> ")
					continue
				}

				if err := conn.WriteMessage(websocket.TextMessage, []byte(line)); err != nil {
					fmt.Printf("%s: %v\n", color.RedString("Send error"), err)
					return err
				}
				fmt.Printf("%s %s\n", color.GreenString("-->"), line)
				fmt.Print("> ")
			}
		}
	},
}

func splitHeader(h string) []string {
	for i := 0; i < len(h); i++ {
		if h[i] == ':' {
			return []string{h[:i], h[i+1:]}
		}
	}
	return nil
}

func init() {
	rootCmd.AddCommand(wsCmd)
	wsCmd.Flags().StringSliceP("header", "H", []string{}, "Custom headers (format: 'Key: Value')")
}
