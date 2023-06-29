package analog

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type testEffect struct{}

func (e *testEffect) Process(input []Signal) []Signal {
	output := make([]Signal, 0, len(input))
	for _, signal := range input {
		switch value := signal.(type) {
		case Float32:
			output = append(output, value+10)
		default:
			output = append(output, value)
		}
	}

	return output
}

func TestPedalPassthrough(t *testing.T) {
	pedal := NewPedal(&testEffect{})
	inputWire := make(Wire, 10)
	outputWire := make(Wire, 10)

	inputSignal := make([]Signal, 10)
	for i := 0; i < len(inputSignal); i++ {
		inputWire <- Float32(i)
	}

	pedal.GetInputJack()[0].Connect(inputWire)
	pedal.GetOutputJack()[0].Connect(outputWire)

	shouldRun := true
	go pedal.Run(&shouldRun)

	for i := range inputSignal {
		result := <-pedal.GetOutputJack()[0].GetWire()
		value, ok := result.(Float32)
		require.True(t, ok)
		require.Equal(t, value, Float32(i))
	}

	shouldRun = false
}

func TestPedalEffect(t *testing.T) {
	pedal := NewPedal(&testEffect{})
	inputWire := make(Wire, 10)
	outputWire := make(Wire, 10)

	inputSignal := make([]Signal, 10)
	for i := 0; i < len(inputSignal); i++ {
		inputWire <- Float32(i)
	}

	pedal.GetInputJack()[0].Connect(inputWire)
	pedal.GetOutputJack()[0].Connect(outputWire)

	pedal.Toggle()
	shouldRun := true
	go pedal.Run(&shouldRun)

	for i := range inputSignal {
		result := <-pedal.GetOutputJack()[0].GetWire()
		value, ok := result.(Float32)
		require.True(t, ok)
		require.Equal(t, value, Float32(i+10))
	}

	shouldRun = false
}
