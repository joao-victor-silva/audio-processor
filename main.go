package main

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
*/
import "C"
import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/joao-victor-silva/audio-processor/audio"
	"github.com/joao-victor-silva/audio-processor/effect"
)

func main() {
	threshold := flag.Float64("threshold", 0.00012, "threshold")
	min := flag.Float64("min", 0.02, "min")
	max := flag.Float64("max", 0.8, "max")
	// max := flag.Float64("max", 1.0, "max")
	// min := flag.Float64("min", -1.0, "min")
	flag.Parse()

	sdlManager, err := audio.NewSDL()
	if err != nil {
		panic(err)
	}
	defer sdlManager.Close()

	fmt.Println("Input devices:")
	sdlManager.ListAudioDevice(true)
	fmt.Println("\nOutput devices:")
	sdlManager.ListAudioDevice(false)

	mic, err := sdlManager.NewAudioDevice(true)
	if err != nil {
		panic("Counldn't open the mic device. " + err.Error())
	}
	defer mic.Close()
	mic.Unpause()

	headphone, err := sdlManager.NewAudioDevice(false)
	if err != nil {
		panic("Counldn't open the headphone device" + err.Error())
	}
	defer headphone.Close()
	headphone.Unpause()

	if mic.AudioFormat() != headphone.AudioFormat() {
		panic("Couldn't use the same audio format for mic and headphones")
	}

	// compressor := effect.Effect{Min: *min, Max: *max, Threshold: *threshold, LogTail: make([]effect.LogReg, 4096), LastLogRegIndex: 0, Samples: 1024}

	fmt.Println("Min: ", *min)
	fmt.Println("Max: ", *max)
	fmt.Println("Threshold: ", *threshold)
	noiseGate := effect.NoiseGate{Threshold: *threshold, LogTail: make([]effect.LogReg, 4096)}
	defer (func() {
		for _, data := range noiseGate.LogTail[noiseGate.LastLogRegIndex:] {
			fmt.Println(data)
		}
		for _, data := range noiseGate.LogTail[:noiseGate.LastLogRegIndex] {
			fmt.Println(data)
		}
	})()
	upwardCompressor := effect.UpwardCompressor{Min: *min, Threshold: *threshold, Factor: 4.0, LogTail: make([]effect.LogReg, 4096)}
	defer (func() {
		for _, data := range upwardCompressor.LogTail[upwardCompressor.LastLogRegIndex:] {
			fmt.Println(data)
		}
		for _, data := range upwardCompressor.LogTail[:upwardCompressor.LastLogRegIndex] {
			fmt.Println(data)
		}
	})()
	downwardCompressor := effect.DownwardCompressor{Max: *max, LogTail: make([]effect.LogReg, 4096)}
	defer (func() {
		for _, data := range downwardCompressor.LogTail[downwardCompressor.LastLogRegIndex:] {
			fmt.Println(data)
		}
		for _, data := range downwardCompressor.LogTail[:downwardCompressor.LastLogRegIndex] {
			fmt.Println(data)
		}
	})()

	raw := audio.NewProcessor("data.bin")
	defer raw.Close()

	copyFromMic := effect.Copy{}
	go copyFromMic.Process(mic, raw)

	processed := audio.NewProcessor("processed-data.bin")
	defer processed.Close()
	// go compressor.Process(raw, processed)
	noiseGateProcessor := audio.NewProcessor("noisegate-data.bin")
	defer noiseGateProcessor.Close()

	upwardCompressorProcessor := audio.NewProcessor("upward-data.bin")
	defer upwardCompressorProcessor.Close()

	go noiseGate.Process(raw, noiseGateProcessor)
	go upwardCompressor.Process(noiseGateProcessor, upwardCompressorProcessor)
	go downwardCompressor.Process(upwardCompressorProcessor, processed)

	copyToHeadphone := effect.Copy{}
	go copyToHeadphone.Process(processed, headphone)

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <-mainThreadSignals
}
