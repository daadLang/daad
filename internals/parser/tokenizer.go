package parser

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

	// Boolean literals
	TRUE  TokenType = "TRUE"  // صحيح
	FALSE TokenType = "FALSE" // خطأ
)

var keywords = map[string]TokenType{
	"إذا": IF,
	"لو":  IF,

	"وإذا": ELIF,
	"ولو":  ELIF,

	"وإلا": ELSE,

	"طالما": WHILE,
	"مادام": WHILE,

	"لكل": FOR,
	"في":  IN,

	//? new feature: reapeat N times equivalent to for _ in range(N)
	"كرر":  REPEAT, // new feature
	"مرات": TIMES,  // new feature

	"أرجع": RETURN,
	"دالة": FUNC,
	"اخرج": BREAK,
	"تابع": CONTINUE,
	"صحيح": TRUE,
	"خطأ":  FALSE,
	"و":    AND,
	"أو":   OR,

	"ليس": NOT,
	"لا":  NOT,
}
