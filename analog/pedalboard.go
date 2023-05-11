package analog

type PedalBoard interface{
	AddPedal(pedal Pedal, index int)
	Toggle(index int)
}

type pedalBoard struct {
	input InputJack
	output OutputJack
	pedals []Pedal
}

func NewPedalBoard() PedalBoard {
	return &pedalBoard{
		input: NewInputJack(),
		output: NewOutputJack(),
	}
}

func (p *pedalBoard) AddPedal(pedal Pedal, index int) {
	var input Jack
	var output Jack

	if (index > len(p.pedals)) {
		panic("Index out of bounds")
	}

	if (index == 0) {
		input = p.input
	} else {
		input = p.pedals[index-1].GetOutputJack()[0]
	}

	if (index == len(p.pedals)) {
		output = p.output
	} else {
		output = p.pedals[index-1].GetInputJack()[0]
	}

	inputWire := make(Wire)
	outputWire := make(Wire)

	output.Connect(outputWire)
	pedal.GetOutputJack()[0].Connect(outputWire)

	pedal.GetInputJack()[0].Connect(inputWire)
	input.Connect(inputWire)
	
	// shift right pedals
	p.pedals = append(p.pedals, nil)
	copy(p.pedals[index+1:], p.pedals[index:])

	p.pedals[index] = pedal
}

func (p *pedalBoard) Toggle(index int) {
	if (index > len(p.pedals)) {
		panic("Index out of bounds")
	}

	p.pedals[index].Toggle()
}
