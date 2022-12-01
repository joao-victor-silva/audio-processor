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
	"math"
	"os"
	"os/signal"

	"github.com/joao-victor-silva/audio-processor/audio"
)


type ProcessData interface {
	Process(input <- chan byte, output chan <- byte, dataType C.SDL_AudioFormat)
}

type Copy struct {}
type Effect struct {
	min float64
	max float64
	threshold float64
}

func main() {
	max := flag.Float64("max", 0.03, "max")
	min := flag.Float64("min", 0.002, "min")
	threshold := flag.Float64("threshold", 0.001, "threshold")
	flag.Parse()

	sdlManager, err := audio.NewSDL()
	if err != nil {
		panic(err)
	}
	defer sdlManager.Close()

	mic, err := sdlManager.NewAudioDevice(true)
	if err != nil {
		panic("Counldn't open the mic device. " + err.Error())
	}
	defer mic.Close()

	headphone, err := sdlManager.NewAudioDevice(false)
	if err != nil {
		panic("Counldn't open the headphone device" + err.Error())
	}
	defer headphone.Close()


	if mic.AudioFormat() != headphone.AudioFormat() {
		panic("Couldn't use the same audio format for mic and headphones")
	}
	
	mic.Unpause()
	headphone.Unpause()

	copyFromRecord := Effect{ min: *min, max: *max, threshold: *threshold }
	go copyFromRecord.Process(mic, headphone, (C.SDL_AudioFormat) (mic.AudioFormat()))

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <- mainThreadSignals
}

func (*Copy) Process(inputDevice audio.AudioDevice , outputDevice audio.AudioDevice, audioFormat C.SDL_AudioFormat) {
	for outputDevice.IsChannelOpen() {
		data := inputDevice.ReadData()
		fmt.Println("value:", data)
		outputDevice.WriteData(data)
	}
}

func (effect *Effect) Process(inputDevice audio.AudioDevice , outputDevice audio.AudioDevice, audioFormat C.SDL_AudioFormat) {
	samples := make([]float64, 64)
	// vol := 0.0
	// for i, _ := range(samples)
	// 
	// rad := 0.1
	// 
	// math.Sqrt(a2 + b2 + c2 + d2 + e2 / 5.0)

	// cycle through index 
	// 0100 1000
	// 0011 0111
	// samples[binary & (2**n - 1)]
	// 8 bits =>  - (2**(n - 1)) ... (2**(n-1) - 1)
	//
	// 0.01, -0.02, 0.05 => x
	//
	// -0.03 +0.03
	//
	//
	// min 0.3
	// max 0.8
	// sample x
	// vol 0.7
	// vol 0.9
	// vol 0.2

	// x 0.2 => y = x * (max / vol)
	// y 0.3
	i := 0
	for outputDevice.IsChannelOpen() {
		data := inputDevice.ReadData()

		samples[i] = float64(data * data) / float64(len(samples))
		volume := 0.0
		for _, sample := range samples {
			volume += sample
		}
		volume = math.Sqrt(volume)
		if (volume < effect.threshold) {
			data = 0.0
		} else if (volume < effect.min) {
			//amp
			data = data * float32(effect.min / volume)
		} else if (volume > effect.max) {
			//reduce
			data = data * float32(effect.max / volume)
		}
		fmt.Println("volume:", volume, "data:", data)
		outputDevice.WriteData(data)
		i += 1
		i = i & (len(samples) - 1)
		// fmt.Println("i:", i)
	}
}
