package lexer

import (
	"testing"

	"github.com/TechGeeks-Club/daad/internals/lexer"
)

// Helper function to compare tokens
func compareTokens(t *testing.T, testName string, expected []lexer.Token, actual []lexer.Token) {
	if len(expected) != len(actual) {
		t.Errorf("%s: Expected %d tokens, got %d tokens", testName, len(expected), len(actual))
		t.Logf("Expected tokens:")
		for i, tok := range expected {
			t.Logf("  [%d] Type: %s, Value: %q", i, tok.Type.String(), tok.Value)
		}
		t.Logf("Actual tokens:")
		for i, tok := range actual {
			t.Logf("  [%d] Type: %s, Value: %q", i, tok.Type.String(), tok.Value)
		}
		return
	}

	for i := 0; i < len(expected); i++ {
		if expected[i].Type != actual[i].Type {
			t.Errorf("%s: Token %d type mismatch. Expected %s, got %s",
				testName, i, expected[i].Type.String(), actual[i].Type.String())
		}
		if expected[i].Value != actual[i].Value {
			t.Errorf("%s: Token %d value mismatch. Expected %q, got %q",
				testName, i, expected[i].Value, actual[i].Value)
		}
	}
}

func TestTokenizeBasicArithmetic(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/basic_arithmetic.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []lexer.Token{
		// # Test basic arithmetic operations
		{Type: lexer.COMMENT, Value: "# Test basic arithmetic operations\n"},

		// متغير = 10 + 5
		{Type: lexer.NAME, Value: "متغير"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.PLUS, Value: "+"},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة = متغير * 2 - 3// means that this line is empty e.g `     \n`
		{Type: lexer.NAME, Value: "نتيجة"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NAME, Value: "متغير"},
		{Type: lexer.MULT, Value: "*"},
		{Type: lexer.NUMBER, Value: "2"},
		{Type: lexer.MINUS, Value: "-"},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قوة = 2 ** 8
		{Type: lexer.NAME, Value: "قوة"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "2"},
		{Type: lexer.POWER, Value: "**"},
		{Type: lexer.NUMBER, Value: "8"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قسمة = 100 / 4
		{Type: lexer.NAME, Value: "قسمة"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "100"},
		{Type: lexer.DIVIDE, Value: "/"},
		{Type: lexer.NUMBER, Value: "4"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// باقي = 17 % 5
		{Type: lexer.NAME, Value: "باقي"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "17"},
		{Type: lexer.MOD, Value: "%"},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		{Type: lexer.EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeBasicArithmetic", expected, tokens)
}

func TestTokenizeControlFlow(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/control_flow.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []lexer.Token{
		// # Test control flow statements
		{Type: lexer.COMMENT, Value: "# Test control flow statements\n"},

		// اذا العدد > 10:
		{Type: lexer.IF, Value: "اذا"},
		{Type: lexer.NAME, Value: "العدد"},
		{Type: lexer.GREATER, Value: ">"},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     اطبع("كبير")
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.NAME, Value: "اطبع"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.STRING, Value: "\"كبير\""},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// والا:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.ELSE, Value: "والا"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     اطبع("صغير")
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.NAME, Value: "اطبع"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.STRING, Value: "\"صغير\""},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// طالما العداد < 100:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.WHILE, Value: "طالما"},
		{Type: lexer.NAME, Value: "العداد"},
		{Type: lexer.LESS, Value: "<"},
		{Type: lexer.NUMBER, Value: "100"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     العداد += 1
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.NAME, Value: "العداد"},
		{Type: lexer.PLUS_ASSIGN, Value: "+="},
		{Type: lexer.NUMBER, Value: "1"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// لكل عنصر في القائمة:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.FOR, Value: "لكل"},
		{Type: lexer.NAME, Value: "عنصر"},
		{Type: lexer.IN, Value: "في"},
		{Type: lexer.NAME, Value: "القائمة"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     اطبع(عنصر)
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.NAME, Value: "اطبع"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.NAME, Value: "عنصر"},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.NEWLINE, Value: "\n"},

		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeControlFlow", expected, tokens)
}

func TestTokenizeFunctions(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/functions.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []lexer.Token{
		// # Test function definitions
		{Type: lexer.COMMENT, Value: "# Test function definitions\n"},

		// دالة جمع(أ, ب) -> عدد:
		{Type: lexer.FUNC, Value: "دالة"},
		{Type: lexer.NAME, Value: "جمع"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.NAME, Value: "أ"},
		{Type: lexer.COMMA, Value: ","},
		{Type: lexer.NAME, Value: "ب"},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.RETTYPE, Value: "->"},
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     ارجع أ + ب
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.RETURN, Value: "ارجع"},
		{Type: lexer.NAME, Value: "أ"},
		{Type: lexer.PLUS, Value: "+"},
		{Type: lexer.NAME, Value: "ب"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// دالة طباعة_رسالة(نص):
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.FUNC, Value: "دالة"},
		{Type: lexer.NAME, Value: "طباعة_رسالة"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.NAME, Value: "نص"},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     اطبع(نص)
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.NAME, Value: "اطبع"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.NAME, Value: "نص"},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     ارجع
		{Type: lexer.RETURN, Value: "ارجع"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة = جمع(5, 10)
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.NAME, Value: "نتيجة"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NAME, Value: "جمع"},
		{Type: lexer.LPAREN, Value: "("},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.COMMA, Value: ","},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.RPAREN, Value: ")"},
		{Type: lexer.NEWLINE, Value: "\n"},

		{Type: lexer.EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeFunctions", expected, tokens)
}

