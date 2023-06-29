package analog

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type testPedal struct {
	inputs []InputJack
	outputs []OutputJack
	isOn bool
}


func (p *testPedal) GetInputJack() []InputJack {
	return p.inputs
}

func (p *testPedal) GetOutputJack() []OutputJack {
	return p.outputs
}

func (p *testPedal) Toggle() {
	p.isOn = !p.isOn
}

func (p *testPedal) Run(shouldRun *bool) {
}


func TestAddPedal(t *testing.T) {
	pedalboard := NewPedalBoard().(*pedalBoard)

	pedal := &testPedal{
		inputs: []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
		isOn: false,
	}
	
	pedalboard.AddPedal(pedal, 0)

	require.Equal(t, 1, len(pedalboard.pedals))

	require.Equal(t, 1, len(pedal.inputs))
	require.Equal(t, pedalboard.input, pedal.inputs[0])

	require.Equal(t, 1, len(pedal.outputs))
	require.Equal(t, pedalboard.output, pedal.outputs[0])
}

func TestAddPedalToInvalidIndex(t *testing.T) {
	pedalboard := NewPedalBoard().(*pedalBoard)

	pedal := &testPedal{
		inputs: []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
		isOn: false,
	}
	
	result := pedalboard.AddPedal(pedal, -1)

	require.Equal(t, errors.New("Index out of bounds"), result)

	result = pedalboard.AddPedal(pedal, 1000)
	require.Equal(t, errors.New("Index out of bounds"), result)
}

func TestUpdateJackConnections(t *testing.T) {
	pedalboard := NewPedalBoard().(*pedalBoard)

	pedalOne := &testPedal{
		inputs: []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
		isOn: false,
	}
	
	pedalboard.AddPedal(pedalOne, 0)

	require.Equal(t, 1, len(pedalboard.pedals))

	require.Equal(t, 1, len(pedalOne.inputs))
	require.Equal(t, 1, len(pedalOne.outputs))
	require.Equal(t, pedalboard.input, pedalOne.inputs[0])
	require.Equal(t, pedalboard.output, pedalOne.outputs[0])

	pedalLast := &testPedal{
		inputs: []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
		isOn: false,
	}
	
	pedalboard.AddPedal(pedalLast, 1)

	require.Equal(t, 2, len(pedalboard.pedals))
	require.Equal(t, pedalOne, pedalboard.pedals[0])
	require.Equal(t, pedalLast, pedalboard.pedals[1])

	require.Equal(t, 1, len(pedalOne.inputs))
	require.Equal(t, 1, len(pedalOne.outputs))
	require.Equal(t, pedalboard.input, pedalOne.inputs[0])

	require.Equal(t, 1, len(pedalLast.inputs))
	require.Equal(t, 1, len(pedalLast.outputs))
	require.Equal(t, pedalOne.outputs[0], pedalLast.inputs[0])
	require.Equal(t, pedalboard.output, pedalLast.outputs[0])


	newPedalFirst := &testPedal{
		inputs: []InputJack{NewInputJack()},
		outputs: []OutputJack{NewOutputJack()},
		isOn: false,
	}
	
	pedalboard.AddPedal(newPedalFirst, 0)

	require.Equal(t, 3, len(pedalboard.pedals))
	require.Equal(t, newPedalFirst, pedalboard.pedals[0])
	require.Equal(t, pedalOne, pedalboard.pedals[1])
	require.Equal(t, pedalLast, pedalboard.pedals[2])

	require.Equal(t, 1, len(pedalOne.inputs))
	require.Equal(t, 1, len(pedalOne.outputs))

	require.Equal(t, 1, len(newPedalFirst.inputs))
	require.Equal(t, 1, len(newPedalFirst.outputs))
	require.Equal(t, pedalboard.input, newPedalFirst.inputs[0])
	require.Equal(t, pedalOne.inputs[0], newPedalFirst.outputs[0])
}
