package interpreter

import ast "github.com/daadLang/daad/internals/ast"

type ValueType int

const (
	IntType ValueType = iota
	FloatType
	CharType
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

func ToValue(v interface{}) Value {
	switch val := v.(type) {
	case int:
		return IntValue{V: val}
	case int64:
		return IntValue{V: int(val)}
	case float64:
		return FloatValue{V: val}
	case float32:
		return FloatValue{V: float64(val)}
	case string:
		return StringValue{V: val}
	case rune:
		return CharValue{V: val}
	case bool:
		return BoolValue{V: val}
	case nil:
		return NoneValue{}
	case Value:
		return val
	default:
		return val
	}
}

type IntValue struct {
	V int
}

func (IntValue) Type() ValueType { return IntType }

type FloatValue struct {
	V float64
}

func (FloatValue) Type() ValueType { return FloatType }

type StringValue struct {
	V string
}

func (StringValue) Type() ValueType { return StringType }
func (s StringValue) Len() int {
	return len(s.V)
}

type CharValue struct {
	V rune
}

func (CharValue) Type() ValueType { return CharType }

type BoolValue struct {
	V bool
}

func (BoolValue) Type() ValueType { return BoolType }

type ListValue struct {
	Elements []Value
}

func (ListValue) Type() ValueType { return ListType }

func (l ListValue) Len() int {
	return len(l.Elements)
}

type TupleValue struct {
	Elements []Value
}

func (TupleValue) Type() ValueType { return TupleType }

func (t TupleValue) Len() int {
	return len(t.Elements)
}

type DictValue struct {
	Entries map[interface{}]Value
}

func (DictValue) Type() ValueType { return DictType }

func (d DictValue) Len() int {
	return len(d.Entries)
}

type NoneValue struct{}

func (NoneValue) Type() ValueType { return NoneType }

type Callable interface {
	Call(args []Value) (Value, error)
}

type FunctionValue struct {
	Name     string
	Params   []string // arg names
	Defaults []Value  // default values for the LAST N args
	Body     []ast.Stmt
	Env      *Env // closure
}

func (*FunctionValue) Type() ValueType { return FunctionType }

func (f *FunctionValue) RequiredCount() int {
	return len(f.Params) - len(f.Defaults)
}

type BuiltinFunction func(args []Value) (Value, error)

type BuiltinValue struct {
	Name string
	Fn   BuiltinFunction
}

func (*BuiltinValue) Type() ValueType { return BuiltinType }
