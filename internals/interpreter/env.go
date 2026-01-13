package interpreter

type Env struct {
	values map[string]any
	parent *Env
}

func NewEnv(parent *Env) *Env {
	return &Env{
		values: make(map[string]any),
		parent: parent,
	}
}

func (e *Env) Get(name string) any {
	if v, ok := e.values[name]; ok {
		return v
	}
	if e.parent != nil {
		return e.parent.Get(name)
	}
	panic("undefined variable: " + name)
}

func (e *Env) Set(name string, val any) {
	e.values[name] = val
}
