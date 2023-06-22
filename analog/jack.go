package analog

type jack struct {
	wire Wire
}

type Jack interface {
	Connect(wire Wire)
	GetWire() Wire
}

type InputJack interface {
	Jack
	ReceiveSignal() Signal
	BufferedReceiveSignal(size int) []Signal
}

type OutputJack interface {
	Jack
	SendSignal(Signal)
	BufferedSendSignal([]Signal)
}

// type InputOutputJack interface {
// 	InputJack
// 	OutputJack
// }

func (j *jack) Connect(wire Wire) {
	open := true

	select {
	case _, open = <-j.wire:
	default:
	}
	if open {
		close(j.wire)
	}

	j.wire = wire
}

func (j *jack) GetWire() Wire {
	return j.wire
}

func (j *jack) ReceiveSignal() Signal {
	return <-j.wire
}

func (j *jack) BufferedReceiveSignal(size int) []Signal {
	result := make([]Signal, 0, size)
	for i := 0; i < size; i++ {
		var value Signal
		select {
		case value = <-j.wire:
		default:
		}
		result = append(result, value)
	}
	return result
}

func (j *jack) SendSignal(signal Signal) {
	j.wire <- signal
}

func (j *jack) BufferedSendSignal(signals []Signal) {
	for _, signal := range signals {
		j.wire <- signal
	}
}

func NewInputJack() InputJack {
	return &jack{wire: make(Wire, 1024)}
}

func NewOutputJack() OutputJack {
	return &jack{wire: make(Wire, 1024)}
}
