package interpreter

type SignalType int

const (
	NoSignal SignalType = iota
	ReturnSignal
	BreakSignal
	ContinueSignal
	ErrorSignal
)

type Signal struct {
	Type  SignalType
	Value Value
	Err   error
}

func NewNoSignal() Signal             { return Signal{Type: NoSignal} }
func NewReturnSignal(v Value) Signal  { return Signal{Type: ReturnSignal, Value: v} }
func NewBreakSignal() Signal          { return Signal{Type: BreakSignal} }
func NewContinueSignal() Signal       { return Signal{Type: ContinueSignal} }
func NewErrorSignal(err error) Signal { return Signal{Type: ErrorSignal, Err: err} }

func (s Signal) IsReturn() bool   { return s.Type == ReturnSignal }
func (s Signal) IsBreak() bool    { return s.Type == BreakSignal }
func (s Signal) IsContinue() bool { return s.Type == ContinueSignal }
func (s Signal) IsError() bool    { return s.Type == ErrorSignal }
