package ast_test

import (
	"testing"

	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
	"github.com/daadLang/daad/internals/parser"
)

func TestParseOOPClassUsesAttributeNodes(t *testing.T) {
	tokens, err := lexer.Tokenize("../examples/oop_class.daad")
	if err != nil {
		t.Fatalf("Failed to tokenize OOP example: %v", err)
	}

	module := parser.NewParser(tokens).Parse()
	if len(module.Body) < 3 {
		t.Fatalf("expected at least 3 top-level statements, got %d", len(module.Body))
	}

	classDef, ok := module.Body[0].(*ast.ClassDefStmt)
	if !ok {
		t.Fatalf("expected first statement to be ClassDefStmt, got %T", module.Body[0])
	}

	if len(classDef.Body) < 2 {
		t.Fatalf("expected class body to have methods, got %d", len(classDef.Body))
	}

	constructor, ok := classDef.Body[0].(*ast.FunctionDefStmt)
	if !ok || constructor.Name != "__بناء__" {
		t.Fatalf("expected first class statement to be __بناء__ method, got %T", classDef.Body[0])
	}

	assignStmt, ok := constructor.Body[0].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected constructor first statement to be ExprStmt, got %T", constructor.Body[0])
	}
	assignExpr, ok := assignStmt.Value.(*ast.Assign)
	if !ok {
		t.Fatalf("expected constructor expression to be Assign, got %T", assignStmt.Value)
	}
	if _, ok := assignExpr.Target.(*ast.Attribute); !ok {
		t.Fatalf("expected constructor assignment target to be Attribute, got %T", assignExpr.Target)
	}

	methodCallStmt, ok := module.Body[len(module.Body)-1].(*ast.ExprStmt)
	if !ok {
		t.Fatalf("expected last statement to be ExprStmt, got %T", module.Body[len(module.Body)-1])
	}
	methodCall, ok := methodCallStmt.Value.(*ast.Call)
	if !ok {
		t.Fatalf("expected last expression to be Call, got %T", methodCallStmt.Value)
	}
	if _, ok := methodCall.Func.(*ast.Attribute); !ok {
		t.Fatalf("expected method call function to be Attribute, got %T", methodCall.Func)
	}
}

func TestParseClassWithSingleInheritance(t *testing.T) {
	tokens, err := lexer.TokenizeString("صنف ابن(اب):\n    دالة مرحبا(ذاتي):\n        ارجع 1\n")
	if err != nil {
		t.Fatalf("Failed to tokenize inheritance class: %v", err)
	}

	module := parser.NewParser(tokens).Parse()
	if len(module.Body) != 1 {
		t.Fatalf("expected one top-level statement, got %d", len(module.Body))
	}

	classDef, ok := module.Body[0].(*ast.ClassDefStmt)
	if !ok {
		t.Fatalf("expected ClassDefStmt, got %T", module.Body[0])
	}

	if classDef.Name != "ابن" {
		t.Fatalf("expected class name ابن, got %s", classDef.Name)
	}

	if classDef.Parent != "اب" {
		t.Fatalf("expected parent class اب, got %s", classDef.Parent)
	}
}
