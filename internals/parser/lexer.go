package parser

import "unicode"

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

func (l *Lexer) goback() int32 {
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

// TODO : complete this69
func (l *Lexer) readString() string {
	start := l.pos
	str_delemiter := l.peek()
	next := l.advance() // ? the first "
	// multi line string
	if next == str_delemiter {
		next := l.advance()
		if next == str_delemiter {
			return l.readMultiLineString()
		}
		return string(str_delemiter) + string(str_delemiter)
	}

	// one line string
	for {
		r := l.peek()
		if r == str_delemiter {
			l.advance()
			break
		} else {
			l.advance()
		}
	}
	return string(l.input[start:l.pos])
}

// TODO:
func (l *Lexer) readMultiLineString() string {
	return ""
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

	if r == '"' || r == '\'' {
		return Token{Type: STRING, Value: l.readString()}

	}

	switch r {
	// Arithmetic operators
	case '+':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: PLUS_ASSIGN, Value: "+="}
		}
		return Token{Type: PLUS, Value: "+"}
	case '-':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: MINUS_ASSIGN, Value: "-="}
		}
		return Token{Type: MINUS, Value: "-"}
	case '*':
		next := l.advance()
		if next == '*' {
			l.advance()
			if l.peek() == '=' {
				l.advance()
				return Token{Type: POWER_ASSIGN, Value: "**="}
			}
			return Token{Type: POWER, Value: "**"}
		}
		if next == '=' {
			l.advance()
			return Token{Type: MULT_ASSIGN, Value: "*="}
		}
		return Token{Type: MULT, Value: "*"}
	case '/':
		next := l.advance()
		if next == '/' {
			l.advance()
			return Token{Type: FLOORDIV, Value: "//"}
		}
		if next == '=' {
			l.advance()
			return Token{Type: DIVIDE_ASSIGN, Value: "/="}
		}
		return Token{Type: DIVIDE, Value: "/"}
	case '%':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: MOD_ASSIGN, Value: "%="}
		}
		return Token{Type: MOD, Value: "%"}

	// Comparison operators
	case '=':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: EQ, Value: "=="}
		}
		return Token{Type: ASSIGN, Value: "="}
	case '!':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: NEQ, Value: "!="}
		}
		return Token{Type: ILLEGAL, Value: "!"}
	case '<':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: LEQ, Value: "<="}
		}
		if next == '<' {
			l.advance()
			return Token{Type: LSHIFT, Value: "<<"}
		}
		return Token{Type: LESS, Value: "<"}
	case '>':
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: GEQ, Value: ">="}
		}
		if next == '>' {
			l.advance()
			return Token{Type: RSHIFT, Value: ">>"}
		}
		return Token{Type: GREATER, Value: ">"}

	// Bitwise operators
	case '&':
		l.advance()
		return Token{Type: BITWISE_AND, Value: "&"}
	case '|':
		l.advance()
		return Token{Type: BITWISE_OR, Value: "|"}
	case '^':
		l.advance()
		return Token{Type: BITWISE_XOR, Value: "^"}
	case '~':
		l.advance()
		return Token{Type: BITWISE_NOT, Value: "~"}

	// Delimiters
	case '(':
		l.advance()
		return Token{Type: LPAREN, Value: "("}
	case ')':
		l.advance()
		return Token{Type: RPAREN, Value: ")"}
	case '[':
		l.advance()
		return Token{Type: LBRACKET, Value: "["}
	case ']':
		l.advance()
		return Token{Type: RBRACKET, Value: "]"}
	case '{':
		l.advance()
		return Token{Type: LBRACE, Value: "{"}
	case '}':
		l.advance()
		return Token{Type: RBRACE, Value: "}"}
	case ',':
		l.advance()
		return Token{Type: COMMA, Value: ","}
	case '.':
		l.advance()
		return Token{Type: DOT, Value: "."}
	case ':':
		l.advance()
		return Token{Type: COLON, Value: ":"}
	case ';':
		l.advance()
		return Token{Type: SEMICOLON, Value: ";"}
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
