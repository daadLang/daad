package ast_test

import (
	"flag"
	"fmt"
	"os"
	"testing"

	"github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
)

var inputFile string

func init() {
	flag.StringVar(&inputFile, "file", "", "Path to the .daad file to parse and print AST")
}

func TestPrintAST(t *testing.T) {
	var testFiles []string

	if inputFile != "" {
		testFiles = append(testFiles, inputFile)
	} else {
		files, err := os.ReadDir("../examples/")
		if err != nil {
			t.Fatalf("Failed to read examples directory: %v", err)
		}
		for _, file := range files {
			if !file.IsDir() && len(file.Name()) > 5 && file.Name()[len(file.Name())-5:] == ".daad" {
				testFiles = append(testFiles, "../examples/"+file.Name())
			}
		}
	}

	for _, file := range testFiles {
		t.Run(file, func(t *testing.T) {
			tokens, err := lexer.Tokenize(file)
			if err != nil {
				t.Fatalf("Failed to tokenize file %s: %v", file, err)
			}

			p := parser.NewParser(tokens)
			module := p.Parse()

			fmt.Printf("AST for %s:\n", file)
			ast.PrintAST(module)
		})
	}
}
