package interpreter_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/daadLang/daad/internals/interpreter"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
)

func rawValueForImportTests(v interpreter.Value) interface{} {
	switch val := v.(type) {
	case interpreter.IntValue:
		return val.V
	case interpreter.FloatValue:
		return val.V
	case interpreter.StringValue:
		return val.V
	case interpreter.BoolValue:
		return val.V
	case interpreter.CharValue:
		return val.V
	case interpreter.NoneValue:
		return nil
	default:
		return v
	}
}

func writeDaadFile(t *testing.T, path string, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("failed to create dir for %s: %v", path, err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("failed to write %s: %v", path, err)
	}
}

func runDaadFile(t *testing.T, mainFile string) *interpreter.Interpreter {
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

func TestInterpreterImportAndFromImport(t *testing.T) {
	tmp := t.TempDir()

	mainFile := filepath.Join(tmp, "main.daad")
	moduleFile := filepath.Join(tmp, "mathlib.daad")

	writeDaadFile(t, moduleFile, "قيمة = 7\n")
	writeDaadFile(t, mainFile, "استورد mathlib\nمن mathlib استورد قيمة باسم ن\nنتيجة = mathlib.قيمة + ن\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("نتيجة"))
	if got != 14 {
		t.Fatalf("expected نتيجة to be 14, got %v", got)
	}
}

func TestInterpreterRelativeFromImport(t *testing.T) {
	tmp := t.TempDir()

	mainFile := filepath.Join(tmp, "pkg", "main.daad")
	helperFile := filepath.Join(tmp, "pkg", "helper.daad")

	writeDaadFile(t, helperFile, "قيمة = 9\n")
	writeDaadFile(t, mainFile, "من . استورد helper\nق = helper.قيمة\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("ق"))
	if got != 9 {
		t.Fatalf("expected ق to be 9, got %v", got)
	}
}

func TestInterpreterFromImportStar(t *testing.T) {
	tmp := t.TempDir()

	mainFile := filepath.Join(tmp, "main.daad")
	moduleFile := filepath.Join(tmp, "lib.daad")

	writeDaadFile(t, moduleFile, "أ = 1\nب = 2\n")
	writeDaadFile(t, mainFile, "من lib استورد *\nج = أ + ب\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("ج"))
	if got != 3 {
		t.Fatalf("expected ج to be 3, got %v", got)
	}
}
