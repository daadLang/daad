package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
)

func (i *Interpreter) execIfStmt(stmt *ast.IfStmt) {
	cond := i.execExpr(stmt.Test)
	condBool, ok := cond.(BoolValue)
	if !ok {
		panic(newRuntimeError("if condition must be a boolean, got %T", cond))
	}
	if condBool.V {
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

	var items []Value
	switch iter := iterable.(type) {
	case ListValue:
		items = iter.Elements
	case TupleValue:
		items = iter.Elements
	default:
		panic(newRuntimeError("for loop iterable must be a list or tuple, got %T", iterable))
	}

	for _, item := range items {
		targetName, ok := stmt.Target.(*ast.Name)
		if !ok {
			panic(newRuntimeError("for loop target must be a variable name"))
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
		condBool, ok := cond.(BoolValue)
		if !ok {
			panic(newRuntimeError("while condition must be a boolean, got %T", cond))
		}
		if !condBool.V {
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

func (i *Interpreter) execAugmentedAssignStmt(stmt *ast.AugmentedAssignStmt) {
	targetName, ok := stmt.Target.(*ast.Name)
	if !ok {
		panic(newRuntimeError("augmented assignment target must be a variable name"))
	}

	currentValue := i.env.Get(targetName.Id)
	operandValue := i.execExpr(stmt.Value)

	binOp := &ast.BinOp{
		Left:  &ast.Constant{Value: extractRawValue(currentValue)},
		Right: &ast.Constant{Value: extractRawValue(operandValue)},
		Op:    augOpToBinOp(stmt.Op),
	}

	result := i.execBinOpExpr(binOp)
	i.env.Set(targetName.Id, result)
}

func augOpToBinOp(op lexer.TokenType) lexer.TokenType {
	switch op {
	case lexer.PLUS_ASSIGN:
		return lexer.PLUS
	case lexer.MINUS_ASSIGN:
		return lexer.MINUS
	case lexer.MULT_ASSIGN:
		return lexer.MULT
	case lexer.DIVIDE_ASSIGN:
		return lexer.DIVIDE
	case lexer.MOD_ASSIGN:
		return lexer.MOD
	case lexer.POWER_ASSIGN:
		return lexer.POWER
	default:
		panic(newRuntimeError("unknown augmented assignment operator: %v", op))
	}
}

func (i *Interpreter) execFunctionDefStmt(stmt *ast.FunctionDefStmt) {
	defaults := make([]Value, len(stmt.Defaults))
	for idx, defaultExpr := range stmt.Defaults {
		defaults[idx] = i.execExpr(defaultExpr)
	}

	funcValue := &FunctionValue{
		Name:     stmt.Name,
		Params:   stmt.Args, // arg names
		Defaults: defaults,  //default values
		Body:     stmt.Body,
		Env:      i.env,
	}
	i.env.Set(stmt.Name, funcValue)
}
