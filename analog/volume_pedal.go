package analog

import (
	"fmt"
	"math/cmplx"

	"github.com/joao-victor-silva/audio-processor/audio"
	"github.com/mjibson/go-dsp/fft"
)

type VolumeEffect struct {
}

func (e *VolumeEffect) Process(signals []Signal) []Signal {
	output := make([]float64, 0, len(signals))
	for _, signal := range signals {
		switch value := signal.(type) {
		case Float32:
			output = append(output, float64(value))
		case *audio.Sample:
			output = append(output, float64(value.Value))
		default:
			panic("Not implemented for types diferent of float32")
		}
	}

	frequency := fft.FFTReal(output)

	average := 0.0
	for _, i := range frequency {
		fmt.Println(cmplx.Abs(i))
		average += cmplx.Abs(i)
	}

	fmt.Println("average:", average / float64(len(frequency)), "\n\n")

	return signals
}
