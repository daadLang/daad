package interpreter

type Env struct {
	values map[string]Value
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		values: make(map[string]Value),
		parent: parent,
	}
}

func (e *Env) Get(name string) Value {
	if v, ok := e.values[name]; ok {
		return v
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	panic("undefined variable: " + name)
}

func (e *Env) Set(name string, val Value) {
	e.values[name] = val
}
