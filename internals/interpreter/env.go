package interpreter

import "github.com/daadLang/daad/internals/ast"

type Env struct {
	values map[string]ast.Var
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		values: make(map[string]ast.Var),
		parent: parent,
	}
}

func (e *Env) Get(name string) ast.Var {
	if v, ok := e.values[name]; ok {
		return v
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	panic("undefined variable: " + name)
}

func (e *Env) Set(name string, val ast.Var) {
	e.values[name] = val
}
