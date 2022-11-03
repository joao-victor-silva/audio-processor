package main

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
*/
import "C"
import (
	"encoding/binary"
	"math"
	"os"
	"os/signal"
	"github.com/joao-victor-silva/audio-processor/audio"
)


type ProcessData interface {
	Process(input <- chan byte, output chan <- byte, dataType C.SDL_AudioFormat)
}

type Copy struct {}
type Effect struct {}

func main() {
	sdlManager, err := audio.NewSDL()
	if err != nil {
		panic(err)
	}
	defer sdlManager.Close()

	var userdata audio.UserData
	userdata.Record = make(chan byte, 1024 * 4)
	userdata.Process = make(chan byte, 1024 * 4)
	userdata.Playback = make(chan byte, 1024 * 4)

	defer close (userdata.Record)

	mic, err := audio.NewAudioDevice(true, &userdata)
	if err != nil {
		panic("Counldn't open the mic device")
	}
	defer mic.Close()

	headphone, err := audio.NewAudioDevice(false, &userdata)
	if err != nil {
		panic("Counldn't open the headphone device")
	}
	defer headphone.Close()


	if mic.AudioFormat() != headphone.AudioFormat() {
		panic("Couldn't use the same audio format for mic and headphones")
	}
	
	mic.Unpause()
	headphone.Unpause()

	copyFromRecord := Copy{}
	copyToPlayback := Copy{}
	go copyFromRecord.Process(userdata.Record, userdata.Process, (C.SDL_AudioFormat) (mic.AudioFormat()))
	go copyToPlayback.Process(userdata.Process, userdata.Playback, (C.SDL_AudioFormat) (headphone.AudioFormat()))

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <- mainThreadSignals
}

func (*Copy) Process(input <- chan byte, output chan <- byte, audioFormat C.SDL_AudioFormat) {
	for data := range input {
		output <- data
	}
}

func (*Effect) Process(input <- chan byte, output chan <- byte, audioFormat C.SDL_AudioFormat) {
	for true {
		if audioFormat == C.AUDIO_F32 {
			binaryData := make([]byte, 4)
			binaryData[0] = <- input
			binaryData[1] = <- input
			binaryData[2] = <- input
			binaryData[3] = <- input

			buffer := math.Float32frombits(binary.LittleEndian.Uint32(binaryData))
			buffer = buffer / 100
			binary.LittleEndian.PutUint32(binaryData, math.Float32bits(buffer))

			output <- binaryData[0]
			output <- binaryData[1]
			output <- binaryData[2]
			output <- binaryData[3]
		}
	}

	// for data := range input {
	// }
}
