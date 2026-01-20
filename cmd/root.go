package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/daadLang/daad/internals/interpreter"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "daad [file]",
	Short: "Daad - Arabic Programming Language",
	Long:  `Daad (ض) is an Arabic programming language with Python-like syntax.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		runInterpreter(args[0])

	},
}

var rootCmdAr = &cobra.Command{
	Use:   "ض [ملف]",
	Short: "ض - لغة برمجة عربية",
	Long:  `ض (daad) هي لغة برمجة عربية بصيغة مشابهة لبايثون.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			return
		}
		runInterpreter(args[0])
	},
}

func Execute() {
	binaryName := filepath.Base(os.Args[0])

	isArabicBinary := false
	for _, r := range binaryName {
		if unicode.Is(unicode.Arabic, r) ||
			(unicode.Is(unicode.Inherited, r) && strings.ContainsRune("ض", r)) {
			isArabicBinary = true
			break
		}
	}

	if isArabicBinary {
		if err := rootCmdAr.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(tokenizeCmd)
	rootCmdAr.AddCommand(tokenizeCmdAr)
	rootCmd.AddCommand(astCmd)
	rootCmdAr.AddCommand(astCmdAr)
}

func runInterpreter(filePath string) {
	// check if file exists
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

	i := interpreter.NewInterpreter()
	i.Run(&module)

}
