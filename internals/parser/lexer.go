package parser

import (
	"bufio"
	"unicode"
)

// ? rune is int32 representing a Unicode code char
type Lexer struct {
	reader      *bufio.Reader
	currentChar rune
	hasChar     bool
	buffer      []rune
	spacesCount int
	openBraces  int // we need to track any type of open braces `({[` to handle spaces , (in case of openBraces > 0 we ignore spaces and newlines)
	isNewLine   bool
	Tokens      []Token
}

func NewLexer(r *bufio.Reader) *Lexer {
	return &Lexer{
		reader:      r,
		currentChar: 0,
		hasChar:     false,
		buffer:      make([]rune, 0, 64),
		spacesCount: 0,
		openBraces:  0,
		isNewLine:   true,
		Tokens:      make([]Token, 0),
	}
}

func (l *Lexer) peek() rune {
	if !l.hasChar {
		r, _, err := l.reader.ReadRune()
		if err != nil {
			return 0
		}
		l.currentChar = r
		l.hasChar = true
	}
	return l.currentChar
}

func (l *Lexer) advance() rune {
	r := l.peek()
	if r != 0 {
		l.buffer = append(l.buffer, r)
		l.hasChar = false
	}
	return r
}

func (l *Lexer) clearBuffer() {
	l.buffer = l.buffer[:0]
}

func (l *Lexer) getBuffer() string {
	return string(l.buffer)
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.peek()) {
		l.advance()
	}
}

// this function will replace skipWhitespace
func (l *Lexer) handleSpaces() (Token, int) {
	if l.openBraces > 0 {
		l.skipWhitespace()
		return Token{}, 0
	} else {
		if l.isNewLine {
			for {
				peek := l.peek()
				if peek == ' ' {
					l.spacesCount++
					l.advance()
				} else if peek == '\t' {
					l.spacesCount += 4
					l.advance()
				} else if peek == '\n' || peek == '\r' {
					// skip in case of empty line e.g `     \n`
					l.spacesCount = 0
					l.advance()
					// stay in isNewLine state and continue to next line
					continue
				} else {
					l.isNewLine = false
					if l.spacesCount > 0 {
						spaces := make([]rune, l.spacesCount)
						for i := range spaces {
							spaces[i] = ' '
						}
						return Token{Type: NAME, Value: string(spaces)}, 1
					}
					break
				}
			}
		} else {
			for {
				peek := l.peek()
				if peek == ' ' || peek == '\t' {
					l.advance()
				} else if peek == '\n' || peek == '\r' {
					l.isNewLine = true
					l.advance()
					return Token{Type: NEWLINE, Value: "\n"}, 1
				} else {
					break
				}
			}
		}

	}
	return Token{}, 0

}

func (l *Lexer) readName() string {
	l.clearBuffer()
	for {
		r := l.peek()
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
			l.advance()
		} else {
			break
		}
	}
	return l.getBuffer()
}

func (l *Lexer) readString() string {
	l.clearBuffer()
	str_delimiter := l.peek()
	l.advance() // skip the first one

	// Read multiline string
	if l.peek() == str_delimiter {
		l.advance() // skip 2
		if l.peek() == str_delimiter {
			l.advance() // skip 3
			for {
				r := l.peek()
				if r == 0 { // EOF ?
					break
				}
				if r == str_delimiter {
					l.advance()
					if l.peek() == str_delimiter {
						l.advance()
						if l.peek() == str_delimiter {
							l.advance()
							break
						}
					}
				} else {
					l.advance()
				}
			}
			return l.getBuffer()
		}
		// Empty string: "" or ''
		return string(str_delimiter) + string(str_delimiter)
	}

	// read until """ (or ''')
	for {
		r := l.peek()
		if r == 0 { // EOF
			break
		}
		if r == str_delimiter {
			l.advance()
			break
		}
		if r == '\\' {
			l.advance()
			if l.peek() != 0 {
				l.advance()
			}
		} else {
			l.advance()
		}
	}
	return l.getBuffer()
}

