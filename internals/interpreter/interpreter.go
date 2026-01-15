package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
)

type Interpreter struct {
	env *Env
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnv(nil),
	}
}

func (i *Interpreter) Run(m *ast.Module) {
	for _, stmt := range m.Body {
		i.execStmt(stmt)
	}
}

func (i *Interpreter) execStmt(stmt ast.Stmt) Signal {
	switch e := stmt.(type) {
	case *ast.ExprStmt:
		i.execExpr(e.Value)
	case *ast.IfStmt:
		i.execIfStmt(e)
	case *ast.ForStmt:
		i.execForStmt(e)
	case *ast.WhileStmt:
		i.execWhileStmt(e)
	case *ast.AssignStmt:
		value := i.execExpr(e.Value)
		i.env.Set(e.Target.Id, value)
	case *ast.ReturnStmt:
		return i.execReturnStmt(e)
	default:
		panic(newRuntimeError("unknown statement: %T", stmt))
	}
	return Signal{SignalType: NoSignal}
}

func (i *Interpreter) execExpr(expr ast.Expr) Value {
	switch e := expr.(type) {
	case *ast.Constant:
		return e.Value

	case *ast.Name:
		return i.env.Get(e.Id)

	case *ast.UnaryOp:
		return i.execUnaryOpExpr(e)

	case *ast.BinOp:
		return i.execBinOpExpr(e)
	}

	panic(newRuntimeError("unknown expression: %T", expr))
}
