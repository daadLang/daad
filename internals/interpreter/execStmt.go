package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
)

// execBlock executes a sequence of statements and returns the first non-NoSignal
// signal it encounters. If all statements execute normally it returns NewNoSignal().
func (i *Interpreter) execBlock(stmts []ast.Stmt) Signal {
	for _, s := range stmts {
		sig := i.execStmt(s)
		if sig.Type != NoSignal {
			return sig
		}
	}
	return NewNoSignal()
}

func (i *Interpreter) execIfStmt(stmt *ast.IfStmt) Signal {
	cond := i.execExpr(stmt.Test)
	condBool, ok := cond.(BoolValue)
	if !ok {
		panic(newRuntimeError("if condition must be a boolean, got %T", cond))
	}
	if condBool.V {
		sig := i.execBlock(stmt.Body)
		if sig.Type != NoSignal {
			return sig
		}
	} else {
		sig := i.execBlock(stmt.Orelse)
		if sig.Type != NoSignal {
			return sig
		}
	}
	return NewNoSignal()
}

func (i *Interpreter) execForStmt(stmt *ast.ForStmt) Signal {
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
		sig := i.execBlock(stmt.Body)
		if sig.IsReturn() || sig.IsError() {
			return sig
		}
		if sig.IsBreak() {
			// consume break: exit loop normally
			return NewNoSignal()
		}
		if sig.IsContinue() {
			// start next iteration
			continue
		}
	}
	// run orelse when loop completes normally
	if len(stmt.Orelse) > 0 {
		sig := i.execBlock(stmt.Orelse)
		if sig.Type != NoSignal {
			return sig
		}
	}

	return NewNoSignal()

}

func (i *Interpreter) execWhileStmt(stmt *ast.WhileStmt) Signal {
	for {
		cond := i.execExpr(stmt.Test)
		condBool, ok := cond.(BoolValue)
		if !ok {
			panic(newRuntimeError("while condition must be a boolean, got %T", cond))
		}
		if !condBool.V {
			break
		}
		sig := i.execBlock(stmt.Body)
		if sig.IsReturn() || sig.IsError() {
			return sig
		}
		if sig.IsBreak() {
			return NewNoSignal()
		}
		if sig.IsContinue() {
			continue
		}
	}

	// run orelse when while completes normally
	if len(stmt.Orelse) > 0 {
		sig := i.execBlock(stmt.Orelse)
		if sig.Type != NoSignal {
			return sig
		}
	}

	return NewNoSignal()
}

func (i *Interpreter) execReturnStmt(stmt *ast.ReturnStmt) Signal {
	value := i.execExpr(stmt.Value)
	return NewReturnSignal(value)
}

func (i *Interpreter) execAugmentedAssignStmt(stmt *ast.AugmentedAssignStmt) Signal {
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
	return NewNoSignal()
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

func (i *Interpreter) execFunctionDefStmt(stmt *ast.FunctionDefStmt) Signal {
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
	return NewNoSignal()
}
