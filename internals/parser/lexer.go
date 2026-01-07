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
		next := l.advance()
		if next == '=' {
			l.advance()
			return Token{Type: EQ, Value: "=="}
		}
		l.goback()
		return Token{Type: ASSIGN, Value: "="}
	case '!':
		next := l.advance()
		if next == '=' {
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
