package effect

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
*/
import "C"
import (
	"encoding/binary"
	"fmt"
	"io"
	"math"
	"time"

	"github.com/joao-victor-silva/audio-processor/audio"
)

type ProcessData interface {
	Process(input <- chan byte, output chan <- byte)
}

type Copy struct{
	File io.Writer
}
type Effect struct {
	Min float64
	Max float64
	Threshold float64
}
//
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

func (effect *Effect) Process(inputDevice audio.AudioDevice , outputDevice audio.AudioDevice) {
	samples := make([]float64, 256)
	i := 0
	for outputDevice.IsChannelOpen() {
		begin := time.Now()
		dataBeforeEffect := inputDevice.ReadData()

		samples[i] = float64(dataBeforeEffect)
		average := 0.0
		for _, sample := range samples {
			average += sample
		}
		average /= float64(len(samples))

		volume := 0.0
		for _, sample := range samples {
			volume += math.Pow(sample - average, 2)
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
		outputDevice.WriteData(dataAfterEffect)
		i += 1
		i = i & (len(samples) - 1)
		end := time.Now()
		elapsed := end.Sub(begin)
		fmt.Println(elapsed.String())
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
		volume += math.Pow(sample - average, 2)
	}
	volume = math.Sqrt(volume) / float64(len(samples))

	return volume
}

func (c *Copy) Process(inputDevice audio.AudioDevice , outputDevice audio.AudioDevice) {
	for outputDevice.IsChannelOpen() {
		begin := time.Now()
		data := inputDevice.ReadData()
		binaryData := make([]byte, 4)
		binary.LittleEndian.PutUint32(binaryData, math.Float32bits(data))
		c.File.Write(binaryData)
		// fmt.Println("value:", data)
		outputDevice.WriteData(data)
		end := time.Now()
		elapsed := end.Sub(begin)
		fmt.Printf("%v\n", elapsed.String())
	}
}
