package interpreter

import (
	"fmt"

	"github.com/daadLang/daad/internals/lexer"
)

func castAdd(left, right Value) Value {
	// str has highest priority
	lStr, lIsStr := left.(StringValue)
	rStr, rIsStr := right.(StringValue)

	if lIsStr || rIsStr {
		leftString := ""
		rightString := ""

		if lIsStr {
			leftString = lStr.V
		} else {
			leftString = fmt.Sprintf("%v", extractRawValue(left))
		}

		if rIsStr {
			rightString = rStr.V
		} else {
			rightString = fmt.Sprintf("%v", extractRawValue(right))
		}

		return StringValue{V: leftString + rightString}
	}

	// bool to int
	if lBool, ok := left.(BoolValue); ok {
		if lBool.V {
			left = IntValue{V: 1}
		} else {
			left = IntValue{V: 0}
		}
	}
	if rBool, ok := right.(BoolValue); ok {
		if rBool.V {
			right = IntValue{V: 1}
		} else {
			right = IntValue{V: 0}
		}
	}

	//  float64 > int
	lInt, lIsInt := left.(IntValue)
	rInt, rIsInt := right.(IntValue)
	lFloat, lIsFloat := left.(FloatValue)
	rFloat, rIsFloat := right.(FloatValue)

	if lIsFloat || rIsFloat {
		if lIsFloat && rIsFloat {
			return FloatValue{V: lFloat.V + rFloat.V}
		} else if lIsFloat && rIsInt {
			return FloatValue{V: lFloat.V + float64(rInt.V)}
		} else if lIsInt && rIsFloat {
			return FloatValue{V: float64(lInt.V) + rFloat.V}
		}
	} else if lIsInt && rIsInt {
		return IntValue{V: lInt.V + rInt.V}
	}

	panic(newUnsupportedOperationError(lexer.PLUS, left, right))
}

func extractRawValue(v Value) interface{} {
	switch val := v.(type) {
	case IntValue:
		return val.V
	case FloatValue:
		return val.V
	case StringValue:
		return val.V
	case CharValue:
		return val.V
	case BoolValue:
		return val.V
	case NoneValue:
		return nil
	default:
		return v
	}
}

func compareEqual(left, right Value) bool {
	lRaw := extractRawValue(left)
	rRaw := extractRawValue(right)

	switch l := left.(type) {
	case IntValue:
		switch r := right.(type) {
		case IntValue:
			return l.V == r.V
		case FloatValue:
			return float64(l.V) == r.V
		}
	case FloatValue:
		switch r := right.(type) {
		case IntValue:
			return l.V == float64(r.V)
		case FloatValue:
			return l.V == r.V
		}
	case StringValue:
		if r, ok := right.(StringValue); ok {
			return l.V == r.V
		}
	case BoolValue:
		if r, ok := right.(BoolValue); ok {
			return l.V == r.V
		}
	case NoneValue:
		_, ok := right.(NoneValue)
		return ok
	}

	return lRaw == rRaw
}

// check if left < right
func compareLess(left, right Value) bool {
	switch l := left.(type) {
	case IntValue:
		switch r := right.(type) {
		case IntValue:
			return l.V < r.V
		case FloatValue:
			return float64(l.V) < r.V
		}
	case FloatValue:
		switch r := right.(type) {
		case IntValue:
			return l.V < float64(r.V)
		case FloatValue:
			return l.V < r.V
		}
	case StringValue:
		if r, ok := right.(StringValue); ok {
			return l.V < r.V
		}
	}

	panic(newTypeError("'<' not supported between '%T' and '%T'", left, right))
}

// handle 'in' operator
func containsValue(container, item Value) bool {
	switch c := container.(type) {
	case ListValue:
		for _, elem := range c.Elements {
			if compareEqual(elem, item) {
				return true
			}
		}
		return false
	case TupleValue:
		for _, elem := range c.Elements {
			if compareEqual(elem, item) {
				return true
			}
		}
		return false
	case StringValue:
		if s, ok := item.(StringValue); ok {
			return len(c.V) > 0 && len(s.V) > 0 &&
				containsString(c.V, s.V)
		}
		if ch, ok := item.(CharValue); ok {
			return containsRune(c.V, ch.V)
		}
		return false
	case DictValue:
		key := extractRawValue(item)
		_, ok := c.Entries[key]
		return ok
	}

	panic(newTypeError("argument of type '%T' is not iterable", container))
}

func containsString(haystack, needle string) bool {
	return len(needle) <= len(haystack) &&
		(haystack == needle ||
			len(haystack) >= len(needle) &&
				searchString(haystack, needle))
}

func searchString(haystack, needle string) bool {
	for i := 0; i <= len(haystack)-len(needle); i++ {
		if haystack[i:i+len(needle)] == needle {
			return true
		}
	}
	return false
}

func containsRune(s string, r rune) bool {
	for _, c := range s {
		if c == r {
			return true
		}
	}
	return false
}

func castNumericOp(left, right Value, op lexer.TokenType) (float64, float64, bool) {
	// bool to int
	if lBool, ok := left.(BoolValue); ok {
		if lBool.V {
			left = IntValue{V: 1}
		} else {
			left = IntValue{V: 0}
		}
	}
	if rBool, ok := right.(BoolValue); ok {
		if rBool.V {
			right = IntValue{V: 1}
		} else {
			right = IntValue{V: 0}
		}
	}

	lInt, lIsInt := left.(IntValue)
	rInt, rIsInt := right.(IntValue)
	lFloat, lIsFloat := left.(FloatValue)
	rFloat, rIsFloat := right.(FloatValue)

	var leftVal, rightVal float64
	var needFloat bool

	if lIsFloat {
		leftVal = lFloat.V
		needFloat = true
	} else if lIsInt {
		leftVal = float64(lInt.V)
	} else {
		panic(newUnsupportedOperationError(op, left, right))
	}

	if rIsFloat {
		rightVal = rFloat.V
		needFloat = true
	} else if rIsInt {
		rightVal = float64(rInt.V)
	} else {
		panic(newUnsupportedOperationError(op, left, right))
	}

	return leftVal, rightVal, needFloat || lIsFloat || rIsFloat
}
