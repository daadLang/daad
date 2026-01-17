package cmd

import (
	"fmt"
	"os"

	"github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
	"github.com/spf13/cobra"
)

var astCmd = &cobra.Command{
	Use:   "ast [file]",
	Short: "Display AST for a Daad source file",
	Long:  `Display AST for a Daad source file`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAst(args[0])
	},
}

var astCmdAr = &cobra.Command{
	Use:   "هيكل [ملف]",
	Short: "عرض الشجرة التركيبية المجردة لملف مصدر Daad",
	Long:  `عرض الشجرة التركيبية المجردة لملف مصدر Daad.`,
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAst(args[0])
	},
}

func runAst(filePath string) {
	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		fmt.Fprintf(os.Stderr, "Error: File '%s' does not exist\n", filePath)
		os.Exit(1)
	}

	// ast the file
	tokens, err := lexer.Tokenize(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error tokenizing file: %v\n", err)
		os.Exit(1)
	}

	p := parser.NewParser(tokens)
	module := p.Parse()

	ast.PrintAST(module)
}
