package main

import (
	"fmt"
	"unicode"
)

type Token struct {
	Type  TokenType
	Value string
}

type TokenType string

const (
	EOF     TokenType = "EOF"
	ILLEGAL TokenType = "ILLEGAL"
	NEWLINE TokenType = "NEWLINE"

	// Literals
	IDENT  TokenType = "IDENT"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"

	// Keywords
	IF       TokenType = "IF"
	ELIF     TokenType = "ELIF"
	ELSE     TokenType = "ELSE"
	WHILE    TokenType = "WHILE"
	FOR      TokenType = "FOR"
	IN       TokenType = "IN"
	RETURN   TokenType = "RETURN"
	FUNC     TokenType = "FUNC"
	BREAK    TokenType = "BREAK"
	CONTINUE TokenType = "CONTINUE"

	// Operators / delimiters
	PLUS   TokenType = "PLUS"
	MINUS  TokenType = "MINUS"
	MULT   TokenType = "MULT"
	DIVIDE TokenType = "DIVIDE"
	POWER  TokenType = "POWER"
	MOD    TokenType = "MOD"
	LPAREN TokenType = "LPAREN"
	RPAREN TokenType = "RPAREN"
	COLON  TokenType = "COLON"
	EQ     TokenType = "EQ"
	NEQ    TokenType = "NEQ"

	// Boolean literals
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"
)

var keywords = map[string]TokenType{
	"إذا":   IF,
	"وإذا":  ELIF,
	"وإلا":  ELSE,
	"طالما": WHILE,
	"لكل":   FOR,
	"في":    IN,
	"أرجع":  RETURN,
	"دالة":  FUNC,
	"اخرج":  BREAK,
	"تابع":  CONTINUE,
	"صحيح":  TRUE,
	"خطأ":   FALSE,
}

func tokenizer_test() {
	source := `دالة تحقق(عدد):
    إذا عدد % 2 == 0:
        أرجع "زوجي"
    وإلا:
        أرجع "فردي"`

	lexer := Lexer{input: []int32(source)}
	for {
		token := lexer.NextToken()
		fmt.Printf("%d -> %q\n", token.Type, token.Value)
		if token.Type == EOF {
			break
		}
	}
}

type Lexer struct {
	input  []int32
	pos    int
	Tokens []Token
}

func (l *Lexer) peek() int32 {
	if l.pos >= len(l.input) {
		return 0
	}
	return l.input[l.pos]
}

func (l *Lexer) advance() int32 {
	r := l.peek()
	l.pos++
	return r
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.advance()
	}
}

func (l *Lexer) readIdentifier() string {
	start := l.pos
	for {
		r := l.peek()
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			l.advance()
		} else {
			break
		}
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) readNumber() string {
	start := l.pos
	for unicode.IsDigit(l.peek()) {
		l.advance()
	}
	return string(l.input[start:l.pos])
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()
	r := l.peek()
	if r == 0 {
		return Token{Type: EOF, Value: ""}
	}

	if unicode.IsLetter(r) {
		value := l.readIdentifier()
		if _, ok := keywords[value]; ok {
			return Token{Type: keywords[value], Value: value}
		}
		return Token{Type: IDENT, Value: value}
	}

	if unicode.IsDigit(r) {
		return Token{Type: NUMBER, Value: l.readNumber()}
	}

	switch r {
	case '+':
		l.advance()
		return Token{Type: PLUS, Value: "+"}
	case '-':
		l.advance()
		return Token{Type: MINUS, Value: "-"}
	case '*':
		l.advance()
		return Token{Type: MULT, Value: "*"}
	case '/':
		l.advance()
		return Token{Type: DIVIDE, Value: "/"}
	case '^':
		l.advance()
		return Token{Type: POWER, Value: "^"}
	case '%':
		l.advance()
		return Token{Type: MOD, Value: "%"}
	case '(':
		l.advance()
		return Token{Type: LPAREN, Value: "("}
	case ')':
		l.advance()
		return Token{Type: RPAREN, Value: ")"}
	case ':':
		l.advance()
		return Token{Type: COLON, Value: ":"}
	case '=':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: EQ, Value: "=="}
		}
		return Token{Type: ILLEGAL, Value: "="}
	case '!':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: NEQ, Value: "!="}
		}
		return Token{Type: ILLEGAL, Value: "!"}
	case '\n':
		l.advance()
		return Token{Type: NEWLINE, Value: "\n"}

	}

	l.advance()
	return Token{Type: ILLEGAL, Value: string(r)}
}

func (l *Lexer) Tokenize() []Token {
	for {
		tok := l.NextToken()
		l.Tokens = append(l.Tokens, tok)
		if tok.Type == EOF {
			break
		}
	}
	return l.Tokens
}

func main() {
	source := `دالة تحقق(عدد):
    إذا عدد % 2 == 0:
        أرجع "زوجي"
    وإلا:
        أرجع "فردي"`

	lexer := Lexer{input: []int32(source)}
	for {
		token := lexer.NextToken()
		fmt.Printf("%s -> %q\n", token.Type, token.Value)
		if token.Type == EOF {
			break
		}
	}
}
