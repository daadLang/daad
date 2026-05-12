package interpreter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daadLang/daad/internals/interpreter"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
)

func getIntValue(v interpreter.Value) int {
	switch val := v.(type) {
	case interpreter.IntValue:
		return val.V
	default:
		return -999
	}
}

func createFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func testRunFile(t *testing.T, mainFile string) *interpreter.Interpreter {
	t.Helper()

	tokens, err := lexer.Tokenize(mainFile)
	if err != nil {
		t.Fatalf("failed to tokenize %s: %v", mainFile, err)
	}

	p := parser.NewParser(tokens)
	module := p.Parse()

	interp := interpreter.NewInterpreter()
	interp.SetSourcePath(mainFile)
	interp.Run(&module)
	return interp
}

func TestNestedImportChaining(t *testing.T) {
	tmp := t.TempDir()

	pkgDir := filepath.Join(tmp, "pkg")
	subDir := filepath.Join(pkgDir, "sub")
	deepDir := filepath.Join(subDir, "deep")

	createFile(t, filepath.Join(pkgDir, "_.daad"), "من . استورد sub\n")
	createFile(t, filepath.Join(subDir, "_.daad"), "من . استورد deep\n")
	createFile(t, filepath.Join(deepDir, "_.daad"), "قيمة = 42\n")

	mainFile := filepath.Join(tmp, "main.daad")
	createFile(t, mainFile, "استورد pkg.sub.deep\nج = pkg.sub.deep.قيمة\n")

	interp := testRunFile(t, mainFile)
	got := getIntValue(interp.GetVar("ج"))
	if got != 42 {
		t.Fatalf("expected ج to be 42, got %v", got)
	}
}

func TestSimpleAndNestedImportCoexist(t *testing.T) {
	tmp := t.TempDir()

	pkgDir := filepath.Join(tmp, "pkg")
	subDir := filepath.Join(pkgDir, "sub")

	createFile(t, filepath.Join(pkgDir, "_.daad"), "من . استورد sub\nx = 10\n")
	createFile(t, filepath.Join(subDir, "_.daad"), "y = 20\n")

	mainFile := filepath.Join(tmp, "main.daad")
	createFile(t, mainFile, "استورد pkg\nاستورد pkg.sub\nأ = pkg.x\nب = pkg.sub.y\nج = أ + ب\n")

	interp := testRunFile(t, mainFile)
	got := getIntValue(interp.GetVar("ج"))
	if got != 30 {
		t.Fatalf("expected ج to be 30, got %v", got)
	}
}

func TestNestedImportWithFunction(t *testing.T) {
	tmp := t.TempDir()

	pkgDir := filepath.Join(tmp, "tools")
	advDir := filepath.Join(pkgDir, "advanced")

	createFile(t, filepath.Join(pkgDir, "_.daad"), "من . استورد advanced\n")
	createFile(t, filepath.Join(advDir, "_.daad"), "دالة حساب(أ, ب):\n    ارجع أ * ب + 10\n")

	mainFile := filepath.Join(tmp, "main.daad")
	createFile(t, mainFile, "استورد tools.advanced\nج = tools.advanced.حساب(3, 4)\n")

	interp := testRunFile(t, mainFile)
	got := getIntValue(interp.GetVar("ج"))
	if got != 22 {
		t.Fatalf("expected ج to be 22, got %v", got)
	}
}

func TestFromImportNestedModule(t *testing.T) {
	tmp := t.TempDir()

	pkgDir := filepath.Join(tmp, "pkg")
	subDir := filepath.Join(pkgDir, "sub")

	createFile(t, filepath.Join(pkgDir, "_.daad"), "من . استورد sub\n")
	createFile(t, filepath.Join(subDir, "_.daad"), "رقم = 55\n")

	mainFile := filepath.Join(tmp, "main.daad")
	createFile(t, mainFile, "من pkg استورد sub باسم deep\nپ = deep.رقم\n")

	interp := testRunFile(t, mainFile)
	got := getIntValue(interp.GetVar("پ"))
	if got != 55 {
		t.Fatalf("expected پ to be 55, got %v", got)
	}
}
