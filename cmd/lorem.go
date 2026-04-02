package cmd

import (
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var loremCmd = &cobra.Command{
	Use:   "lorem",
	Short: "Generate lorem ipsum text",
	Long:  `Generate placeholder text (lorem ipsum).`,
}

var loremWords = []string{
	"lorem", "ipsum", "dolor", "sit", "amet", "consectetur", "adipiscing", "elit",
	"sed", "do", "eiusmod", "tempor", "incididunt", "ut", "labore", "et", "dolore",
	"magna", "aliqua", "enim", "ad", "minim", "veniam", "quis", "nostrud",
	"exercitation", "ullamco", "laboris", "nisi", "aliquip", "ex", "ea", "commodo",
	"consequat", "duis", "aute", "irure", "in", "reprehenderit", "voluptate",
	"velit", "esse", "cillum", "fugiat", "nulla", "pariatur", "excepteur", "sint",
	"occaecat", "cupidatat", "non", "proident", "sunt", "culpa", "qui", "officia",
	"deserunt", "mollit", "anim", "id", "est", "laborum",
}

var loremTextCmd = &cobra.Command{
	Use:   "text",
	Short: "Generate lorem ipsum paragraphs",
	RunE: func(cmd *cobra.Command, args []string) error {
		paragraphs, _ := cmd.Flags().GetInt("paragraphs")
		wordsPerParagraph, _ := cmd.Flags().GetInt("words")

		rand.Seed(time.Now().UnixNano())

		for i := 0; i < paragraphs; i++ {
			if i > 0 {
				fmt.Println()
			}
			fmt.Println(generateParagraph(wordsPerParagraph))
		}

		return nil
	},
}

var loremWordsCmd = &cobra.Command{
	Use:   "words",
	Short: "Generate lorem ipsum words",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")

		rand.Seed(time.Now().UnixNano())

		words := make([]string, count)
		for i := 0; i < count; i++ {
			words[i] = loremWords[rand.Intn(len(loremWords))]
		}

		fmt.Println(strings.Join(words, " "))
		return nil
	},
}

var loremSentencesCmd = &cobra.Command{
	Use:   "sentences",
	Short: "Generate lorem ipsum sentences",
	RunE: func(cmd *cobra.Command, args []string) error {
		count, _ := cmd.Flags().GetInt("count")
		wordsPerSentence, _ := cmd.Flags().GetInt("words")

		rand.Seed(time.Now().UnixNano())

		for i := 0; i < count; i++ {
			fmt.Println(generateSentence(wordsPerSentence))
		}

		return nil
	},
}

func generateParagraph(words int) string {
	sentences := words / 10
	if sentences < 3 {
		sentences = 3
	}

	var result []string
	for i := 0; i < sentences; i++ {
		result = append(result, generateSentence(10+rand.Intn(10)))
	}

	return strings.Join(result, " ")
}

func generateSentence(words int) string {
	rand.Seed(time.Now().UnixNano())

	var result []string
	for i := 0; i < words; i++ {
		word := loremWords[rand.Intn(len(loremWords))]
		if i == 0 {
			word = strings.Title(word)
		}
		result = append(result, word)
	}

	return strings.Join(result, " ") + "."
}

func init() {
	rootCmd.AddCommand(loremCmd)
	loremCmd.AddCommand(loremTextCmd)
	loremCmd.AddCommand(loremWordsCmd)
	loremCmd.AddCommand(loremSentencesCmd)

	loremTextCmd.Flags().IntP("paragraphs", "p", 3, "Number of paragraphs")
	loremTextCmd.Flags().IntP("words", "w", 50, "Words per paragraph")

	loremWordsCmd.Flags().IntP("count", "c", 10, "Number of words")

	loremSentencesCmd.Flags().IntP("count", "c", 3, "Number of sentences")
	loremSentencesCmd.Flags().IntP("words", "w", 10, "Words per sentence")
}
