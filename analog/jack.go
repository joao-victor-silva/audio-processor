package analog

type jack struct {
	wire Wire
}

type Jack interface {
	Connect(wire Wire)
}

type InputJack interface {
	Jack
	ReceiveSignal() Signal
}

type OutputJack interface {
	Jack
	SendSignal(Signal)
}

// type InputOutputJack interface {
// 	InputJack
// 	OutputJack
// }

func (j *jack) Connect(wire Wire) {
	_, open := <- j.wire
	if open {
		close(j.wire)
	}

	j.wire = wire
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
