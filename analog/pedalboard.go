package analog

import "errors"

type PedalBoard interface{
	AddPedal(pedal Pedal, index int) error
	Toggle(index int) error
	InputConnect(wire Wire)
	OutputConnect(wire Wire)
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

func (p *pedalBoard) AddPedal(pedal Pedal, index int) error {
	var input Jack
	var output Jack

	if (index > len(p.pedals) || index < 0) {
		return errors.New("Index out of bounds")
	}

	if (index == 0) {
		input = p.input
	} else {
		input = p.pedals[index-1].GetOutputJack()[0]
	}

	if (index == len(p.pedals)) {
		output = p.output
	} else {
		output = p.pedals[index].GetInputJack()[0]
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

	return nil
}

func (p *pedalBoard) Toggle(index int) error {
	if (index > len(p.pedals) || index < 0) {
		return errors.New("Index out of bounds")
	}

	p.pedals[index].Toggle()

	return nil
}


func (p *pedalBoard) InputConnect(wire Wire) {
	p.input.Connect(wire)
}

func (p *pedalBoard) OutputConnect(wire Wire) {
	p.output.Connect(wire)
}
