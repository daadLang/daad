package interpreter

import (
	ast "github.com/daadLang/daad/internals/ast"
	"github.com/daadLang/daad/internals/lexer"
)

func (i *Interpreter) execConstExpr(e *ast.Constant) Value {
	switch v := e.Value.(type) {
	case int:
		return IntValue{V: v}
	case float64:
		return FloatValue{V: v}
	case string:
		return StringValue{V: v}
	case rune:
		return CharValue{V: v}
	case bool:
		return BoolValue{V: v}
	case nil:
		return NoneValue{}
	default:
		panic(newRuntimeError("unknown constant type: %T", e.Value))
	}
}

func (i *Interpreter) execBinOpExpr(e *ast.BinOp) Value {
	left := i.execExpr(e.Left)
	right := i.execExpr(e.Right)

	switch e.Op {
	case lexer.PLUS:
		return castAdd(left, right)

	case lexer.MINUS:
		leftVal, rightVal, isFloat := castNumericOp(left, right, lexer.MINUS)
		if isFloat {
			return FloatValue{V: leftVal - rightVal}
		}
		return IntValue{V: int(leftVal) - int(rightVal)}

	case lexer.MULT:
		leftVal, rightVal, isFloat := castNumericOp(left, right, lexer.MULT)
		if isFloat {
			return FloatValue{V: leftVal * rightVal}
		}
		return IntValue{V: int(leftVal) * int(rightVal)}

	case lexer.DIVIDE:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.DIVIDE)
		return FloatValue{V: leftVal / rightVal}

	case lexer.FLOORDIV:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.FLOORDIV)
		return IntValue{V: int(leftVal / rightVal)}

	case lexer.MOD:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.MOD)
		return IntValue{V: int(leftVal) % int(rightVal)}

	case lexer.POWER:
		leftVal, rightVal, isFloat := castNumericOp(left, right, lexer.POWER)
		result := 1.0
		for j := 0; j < int(rightVal); j++ {
			result *= leftVal
		}
		if isFloat {
			return FloatValue{V: result}
		}
		return IntValue{V: int(result)}

	case lexer.BITWISE_AND:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.BITWISE_AND)
		return IntValue{V: int(leftVal) & int(rightVal)}

	case lexer.BITWISE_OR:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.BITWISE_OR)
		return IntValue{V: int(leftVal) | int(rightVal)}

	case lexer.BITWISE_XOR:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.BITWISE_XOR)
		return IntValue{V: int(leftVal) ^ int(rightVal)}

	case lexer.LSHIFT:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.LSHIFT)
		return IntValue{V: int(leftVal) << uint(rightVal)}

	case lexer.RSHIFT:
		leftVal, rightVal, _ := castNumericOp(left, right, lexer.RSHIFT)
		return IntValue{V: int(leftVal) >> uint(rightVal)}

	default:
		panic(newRuntimeError("unknown binary operator: %v", e.Op))
	}
}

func (i *Interpreter) execUnaryOpExpr(e *ast.UnaryOp) Value {
	operand := i.execExpr(e.Expr)

	switch e.Op {
	case lexer.NOT:
		boolVal, ok := operand.(BoolValue)
		if !ok {
			panic(newTypeError("unary NOT requires boolean operand, got %T", operand))
		}
		return BoolValue{V: !boolVal.V}

	case lexer.MINUS:
		switch v := operand.(type) {
		case IntValue:
			return IntValue{V: -v.V}
		case FloatValue:
			return FloatValue{V: -v.V}
		default:
			panic(newTypeError("unary MINUS requires numeric operand, got %T", operand))
		}

	case lexer.PLUS:
		switch v := operand.(type) {
		case IntValue:
			return v
		case FloatValue:
			return v
		default:
			panic(newTypeError("unary PLUS requires numeric operand, got %T", operand))
		}

	case lexer.BITWISE_NOT:
		intVal, ok := operand.(IntValue)
		if !ok {
			panic(newTypeError("unary BITWISE_NOT requires integer operand, got %T", operand))
		}
		return IntValue{V: ^intVal.V}

	default:
		panic(newRuntimeError("unknown unary operator: %v", e.Op))
	}
}

func (i *Interpreter) execBoolOpExpr(e *ast.BoolOp) Value {
	left := i.execExpr(e.Left)

	switch e.Op {
	case lexer.AND:
		leftBool, ok := left.(BoolValue)
		if !ok {
			panic(newTypeError("AND requires boolean operands, got %T", left))
		}
		if !leftBool.V {
			return BoolValue{V: false}
		}
		right := i.execExpr(e.Right)
		rightBool, ok := right.(BoolValue)
		if !ok {
			panic(newTypeError("AND requires boolean operands, got %T", right))
		}
		return BoolValue{V: rightBool.V}

	case lexer.OR:
		leftBool, ok := left.(BoolValue)
		if !ok {
			panic(newTypeError("OR requires boolean operands, got %T", left))
		}
		if leftBool.V {
			return BoolValue{V: true}
		}
		right := i.execExpr(e.Right)
		rightBool, ok := right.(BoolValue)
		if !ok {
			panic(newTypeError("OR requires boolean operands, got %T", right))
		}
		return BoolValue{V: rightBool.V}

	default:
		panic(newRuntimeError("unknown boolean operator: %v", e.Op))
	}
}

