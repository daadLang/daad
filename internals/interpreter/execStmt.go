package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
)

func (i *Interpreter) execIfStmt(stmt *ast.IfStmt) {
	cond := i.execExpr(stmt.Test)
	condBool, ok := cond.(bool)
	if !ok {
		panic("if condition must be a boolean")
	}
	if condBool {
		for _, s := range stmt.Body {
			i.execStmt(s)
		}
	} else {
		for _, s := range stmt.Orelse {
			i.execStmt(s)
		}
	}
}

func (i *Interpreter) execForStmt(stmt *ast.ForStmt) {
	iterable := i.execExpr(stmt.Iter)
	iterableSlice, ok := iterable.([]any)
	if !ok {
		panic("for loop iterable must be a slice")
	}
	for _, item := range iterableSlice {
		targetName, ok := stmt.Target.(*ast.Name)
		if !ok {
			panic("for loop target must be a variable name")
		}
		i.env.Set(targetName.Id, item)
		for _, s := range stmt.Body {
			i.execStmt(s)
		}
	}
}

func (i *Interpreter) execWhileStmt(stmt *ast.WhileStmt) {
	for {
		cond := i.execExpr(stmt.Test)
		condBool, ok := cond.(bool)
		if !ok {
			panic("while condition must be a boolean")
		}
		if !condBool {
			break
		}
		for _, s := range stmt.Body {
			i.execStmt(s)
		}
	}
}

func (i *Interpreter) execReturnStmt(stmt *ast.ReturnStmt) {
	value := i.execExpr(stmt.Value)
	// Using panic to simulate return behavior
	panic(returnValue{value})
}

type returnValue struct {
	value any
}
