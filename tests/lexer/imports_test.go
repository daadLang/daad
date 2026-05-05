package lexer_test

import (
	"testing"

	"github.com/daadLang/daad/internals/lexer"
)

func countTokens(tokens []lexer.Token, typ lexer.TokenType) int {
	count := 0
	for _, tok := range tokens {
		if tok.Type == typ {
			count++
		}
	}
	return count
}

func TestTokenizeImportsFixture(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/imports.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize import fixture: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("expected import fixture to produce tokens")
	}

	if tokens[0].Type != lexer.COMMENT {
		t.Fatalf("expected first token to be COMMENT, got %s", tokens[0].Type.String())
	}

	if countTokens(tokens, lexer.FROM) != 7 {
		t.Fatalf("expected 7 FROM tokens, got %d", countTokens(tokens, lexer.FROM))
	}

	if countTokens(tokens, lexer.IMPORT) != 11 {
		t.Fatalf("expected 11 IMPORT tokens, got %d", countTokens(tokens, lexer.IMPORT))
	}

	if countTokens(tokens, lexer.AS) != 5 {
		t.Fatalf("expected 5 AS tokens, got %d", countTokens(tokens, lexer.AS))
	}

	if countTokens(tokens, lexer.MULT) != 2 {
		t.Fatalf("expected 2 MULT tokens for star imports, got %d", countTokens(tokens, lexer.MULT))
	}

	if countTokens(tokens, lexer.DOT) != 11 {
		t.Fatalf("expected 11 DOT tokens for dotted and relative imports, got %d", countTokens(tokens, lexer.DOT))
	}
}
