package interpreter

type Env struct {
	values   map[string]Value
	builtins map[string]Value
	parent   *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		values:   make(map[string]Value),
		builtins: make(map[string]Value),
		parent:   parent,
	}
}

func (e *Env) Get(name string) Value {
	if v, ok := e.values[name]; ok {
		return v
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	if v, ok := e.builtins[name]; ok {
		return v
	}
	panic("undefined variable: " + name)
}

func (e *Env) Set(name string, val Value) {
	e.values[name] = val
}

func (e *Env) SetBuiltin(name string, val Value) {
	e.builtins[name] = val
}

func (e *Env) IsBuiltin(name string) bool {
	if _, ok := e.builtins[name]; ok {
		return true
	}
	if e.parent != nil {
		return e.parent.IsBuiltin(name)
	}
	return false
}
