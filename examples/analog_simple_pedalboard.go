package main

import (
	"encoding/binary"
	"math"
	"os"
	_ "os"
	"os/signal"

	"github.com/joao-victor-silva/audio-processor/analog"
	"github.com/joao-victor-silva/audio-processor/audio"
)

func main() {
	sdlManager, err := audio.NewSDL()
	if err != nil {
		panic(err)
	}
	defer sdlManager.Close()

	mic, err := sdlManager.NewAudioDevice(0, true)
	if err != nil {
		panic("Counldn't open the mic device. " + err.Error())
	}
	defer mic.Close()
	mic.Unpause()

	headphone, err := sdlManager.NewAudioDevice(0, false)
	if err != nil {
		panic("Counldn't open the headphone device" + err.Error())
	}
	defer headphone.Close()
	headphone.Unpause()

	if mic.AudioFormat() != headphone.AudioFormat() {
		panic("Couldn't use the same audio format for mic and headphones")
	}

	// Add pedal in position to rewire other pedals
	pedalBoard := analog.NewPedalBoard()

	recordEffect := analog.NewRecorderEffect("raw-input.bin", 1024)
	defer recordEffect.Close()

	recordPedal := analog.NewPedal(recordEffect)
	pedalBoard.AddPedal(recordPedal, 0)

	micOutput := analog.NewOutputJack()
	go func(device audio.AudioDevice, jack analog.OutputJack) {
		for {
			data := device.ReadData()
			jack.SendSignal(&data)
		}
	}(mic, micOutput)

	headphoneInput := analog.NewInputJack()
	go func(device audio.AudioDevice, jack analog.InputJack) {
		for {
			signal := jack.ReceiveSignal()
			device.WriteData(audio.Sample{Value: math.Float32frombits(
				binary.LittleEndian.Uint32(signal.ToBytes()),
			)})
		}
	}(headphone, headphoneInput)

	pedalBoard.InputConnect(micOutput.GetWire())
	pedalBoard.OutputConnect(headphoneInput.GetWire())

	shouldRun := true

	recordPedal.Toggle()
	go recordPedal.Run(&shouldRun)

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <-mainThreadSignals

	// effectOnePedal := analog.NewDummyPedal()
	//
	// effectTwoPedal := analog.NewDummyPedal()
	//
	// // Add pedal in container
	// // Connect raw input to pedal input jack (mic -> pedal)
	// // Connect pedal output jack to raw output (pedal -> headphone)
	// pedalBoard.AddPedal(effectOnePedal, 0)
	//
	// // Add pedal in container
	// // Disconnect effectOnePedal output from raw output
	// // Connect effectOnePedal output to pedal input jack (pedalOne -> pedalTwo)
	// // Connect pedal output jack to raw output (pedalTwo -> headphone)
	// pedalBoard.AddPedal(effectTwoPedal, 1)
}
