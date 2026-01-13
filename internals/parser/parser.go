package parser

import (
	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
)

type Parser struct {
	Tokens []lexer.Token
	Pos    int
	Module ast.Module
}

func NewParser(tokens []lexer.Token) *Parser {
	return &Parser{
		Tokens: tokens,
		Pos:    0,
		Module: ast.Module{Body: []ast.Stmt{}},
	}
}
