package ast_test

import (
	"flag"
	"fmt"
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
	if inputFile == "" {
		t.Skip("No input file provided. Use -file flag to specify a .daad file")
	}

	tokens, err := lexer.Tokenize(inputFile)
	if err != nil {
		t.Fatalf("Failed to tokenize file %s: %v", inputFile, err)
	}

	p := parser.NewParser(tokens)
	module := p.Parse()

	fmt.Printf("AST for %s:\n", inputFile)
	ast.PrintAST(module)
}