func (l *Lexer) readComment() string {
	l.clearBuffer()
	for {
		r := l.peek()
		if r == 0 || r == '\n' || r == '\r' {
			// Consume the newline at the end of the comment
			if r == '\n' || r == '\r' {
				l.advance()
			}
			break
		}
		l.advance()
	}
	return l.getBuffer()
}

func (l *Lexer) readNumber() string {
	l.clearBuffer()
	decimalPointSeen := false
	for {
		peek := l.peek()
		if unicode.IsDigit(peek) {
			l.advance()
		} else if peek == '.' && !decimalPointSeen {
			decimalPointSeen = true
			l.advance()
		} else {
			break
		}
	}
	return l.getBuffer()
}

// أ,إو ...  to ا
func (l *Lexer) simplifyKeyword(keyword []rune) string {
	for index, char := range keyword {
		if char == 'أ' || char == 'إ' || char == 'ؤ' || char == 'ء' || char == 'ى' {
			keyword[index] = 'ا'
		}
	}
	return string(keyword)
}

func (l *Lexer) NextToken() Token {
	tok, handled := l.handleSpaces()
	if handled > 0 {
		l.spacesCount = 0
		return tok
	}
	r := l.peek()
	if r == 0 {
		return Token{Type: EOF, Value: ""}
	}

	if unicode.IsLetter(r) {
		value := l.readName()
		simplified := l.simplifyKeyword([]rune(value))
		if _, ok := keywords[simplified]; ok {

			return Token{Type: keywords[simplified], Value: value}
		}
		return Token{Type: NAME, Value: value}
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
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: PLUS_ASSIGN, Value: "+="}
		}
		return Token{Type: PLUS, Value: "+"}
	case '-':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: MINUS_ASSIGN, Value: "-="}
		}
		if l.peek() == '>' {
			l.advance()
			return Token{Type: RETTYPE, Value: "->"}
		}
		return Token{Type: MINUS, Value: "-"}
	case '*':
		l.advance()
		if l.peek() == '*' {
			l.advance()
			if l.peek() == '=' {
				l.advance()
				return Token{Type: POWER_ASSIGN, Value: "**="}
			}
			return Token{Type: POWER, Value: "**"}
		}
		if l.peek() == '=' {
			l.advance()
			return Token{Type: MULT_ASSIGN, Value: "*="}
		}
		return Token{Type: MULT, Value: "*"}
	case '/':
		l.advance()
		if l.peek() == '/' {
			l.advance()
			return Token{Type: FLOORDIV, Value: "//"}
		}
		if l.peek() == '=' {
			l.advance()
			return Token{Type: DIVIDE_ASSIGN, Value: "/="}
		}
		return Token{Type: DIVIDE, Value: "/"}
	case '%':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: MOD_ASSIGN, Value: "%="}
		}
		return Token{Type: MOD, Value: "%"}

	// Comparison operators
	case '=':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: EQ, Value: "=="}
		}
		return Token{Type: ASSIGN, Value: "="}
	case '!':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: NEQ, Value: "!="}
		}
		return Token{Type: ILLEGAL, Value: "!"}
	case '<':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: LEQ, Value: "<="}
		}
		if l.peek() == '<' {
			l.advance()
			return Token{Type: LSHIFT, Value: "<<"}
		}
		return Token{Type: LESS, Value: "<"}
	case '>':
		l.advance()
		if l.peek() == '=' {
			l.advance()
			return Token{Type: GEQ, Value: ">="}
		}
		if l.peek() == '>' {
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
		l.openBraces++
		return Token{Type: LPAREN, Value: "("}
	case ')':
		l.advance()
		l.openBraces--
		return Token{Type: RPAREN, Value: ")"}
	case '[':
		l.advance()
		l.openBraces++
		return Token{Type: LBRACKET, Value: "["}
	case ']':
		l.advance()
		l.openBraces--
		return Token{Type: RBRACKET, Value: "]"}
	case '{':
		l.advance()
		l.openBraces++
		return Token{Type: LBRACE, Value: "{"}
	case '}':
		l.advance()
		l.openBraces--
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
	case '#':
		comment := l.readComment()
		l.isNewLine = true
		return Token{Type: COMMENT, Value: comment}
	}

	l.advance()
	return Token{Type: ILLEGAL, Value: string(r)}
}