func TestTokenizeStrings(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/strings.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []lexer.Token{
		// # Test string literals
		{Type: lexer.COMMENT, Value: "# Test string literals\n"},

		// اسم = "محمد"
		{Type: lexer.NAME, Value: "اسم"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.STRING, Value: "\"محمد\""},
		{Type: lexer.NEWLINE, Value: "\n"},

		// رسالة = 'مرحبا بك'
		{Type: lexer.NAME, Value: "رسالة"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.STRING, Value: "'مرحبا بك'"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نص_متعدد = """..."""
		{Type: lexer.NAME, Value: "نص_متعدد"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.STRING, Value: "\"\"\"\nهذا نص\nمتعدد الأسطر\nيحتوي على عدة أسطر\n\"\"\""},
		{Type: lexer.NEWLINE, Value: "\n"},

		// فارغ = ""
		{Type: lexer.NAME, Value: "فارغ"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.STRING, Value: "\"\""},
		{Type: lexer.NEWLINE, Value: "\n"},

		// فارغ2 = ''
		{Type: lexer.NAME, Value: "فارغ2"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.STRING, Value: "''"},
		{Type: lexer.NEWLINE, Value: "\n"},

		{Type: lexer.EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeStrings", expected, tokens)
}

func TestTokenizeOperators(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/operators.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []lexer.Token{
		// # Test all operators
		{Type: lexer.COMMENT, Value: "# Test all operators\n"},

		// # Comparison
		{Type: lexer.COMMENT, Value: "# Comparison\n"},

		// نتيجة1 = 5 == 5
		{Type: lexer.NAME, Value: "نتيجة1"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.EQ, Value: "=="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة2 = 10 != 5
		{Type: lexer.NAME, Value: "نتيجة2"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.NEQ, Value: "!="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة3 = 7 < 10
		{Type: lexer.NAME, Value: "نتيجة3"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "7"},
		{Type: lexer.LESS, Value: "<"},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة4 = 15 > 10
		{Type: lexer.NAME, Value: "نتيجة4"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "15"},
		{Type: lexer.GREATER, Value: ">"},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة5 = 5 <= 5
		{Type: lexer.NAME, Value: "نتيجة5"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.LEQ, Value: "<="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// نتيجة6 = 20 >= 15
		{Type: lexer.NAME, Value: "نتيجة6"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "20"},
		{Type: lexer.GEQ, Value: ">="},
		{Type: lexer.NUMBER, Value: "15"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// # Logical
		{Type: lexer.COMMENT, Value: "# Logical\n"},

		// صحيح و خطأ
		{Type: lexer.TRUE, Value: "صحيح"},
		{Type: lexer.AND, Value: "و"},
		{Type: lexer.FALSE, Value: "خطأ"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// صحيح او خطأ
		{Type: lexer.TRUE, Value: "صحيح"},
		{Type: lexer.OR, Value: "او"},
		{Type: lexer.FALSE, Value: "خطأ"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// ليس صحيح
		{Type: lexer.NOT, Value: "ليس"},
		{Type: lexer.TRUE, Value: "صحيح"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// # Bitwise
		{Type: lexer.COMMENT, Value: "# Bitwise\n"},

		// قيمة1 = 5 & 3
		{Type: lexer.NAME, Value: "قيمة1"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.BITWISE_AND, Value: "&"},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قيمة2 = 5 | 3
		{Type: lexer.NAME, Value: "قيمة2"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.BITWISE_OR, Value: "|"},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قيمة3 = 5 ^ 3
		{Type: lexer.NAME, Value: "قيمة3"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.BITWISE_XOR, Value: "^"},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قيمة4 = ~5
		{Type: lexer.NAME, Value: "قيمة4"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.BITWISE_NOT, Value: "~"},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قيمة5 = 2 << 3
		{Type: lexer.NAME, Value: "قيمة5"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "2"},
		{Type: lexer.LSHIFT, Value: "<<"},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// قيمة6 = 16 >> 2
		{Type: lexer.NAME, Value: "قيمة6"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "16"},
		{Type: lexer.RSHIFT, Value: ">>"},
		{Type: lexer.NUMBER, Value: "2"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// # Assignment
		{Type: lexer.COMMENT, Value: "# Assignment\n"},

		// عدد = 10
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.ASSIGN, Value: "="},
		{Type: lexer.NUMBER, Value: "10"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// عدد += 5
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.PLUS_ASSIGN, Value: "+="},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// عدد -= 3
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.MINUS_ASSIGN, Value: "-="},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// عدد *= 2
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.MULT_ASSIGN, Value: "*="},
		{Type: lexer.NUMBER, Value: "2"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// عدد /= 4
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.DIVIDE_ASSIGN, Value: "/="},
		{Type: lexer.NUMBER, Value: "4"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// عدد %= 3
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.MOD_ASSIGN, Value: "%="},
		{Type: lexer.NUMBER, Value: "3"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// عدد **= 2
		{Type: lexer.NAME, Value: "عدد"},
		{Type: lexer.POWER_ASSIGN, Value: "**="},
		{Type: lexer.NUMBER, Value: "2"},
		{Type: lexer.NEWLINE, Value: "\n"},

		{Type: lexer.EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeOperators", expected, tokens)
}

func TestTokenizeKeywords(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/keywords.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []lexer.Token{
		// # Test all keywords
		{Type: lexer.COMMENT, Value: "# Test all keywords\n"},

		// إذا صحيح:
		{Type: lexer.IF, Value: "إذا"},
		{Type: lexer.TRUE, Value: "صحيح"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     ارجع 1
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.RETURN, Value: "ارجع"},
		{Type: lexer.NUMBER, Value: "1"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// وإذا خطأ:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.ELIF, Value: "وإذا"},
		{Type: lexer.FALSE, Value: "خطأ"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     ارجع 0
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.RETURN, Value: "ارجع"},
		{Type: lexer.NUMBER, Value: "0"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// وإلا:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.ELSE, Value: "وإلا"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     ارجع -1
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.RETURN, Value: "ارجع"},
		{Type: lexer.MINUS, Value: "-"},
		{Type: lexer.NUMBER, Value: "1"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// طالما خطأ:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.WHILE, Value: "طالما"},
		{Type: lexer.FALSE, Value: "خطأ"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     اخرج
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.BREAK, Value: "اخرج"},
		{Type: lexer.NEWLINE, Value: "\n"},

		// كرر 5 مرات:
		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.REPEAT, Value: "كرر"},
		{Type: lexer.NUMBER, Value: "5"},
		{Type: lexer.TIMES, Value: "مرات"},
		{Type: lexer.COLON, Value: ":"},
		{Type: lexer.NEWLINE, Value: "\n"},

		//     تابع
		{Type: lexer.INDENT, Value: ""},
		{Type: lexer.CONTINUE, Value: "تابع"},
		{Type: lexer.NEWLINE, Value: "\n"},

		{Type: lexer.DEDENT, Value: ""},
		{Type: lexer.EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeKeywords", expected, tokens)
}
