package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
)

type Interpreter struct {
	env *Env
}

func NewInterpreter() *Interpreter {
	env := NewEnv(nil)
	RegisterBuiltins(env)
	return &Interpreter{
		env: env,
	}
}

func (i *Interpreter) SetVar(name string, value Value) {
	i.env.Set(name, value)
}

func (i *Interpreter) GetVar(name string) Value {
	return i.env.Get(name)
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
		return NewNoSignal()
	case *ast.IfStmt:
		return i.execIfStmt(e)
	case *ast.ForStmt:
		return i.execForStmt(e)
	case *ast.RepeatStmt:
		return i.execRepeatStmt(e)
	case *ast.WhileStmt:
		return i.execWhileStmt(e)
	case *ast.AssignStmt:
		value := i.execExpr(e.Value)
		i.env.Set(e.Target.Id, value)
		return NewNoSignal()
	case *ast.AugmentedAssignStmt:
		return i.execAugmentedAssignStmt(e)
	case *ast.FunctionDefStmt:
		return i.execFunctionDefStmt(e)
	case *ast.ReturnStmt:
		return i.execReturnStmt(e)
	case *ast.BreakStmt:
		return NewBreakSignal()
	case *ast.ContinueStmt:
		return NewContinueSignal()

	default:
		panic(newRuntimeError("unknown statement: %T", stmt))
	}
}

func (i *Interpreter) execExpr(expr ast.Expr) Value {
	switch e := expr.(type) {
	case *ast.Constant:
		return i.execConstExpr(e)

	case *ast.Name:
		return i.env.Get(e.Id)

	case *ast.UnaryOp:
		return i.execUnaryOpExpr(e)

	case *ast.BinOp:
		return i.execBinOpExpr(e)

	case *ast.BoolOp:
		return i.execBoolOpExpr(e)

	case *ast.Compare:
		return i.execCompareExpr(e)

	case *ast.Assign:
		return i.execAssignExpr(e)

	case *ast.Call:
		return i.execCallExpr(e)

	case *ast.Subscript:
		return i.execSubscriptExpr(e)

	case *ast.List:
		return i.execListExpr(e)

	case *ast.Dict:
		return i.execDictExpr(e)

	case *ast.Tuple:
		return i.execTupleExpr(e)
	}

	panic(newRuntimeError("unknown expression: %T", expr))
}
