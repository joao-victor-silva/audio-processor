package examples

import "github.com/joao-victor-silva/audio-processor/analog"

func main() {
	// Add pedal in position to rewire other pedals
	pedalBoard := analog.NewPedalBoard()

	effectOnePedal := analog.NewDummyPedal()

	effectTwoPedal := analog.NewDummyPedal()

	// Add pedal in container
	// Connect raw input to pedal input jack (mic -> pedal)
	// Connect pedal output jack to raw output (pedal -> headphone)
	pedalBoard.AddPedal(effectOnePedal, 0)

	// Add pedal in container
	// Disconnect effectOnePedal output from raw output
	// Connect effectOnePedal output to pedal input jack (pedalOne -> pedalTwo)
	// Connect pedal output jack to raw output (pedalTwo -> headphone)
	pedalBoard.AddPedal(effectTwoPedal, 1)
}
