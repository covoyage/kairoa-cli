package cmd

import (
	"encoding/json"
	"fmt"

	"github.com/spf13/cobra"
)

var keycodeCmd = &cobra.Command{
	Use:   "keycode",
	Short: "Keyboard key code reference",
	Long:  `Show keyboard key codes and event properties.`,
}

var keycodeLookupCmd = &cobra.Command{
	Use:   "lookup [key]",
	Short: "Look up key code information",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		key := args[0]

		keyInfo := getKeyInfo(key)

		format, _ := cmd.Flags().GetString("format")
		if format == "json" {
			jsonData, _ := json.MarshalIndent(keyInfo, "", "  ")
			fmt.Println(string(jsonData))
		} else {
			fmt.Printf("Key: %s\n", keyInfo.Key)
			fmt.Printf("Code: %s\n", keyInfo.Code)
			fmt.Printf("KeyCode: %d\n", keyInfo.KeyCode)
			fmt.Printf("Which: %d\n", keyInfo.Which)
			fmt.Printf("Location: %d\n", keyInfo.Location)
			fmt.Printf("Modifiers: Alt=%v, Ctrl=%v, Shift=%v, Meta=%v\n",
				keyInfo.AltKey, keyInfo.CtrlKey, keyInfo.ShiftKey, keyInfo.MetaKey)
		}

		return nil
	},
}

var keycodeListCmd = &cobra.Command{
	Use:   "list",
	Short: "List common key codes",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Common Key Codes:")
		fmt.Println()
		fmt.Println("Letters:")
		for i := 0; i < 26; i++ {
			key := string(rune('A' + i))
			code := 65 + i
			fmt.Printf("  %s: KeyCode=%d, Code=Key%s\n", key, code, key)
		}
		fmt.Println()
		fmt.Println("Numbers:")
		for i := 0; i < 10; i++ {
			key := string(rune('0' + i))
			code := 48 + i
			fmt.Printf("  %s: KeyCode=%d, Code=Digit%s\n", key, code, key)
		}
		fmt.Println()
		fmt.Println("Function Keys:")
		for i := 1; i <= 12; i++ {
			code := 111 + i
			fmt.Printf("  F%d: KeyCode=%d, Code=F%d\n", i, code, i)
		}
		fmt.Println()
		fmt.Println("Special Keys:")
		specialKeys := map[string]int{
			"Enter": 13, "Escape": 27, "Space": 32, "Backspace": 8,
			"Tab": 9, "Shift": 16, "Ctrl": 17, "Alt": 18,
			"CapsLock": 20, "Pause": 19, "Insert": 45, "Delete": 46,
			"Home": 36, "End": 35, "PageUp": 33, "PageDown": 34,
			"ArrowLeft": 37, "ArrowUp": 38, "ArrowRight": 39, "ArrowDown": 40,
		}
		for name, code := range specialKeys {
			fmt.Printf("  %s: KeyCode=%d\n", name, code)
		}

		return nil
	},
}

type KeyInfo struct {
	Key       string `json:"key"`
	Code      string `json:"code"`
	KeyCode   int    `json:"keyCode"`
	Which     int    `json:"which"`
	Location  int    `json:"location"`
	AltKey    bool   `json:"altKey"`
	CtrlKey   bool   `json:"ctrlKey"`
	ShiftKey  bool   `json:"shiftKey"`
	MetaKey   bool   `json:"metaKey"`
}

func getKeyInfo(key string) KeyInfo {
	// Simple key code mapping
	keyCodes := map[string]int{
		"a": 65, "b": 66, "c": 67, "d": 68, "e": 69, "f": 70,
		"g": 71, "h": 72, "i": 73, "j": 74, "k": 75, "l": 76,
		"m": 77, "n": 78, "o": 79, "p": 80, "q": 81, "r": 82,
		"s": 83, "t": 84, "u": 85, "v": 86, "w": 87, "x": 88,
		"y": 89, "z": 90,
		"0": 48, "1": 49, "2": 50, "3": 51, "4": 52,
		"5": 53, "6": 54, "7": 55, "8": 56, "9": 57,
		"Enter": 13, "Escape": 27, "Space": 32, "Backspace": 8,
		"Tab": 9, "Shift": 16, "Control": 17, "Alt": 18,
	}

	keyCode := keyCodes[key]
	if keyCode == 0 {
		keyCode = -1
	}

	return KeyInfo{
		Key:      key,
		Code:     "Key" + key,
		KeyCode:  keyCode,
		Which:    keyCode,
		Location: 0,
	}
}

func init() {
	rootCmd.AddCommand(keycodeCmd)
	keycodeCmd.AddCommand(keycodeLookupCmd)
	keycodeCmd.AddCommand(keycodeListCmd)

	keycodeLookupCmd.Flags().StringP("format", "f", "text", "Output format (text, json)")
}
