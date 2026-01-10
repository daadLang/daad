package lexer

import (
	"testing"
)

// Helper function to compare tokens
func compareTokens(t *testing.T, testName string, expected []Token, actual []Token) {
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
	tokens, err := Tokenize("../../tests/basic_arithmetic.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []Token{
		// # Test basic arithmetic operations
		{Type: COMMENT, Value: "# Test basic arithmetic operations\n"},

		// متغير = 10 + 5
		{Type: NAME, Value: "متغير"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "10"},
		{Type: PLUS, Value: "+"},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة = متغير * 2 - 3// means that this line is empty e.g `     \n`
		{Type: NAME, Value: "نتيجة"},
		{Type: ASSIGN, Value: "="},
		{Type: NAME, Value: "متغير"},
		{Type: MULT, Value: "*"},
		{Type: NUMBER, Value: "2"},
		{Type: MINUS, Value: "-"},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// قوة = 2 ** 8
		{Type: NAME, Value: "قوة"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "2"},
		{Type: POWER, Value: "**"},
		{Type: NUMBER, Value: "8"},
		{Type: NEWLINE, Value: "\n"},

		// قسمة = 100 / 4
		{Type: NAME, Value: "قسمة"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "100"},
		{Type: DIVIDE, Value: "/"},
		{Type: NUMBER, Value: "4"},
		{Type: NEWLINE, Value: "\n"},

		// باقي = 17 % 5
		{Type: NAME, Value: "باقي"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "17"},
		{Type: MOD, Value: "%"},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		{Type: EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeBasicArithmetic", expected, tokens)
}

func TestTokenizeControlFlow(t *testing.T) {
	tokens, err := Tokenize("../../tests/control_flow.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []Token{
		// # Test control flow statements
		{Type: COMMENT, Value: "# Test control flow statements\n"},

		// اذا العدد > 10:
		{Type: IF, Value: "اذا"},
		{Type: NAME, Value: "العدد"},
		{Type: GREATER, Value: ">"},
		{Type: NUMBER, Value: "10"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     اطبع("كبير")
		{Type: INDENT, Value: ""},
		{Type: NAME, Value: "اطبع"},
		{Type: LPAREN, Value: "("},
		{Type: STRING, Value: "\"كبير\""},
		{Type: RPAREN, Value: ")"},
		{Type: NEWLINE, Value: "\n"},

		// والا:
		{Type: DEDENT, Value: ""},
		{Type: ELSE, Value: "والا"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     اطبع("صغير")
		{Type: INDENT, Value: ""},
		{Type: NAME, Value: "اطبع"},
		{Type: LPAREN, Value: "("},
		{Type: STRING, Value: "\"صغير\""},
		{Type: RPAREN, Value: ")"},
		{Type: NEWLINE, Value: "\n"},

		// طالما العداد < 100:
		{Type: DEDENT, Value: ""},
		{Type: WHILE, Value: "طالما"},
		{Type: NAME, Value: "العداد"},
		{Type: LESS, Value: "<"},
		{Type: NUMBER, Value: "100"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     العداد += 1
		{Type: INDENT, Value: ""},
		{Type: NAME, Value: "العداد"},
		{Type: PLUS_ASSIGN, Value: "+="},
		{Type: NUMBER, Value: "1"},
		{Type: NEWLINE, Value: "\n"},

		// لكل عنصر في القائمة:
		{Type: DEDENT, Value: ""},
		{Type: FOR, Value: "لكل"},
		{Type: NAME, Value: "عنصر"},
		{Type: IN, Value: "في"},
		{Type: NAME, Value: "القائمة"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     اطبع(عنصر)
		{Type: INDENT, Value: ""},
		{Type: NAME, Value: "اطبع"},
		{Type: LPAREN, Value: "("},
		{Type: NAME, Value: "عنصر"},
		{Type: RPAREN, Value: ")"},
		{Type: NEWLINE, Value: "\n"},

		{Type: DEDENT, Value: ""},
		{Type: EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeControlFlow", expected, tokens)
}

func TestTokenizeFunctions(t *testing.T) {
	tokens, err := Tokenize("../../tests/functions.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []Token{
		// # Test function definitions
		{Type: COMMENT, Value: "# Test function definitions\n"},

		// دالة جمع(أ, ب) -> عدد:
		{Type: FUNC, Value: "دالة"},
		{Type: NAME, Value: "جمع"},
		{Type: LPAREN, Value: "("},
		{Type: NAME, Value: "أ"},
		{Type: COMMA, Value: ","},
		{Type: NAME, Value: "ب"},
		{Type: RPAREN, Value: ")"},
		{Type: RETTYPE, Value: "->"},
		{Type: NAME, Value: "عدد"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     ارجع أ + ب
		{Type: INDENT, Value: ""},
		{Type: RETURN, Value: "ارجع"},
		{Type: NAME, Value: "أ"},
		{Type: PLUS, Value: "+"},
		{Type: NAME, Value: "ب"},
		{Type: NEWLINE, Value: "\n"},

		// دالة طباعة_رسالة(نص):
		{Type: DEDENT, Value: ""},
		{Type: FUNC, Value: "دالة"},
		{Type: NAME, Value: "طباعة_رسالة"},
		{Type: LPAREN, Value: "("},
		{Type: NAME, Value: "نص"},
		{Type: RPAREN, Value: ")"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     اطبع(نص)
		{Type: INDENT, Value: ""},
		{Type: NAME, Value: "اطبع"},
		{Type: LPAREN, Value: "("},
		{Type: NAME, Value: "نص"},
		{Type: RPAREN, Value: ")"},
		{Type: NEWLINE, Value: "\n"},

		//     ارجع
		{Type: RETURN, Value: "ارجع"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة = جمع(5, 10)
		{Type: DEDENT, Value: ""},
		{Type: NAME, Value: "نتيجة"},
		{Type: ASSIGN, Value: "="},
		{Type: NAME, Value: "جمع"},
		{Type: LPAREN, Value: "("},
		{Type: NUMBER, Value: "5"},
		{Type: COMMA, Value: ","},
		{Type: NUMBER, Value: "10"},
		{Type: RPAREN, Value: ")"},
		{Type: NEWLINE, Value: "\n"},

		{Type: EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeFunctions", expected, tokens)
}

func TestTokenizeStrings(t *testing.T) {
	tokens, err := Tokenize("../../tests/strings.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []Token{
		// # Test string literals
		{Type: COMMENT, Value: "# Test string literals\n"},

		// اسم = "محمد"
		{Type: NAME, Value: "اسم"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "\"محمد\""},
		{Type: NEWLINE, Value: "\n"},

		// رسالة = 'مرحبا بك'
		{Type: NAME, Value: "رسالة"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "'مرحبا بك'"},
		{Type: NEWLINE, Value: "\n"},

		// نص_متعدد = """..."""
		{Type: NAME, Value: "نص_متعدد"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "\"\"\"\nهذا نص\nمتعدد الأسطر\nيحتوي على عدة أسطر\n\"\"\""},
		{Type: NEWLINE, Value: "\n"},

		// فارغ = ""
		{Type: NAME, Value: "فارغ"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "\"\""},
		{Type: NEWLINE, Value: "\n"},

		// فارغ2 = ''
		{Type: NAME, Value: "فارغ2"},
		{Type: ASSIGN, Value: "="},
		{Type: STRING, Value: "''"},
		{Type: NEWLINE, Value: "\n"},

		{Type: EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeStrings", expected, tokens)
}

func TestTokenizeOperators(t *testing.T) {
	tokens, err := Tokenize("../../tests/operators.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []Token{
		// # Test all operators
		{Type: COMMENT, Value: "# Test all operators\n"},

		// # Comparison
		{Type: COMMENT, Value: "# Comparison\n"},

		// نتيجة1 = 5 == 5
		{Type: NAME, Value: "نتيجة1"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "5"},
		{Type: EQ, Value: "=="},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة2 = 10 != 5
		{Type: NAME, Value: "نتيجة2"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "10"},
		{Type: NEQ, Value: "!="},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة3 = 7 < 10
		{Type: NAME, Value: "نتيجة3"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "7"},
		{Type: LESS, Value: "<"},
		{Type: NUMBER, Value: "10"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة4 = 15 > 10
		{Type: NAME, Value: "نتيجة4"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "15"},
		{Type: GREATER, Value: ">"},
		{Type: NUMBER, Value: "10"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة5 = 5 <= 5
		{Type: NAME, Value: "نتيجة5"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "5"},
		{Type: LEQ, Value: "<="},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		// نتيجة6 = 20 >= 15
		{Type: NAME, Value: "نتيجة6"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "20"},
		{Type: GEQ, Value: ">="},
		{Type: NUMBER, Value: "15"},
		{Type: NEWLINE, Value: "\n"},

		// # Logical
		{Type: COMMENT, Value: "# Logical\n"},

		// صحيح و خطأ
		{Type: TRUE, Value: "صحيح"},
		{Type: AND, Value: "و"},
		{Type: FALSE, Value: "خطأ"},
		{Type: NEWLINE, Value: "\n"},

		// صحيح او خطأ
		{Type: TRUE, Value: "صحيح"},
		{Type: OR, Value: "او"},
		{Type: FALSE, Value: "خطأ"},
		{Type: NEWLINE, Value: "\n"},

		// ليس صحيح
		{Type: NOT, Value: "ليس"},
		{Type: TRUE, Value: "صحيح"},
		{Type: NEWLINE, Value: "\n"},

		// # Bitwise
		{Type: COMMENT, Value: "# Bitwise\n"},

		// قيمة1 = 5 & 3
		{Type: NAME, Value: "قيمة1"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "5"},
		{Type: BITWISE_AND, Value: "&"},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// قيمة2 = 5 | 3
		{Type: NAME, Value: "قيمة2"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "5"},
		{Type: BITWISE_OR, Value: "|"},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// قيمة3 = 5 ^ 3
		{Type: NAME, Value: "قيمة3"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "5"},
		{Type: BITWISE_XOR, Value: "^"},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// قيمة4 = ~5
		{Type: NAME, Value: "قيمة4"},
		{Type: ASSIGN, Value: "="},
		{Type: BITWISE_NOT, Value: "~"},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		// قيمة5 = 2 << 3
		{Type: NAME, Value: "قيمة5"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "2"},
		{Type: LSHIFT, Value: "<<"},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// قيمة6 = 16 >> 2
		{Type: NAME, Value: "قيمة6"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "16"},
		{Type: RSHIFT, Value: ">>"},
		{Type: NUMBER, Value: "2"},
		{Type: NEWLINE, Value: "\n"},

		// # Assignment
		{Type: COMMENT, Value: "# Assignment\n"},

		// عدد = 10
		{Type: NAME, Value: "عدد"},
		{Type: ASSIGN, Value: "="},
		{Type: NUMBER, Value: "10"},
		{Type: NEWLINE, Value: "\n"},

		// عدد += 5
		{Type: NAME, Value: "عدد"},
		{Type: PLUS_ASSIGN, Value: "+="},
		{Type: NUMBER, Value: "5"},
		{Type: NEWLINE, Value: "\n"},

		// عدد -= 3
		{Type: NAME, Value: "عدد"},
		{Type: MINUS_ASSIGN, Value: "-="},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// عدد *= 2
		{Type: NAME, Value: "عدد"},
		{Type: MULT_ASSIGN, Value: "*="},
		{Type: NUMBER, Value: "2"},
		{Type: NEWLINE, Value: "\n"},

		// عدد /= 4
		{Type: NAME, Value: "عدد"},
		{Type: DIVIDE_ASSIGN, Value: "/="},
		{Type: NUMBER, Value: "4"},
		{Type: NEWLINE, Value: "\n"},

		// عدد %= 3
		{Type: NAME, Value: "عدد"},
		{Type: MOD_ASSIGN, Value: "%="},
		{Type: NUMBER, Value: "3"},
		{Type: NEWLINE, Value: "\n"},

		// عدد **= 2
		{Type: NAME, Value: "عدد"},
		{Type: POWER_ASSIGN, Value: "**="},
		{Type: NUMBER, Value: "2"},
		{Type: NEWLINE, Value: "\n"},

		{Type: EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeOperators", expected, tokens)
}

func TestTokenizeKeywords(t *testing.T) {
	tokens, err := Tokenize("../../tests/keywords.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize file: %v", err)
	}

	expected := []Token{
		// # Test all keywords
		{Type: COMMENT, Value: "# Test all keywords\n"},

		// إذا صحيح:
		{Type: IF, Value: "إذا"},
		{Type: TRUE, Value: "صحيح"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     ارجع 1
		{Type: INDENT, Value: ""},
		{Type: RETURN, Value: "ارجع"},
		{Type: NUMBER, Value: "1"},
		{Type: NEWLINE, Value: "\n"},

		// وإذا خطأ:
		{Type: DEDENT, Value: ""},
		{Type: ELIF, Value: "وإذا"},
		{Type: FALSE, Value: "خطأ"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     ارجع 0
		{Type: INDENT, Value: ""},
		{Type: RETURN, Value: "ارجع"},
		{Type: NUMBER, Value: "0"},
		{Type: NEWLINE, Value: "\n"},

		// وإلا:
		{Type: DEDENT, Value: ""},
		{Type: ELSE, Value: "وإلا"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     ارجع -1
		{Type: INDENT, Value: ""},
		{Type: RETURN, Value: "ارجع"},
		{Type: MINUS, Value: "-"},
		{Type: NUMBER, Value: "1"},
		{Type: NEWLINE, Value: "\n"},

		// طالما خطأ:
		{Type: DEDENT, Value: ""},
		{Type: WHILE, Value: "طالما"},
		{Type: FALSE, Value: "خطأ"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     اخرج
		{Type: INDENT, Value: ""},
		{Type: BREAK, Value: "اخرج"},
		{Type: NEWLINE, Value: "\n"},

		// كرر 5 مرات:
		{Type: DEDENT, Value: ""},
		{Type: REPEAT, Value: "كرر"},
		{Type: NUMBER, Value: "5"},
		{Type: TIMES, Value: "مرات"},
		{Type: COLON, Value: ":"},
		{Type: NEWLINE, Value: "\n"},

		//     تابع
		{Type: INDENT, Value: ""},
		{Type: CONTINUE, Value: "تابع"},
		{Type: NEWLINE, Value: "\n"},

		{Type: DEDENT, Value: ""},
		{Type: EOF, Value: ""},
	}

	compareTokens(t, "TestTokenizeKeywords", expected, tokens)
}
