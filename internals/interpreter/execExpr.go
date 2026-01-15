package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
)

func (i *Interpreter) execBinOpExpr(e *ast.BinOp) Value {
	left := i.execExpr(e.Left)
	right := i.execExpr(e.Right)

	switch e.Op {
	case lexer.PLUS:
		return castAdd(left, right)

	case lexer.MINUS:
		leftVal, rightVal, isFloat := castNumericOp(left, right, "-")
		if isFloat {
			return leftVal - rightVal
		}
		return int(leftVal) - int(rightVal)

	case lexer.MULT:
		leftVal, rightVal, isFloat := castNumericOp(left, right, "*")
		if isFloat {
			return leftVal * rightVal
		}
		return int(leftVal) * int(rightVal)

	case lexer.DIVIDE:
		leftVal, rightVal, _ := castNumericOp(left, right, "/")
		return leftVal / rightVal

	default:
		panic(newRuntimeError("unknown binary operator: %v", e.Op))
	}
}

func (i *Interpreter) execUnaryOpExpr(e *ast.UnaryOp) Value {
	switch e.Op {
	case lexer.NOT:
		return !i.execExpr(e.Expr).(bool)
	default:
		panic(newRuntimeError("unknown unary operator: %v", e.Op))
	}
}
