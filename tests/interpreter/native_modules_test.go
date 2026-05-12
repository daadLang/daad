package interpreter_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func TestNativeMathModule(t *testing.T) {
	mainFile := filepath.Join(t.TempDir(), "main.daad")
	writeDaadFile(t, mainFile, "استورد رياضيات\nنتيجة = رياضيات.جذر(9) + رياضيات.مطلق(-3)\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("نتيجة"))
	if got != 6.0 {
		t.Fatalf("expected نتيجة to be 6.0, got %v", got)
	}
}

func TestNativeRandomModuleSeed(t *testing.T) {
	mainFile := filepath.Join(t.TempDir(), "main.daad")
	writeDaadFile(t, mainFile, "استورد عشوائي\nعشوائي.اضبط_البذرة(42)\nأ = عشوائي.عدد_عشوائي(1, 100)\nعشوائي.اضبط_البذرة(42)\nب = عشوائي.عدد_عشوائي(1, 100)\n")

	interp := runDaadFile(t, mainFile)
	first := rawValueForImportTests(interp.GetVar("أ"))
	second := rawValueForImportTests(interp.GetVar("ب"))
	if first != second {
		t.Fatalf("expected deterministic values with same seed, got %v and %v", first, second)
	}
}

func TestNativeTimeModule(t *testing.T) {
	mainFile := filepath.Join(t.TempDir(), "main.daad")
	writeDaadFile(t, mainFile, "استورد وقت\nالان_ث = وقت.الان()\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("الان_ث"))
	value, ok := got.(int)
	if !ok || value <= 0 {
		t.Fatalf("expected الان_ث to be a positive int, got %v", got)
	}
}

func TestNativePathModule(t *testing.T) {
	tmp := t.TempDir()
	mainFile := filepath.Join(tmp, "main.daad")
	writeDaadFile(t, mainFile, "استورد مسار\nنتيجة = مسار.ربط_المسار(\"a\", \"b\", \"c.txt\")\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("نتيجة"))
	want := filepath.Join("a", "b", "c.txt")
	if got != want {
		t.Fatalf("expected نتيجة to be %v, got %v", want, got)
	}
}

func TestNativeOSModuleExists(t *testing.T) {
	tmp := t.TempDir()
	filePath := filepath.Join(tmp, "check.txt")
	if err := os.WriteFile(filePath, []byte("ok"), 0o644); err != nil {
		t.Fatalf("failed to write temp file: %v", err)
	}

	mainFile := filepath.Join(tmp, "main.daad")
	writeDaadFile(t, mainFile, fmt.Sprintf("استورد نظام\nنتيجة = نظام.موجود(%q)\n", filePath))

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("نتيجة"))
	if got != true {
		t.Fatalf("expected نتيجة to be true, got %v", got)
	}
}
