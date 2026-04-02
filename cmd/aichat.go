package cmd

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var aiChatCmd = &cobra.Command{
	Use:   "aichat",
	Short: "AI chat client",
	Long:  `Chat with AI models (OpenAI, Anthropic, etc.).`,
}

var aiChatSendCmd = &cobra.Command{
	Use:   "send [message]",
	Short: "Send a message to AI",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		message := strings.Join(args, " ")
		apiKey, _ := cmd.Flags().GetString("api-key")
		baseURL, _ := cmd.Flags().GetString("base-url")
		model, _ := cmd.Flags().GetString("model")
		stream, _ := cmd.Flags().GetBool("stream")

		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		if apiKey == "" {
			return fmt.Errorf("API key required. Set OPENAI_API_KEY or use --api-key")
		}

		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}

		return sendChatRequest(baseURL, apiKey, model, message, stream)
	},
}

var aiChatInteractiveCmd = &cobra.Command{
	Use:   "interactive",
	Short: "Interactive chat session",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		baseURL, _ := cmd.Flags().GetString("base-url")
		model, _ := cmd.Flags().GetString("model")

		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		if apiKey == "" {
			return fmt.Errorf("API key required. Set OPENAI_API_KEY or use --api-key")
		}

		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}

		fmt.Println("AI Chat Interactive Mode")
		fmt.Println("Type 'exit' or 'quit' to exit")
		fmt.Println()

		var messages []map[string]string
		messages = append(messages, map[string]string{
			"role":    "system",
			"content": "You are a helpful assistant.",
		})

		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("You: ")
			if !scanner.Scan() {
				break
			}

			input := strings.TrimSpace(scanner.Text())
			if input == "" {
				continue
			}
			if input == "exit" || input == "quit" {
				break
			}

			messages = append(messages, map[string]string{
				"role":    "user",
				"content": input,
			})

			response, err := sendChatRequestWithHistory(baseURL, apiKey, model, messages)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}

			fmt.Printf("AI: %s\n\n", response)

			messages = append(messages, map[string]string{
				"role":    "assistant",
				"content": response,
			})
		}

		return nil
	},
}

var aiChatModelsCmd = &cobra.Command{
	Use:   "models",
	Short: "List available models",
	RunE: func(cmd *cobra.Command, args []string) error {
		apiKey, _ := cmd.Flags().GetString("api-key")
		baseURL, _ := cmd.Flags().GetString("base-url")

		if apiKey == "" {
			apiKey = os.Getenv("OPENAI_API_KEY")
		}

		if baseURL == "" {
			baseURL = "https://api.openai.com/v1"
		}

		if apiKey == "" {
			// Show common models without API
			fmt.Println("Common OpenAI Models:")
			fmt.Println("  gpt-4")
			fmt.Println("  gpt-4-turbo")
			fmt.Println("  gpt-3.5-turbo")
			fmt.Println("  gpt-4o")
			fmt.Println("  gpt-4o-mini")
			return nil
		}

		return listModels(baseURL, apiKey)
	},
}

func sendChatRequest(baseURL, apiKey, model, message string, stream bool) error {
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
		"stream": stream,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if len(result.Choices) > 0 {
		fmt.Println(result.Choices[0].Message.Content)
	}

	return nil
}

func sendChatRequestWithHistory(baseURL, apiKey, model string, messages []map[string]string) (string, error) {
	if model == "" {
		model = "gpt-3.5-turbo"
	}

	payload := map[string]interface{}{
		"model":    model,
		"messages": messages,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", baseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) > 0 {
		return result.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no response from AI")
}

func listModels(baseURL, apiKey string) error {
	req, err := http.NewRequest("GET", baseURL+"/models", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error: %s - %s", resp.Status, string(body))
	}

	var result struct {
		Data []struct {
			ID string `json:"id"`
		} `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	fmt.Println("Available Models:")
	for _, model := range result.Data {
		fmt.Printf("  %s\n", model.ID)
	}

	return nil
}

func init() {
	rootCmd.AddCommand(aiChatCmd)
	aiChatCmd.AddCommand(aiChatSendCmd)
	aiChatCmd.AddCommand(aiChatInteractiveCmd)
	aiChatCmd.AddCommand(aiChatModelsCmd)

	aiChatCmd.PersistentFlags().StringP("api-key", "k", "", "API key (or set OPENAI_API_KEY env)")
	aiChatCmd.PersistentFlags().StringP("base-url", "u", "", "API base URL")
	aiChatCmd.PersistentFlags().StringP("model", "m", "gpt-3.5-turbo", "Model name")

	aiChatSendCmd.Flags().BoolP("stream", "s", false, "Stream response")
}
