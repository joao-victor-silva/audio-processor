package effect

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
*/
import "C"
import (
	_ "encoding/binary"
	"fmt"
	_ "io"
	"math"
	_ "os"

	"github.com/joao-victor-silva/audio-processor/audio"
)

type Float64 float64
type Float32 float32

func (f Float64) String() string {
	return fmt.Sprintf("%.4f", f*1000)
}

func (f Float32) String() string {
	return fmt.Sprintf("%.4f", f*1000)
}

type LogReg struct {
	Volume  Float64
	Average Float64
	Delta   Float64
	Factor  Float64
	State   string
	Before  Float32
	After   Float32
}

func (r LogReg) String() string {
	return fmt.Sprint("v: ", r.Volume, " avg: ", r.Average, " d: ", r.Delta, " f: ", r.Factor, " state: ", r.State, " b: ", r.Before, " a: ", r.After)
}

type Effect struct {
	Min       float64
	Max       float64
	Threshold float64
	LogTail   []LogReg
}

type Copy struct{}

// type GainEffect struct {
// 	Gain float64
// }
//
// func (g *GainEffect) Process(inputDevice, outputDevice audio.AudioDevice) {
// 	samples := make([]float64, 256)
// 	i := 0
// 	for outputDevice.IsChannelOpen() {
// 		samples[i] = float64(inputDevice.ReadData())
//
// 		average := 0.0
// 		for _, sample := range samples {
// 			average += sample
// 		}
// 		average /= float64(len(samples))
//
// 		volume := getVolume(samples)
//
// 		if (samples[i] > 0) {
// 			samples[i] += g.Gain * volume
// 		} else {
// 			samples[i] -= g.Gain * volume
// 		}
// 	}
// }

func (effect *Effect) Process(inputDevice audio.AudioProcessor, outputDevice audio.AudioProcessor) {
	samples := make([]float64, 256)
	i := 0
	j := 0
	for outputDevice.IsChannelOpen() {
		// begin := time.Now()
		dataBeforeEffect := inputDevice.ReadData()

		samples[i] = float64(dataBeforeEffect)
		average := 0.0
		for _, sample := range samples {
			average += sample
		}
		average /= float64(len(samples))

		volume := 0.0
		for _, sample := range samples {
			volume += math.Pow(sample-average, 2)
		}
		volume = math.Sqrt(volume) / float64(len(samples))

		var dataAfterEffect float32
		var delta float64
		var factor float64
		var state string

		switch {
		case volume < effect.Threshold:
			dataAfterEffect = float32(average)
			state = "Threshold"
		case volume < effect.Min:
			//amp
			delta = float64(dataBeforeEffect) - average
			factor = effect.Min / volume

			dataAfterEffect = float32(average + (delta * factor))
			state = "Below min"
		case volume < effect.Max:
			state = "None"
			dataAfterEffect = dataBeforeEffect
		default:
			//reduce
			delta = float64(dataBeforeEffect) - average
			factor = effect.Max / volume

			dataAfterEffect = float32(average + (delta * factor))
			state = "Above max"
		}
		_ = state
		// fmt.Println("v:", volume, "avg:", average, "d:", delta, "f:", factor, "s:", state, "b:", dataBeforeEffect, "a:", dataAfterEffect)
		logReg := LogReg{Volume: Float64(volume), Average: Float64(average), Delta: Float64(delta), Factor: Float64(factor), State: state, Before: Float32(dataBeforeEffect), After: Float32(dataAfterEffect)}
		effect.LogTail[j] = logReg
		j += 1
		j = j & (len(effect.LogTail) - 1)

		outputDevice.WriteData(dataAfterEffect)
		i += 1
		i = i & (len(samples) - 1)

		// end := time.Now()
		// elapsed := end.Sub(begin)
		// fmt.Println(elapsed.String())
		// fmt.Println("i:", i)
	}
}

func getVolume(samples []float64) float64 {
	average := 0.0
	for _, sample := range samples {
		average += sample
	}
	average /= float64(len(samples))

	volume := 0.0
	for _, sample := range samples {
		volume += math.Pow(sample-average, 2)
	}
	volume = math.Sqrt(volume) / float64(len(samples))

	return volume
}

func (c *Copy) Process(inputDevice audio.AudioProcessor, outputDevice audio.AudioProcessor) {
	for outputDevice.IsChannelOpen() {
		// begin := time.Now()
		data := inputDevice.ReadData()
		// binaryData := make([]byte, 4)
		// binary.LittleEndian.PutUint32(binaryData, math.Float32bits(data))
		// c.File.Write(binaryData)
		// fmt.Println("value:", data)
		outputDevice.WriteData(data)
		// end := time.Now()
		// elapsed := end.Sub(begin)
		// fmt.Printf("%v\n", elapsed.String())
	}
}
