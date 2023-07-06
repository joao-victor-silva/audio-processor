package analog

import "errors"

type connectedPedal struct {
	Pedal Pedal
	ShouldRun *bool
}

type PedalBoard interface{
	AddPedal(pedal Pedal, index int) error
	RemovePedal(index int) error
	Toggle(index int) error
	InputConnect(wire Wire)
	OutputConnect(wire Wire)
}

type pedalBoard struct {
	input InputJack
	output OutputJack
	pedals []*connectedPedal
}

func NewPedalBoard() PedalBoard {
	return &pedalBoard{
		input: NewInputJack(),
		output: NewOutputJack(),
	}
}

func (p *pedalBoard) AddPedal(pedal Pedal, index int) error {
	var input Jack
	var output Jack

	if (index > len(p.pedals) || index < 0) {
		return errors.New("Index out of bounds")
	}

	if (index == 0) {
		input = p.input
	} else {
		input = p.pedals[index-1].Pedal.GetOutputJack()[0]
	}

	if (index == len(p.pedals)) {
		output = p.output
	} else {
		output = p.pedals[index].Pedal.GetInputJack()[0]
	}

	inputWire := make(Wire, 1024)
	outputWire := make(Wire, 1024)

	output.Connect(outputWire)
	pedal.GetOutputJack()[0].Connect(outputWire)

	pedal.GetInputJack()[0].Connect(inputWire)
	input.Connect(inputWire)
	
	// shift right pedals
	p.pedals = append(p.pedals, nil)
	copy(p.pedals[index+1:], p.pedals[index:])

	shouldRun := true
	connectedPedal := connectedPedal{Pedal: pedal, ShouldRun: &shouldRun}
	p.pedals[index] = &connectedPedal
	go pedal.Run(connectedPedal.ShouldRun)

	return nil
}


func (p *pedalBoard) RemovePedal(index int) error {
	var input Jack
	var output Jack

	if (index > len(p.pedals) || index < 0) {
		return errors.New("Index out of bounds")
	}

	*p.pedals[index].ShouldRun = false

	if (index == 0) {
		input = p.input
	} else {
		input = p.pedals[index-1].Pedal.GetOutputJack()[0]
	}

	if (index == len(p.pedals)) {
		output = p.output
	} else {
		output = p.pedals[index].Pedal.GetInputJack()[0]
	}

	wire := make(Wire, 1024)

	output.Connect(wire)
	input.Connect(wire)
	
	// shift left pedals
	copy(p.pedals[index:], p.pedals[index+1:])
	p.pedals = p.pedals[:len(p.pedals) - 1]

	return nil
}

func (p *pedalBoard) Toggle(index int) error {
	if (index > len(p.pedals) || index < 0) {
		return errors.New("Index out of bounds")
	}

	p.pedals[index].Pedal.Toggle()

	return nil
}


func (p *pedalBoard) InputConnect(wire Wire) {
	p.input.Connect(wire)
	if len(p.pedals) > 0 {
		p.pedals[0].Pedal.GetInputJack()[0].Connect(wire)
	}
}

func (p *pedalBoard) OutputConnect(wire Wire) {
	p.output.Connect(wire)
	if len(p.pedals) > 0 {
		p.pedals[len(p.pedals) - 1].Pedal.GetOutputJack()[0].Connect(wire)
	}
}
