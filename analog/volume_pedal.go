package analog

import (
	"fmt"
	"math"
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


// Loudness
// Input -> i0
// Normalize i0 amplitude using K-filter or A-weighting (avoid different loudness perception by frequency)
// Mean square of normalized frequencies
// Lk = 10 * log10( mean square )

func Loudness(samples []complex128) float64 {
	// normalize samples amplitudes
	// Nomalize()

	lk := 10 * math.Log10(MeanSquare(make([]float64, 10)))

	return lk
}


func Normalize(samples []complex128) {

}

func MeanSquare(data []float64) float64 {
	mean := float64(0.0)
	length := len(data)
	for _, value := range data {
		mean += value / float64(length)
	}
	return mean
}


