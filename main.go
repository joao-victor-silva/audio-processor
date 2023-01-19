package main

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
*/
import "C"
import (
	"fmt"
	"flag"
	"os"
	"os/signal"

	"github.com/joao-victor-silva/audio-processor/audio"
	"github.com/joao-victor-silva/audio-processor/effect"
)

func main() {
	threshold := flag.Float64("threshold", 0.001, "threshold")
	min := flag.Float64("min", 0.002, "min")
	max := flag.Float64("max", 0.03, "max")
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

	file, _ := os.OpenFile("data.bin", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	input_output := customAudioDevice{channelIsOpen: true, sampleCount: 0, file: file};
	defer input_output.Close()


	headphone, err := sdlManager.NewAudioDevice(false)
	if err != nil {
		panic("Counldn't open the headphone device" + err.Error())
	}
	defer headphone.Close()
	headphone.Unpause()


	if mic.AudioFormat() != headphone.AudioFormat() {
		panic("Couldn't use the same audio format for mic and headphones")
	}

	// copyEffect := effect.Copy{}
	// go copyEffect.Process(&input_output, headphone)
	copyFromRecord := effect.Effect{ Min: *min, Max: *max, Threshold: *threshold }
	// go copyFromRecord.Process(&input_output, headphone)
	// go copyFromRecord.Process(mic, headphone)
	go copyFromRecord.Process(mic, &input_output)

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <- mainThreadSignals
}
