package interpreter

import "fmt"

func castAdd(left, right Value) Value {
	lStr, lIsStr := left.(string)
	rStr, rIsStr := right.(string)

	if lIsStr || rIsStr {
		if !lIsStr {
			lStr = fmt.Sprintf("%v", left)
		}
		if !rIsStr {
			rStr = fmt.Sprintf("%v", right)
		}
		return lStr + rStr
	}

	// Convert bool to int
	if lBool, ok := left.(bool); ok {
		if lBool {
			left = 1
		} else {
			left = 0
		}
	}
	if rBool, ok := right.(bool); ok {
		if rBool {
			right = 1
		} else {
			right = 0
		}
	}

	lInt, lIsInt := left.(int)
	rInt, rIsInt := right.(int)
	lFloat, lIsFloat := left.(float64)
	rFloat, rIsFloat := right.(float64)

	if lIsFloat || rIsFloat {
		if lIsFloat && rIsFloat {
			return lFloat + rFloat
		} else if lIsFloat && rIsInt {
			return lFloat + float64(rInt)
		} else if lIsInt && rIsFloat {
			return float64(lInt) + rFloat
		}
	} else if lIsInt && rIsInt {
		return lInt + rInt
	}

	panic(newUnsupportedOperationError("+", left, right))
}

func castNumericOp(left, right Value, op string) (float64, float64, bool) {
	// Convert bool to int
	if lBool, ok := left.(bool); ok {
		if lBool {
			left = 1
		} else {
			left = 0
		}
	}
	if rBool, ok := right.(bool); ok {
		if rBool {
			right = 1
		} else {
			right = 0
		}
	}

	lInt, lIsInt := left.(int)
	rInt, rIsInt := right.(int)
	lFloat, lIsFloat := left.(float64)
	rFloat, rIsFloat := right.(float64)

	var leftVal, rightVal float64
	var needFloat bool

	if lIsFloat {
		leftVal = lFloat
		needFloat = true
	} else if lIsInt {
		leftVal = float64(lInt)
	} else {
		panic(newUnsupportedOperationError(op, left, right))
	}

	if rIsFloat {
		rightVal = rFloat
		needFloat = true
	} else if rIsInt {
		rightVal = float64(rInt)
	} else {
		panic(newUnsupportedOperationError(op, left, right))
	}

	return leftVal, rightVal, needFloat || lIsFloat || rIsFloat
}
