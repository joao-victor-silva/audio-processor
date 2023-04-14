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

	inputId := flag.Int("inputId", 0, "inputId")
	outputId := flag.Int("outputId", 0, "outputId")
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

	mic, err := sdlManager.NewAudioDevice(*inputId, true)
	if err != nil {
		panic("Counldn't open the mic device. " + err.Error())
	}
	defer mic.Close()
	mic.Unpause()

	headphone, err := sdlManager.NewAudioDevice(*outputId, false)
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
	downwardCompressor := effect.DownwardCompressor{Max: *max, MaxExpected: 0.00200, Factor: 0.95, LogTail: make([]effect.LogReg, 4096)}
	defer (func() {
		for _, data := range downwardCompressor.LogTail[downwardCompressor.LastLogRegIndex:] {
			fmt.Println(data)
		}
		for _, data := range downwardCompressor.LogTail[:downwardCompressor.LastLogRegIndex] {
			fmt.Println(data)
		}
	})()
	levelLogger := effect.LevelLogger{}
	defer levelLogger.Print()

	levelNormalizer := effect.LevelNormalizer{Min: 0.000007, Max: 0.00300, Dynamic: false}

	raw := audio.NewProcessor("data.bin")
	defer raw.Close()

	normalized := audio.NewProcessor("normalized.bin")
	defer normalized.Close()

	processed := audio.NewProcessor("processed-data.bin")
	defer processed.Close()

	copyFromMic := effect.Copy{}

	noiseGateProcessor := audio.NewProcessor("noisegate-data.bin")
	defer noiseGateProcessor.Close()

	go copyFromMic.Process(mic, raw)
	go levelNormalizer.Process(raw, normalized)
	go noiseGate.Process(normalized, noiseGateProcessor)
	go levelLogger.Process(noiseGateProcessor, processed)

	// // go compressor.Process(raw, processed)
	//
	// upwardCompressorProcessor := audio.NewProcessor("upward-data.bin")
	// defer upwardCompressorProcessor.Close()
	//
	// go upwardCompressor.Process(noiseGateProcessor, upwardCompressorProcessor)
	// go downwardCompressor.Process(upwardCompressorProcessor, processed)

	copyToHeadphone := effect.Copy{}
	go copyToHeadphone.Process(processed, headphone)

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <-mainThreadSignals
}
