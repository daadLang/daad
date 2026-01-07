package parser

type TokenType int

type Token struct {
	Type  TokenType
	Value string
}

const (
	EOF TokenType = iota
	ILLEGAL
	NEWLINE // \n

	// Literals
	IDENT  // variable names
	NUMBER // 12345
	STRING // "hello"

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

	// Operators / delimiters
	PLUS   // +
	MINUS  // -
	MULT   // *
	DIVIDE // /
	POWER  // ^
	MOD    // %
	LPAREN // (
	RPAREN // )
	COLON  // :
	EQ     // ==
	NEQ    // !=
	ASSIGN // =

	NOT // ليس
	AND // و
	OR  //  او

	// bool
	TRUE
	FALSE
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
