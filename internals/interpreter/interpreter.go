package interpreter

import (
	"fmt"

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

func (i *Interpreter) execStmt(stmt ast.Stmt) {
	switch e := stmt.(type) {
	case *ast.ExprStmt:
		i.execExpr(e.Value)
	case *ast.IfStmt:
		i.execIfStmt(e)
	case *ast.ForStmt:
		i.execForStmt(e)
	case *ast.WhileStmt:
		i.execWhileStmt(e)
	case *ast.ReturnStmt:
		i.execReturnStmt(e)
	default:
		panic("unknown statement: " + fmt.Sprintf("%T", stmt))
	}
}

func (i *Interpreter) execExpr(expr ast.Expr) any {
	switch e := expr.(type) {
	case *ast.Constant:
		return e.Value

	case *ast.Name:
		return i.env.Get(e.Id)
	}
	panic("unknown expression: " + fmt.Sprintf("%T", expr))
}
