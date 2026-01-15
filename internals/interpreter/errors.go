package interpreter

import "fmt"

type RuntimeError struct {
	Message string
}

func (e *RuntimeError) Error() string {
	return fmt.Sprintf("RuntimeError: %s", e.Message)
}

func newRuntimeError(format string, args ...interface{}) *RuntimeError {
	return &RuntimeError{
		Message: fmt.Sprintf(format, args...),
	}
}

type TypeError struct {
	Message string
}

func (e *TypeError) Error() string {
	return fmt.Sprintf("TypeError: %s", e.Message)
}

func newTypeError(format string, args ...interface{}) *TypeError {
	return &TypeError{
		Message: fmt.Sprintf(format, args...),
	}
}

type UnsupportedOperationError struct {
	Op    string
	Left  interface{}
	Right interface{}
}

func (e *UnsupportedOperationError) Error() string {
	return fmt.Sprintf("TypeError: unsupported operand type(s) for %s: '%T' and '%T'", e.Op, e.Left, e.Right)
}

func newUnsupportedOperationError(op string, left, right interface{}) *UnsupportedOperationError {
	return &UnsupportedOperationError{
		Op:    op,
		Left:  left,
		Right: right,
	}
}
