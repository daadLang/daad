package parser_test

import (
	"testing"

	"github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
)

func mustTokenize(t *testing.T, path string) []lexer.Token {
	t.Helper()

	tokens, err := lexer.Tokenize(path)
	if err != nil {
		t.Fatalf("failed to tokenize %s: %v", path, err)
	}
	return tokens
}

func mustParse(t *testing.T, path string) ast.Module {
	t.Helper()

	tokens := mustTokenize(t, path)
	p := parser.NewParser(tokens)
	return p.Parse()
}

func assertAlias(t *testing.T, got ast.Alias, wantName string, wantAs *string) {
	t.Helper()

	if got.Name != wantName {
		t.Fatalf("expected alias name %q, got %q", wantName, got.Name)
	}

	if wantAs == nil {
		if got.AsName != nil {
			t.Fatalf("expected alias %q to have no as-name, got %q", wantName, *got.AsName)
		}
		return
	}

	if got.AsName == nil {
		t.Fatalf("expected alias %q to have as-name %q, got nil", wantName, *wantAs)
	}

	if *got.AsName != *wantAs {
		t.Fatalf("expected alias %q to have as-name %q, got %q", wantName, *wantAs, *got.AsName)
	}
}

func TestParseImportsFixture(t *testing.T) {
	module := mustParse(t, "../examples/imports.daad")

	if len(module.Body) != 11 {
		t.Fatalf("expected 11 statements, got %d", len(module.Body))
	}

	// 1) استورد مكتبة
	if stmt, ok := module.Body[0].(*ast.ImportStmt); !ok {
		t.Fatalf("statement 0: expected *ast.ImportStmt, got %T", module.Body[0])
	} else {
		if stmt.Module != "" {
			t.Fatalf("statement 0: expected empty Module, got %q", stmt.Module)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 0: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "مكتبة", nil)
	}

	// 2) استورد مكتبة.فرعية باسم فرعية
	if stmt, ok := module.Body[1].(*ast.ImportStmt); !ok {
		t.Fatalf("statement 1: expected *ast.ImportStmt, got %T", module.Body[1])
	} else {
		alias := "فرعية"
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 1: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "مكتبة.فرعية", &alias)
	}

	// 3) استورد مكتبة.فرعية, أخرى.حزمة ك حزمة
	if stmt, ok := module.Body[2].(*ast.ImportStmt); !ok {
		t.Fatalf("statement 2: expected *ast.ImportStmt, got %T", module.Body[2])
	} else {
		alias := "حزمة"
		if len(stmt.Names) != 2 {
			t.Fatalf("statement 2: expected 2 imported names, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "مكتبة.فرعية", nil)
		assertAlias(t, stmt.Names[1], "أخرى.حزمة", &alias)
	}

	// 4) استورد أداة كـ أداة_بديلة
	if stmt, ok := module.Body[3].(*ast.ImportStmt); !ok {
		t.Fatalf("statement 3: expected *ast.ImportStmt, got %T", module.Body[3])
	} else {
		alias := "أداة_بديلة"
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 3: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "أداة", &alias)
	}

	// 5) من حزمة استورد اسم
	if stmt, ok := module.Body[4].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 4: expected *ast.FromImportStmt, got %T", module.Body[4])
	} else {
		if stmt.Module != "حزمة" {
			t.Fatalf("statement 4: expected module %q, got %q", "حزمة", stmt.Module)
		}
		if stmt.Level != 0 {
			t.Fatalf("statement 4: expected level 0, got %d", stmt.Level)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 4: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "اسم", nil)
	}

	// 6) من حزمة استورد اسم باسم alias, ثاني
	if stmt, ok := module.Body[5].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 5: expected *ast.FromImportStmt, got %T", module.Body[5])
	} else {
		alias := "اسم_بديل"
		if stmt.Module != "حزمة" {
			t.Fatalf("statement 5: expected module %q, got %q", "حزمة", stmt.Module)
		}
		if len(stmt.Names) != 2 {
			t.Fatalf("statement 5: expected 2 imported names, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "اسم", &alias)
		assertAlias(t, stmt.Names[1], "ثاني", nil)
	}

	// 7) من حزمة.فرعية استورد *
	if stmt, ok := module.Body[6].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 6: expected *ast.FromImportStmt, got %T", module.Body[6])
	} else {
		if stmt.Module != "حزمة.فرعية" {
			t.Fatalf("statement 6: expected module %q, got %q", "حزمة.فرعية", stmt.Module)
		}
		if stmt.Level != 0 {
			t.Fatalf("statement 6: expected level 0, got %d", stmt.Level)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 6: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "*", nil)
	}

	// 8) من . استورد محلي
	if stmt, ok := module.Body[7].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 7: expected *ast.FromImportStmt, got %T", module.Body[7])
	} else {
		if stmt.Module != "" {
			t.Fatalf("statement 7: expected empty module, got %q", stmt.Module)
		}
		if stmt.Level != 1 {
			t.Fatalf("statement 7: expected level 1, got %d", stmt.Level)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 7: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "محلي", nil)
	}

	// 9) من . استورد *
	if stmt, ok := module.Body[8].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 8: expected *ast.FromImportStmt, got %T", module.Body[8])
	} else {
		if stmt.Module != "" {
			t.Fatalf("statement 8: expected empty module, got %q", stmt.Module)
		}
		if stmt.Level != 1 {
			t.Fatalf("statement 8: expected level 1, got %d", stmt.Level)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 8: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "*", nil)
	}

	// 10) من .. استورد أعلى
	if stmt, ok := module.Body[9].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 9: expected *ast.FromImportStmt, got %T", module.Body[9])
	} else {
		if stmt.Module != "" {
			t.Fatalf("statement 9: expected empty module, got %q", stmt.Module)
		}
		if stmt.Level != 2 {
			t.Fatalf("statement 9: expected level 2, got %d", stmt.Level)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 9: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "أعلى", nil)
	}

	// 11) من ..حزمة.فرعية استورد اسم2 باسم اسم_مستعار
	if stmt, ok := module.Body[10].(*ast.FromImportStmt); !ok {
		t.Fatalf("statement 10: expected *ast.FromImportStmt, got %T", module.Body[10])
	} else {
		alias := "اسم_مستعار"
		if stmt.Module != "حزمة.فرعية" {
			t.Fatalf("statement 10: expected module %q, got %q", "حزمة.فرعية", stmt.Module)
		}
		if stmt.Level != 2 {
			t.Fatalf("statement 10: expected level 2, got %d", stmt.Level)
		}
		if len(stmt.Names) != 1 {
			t.Fatalf("statement 10: expected 1 imported name, got %d", len(stmt.Names))
		}
		assertAlias(t, stmt.Names[0], "اسم2", &alias)
	}
}
