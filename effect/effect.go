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
	"time"

	"github.com/joao-victor-silva/audio-processor/audio"
	// "golang.org/x/text/cases"
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
	Timestamp time.Duration
	Volume    Float64
	Average   Float64
	Delta     Float64
	Factor    Float64
	State     string
	Before    Float32
	After     Float32
}

func (r LogReg) String() string {
	return fmt.Sprint(r.Timestamp, "-> state: ", r.State,  " volume: ", r.Volume, " factor: ", r.Factor, " before: ", r.Before, " after: ", r.After)
}

type Effect struct {
	Min             float64
	Max             float64
	Threshold       float64
	LogTail         []LogReg
	LastLogRegIndex int
	Samples         int
	Attack          int
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

type CompressorState int
const (
	Attack CompressorState = iota
	Compressing
	Release
	Threshold

	Nop
)

type DownwardCompressor struct {
	Max float64
	Factor float64
	LogTail         []LogReg
	LastLogRegIndex int
}

func (c *DownwardCompressor) Process(inputDevice audio.AudioProcessor, outputDevice audio.AudioProcessor) {
	startTime := time.Now()
	for outputDevice.IsChannelOpen() {
		sample := inputDevice.ReadData()
		dataBeforeEffect, volume := sample.Value, sample.Volume

		if (volume < c.Max) {
			outputDevice.WriteData(sample)
			continue
		}

		factor := c.Max / volume

		// TODO: Decide log strategy
		dataAfterEffect := float32((float64(dataBeforeEffect) * factor))

		logReg := LogReg{Timestamp: time.Now().Sub(startTime), Volume: Float64(volume), State: "3 - Downward Compressor", Factor: Float64(factor), Before: Float32(dataBeforeEffect), After: Float32(dataAfterEffect)}
		c.LogTail[c.LastLogRegIndex] = logReg
		c.LastLogRegIndex += 1
		c.LastLogRegIndex = c.LastLogRegIndex & (len(c.LogTail) - 1)
		outputDevice.WriteData(audio.Sample{Value: dataAfterEffect, Volume: c.Max})
	}
}

type UpwardCompressor struct {
	Min float64
	Threshold float64
	Factor float64
	LogTail         []LogReg
	LastLogRegIndex int
}

func (c *UpwardCompressor) Process(inputDevice audio.AudioProcessor, outputDevice audio.AudioProcessor) {
	startTime := time.Now()
	for outputDevice.IsChannelOpen() {
		sample := inputDevice.ReadData()
		dataBeforeEffect, volume := sample.Value, sample.Volume

		if (volume > c.Min || volume <= c.Threshold) {
			outputDevice.WriteData(sample)
			continue
		}


		// factor := c.Min / c.Threshold
		factor := c.Factor

		delta := (c.Min - volume) / (c.Min - c.Threshold) // 0.0 -> 1.0
		// Options
		// delta := ((c.Min - volume) - c.Threshold) / c.Min  // 0.0 -> 1.0
		// delta := (c.Min - volume) - c.Threshold/ (c.Min - c.Threshold) // 0.0 -> 1.0
		// dataAfter := (1.0 + (delta * factor)) * float64(dataBeforeEffect) // Pode passar do minimo

		// dataAfter := dataBeforeEffect * (float32(c.Factor) * float32(delta))


		// volume -> c.Min | factor -> 0.0 (algo somado)
		// volume -> c.Threshold | factor -> 100.0 (algo somado)

		// TODO: Decide log strategy
		dataAfterEffect := float32((1.0 + (delta * factor)) * float64(dataBeforeEffect))

		logReg := LogReg{Timestamp: time.Now().Sub(startTime), Volume: Float64(volume), State: "2 - Upward Compressor", Factor: Float64(factor), Before: Float32(dataBeforeEffect), After: Float32(dataAfterEffect)}
		c.LogTail[c.LastLogRegIndex] = logReg
		c.LastLogRegIndex += 1
		c.LastLogRegIndex = c.LastLogRegIndex & (len(c.LogTail) - 1)
		outputDevice.WriteData(audio.Sample{Value: dataAfterEffect, Volume: c.Min})
	}
}

type NoiseGate struct {
	Threshold float64
	LogTail         []LogReg
	LastLogRegIndex int
}

func (n *NoiseGate) Process(inputDevice audio.AudioProcessor, outputDevice audio.AudioProcessor) {
	startTime := time.Now()

	for outputDevice.IsChannelOpen() {
		sample := inputDevice.ReadData()
		dataBeforeEffect, volume := sample.Value, sample.Volume

		if (volume > n.Threshold) {
			outputDevice.WriteData(sample)
			continue
		}

		logReg := LogReg{Timestamp: time.Now().Sub(startTime), Volume: Float64(volume), State: "1 - Noise gate", Before: Float32(dataBeforeEffect), After: Float32(0.0)}
		n.LogTail[n.LastLogRegIndex] = logReg
		n.LastLogRegIndex += 1
		n.LastLogRegIndex = n.LastLogRegIndex & (len(n.LogTail) - 1)
		outputDevice.WriteData(audio.Sample{Value: float32(0.0), Volume: float64(0.0)})
	}
}

func (effect *Effect) Process(inputDevice audio.AudioProcessor, outputDevice audio.AudioProcessor) {
	samples := make([]float64, effect.Samples)
	i := 0
	effect.LastLogRegIndex = 0
	startTime := time.Now()

	// samplesWithAttck := 0
	// alpha := 1 / effect.Attack
	effectState := Nop


	for outputDevice.IsChannelOpen() {
		// begin := time.Now()
		sample := inputDevice.ReadData()
		dataBeforeEffect, volume := sample.Value, sample.Volume

		samples[i] = float64(dataBeforeEffect)
		average := 0.0
		for _, sample := range samples {
			average += sample
		}
		average /= float64(len(samples))

		var dataAfterEffect float32
		var delta float64
		var factor float64
		var state string

		var desiredState CompressorState	

		switch {
		case volume < effect.Threshold:
			dataAfterEffect = float32(average)
			state = "Threshold"
			desiredState = Nop
		case volume < effect.Min:
			//amp
			delta = float64(dataBeforeEffect) - average
			factor = effect.Min / volume

			dataAfterEffect = float32(average + (delta * factor))
			state = "Below min"
			desiredState = Compressing
		case volume < effect.Max:
			state = "None"
			dataAfterEffect = dataBeforeEffect
			desiredState = Compressing
		default:
			//reduce
			delta = float64(dataBeforeEffect) - average
			factor = effect.Max / volume

			dataAfterEffect = float32(average + (delta * factor))
			state = "Above max"
			desiredState = Compressing
		}
		_ = state


		switch {
		case effectState == Nop && desiredState != Nop:
			effectState = Attack
		case effectState != Nop && desiredState == Nop:
			effectState = Release
		}

		switch {
		case effectState == Attack:
			// Do the attach with the alpha blending and check need to change state for Compressing
		case effectState == Compressing:
			// Compressing (a.k.a do nothing)
		case effectState == Release:
			// TODO
		}

		// fmt.Println("v:", volume, "avg:", average, "d:", delta, "f:", factor, "s:", state, "b:", dataBeforeEffect, "a:", dataAfterEffect)
		logReg := LogReg{Timestamp: time.Now().Sub(startTime), Volume: Float64(volume), Average: Float64(average), Delta: Float64(delta), Factor: Float64(factor), State: state, Before: Float32(dataBeforeEffect), After: Float32(dataAfterEffect)}
		effect.LogTail[effect.LastLogRegIndex] = logReg
		effect.LastLogRegIndex += 1
		effect.LastLogRegIndex = effect.LastLogRegIndex & (len(effect.LogTail) - 1)

		outputDevice.WriteData(audio.Sample{Value: dataAfterEffect, Volume: volume})
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
