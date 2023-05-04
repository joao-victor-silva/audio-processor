package analog

type jack struct {
	wire Wire
}

type InputJack interface {
	ReceiveSignal() Signal
}

type OutputJack interface {
	SendSignal(Signal)
}

func (j *jack) ReceiveSignal() Signal {
	return <- j.wire
}

func (j *jack) SendSignal(signal Signal) {
	j.wire <- signal
}

func NewInputJack() InputJack {
	return &jack{wire: make(Wire)}
}

func NewOutputJack() OutputJack {
	return &jack{wire: make(Wire)}
}
