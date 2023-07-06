package analog

import (
	"math"
	"encoding/binary"

	"github.com/joao-victor-silva/audio-processor/audio"
)

func ReadFromAudioDevice(device audio.AudioDevice, jack OutputJack) {
	for {
		data := device.ReadData()
		jack.SendSignal(&data)
	}
}

func WriteInAudioDevice(device audio.AudioDevice, jack InputJack) {
	for {
		signal := jack.ReceiveSignal()
		device.WriteData(audio.Sample{Value: math.Float32frombits(
			binary.LittleEndian.Uint32(signal.ToBytes()),
		)})
	}
}
