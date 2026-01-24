package interpreter

type Adapter interface{}

type FuncAdapter interface {
	Adapter
	adapter()
}

type FuncArgsAdapter struct {
	FuncAdapter

	args []Value
}

func (*FuncArgsAdapter) adapter() {}

type FuncKwargsAdapter struct {
	FuncAdapter

	kwargs map[string]Value
}

func (*FuncKwargsAdapter) adapter() {}
