package interpreter

import ast "github.com/daadLang/daad/internals/ast"

type ValueType int

const (
	IntType ValueType = iota
	FloatType
	StringType
	BoolType
	ListType
	TupleType
	DictType
	FunctionType
	BuiltinType
	NoneType
)

type Value interface{}

type NumberValue struct {
	V float64
}

func (NumberValue) Type() ValueType { return IntType }

type StringValue struct {
	V string
}

func (StringValue) Type() ValueType { return StringType }
func (s StringValue) Len() int {
	return len(s.V)
}

type BoolValue struct {
	V bool
}

func (BoolValue) Type() ValueType { return BoolType }

type NoneValue struct{}

func (NoneValue) Type() ValueType { return NoneType }

type Callable interface {
	Call(args []Value) (Value, error)
}

type FunctionValue struct {
	Name     string
	Args     []string
	Defaults []ast.Expr
	Body     []ast.Stmt
	Env      *Env // closure
}

func (*FunctionValue) Type() ValueType { return FunctionType }

func (f *FunctionValue) Call(args []Value) (Value, error) {
	// TODO: create new env , insert the args, execute body, handle return
	return nil, nil
}

type BuiltinFunction func(args []Value) (Value, error)

type BuiltinValue struct {
	Name string
	Fn   BuiltinFunction
}

func (*BuiltinValue) Type() ValueType { return BuiltinType }

func (b *BuiltinValue) Call(args []Value) (Value, error) {
	return b.Fn(args)
}
