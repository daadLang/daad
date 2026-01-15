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
	iterableSlice, ok := iterable.([]Value)
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

func (i *Interpreter) execReturnStmt(stmt *ast.ReturnStmt) Signal {
	value := i.execExpr(stmt.Value)
	return Signal{
		SignalType: ReturnSignal,
		Value:      value,
	}
}

func (i *Interpreter) execFunctionDefStmt(stmt *ast.FunctionDefStmt) {
	funcValue := &FunctionValue{
		Name:     stmt.Name,
		Args:     stmt.Args,
		Defaults: stmt.Defaults,
		Body:     stmt.Body,
		Env:      i.env,
	}
	i.env.Set(stmt.Name, funcValue)
}

func (i *Interpreter) execFunctionCallExpr(expr *ast.FunctionCallExpr) Value {
	funcValue := i.env.Get(expr.Name)
	if funcValue == nil {
		panic("undefined function: " + expr.Name)
	}
	fv, ok := funcValue.(*FunctionValue)
	if !ok {
		panic("called object is not a function")
	}

	// Evaluate all argument expressions to get their values
	args := make([]Value, len(expr.Args))
	for idx, argExpr := range expr.Args {
		args[idx] = i.execExpr(argExpr)
	}

	result, err := fv.Call(args)
	if err != nil {
		panic(err.Error())
	}
	return result
}
