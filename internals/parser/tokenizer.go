package parser

import (
	"bufio"
	"os"
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
	NAME    TokenType = "NAME"
	NUMBER  TokenType = "NUMBER"
	STRING  TokenType = "STRING"
	COMMENT TokenType = "COMMENT"

	INDENT TokenType = "INDENT"
	DEDENT TokenType = "DEDENT"

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
	REPEAT   TokenType = "REPEAT" // new feature
	TIMES    TokenType = "TIMES"  // new feature

	// Operators

	// Arithmetic operators
	PLUS     TokenType = "PLUS"     // +
	MINUS    TokenType = "MINUS"    // -
	MULT     TokenType = "MULT"     // *
	DIVIDE   TokenType = "DIVIDE"   // /
	FLOORDIV TokenType = "FLOORDIV" // //
	POWER    TokenType = "POWER"    // **
	MOD      TokenType = "MOD"      // %

	// Assignment operators
	ASSIGN        TokenType = "ASSIGN"        // =
	PLUS_ASSIGN   TokenType = "PLUS_ASSIGN"   // +=
	MINUS_ASSIGN  TokenType = "MINUS_ASSIGN"  // -=
	MULT_ASSIGN   TokenType = "MULT_ASSIGN"   // *=
	DIVIDE_ASSIGN TokenType = "DIVIDE_ASSIGN" // /=
	MOD_ASSIGN    TokenType = "MOD_ASSIGN"    // %=
	POWER_ASSIGN  TokenType = "POWER_ASSIGN"  // **=

	// Comparison operators
	EQ      TokenType = "EQ"      // ==
	NEQ     TokenType = "NEQ"     // !=
	LESS    TokenType = "LESS"    // <
	GREATER TokenType = "GREATER" // >
	LEQ     TokenType = "LEQ"     // <=
	GEQ     TokenType = "GEQ"     // >=

	// Logical operators
	AND TokenType = "AND" // و
	OR  TokenType = "OR"  // أو
	NOT TokenType = "NOT" // ليس , لا

	// Bitwise operators
	BITWISE_AND TokenType = "BITWISE_AND" // &
	BITWISE_OR  TokenType = "BITWISE_OR"  // |
	BITWISE_XOR TokenType = "BITWISE_XOR" // ^
	BITWISE_NOT TokenType = "BITWISE_NOT" // ~
	LSHIFT      TokenType = "LSHIFT"      // <<
	RSHIFT      TokenType = "RSHIFT"      // >>

	// Delimiters
	LPAREN    TokenType = "LPAREN"    // (
	RPAREN    TokenType = "RPAREN"    // )
	LBRACKET  TokenType = "LBRACKET"  // [
	RBRACKET  TokenType = "RBRACKET"  // ]
	LBRACE    TokenType = "LBRACE"    // {
	RBRACE    TokenType = "RBRACE"    // }
	COMMA     TokenType = "COMMA"     // ,
	DOT       TokenType = "DOT"       // .
	COLON     TokenType = "COLON"     // :
	SEMICOLON TokenType = "SEMICOLON" // ;
	RETTYPE   TokenType = "RETTYPE"   // ->

	// Boolean literals
	TRUE  TokenType = "TRUE"  // صحيح
	FALSE TokenType = "FALSE" // خطأ
)

// Keywords map - simplified forms (after normalization by simplifyKeyword)
// أ,إ,ؤ,ء,ى are normalized to ا
var keywords = map[string]TokenType{
	"اذا": IF, // إذا → اذا
	"لو":  IF,

	"واذا": ELIF, // وإذا → واذا
	"ولو":  ELIF,

	"والا": ELSE, // وإلا → والا

	"طالما": WHILE,
	"مادام": WHILE,

	"لكل": FOR,
	"في":  IN,

	"كرر":  REPEAT, // repeat N times
	"مرات": TIMES,

	"ارجع": RETURN, // أرجع → ارجع
	"دالة": FUNC,
	"اخرج": BREAK,
	"تابع": CONTINUE,
	"صحيح": TRUE,
	"خطا":  FALSE, // خطأ → خطا
	"و":    AND,
	"او":   OR, // أو → او

	"ليس": NOT,
	"لا":  NOT,
}

func Tokenize(filePath string) ([]Token, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)

	lexer := &Lexer{
		reader:  reader,
		hasChar: false,
		buffer:  make([]rune, 0, 256),
		Tokens:  []Token{},
	}

	for {
		tok := lexer.NextToken()
		lexer.Tokens = append(lexer.Tokens, tok)
		if tok.Type == EOF {
			break
		}
	}

	return lexer.Tokens, nil
}
