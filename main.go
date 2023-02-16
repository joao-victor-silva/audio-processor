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

	compressor := effect.Effect{Min: *min, Max: *max, Threshold: *threshold, LogTail: make([]effect.LogReg, 4096), LastLogRegIndex: 0}
	defer (func() {
		for _, data := range compressor.LogTail[compressor.LastLogRegIndex:] {
			fmt.Println(data)
		}
		for _, data := range compressor.LogTail[:compressor.LastLogRegIndex] {
			fmt.Println(data)
		}
	})()

	raw := audio.NewProcessor("data.bin")
	defer raw.Close()

	copyFromMic := effect.Copy{}
	go copyFromMic.Process(mic, raw)

	processed := audio.NewProcessor("processed-data.bin")
	defer processed.Close()
	go compressor.Process(raw, processed)

	copyToHeadphone := effect.Copy{}
	go copyToHeadphone.Process(processed, headphone)

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <-mainThreadSignals
}
