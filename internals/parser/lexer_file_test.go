package parser

import (
	"testing"
)

func TestTokenizeBasicArithmetic(t *testing.T) {
	tokens, err := Tokenize("../../tests/basic_arithmetic.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	if len(tokens) == 0 {
		t.Fatal("Expected tokens but got none")
	}

	if tokens[len(tokens)-1].Type != EOF {
		t.Errorf("Expected EOF as last token, got %v", tokens[len(tokens)-1].Type)
	}

	hasNumber := false
	hasPlus := false
	hasAssign := false
	for _, tok := range tokens {
		if tok.Type == NUMBER {
			hasNumber = true
		}
		if tok.Type == PLUS {
			hasPlus = true
		}
		if tok.Type == ASSIGN {
			hasAssign = true
		}
	}

	if !hasNumber {
		t.Error("Expected to find NUMBER tokens")
	}
	if !hasPlus {
		t.Error("Expected to find PLUS token")
	}
	if !hasAssign {
		t.Error("Expected to find ASSIGN token")
	}
}

func TestTokenizeControlFlow(t *testing.T) {
	tokens, err := Tokenize("../../tests/control_flow.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	hasIf := false
	hasElse := false
	hasWhile := false
	hasFor := false
	hasIn := false

	for _, tok := range tokens {
		if tok.Type == IF {
			hasIf = true
		}
		if tok.Type == ELSE {
			hasElse = true
		}
		if tok.Type == WHILE {
			hasWhile = true
		}
		if tok.Type == FOR {
			hasFor = true
		}
		if tok.Type == IN {
			hasIn = true
		}
	}

	if !hasIf {
		t.Error("Expected to find IF keyword")
	}
	if !hasElse {
		t.Error("Expected to find ELSE keyword")
	}
	if !hasWhile {
		t.Error("Expected to find WHILE keyword")
	}
	if !hasFor {
		t.Error("Expected to find FOR keyword")
	}
	if !hasIn {
		t.Error("Expected to find IN keyword")
	}
}

func TestTokenizeFunctions(t *testing.T) {
	tokens, err := Tokenize("../../tests/functions.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	hasFunc := false
	hasReturn := false
	hasRetType := false
	hasLParen := false
	hasRParen := false

	for _, tok := range tokens {
		if tok.Type == FUNC {
			hasFunc = true
		}
		if tok.Type == RETURN {
			hasReturn = true
		}
		if tok.Type == RETTYPE {
			hasRetType = true
		}
		if tok.Type == LPAREN {
			hasLParen = true
		}
		if tok.Type == RPAREN {
			hasRParen = true
		}
	}

	if !hasFunc {
		t.Error("Expected to find FUNC keyword")
	}
	if !hasReturn {
		t.Error("Expected to find RETURN keyword")
	}
	if !hasRetType {
		t.Error("Expected to find RETTYPE (->)")
	}
	if !hasLParen {
		t.Error("Expected to find LPAREN")
	}
	if !hasRParen {
		t.Error("Expected to find RPAREN")
	}
}

func TestTokenizeStrings(t *testing.T) {
	tokens, err := Tokenize("../../tests/strings.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	// Count string tokens
	stringCount := 0
	for _, tok := range tokens {
		if tok.Type == STRING {
			stringCount++
		}
	}

	if stringCount < 5 {
		t.Errorf("Expected at least 5 STRING tokens, got %d", stringCount)
	}
}

func TestTokenizeOperators(t *testing.T) {
	tokens, err := Tokenize("../../tests/operators.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	operatorTests := map[TokenType]string{
		EQ:            "==",
		NEQ:           "!=",
		LESS:          "<",
		GREATER:       ">",
		LEQ:           "<=",
		GEQ:           ">=",
		AND:           "و",
		OR:            "او",
		NOT:           "ليس",
		BITWISE_AND:   "&",
		BITWISE_OR:    "|",
		BITWISE_XOR:   "^",
		BITWISE_NOT:   "~",
		LSHIFT:        "<<",
		RSHIFT:        ">>",
		PLUS_ASSIGN:   "+=",
		MINUS_ASSIGN:  "-=",
		MULT_ASSIGN:   "*=",
		DIVIDE_ASSIGN: "/=",
		MOD_ASSIGN:    "%=",
		POWER_ASSIGN:  "**=",
	}

	foundOperators := make(map[TokenType]bool)
	for _, tok := range tokens {
		if _, exists := operatorTests[tok.Type]; exists {
			foundOperators[tok.Type] = true
		}
	}

	for opType, opName := range operatorTests {
		if !foundOperators[opType] {
			t.Errorf("Expected to find operator %s (%s)", opType, opName)
		}
	}
}

func TestTokenizeKeywords(t *testing.T) {
	tokens, err := Tokenize("../../tests/keywords.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	keywordTests := []TokenType{
		IF, ELIF, ELSE, WHILE, REPEAT, TIMES,
		RETURN, BREAK, CONTINUE, TRUE, FALSE,
	}

	foundKeywords := make(map[TokenType]bool)
	for _, tok := range tokens {
		for _, kwType := range keywordTests {
			if tok.Type == kwType {
				foundKeywords[kwType] = true
			}
		}
	}

	for _, kwType := range keywordTests {
		if !foundKeywords[kwType] {
			t.Errorf("Expected to find keyword %s", kwType)
		}
	}
}

func TestTokenizeFileNotFound(t *testing.T) {
	_, err := Tokenize("../../tests/nonexistent.daad")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestTokenizeEmptyFile(t *testing.T) {
	tokens, err := Tokenize("../../tests/basic_arithmetic.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize: %v", err)
	}

	// At minimum, should have EOF token
	hasEOF := false
	for _, tok := range tokens {
		if tok.Type == EOF {
			hasEOF = true
		}
	}

	if !hasEOF {
		t.Error("Expected EOF token in output")
	}
}

func TestTokenDelimiters(t *testing.T) {
	tokens, err := Tokenize("../../tests/functions.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	delimiterTests := []TokenType{
		LPAREN, RPAREN, COMMA, COLON,
	}

	foundDelimiters := make(map[TokenType]bool)
	for _, tok := range tokens {
		for _, delimType := range delimiterTests {
			if tok.Type == delimType {
				foundDelimiters[delimType] = true
			}
		}
	}

	for _, delimType := range delimiterTests {
		if !foundDelimiters[delimType] {
			t.Errorf("Expected to find delimiter %s", delimType)
		}
	}
}
