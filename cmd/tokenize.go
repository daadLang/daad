package cmd

import (
	"fmt"
	"os"

	"github.com/TechGeeks-Club/daad/internals/lexer"
	"github.com/spf13/cobra"
)

var tokenizeCmd = &cobra.Command{
	Use:   "tokenize [file]",
	Short: "Tokenize a Daad source file",
	Long:  `Tokenize a Daad source file and display all tokens.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTokenize(args[0])
	},
}

var tokenizeCmdAr = &cobra.Command{
	Use:   "رمز [ملف]",
	Short: "تحليل ملف مصدر Daad إلى رموز",
	Long:  `تحليل ملف مصدر Daad وعرض جميع الرموز.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runTokenize(args[0])
	},
}

func runTokenize(filePath string) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", filePath)
		os.Exit(1)
	}

	// Tokenize the file
	tokens, err := lexer.Tokenize(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error tokenizing file: %v\n", err)
		os.Exit(1)
	}

	// Print tokens
	fmt.Printf("Tokenizing: %s\n", filePath)
	fmt.Println("---")

	for i, token := range tokens {
		if token.Type == lexer.EOF {
			fmt.Printf("%4d: %-15s %q\n", i, token.Type.String(), token.Value)
			break
		}
		fmt.Printf("%4d: %-15s %q\n", i, token.Type.String(), token.Value)
	}

	fmt.Println("---")
	fmt.Printf("Total tokens: %d\n", len(tokens))
}
