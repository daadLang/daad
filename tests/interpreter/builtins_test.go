package interpreter_test

import (
	"path/filepath"
	"testing"

	"github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/interpreter"
)

func TestBuiltinShadowingInModuleScope(t *testing.T) {
	module := &ast.Module{Body: []ast.Stmt{
		&ast.ExprStmt{Value: &ast.Assign{Target: &ast.Name{Id: "طول"}, Value: &ast.Constant{Value: 9}}},
		&ast.ExprStmt{Value: &ast.Assign{Target: &ast.Name{Id: "نتيجة"}, Value: &ast.Name{Id: "طول"}}},
	}}

	interp := interpreter.NewInterpreter()
	interp.Run(module)

	got := rawValueForImportTests(interp.GetVar("نتيجة"))
	if got != 9 {
		t.Fatalf("expected shadowed builtin value 9, got %v", got)
	}
}

func TestBuiltinStillCallableWhenNotShadowed(t *testing.T) {
	module := &ast.Module{Body: []ast.Stmt{
		&ast.ExprStmt{Value: &ast.Assign{
			Target: &ast.Name{Id: "lenValue"},
			Value: &ast.Call{
				Func: &ast.Name{Id: "طول"},
				Args: []ast.Expr{
					&ast.List{Elements: []ast.Expr{
						&ast.Constant{Value: 1},
						&ast.Constant{Value: 2},
						&ast.Constant{Value: 3},
					}},
				},
			},
		}},
	}}

	interp := interpreter.NewInterpreter()
	interp.Run(module)

	got := rawValueForImportTests(interp.GetVar("lenValue"))
	if got != 3 {
		t.Fatalf("expected builtin طول() to return 3, got %v", got)
	}
}

func TestImportedModuleExportsShadowedBuiltin(t *testing.T) {
	tmp := t.TempDir()

	mainFile := filepath.Join(tmp, "main.daad")
	moduleFile := filepath.Join(tmp, "lib.daad")

	writeDaadFile(t, moduleFile, "طول = 11\nقيمة = 4\n")
	writeDaadFile(t, mainFile, "استورد lib\nمن lib استورد طول باسم قيمة_طول\nمجموع = lib.طول + قيمة_طول + lib.قيمة\n")

	interp := runDaadFile(t, mainFile)
	got := rawValueForImportTests(interp.GetVar("مجموع"))
	if got != 26 {
		t.Fatalf("expected imported shadowed builtin to be exported and usable, got %v", got)
	}

	mod, ok := interp.GetVar("lib").(*interpreter.ModuleValue)
	if !ok {
		t.Fatalf("expected lib to be a ModuleValue, got %T", interp.GetVar("lib"))
	}
	shadowed, ok := mod.Attributes["طول"]
	if !ok {
		t.Fatalf("expected module export to include shadowed builtin name طول")
	}
	if rawValueForImportTests(shadowed) != 11 {
		t.Fatalf("expected lib.طول to be 11, got %v", rawValueForImportTests(shadowed))
	}
}