func (i *Interpreter) execCompareExpr(e *ast.Compare) Value {
	left := i.execExpr(e.Left)
	right := i.execExpr(e.Comparator)

	switch e.Op {
	case lexer.EQ:
		return BoolValue{V: compareEqual(left, right)}
	case lexer.NEQ:
		return BoolValue{V: !compareEqual(left, right)}
	case lexer.LESS:
		return BoolValue{V: compareLess(left, right)}
	case lexer.GREATER:
		return BoolValue{V: compareLess(right, left)}
	case lexer.LEQ:
		return BoolValue{V: compareLess(left, right) || compareEqual(left, right)}
	case lexer.GEQ:
		return BoolValue{V: compareLess(right, left) || compareEqual(left, right)}
	case lexer.IN:
		return BoolValue{V: containsValue(right, left)}
	default:
		panic(newRuntimeError("unknown comparison operator: %v", e.Op))
	}
}

func (i *Interpreter) execAssignExpr(e *ast.Assign) Value {
	value := i.execExpr(e.Value)

	switch target := e.Target.(type) {
	case *ast.Name:
		i.env.Set(target.Id, value)
	default:
		panic(newRuntimeError("invalid assignment target: %T", e.Target))
	}

	return value
}

func (i *Interpreter) execSubscriptExpr(e *ast.Subscript) Value {
	container := i.execExpr(e.Value)
	index := i.execExpr(e.Index)

	switch c := container.(type) {
	case ListValue:
		idx, ok := index.(IntValue)
		if !ok {
			panic(newTypeError("list indices must be integers, got %T", index))
		}
		if idx.V < 0 || idx.V >= len(c.Elements) {
			panic(newRuntimeError("list index out of range: %d", idx.V))
		}
		return c.Elements[idx.V]

	case StringValue:
		idx, ok := index.(IntValue)
		if !ok {
			panic(newTypeError("string indices must be integers, got %T", index))
		}
		runes := []rune(c.V)
		if idx.V < 0 || idx.V >= len(runes) {
			panic(newRuntimeError("string index out of range: %d", idx.V))
		}
		return CharValue{V: runes[idx.V]}

	case DictValue:
		key := extractRawValue(index)
		if val, ok := c.Entries[key]; ok {
			return val
		}
		panic(newRuntimeError("key not found in dict"))

	case TupleValue:
		idx, ok := index.(IntValue)
		if !ok {
			panic(newTypeError("tuple indices must be integers, got %T", index))
		}
		if idx.V < 0 || idx.V >= len(c.Elements) {
			panic(newRuntimeError("tuple index out of range: %d", idx.V))
		}
		return c.Elements[idx.V]

	default:
		panic(newTypeError("'%T' is not subscriptable", container))
	}
}

func (i *Interpreter) execListExpr(e *ast.List) Value {
	elements := make([]Value, len(e.Elements))
	for idx, elem := range e.Elements {
		elements[idx] = i.execExpr(elem)
	}
	return ListValue{Elements: elements}
}

func (i *Interpreter) execDictExpr(e *ast.Dict) Value {
	entries := make(map[interface{}]Value)
	for idx := range e.Keys {
		key := extractRawValue(i.execExpr(e.Keys[idx]))
		value := i.execExpr(e.Values[idx])
		entries[key] = value
	}
	return DictValue{Entries: entries}
}

func (i *Interpreter) execTupleExpr(e *ast.Tuple) Value {
	elements := make([]Value, len(e.Elements))
	for idx, elem := range e.Elements {
		elements[idx] = i.execExpr(elem)
	}
	return TupleValue{Elements: elements}
}

