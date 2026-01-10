package lexer

type Token struct {
	Type  TokenType
	Value string
}

type TokenType int

const (
	EOF TokenType = iota
	ILLEGAL
	NEWLINE

	// Literals
	NAME
	NUMBER
	STRING
	COMMENT

	INDENT
	DEDENT

	// Keywords
	IF
	ELIF
	ELSE
	WHILE
	FOR
	IN
	RETURN
	FUNC
	BREAK
	CONTINUE
	REPEAT // new feature
	TIMES  // new feature

	// Operators

	// Arithmetic operators
	PLUS     // +
	MINUS    // -
	MULT     // *
	DIVIDE   // /
	FLOORDIV // //
	POWER    // **
	MOD      // %

	// Assignment operators
	ASSIGN        // =
	PLUS_ASSIGN   // +=
	MINUS_ASSIGN  // -=
	MULT_ASSIGN   // *=
	DIVIDE_ASSIGN // /=
	MOD_ASSIGN    // %=
	POWER_ASSIGN  // **=

	// Comparison operators
	EQ      // ==
	NEQ     // !=
	LESS    // <
	GREATER // >
	LEQ     // <=
	GEQ     // >=

	// Logical operators
	AND // و
	OR  // أو
	NOT // ليس , لا

	// Bitwise operators
	BITWISE_AND // &
	BITWISE_OR  // |
	BITWISE_XOR // ^
	BITWISE_NOT // ~
	LSHIFT      // <<
	RSHIFT      // >>

	// Delimiters
	LPAREN    // (
	RPAREN    // )
	LBRACKET  // [
	RBRACKET  // ]
	LBRACE    // {
	RBRACE    // }
	COMMA     // ,
	DOT       // .
	COLON     // :
	SEMICOLON // ;
	RETTYPE   // ->

	// Boolean literals
	TRUE  // صحيح
	FALSE // خطأ
)

var tokenTypeNames = map[TokenType]string{
	EOF:     "EOF",
	ILLEGAL: "ILLEGAL",
	NEWLINE: "NEWLINE",

	NAME:    "NAME",
	NUMBER:  "NUMBER",
	STRING:  "STRING",
	COMMENT: "COMMENT",

	INDENT: "INDENT",
	DEDENT: "DEDENT",

	IF:       "IF",
	ELIF:     "ELIF",
	ELSE:     "ELSE",
	WHILE:    "WHILE",
	FOR:      "FOR",
	IN:       "IN",
	RETURN:   "RETURN",
	FUNC:     "FUNC",
	BREAK:    "BREAK",
	CONTINUE: "CONTINUE",
	REPEAT:   "REPEAT",
	TIMES:    "TIMES",

	PLUS:     "PLUS",
	MINUS:    "MINUS",
	MULT:     "MULT",
	DIVIDE:   "DIVIDE",
	FLOORDIV: "FLOORDIV",
	POWER:    "POWER",
	MOD:      "MOD",

	ASSIGN:        "ASSIGN",
	PLUS_ASSIGN:   "PLUS_ASSIGN",
	MINUS_ASSIGN:  "MINUS_ASSIGN",
	MULT_ASSIGN:   "MULT_ASSIGN",
	DIVIDE_ASSIGN: "DIVIDE_ASSIGN",
	MOD_ASSIGN:    "MOD_ASSIGN",
	POWER_ASSIGN:  "POWER_ASSIGN",

	EQ:      "EQ",
	NEQ:     "NEQ",
	LESS:    "LESS",
	GREATER: "GREATER",
	LEQ:     "LEQ",
	GEQ:     "GEQ",

	AND: "AND",
	OR:  "OR",
	NOT: "NOT",

	BITWISE_AND: "BITWISE_AND",
	BITWISE_OR:  "BITWISE_OR",
	BITWISE_XOR: "BITWISE_XOR",
	BITWISE_NOT: "BITWISE_NOT",
	LSHIFT:      "LSHIFT",
	RSHIFT:      "RSHIFT",

	LPAREN:    "LPAREN",
	RPAREN:    "RPAREN",
	LBRACKET:  "LBRACKET",
	RBRACKET:  "RBRACKET",
	LBRACE:    "LBRACE",
	RBRACE:    "RBRACE",
	COMMA:     "COMMA",
	DOT:       "DOT",
	COLON:     "COLON",
	SEMICOLON: "SEMICOLON",
	RETTYPE:   "RETTYPE",

	TRUE:  "TRUE",
	FALSE: "FALSE",
}

func (t TokenType) String() string {
	if name, ok := tokenTypeNames[t]; ok {
		return name
	}
	return "UNKNOWN"
}

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
