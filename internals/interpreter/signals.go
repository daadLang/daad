package interpreter

const (
	NoSignal = iota
	ReturnSignal
	BreakSignal
	ContinueSignal
	ErrorSignal
)

type Signal struct {
	SignalType int
	Value      Value
}
