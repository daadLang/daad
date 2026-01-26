package interpreter_test

import (
	"testing"

	"github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/interpreter"
	"github.com/daadLang/daad/internals/lexer"
)

func TestReturnSignal(t *testing.T) {
	module := &ast.Module{Body: []ast.Stmt{
		&ast.FunctionDefStmt{
			Name: "foo",
			Args: []string{},
			Body: []ast.Stmt{
				&ast.ReturnStmt{Value: &ast.Constant{Value: 42}},
			},
		},
		&ast.ExprStmt{Value: &ast.Assign{
			Target: &ast.Name{Id: "x"},
			Value:  &ast.Call{Func: &ast.Name{Id: "foo"}},
		}},
	}}

	interp := interpreter.NewInterpreter()
	interp.Run(module)
	v := interp.GetVar("x")
	iv, ok := v.(interpreter.IntValue)
	if !ok {
		t.Fatalf("expected int value, got %T", v)
	}
	if iv.V != 42 {
		t.Fatalf("expected 42, got %d", iv.V)
	}
}

func TestBreakContinueSignals(t *testing.T) {
	module := &ast.Module{Body: []ast.Stmt{
		&ast.ExprStmt{Value: &ast.Assign{Target: &ast.Name{Id: "x"}, Value: &ast.Constant{Value: 0}}},
		&ast.ForStmt{
			Target: &ast.Name{Id: "i"},
			Iter: &ast.List{Elements: []ast.Expr{
				&ast.Constant{Value: 0},
				&ast.Constant{Value: 1},
				&ast.Constant{Value: 2},
				&ast.Constant{Value: 3},
				&ast.Constant{Value: 4},
			}},
			Body: []ast.Stmt{
				&ast.IfStmt{
					Test: &ast.Compare{Left: &ast.Name{Id: "i"}, Op: lexer.EQ, Comparator: &ast.Constant{Value: 2}},
					Body: []ast.Stmt{&ast.ContinueStmt{}},
				},
				&ast.IfStmt{
					Test: &ast.Compare{Left: &ast.Name{Id: "i"}, Op: lexer.EQ, Comparator: &ast.Constant{Value: 4}},
					Body: []ast.Stmt{&ast.BreakStmt{}},
				},
				&ast.ExprStmt{Value: &ast.Assign{Target: &ast.Name{Id: "x"}, Value: &ast.BinOp{Left: &ast.Name{Id: "x"}, Op: lexer.PLUS, Right: &ast.Name{Id: "i"}}}},
			},
		},
	}}

	interp := interpreter.NewInterpreter()
	interp.Run(module)
	v := interp.GetVar("x")
	iv, ok := v.(interpreter.IntValue)
	if !ok {
		t.Fatalf("expected int value, got %T", v)
	}
	if iv.V != 4 {
		t.Fatalf("expected 4, got %d", iv.V)
	}
}
