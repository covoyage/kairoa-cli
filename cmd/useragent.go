package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
)

var useragentCmd = &cobra.Command{
	Use:   "useragent",
	Short: "User-Agent parser",
	Long:  `Parse and analyze User-Agent strings.`,
}

var useragentParseCmd = &cobra.Command{
	Use:   "parse [user-agent]",
	Short: "Parse User-Agent string",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ua := args[0]

		info := parseUserAgent(ua)

		fmt.Printf("Browser: %s\n", info.Browser)
		if info.BrowserVersion != "" {
			fmt.Printf("Browser Version: %s\n", info.BrowserVersion)
		}
		fmt.Printf("OS: %s\n", info.OS)
		if info.OSVersion != "" {
			fmt.Printf("OS Version: %s\n", info.OSVersion)
		}
		if info.Device != "" {
			fmt.Printf("Device: %s\n", info.Device)
		}
		fmt.Printf("Device Type: %s\n", info.DeviceType)

		return nil
	},
}

var useragentListCmd = &cobra.Command{
	Use:   "list",
	Short: "List common User-Agent strings",
	RunE: func(cmd *cobra.Command, args []string) error {
		uas := []struct {
			Name  string
			Value string
		}{
			{"Chrome (Windows)", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
			{"Firefox (Windows)", "Mozilla/5.0 (Windows NT 10.0; Win64; x64; rv:121.0) Gecko/20100101 Firefox/121.0"},
			{"Safari (macOS)", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Safari/605.1.15"},
			{"Edge (Windows)", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 Edg/120.0.0.0"},
			{"Chrome (Android)", "Mozilla/5.0 (Linux; Android 13; SM-G991B) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36"},
			{"Safari (iOS)", "Mozilla/5.0 (iPhone; CPU iPhone OS 17_2 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) Version/17.2 Mobile/15E148 Safari/604.1"},
			{"Chrome (Linux)", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36"},
			{"Googlebot", "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"},
		}

		fmt.Println("Common User-Agent Strings:")
		fmt.Println()

		for _, ua := range uas {
			fmt.Printf("%s:\n  %s\n\n", ua.Name, ua.Value)
		}

		return nil
	},
}

type UserAgentInfo struct {
	Browser        string
	BrowserVersion string
	OS             string
	OSVersion      string
	Device         string
	DeviceType     string
}

func parseUserAgent(ua string) UserAgentInfo {
	info := UserAgentInfo{
		Browser:    "Unknown",
		OS:         "Unknown",
		DeviceType: "desktop",
	}

	ua = strings.TrimSpace(ua)

	// Browser detection
	if strings.Contains(ua, "Chrome") && !strings.Contains(ua, "Edg") && !strings.Contains(ua, "OPR") {
		info.Browser = "Chrome"
		re := regexp.MustCompile(`Chrome/(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.BrowserVersion = matches[1]
		}
	} else if strings.Contains(ua, "Firefox") {
		info.Browser = "Firefox"
		re := regexp.MustCompile(`Firefox/(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.BrowserVersion = matches[1]
		}
	} else if strings.Contains(ua, "Safari") && !strings.Contains(ua, "Chrome") {
		info.Browser = "Safari"
		re := regexp.MustCompile(`Version/(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.BrowserVersion = matches[1]
		}
	} else if strings.Contains(ua, "Edg") {
		info.Browser = "Edge"
		re := regexp.MustCompile(`Edg/(\d+\.\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.BrowserVersion = matches[1]
		}
	}

	// OS detection
	if strings.Contains(ua, "Windows NT 10.0") {
		info.OS = "Windows"
		info.OSVersion = "10"
	} else if strings.Contains(ua, "Windows NT 6.3") {
		info.OS = "Windows"
		info.OSVersion = "8.1"
	} else if strings.Contains(ua, "Mac OS X") {
		info.OS = "macOS"
		re := regexp.MustCompile(`Mac OS X (\d+[._]\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.OSVersion = strings.ReplaceAll(matches[1], "_", ".")
		}
	} else if strings.Contains(ua, "Linux") {
		info.OS = "Linux"
	} else if strings.Contains(ua, "Android") {
		info.OS = "Android"
		re := regexp.MustCompile(`Android (\d+\.?\d*)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.OSVersion = matches[1]
		}
		info.DeviceType = "mobile"
	} else if strings.Contains(ua, "iPhone") || strings.Contains(ua, "iPad") {
		info.OS = "iOS"
		re := regexp.MustCompile(`OS (\d+_\d+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.OSVersion = strings.ReplaceAll(matches[1], "_", ".")
		}
		if strings.Contains(ua, "iPad") {
			info.DeviceType = "tablet"
		} else {
			info.DeviceType = "mobile"
		}
	}

	// Device detection
	if strings.Contains(ua, "iPhone") {
		info.Device = "iPhone"
	} else if strings.Contains(ua, "iPad") {
		info.Device = "iPad"
	} else if strings.Contains(ua, "Android") {
		re := regexp.MustCompile(`Android [^;]+; ([^;)]+)`)
		if matches := re.FindStringSubmatch(ua); len(matches) > 1 {
			info.Device = strings.TrimSpace(matches[1])
		}
	}

	// Bot detection
	if strings.Contains(ua, "bot") || strings.Contains(ua, "crawler") || strings.Contains(ua, "spider") {
		info.DeviceType = "bot"
	}

	return info
}

func init() {
	rootCmd.AddCommand(useragentCmd)
	useragentCmd.AddCommand(useragentParseCmd)
	useragentCmd.AddCommand(useragentListCmd)
}
