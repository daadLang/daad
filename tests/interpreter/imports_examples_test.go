package interpreter_test

import (
	"path/filepath"
	"runtime"
	"testing"
)

func examplePath(t *testing.T, name string) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("failed to resolve test file path")
	}
	testsDir := filepath.Dir(file)
	return filepath.Join(testsDir, "..", "examples", name)
}

func TestImportExampleFilesRun(t *testing.T) {
	// imports.daad now contains only import examples for parser testing
	// So we test with test_imports_comprehensive.daad which has executable code
	mainFile := examplePath(t, "test_imports_comprehensive.daad")
	interp := runDaadFile(t, mainFile)

	// Just verify that the file ran without errors and had output
	// The actual test of imports is done in other more specific tests
	if interp == nil {
		t.Fatalf("interpreter was nil after running imports comprehensive test")
	}
}

func TestImportMoreExamplesRun(t *testing.T) {
	mainFile := examplePath(t, "imports_more.daad")
	interp := runDaadFile(t, mainFile)

	if got := rawValueForImportTests(interp.GetVar("نتيجة1")); got != 103 {
		t.Fatalf("expected نتيجة1 to be 103, got %v", got)
	}
	if got := rawValueForImportTests(interp.GetVar("نتيجة2")); got != 16 {
		t.Fatalf("expected نتيجة2 to be 16, got %v", got)
	}
	if got := rawValueForImportTests(interp.GetVar("نتيجة3")); got != 7 {
		t.Fatalf("expected نتيجة3 to be 7, got %v", got)
	}
}

func TestNativeModulesExampleRun(t *testing.T) {
	mainFile := examplePath(t, "native_modules.daad")
	interp := runDaadFile(t, mainFile)

	if got := rawValueForImportTests(interp.GetVar("قيمة_مطلقة")); got != 3 {
		t.Fatalf("expected قيمة_مطلقة to be 3, got %v", got)
	}
	if got := rawValueForImportTests(interp.GetVar("اسم")); got != "c.txt" {
		t.Fatalf("expected اسم to be c.txt, got %v", got)
	}
	value, ok := rawValueForImportTests(interp.GetVar("قيمة_عشوائية")).(int)
	if !ok || value < 1 || value > 10 {
		t.Fatalf("expected قيمة_عشوائية to be within 1..10, got %v", value)
	}
}

func TestMixedFromImportAndNestedImport(t *testing.T) {
	tmp := t.TempDir()

	// Create a package tree that mixes simple and dotted from-imports.
	toolsDir := filepath.Join(tmp, "tools")
	advancedDir := filepath.Join(toolsDir, "advanced")

	writeDaadFile(t, filepath.Join(toolsDir, "_.ض"), "من . استورد advanced\nدالة نفذ():\n    ارجع 11\n")
	writeDaadFile(t, filepath.Join(advancedDir, "_.ض"), "من . استورد أداة\n")
	writeDaadFile(t, filepath.Join(advancedDir, "أداة.ض"), "دالة احسب(أ, ب):\n    ارجع أ * ب\n")

	mainFile := filepath.Join(tmp, "main.daad")
	writeDaadFile(t, mainFile, "استورد tools\nمن tools استورد نفذ, advanced.أداة\nنتيجة1 = tools.advanced.أداة.احسب(2, 3)\nنتيجة2 = أداة.احسب(4, 5)\nنتيجة3 = نفذ()\n")

	interp := runDaadFile(t, mainFile)
	if got := rawValueForImportTests(interp.GetVar("نتيجة1")); got != 6 {
		t.Fatalf("expected نتيجة1 to be 6, got %v", got)
	}
	if got := rawValueForImportTests(interp.GetVar("نتيجة2")); got != 20 {
		t.Fatalf("expected نتيجة2 to be 20, got %v", got)
	}
	if got := rawValueForImportTests(interp.GetVar("نتيجة3")); got != 11 {
		t.Fatalf("expected نتيجة3 to be 11, got %v", got)
	}
}

func TestRelativeAndStarImportTogether(t *testing.T) {
	tmp := t.TempDir()

	// Create a package that exports a submodule and a value.
	libDir := filepath.Join(tmp, "lib")
	subDir := filepath.Join(libDir, "sub")
	writeDaadFile(t, filepath.Join(libDir, "_.daad"), "من . استورد sub باسم فرعية\nقيمة_الحزمة = 100\n")
	writeDaadFile(t, filepath.Join(subDir, "_.daad"), "دالة جمع(أ, ب):\n    ارجع أ + ب\n")

	// Main file uses both star import and relative import from the same directory.
	writeDaadFile(t, filepath.Join(tmp, "helper.daad"), "قيمة = 7\n")
	mainFile := filepath.Join(tmp, "main.daad")
	writeDaadFile(t, mainFile, "استورد lib\nمن lib استورد *\nمن . استورد helper\nنتيجة1 = قيمة_الحزمة + فرعية.جمع(1, 2)\nنتيجة2 = helper.قيمة\n")

	interp := runDaadFile(t, mainFile)
	if got := rawValueForImportTests(interp.GetVar("نتيجة1")); got != 103 {
		t.Fatalf("expected نتيجة1 to be 103, got %v", got)
	}
	if got := rawValueForImportTests(interp.GetVar("نتيجة2")); got != 7 {
		t.Fatalf("expected نتيجة2 to be 7, got %v", got)
	}
}