func (i *Interpreter) execCallExpr(expr *ast.Call) Value {
	funcValue := i.execExpr(expr.Func)
	if funcValue == nil {
		panic(newRuntimeError("undefined function"))
	}

	// Evaluate positional arguments
	posArgs := make([]Value, len(expr.Args))
	for idx, argExpr := range expr.Args {
		posArgs[idx] = i.execExpr(argExpr)
	}

	// Evaluate keyword arguments
	kwArgs := make(map[string]Value)
	for _, kw := range expr.Kwargs {
		kwArgs[kw.Name] = i.execExpr(kw.Value)
	}

	switch fn := funcValue.(type) {
	case *FunctionValue:
		return i.callFunction(fn, posArgs, kwArgs)
	case *BuiltinValue:
		result, err := fn.Fn(posArgs, kwArgs)
		if err != nil {
			panic(newRuntimeError("%s", err.Error()))
		}
		return result
	case *ClassValue:
		// instantiate object
		obj := &ObjectValue{Class: fn, Attributes: map[string]Value{}}

		// copy non-callable class attributes as initial instance attributes
		for name, v := range fn.Attributes {
			switch v.(type) {
			case *FunctionValue:
				// methods are accessed via attribute lookup and will be bound
			default:
				obj.Attributes[name] = v
			}
		}

		// run constructor if present
		if ctorVal, ok := fn.Attributes["__بناء__"]; ok {
			if ctorFn, ok := ctorVal.(*FunctionValue); ok {
				// call constructor with instance as first arg
				args := make([]Value, 0, len(posArgs)+1)
				args = append(args, obj)
				args = append(args, posArgs...)
				// callFunction returns Value
				i.callFunction(ctorFn, args, kwArgs)
			}
		}

		return obj
	default:
		panic(newTypeError("'%T' is not callable", funcValue))
	}
}

// Attribute access: obj.attr
func (i *Interpreter) execAttributeExpr(e *ast.Attribute) Value {
	container := i.execExpr(e.Value)

	switch c := container.(type) {
	case *ObjectValue:
		// check instance attributes first
		if v, ok := c.Attributes[e.Attr]; ok {
			return v
		}
		// then class attributes (methods or defaults)
		if v, ok := c.Class.Attributes[e.Attr]; ok {
			// if it's a function, bind it to the instance
			if fn, ok := v.(*FunctionValue); ok {
				// return a builtin that will call the function with instance as first arg
				bound := &BuiltinValue{Name: fn.Name, Variadic: false}
				bound.Fn = func(args []Value, kwargs map[string]Value) (Value, error) {
					all := make([]Value, 0, len(args)+1)
					all = append(all, c)
					all = append(all, args...)
					return i.callFunction(fn, all, kwargs), nil
				}
				return bound
			}
			return v
		}
		panic(newRuntimeError("attribute '%s' not found on instance of '%s'", e.Attr, c.Class.Name))

	case *ClassValue:
		if v, ok := c.Attributes[e.Attr]; ok {
			return v
		}
		panic(newRuntimeError("attribute '%s' not found on class '%s'", e.Attr, c.Name))

	default:
		panic(newTypeError("'%T' has no attributes", container))
	}
}

// ===========================================================================
// ===========================================================================
func (i *Interpreter) callFunction(fn *FunctionValue, posArgs []Value, kwArgs map[string]Value) Value {
	// build a map of param name -> index for quick lookup
	paramIndex := make(map[string]int)
	for idx, name := range fn.Params {
		paramIndex[name] = idx
	}

	// create array to hold final argument values (nil means not set)
	finalArgs := make([]Value, len(fn.Params))

	// fill in positional arguments
	for idx, val := range posArgs {
		if idx >= len(fn.Params) {
			panic(newRuntimeError("%s() takes %d argument(s) but %d positional were given",
				fn.Name, len(fn.Params), len(posArgs)))
		}
		finalArgs[idx] = val
	}

	// fill in keyword arguments
	for name, val := range kwArgs {
		idx, exists := paramIndex[name]
		if !exists {
			panic(newRuntimeError("%s() got unexpected keyword argument '%s'", fn.Name, name))
		}
		if finalArgs[idx] != nil {
			panic(newRuntimeError("%s() got multiple values for argument '%s'", fn.Name, name))
		}
		finalArgs[idx] = val
	}

	// fill in defaults for any remaining nil values
	defaultsStart := len(fn.Params) - len(fn.Defaults)
	for idx := range finalArgs {
		if finalArgs[idx] == nil {
			defaultIdx := idx - defaultsStart
			if defaultIdx >= 0 && defaultIdx < len(fn.Defaults) {
				finalArgs[idx] = fn.Defaults[defaultIdx]
			} else {
				// No value and no default - error
				panic(newRuntimeError("%s() missing required argument '%s'",
					fn.Name, fn.Params[idx]))
			}
		}
	}

	funcEnv := NewEnv(fn.Env)
	for idx, paramName := range fn.Params {
		funcEnv.Set(paramName, finalArgs[idx])
	}

	parentEnv := i.env
	i.env = funcEnv

	// execute function body
	var result Value = NoneValue{}
	sig := i.execBlock(fn.Body)
	if sig.IsReturn() {
		result = sig.Value
	}

	i.env = parentEnv

	return result
}
