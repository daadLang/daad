package parser

import (
	"testing"
)

func TestLexerIdentifiers(t *testing.T) {
	input := "متغير عدد_الطلاب اسم123"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: IDENT, Value: "متغير"},
		{Type: IDENT, Value: "عدد_الطلاب"},
		{Type: IDENT, Value: "اسم123"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerNumbers(t *testing.T) {
	input := "42 123 0 999"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: NUMBER, Value: "42"},
		{Type: NUMBER, Value: "123"},
		{Type: NUMBER, Value: "0"},
		{Type: NUMBER, Value: "999"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerKeywords(t *testing.T) {
	input := "إذا وإذا وإلا طالما لكل في أرجع دالة اخرج تابع صحيح خطأ"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	// Keywords are normalized: إ→ا, أ→ا, ؤ→ا, ء→ا, ى→ا
	expected := []Token{
		{Type: IF, Value: "اذا"},    // إذا normalized
		{Type: ELIF, Value: "واذا"}, // وإذا normalized
		{Type: ELSE, Value: "والا"}, // وإلا normalized
		{Type: WHILE, Value: "طالما"},
		{Type: FOR, Value: "لكل"},
		{Type: IN, Value: "في"},
		{Type: RETURN, Value: "ارجع"}, // أرجع normalized
		{Type: FUNC, Value: "دالة"},
		{Type: BREAK, Value: "اخرج"},
		{Type: CONTINUE, Value: "تابع"},
		{Type: TRUE, Value: "صحيح"},
		{Type: FALSE, Value: "خطا"}, // خطأ normalized
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerArithmeticOperators(t *testing.T) {
	input := "+ - ** * // / %"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: PLUS, Value: "+"},
		{Type: MINUS, Value: "-"},
		{Type: POWER, Value: "**"},
		{Type: MULT, Value: "*"},
		{Type: FLOORDIV, Value: "//"},
		{Type: DIVIDE, Value: "/"},
		{Type: MOD, Value: "%"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerComparisonOperators(t *testing.T) {
	input := "== != < > <= >= << >>"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: EQ, Value: "=="},
		{Type: NEQ, Value: "!="},
		{Type: LESS, Value: "<"},
		{Type: GREATER, Value: ">"},
		{Type: LEQ, Value: "<="},
		{Type: GEQ, Value: ">="},
		{Type: LSHIFT, Value: "<<"},
		{Type: RSHIFT, Value: ">>"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerAssignmentOperators(t *testing.T) {
	input := "س = 5"
	input2 := "س += 5"
	input3 := "س -= 5"
	input4 := "س *= 5"
	input5 := "س /= 5"
	input6 := "س %= 5"
	input7 := "س **= 5"

	tests := []struct {
		input    string
		expected TokenType
	}{
		{input, ASSIGN},
		{input2, PLUS_ASSIGN},
		{input3, MINUS_ASSIGN},
		{input4, MULT_ASSIGN},
		{input5, DIVIDE_ASSIGN},
		{input6, MOD_ASSIGN},
		{input7, POWER_ASSIGN},
	}

	for _, test := range tests {
		lexer := Lexer{input: []int32(test.input)}
		lexer.Tokenize()
		found := false
		for _, tok := range lexer.Tokens {
			if tok.Type == test.expected {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find token type %s in input '%s'", test.expected, test.input)
		}
	}
}

func TestLexerAssignmentOperatorsSeparate(t *testing.T) {
	input := "="
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: ASSIGN, Value: "="},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerBitwiseOperators(t *testing.T) {
	input := "& | ^ ~ << >>"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: BITWISE_AND, Value: "&"},
		{Type: BITWISE_OR, Value: "|"},
		{Type: BITWISE_XOR, Value: "^"},
		{Type: BITWISE_NOT, Value: "~"},
		{Type: LSHIFT, Value: "<<"},
		{Type: RSHIFT, Value: ">>"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerDelimiters(t *testing.T) {
	input := "( ) [ ] { } , . : ;"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: LPAREN, Value: "("},
		{Type: RPAREN, Value: ")"},
		{Type: LBRACKET, Value: "["},
		{Type: RBRACKET, Value: "]"},
		{Type: LBRACE, Value: "{"},
		{Type: RBRACE, Value: "}"},
		{Type: COMMA, Value: ","},
		{Type: DOT, Value: "."},
		{Type: COLON, Value: ":"},
		{Type: SEMICOLON, Value: ";"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerStrings(t *testing.T) {
	input := `"مرحبا"`
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: STRING, Value: `"مرحبا"`},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerMultiLineString(t *testing.T) {
	input := `"""
مرحبا
كيف حالك
أهلا وسهلا
"""`
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	if len(lexer.Tokens) < 2 {
		t.Fatal("Expected at least 2 tokens (STRING and EOF)")
	}

	if lexer.Tokens[0].Type != STRING {
		t.Errorf("Expected STRING token, got %v", lexer.Tokens[0].Type)
	}

	// Check that it's a multi-line string starting with """
	if len(lexer.Tokens[0].Value) < 3 || lexer.Tokens[0].Value[:3] != `"""` {
		t.Errorf("Expected multi-line string to start with triple quotes, got: %q", lexer.Tokens[0].Value)
	}
}

func TestLexerComplexExpression(t *testing.T) {
	input := "س=5+3*2"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: IDENT, Value: "س"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "5"},
		{Type: PLUS, Value: "+"},
		{Type: NUMBER, Value: "3"},
		{Type: MULT, Value: "*"},
		{Type: NUMBER, Value: "2"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerFunctionDefinition(t *testing.T) {
	input := "دالة تحقق(عدد):"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: FUNC, Value: "دالة"},
		{Type: IDENT, Value: "تحقق"},
		{Type: LPAREN, Value: "("},
		{Type: IDENT, Value: "عدد"},
		{Type: RPAREN, Value: ")"},
		{Type: COLON, Value: ":"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerIfStatement(t *testing.T) {
	input := "إذا س>10:"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: IF, Value: "اذا"}, // إذا normalized to اذا
		{Type: IDENT, Value: "س"},
		{Type: GREATER, Value: ">"},
		{Type: NUMBER, Value: "10"},
		{Type: COLON, Value: ":"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerLogicalOperators(t *testing.T) {
	input := "س و ص أو ليس ع"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: IDENT, Value: "س"},
		{Type: AND, Value: "و"},
		{Type: IDENT, Value: "ص"},
		{Type: OR, Value: "او"}, // أو normalized
		{Type: NOT, Value: "ليس"},
		{Type: IDENT, Value: "ع"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerKeywordNormalization(t *testing.T) {
	// Test that different Arabic character variants are normalized
	// إ, أ, ؤ, ء, ى should all normalize to ا
	tests := []struct {
		input      string
		expected   TokenType
		normalized string
	}{
		{"إذا", IF, "اذا"},
		{"اذا", IF, "اذا"},
		{"أرجع", RETURN, "ارجع"},
		{"ارجع", RETURN, "ارجع"},
		{"أو", OR, "او"},
		{"او", OR, "او"},
		{"خطأ", FALSE, "خطا"},
		{"خطا", FALSE, "خطا"},
	}

	for _, test := range tests {
		lexer := Lexer{input: []int32(test.input)}
		lexer.Tokenize()

		if len(lexer.Tokens) < 2 {
			t.Fatalf("Expected at least 2 tokens for input %q", test.input)
		}

		if lexer.Tokens[0].Type != test.expected {
			t.Errorf("Input %q: expected type %v, got %v", test.input, test.expected, lexer.Tokens[0].Type)
		}

		if lexer.Tokens[0].Value != test.normalized {
			t.Errorf("Input %q: expected normalized value %q, got %q", test.input, test.normalized, lexer.Tokens[0].Value)
		}
	}
}

func TestLexerWhileLoop(t *testing.T) {
	input := "طالما س<5:"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: WHILE, Value: "طالما"},
		{Type: IDENT, Value: "س"},
		{Type: LESS, Value: "<"},
		{Type: NUMBER, Value: "5"},
		{Type: COLON, Value: ":"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerForLoop(t *testing.T) {
	input := "لكل عنصر في قائمة:"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: FOR, Value: "لكل"},
		{Type: IDENT, Value: "عنصر"},
		{Type: IN, Value: "في"},
		{Type: IDENT, Value: "قائمة"},
		{Type: COLON, Value: ":"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}

func TestLexerArrayList(t *testing.T) {
	input := "[1, 2, 3]"
	lexer := Lexer{input: []int32(input)}
	lexer.Tokenize()

	expected := []Token{
		{Type: LBRACKET, Value: "["},
		{Type: NUMBER, Value: "1"},
		{Type: COMMA, Value: ","},
		{Type: NUMBER, Value: "2"},
		{Type: COMMA, Value: ","},
		{Type: NUMBER, Value: "3"},
		{Type: RBRACKET, Value: "]"},
		{Type: EOF, Value: ""},
	}

	if len(lexer.Tokens) != len(expected) {
		t.Fatalf("Expected %d tokens, got %d", len(expected), len(lexer.Tokens))
	}

	for i, tok := range lexer.Tokens {
		if tok.Type != expected[i].Type || tok.Value != expected[i].Value {
			t.Errorf("Token %d: expected %v, got %v", i, expected[i], tok)
		}
	}
}
