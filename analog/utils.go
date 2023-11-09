package analog

import (
	"math"
	"encoding/binary"

	"github.com/joao-victor-silva/audio-processor/audio"
)

func ReadFromAudioDevice(device audio.AudioDevice, jack OutputJack) {
	i := 1.0
	// c.sampleCount = c.sampleCount + 1
	// if (c.sampleCount >= 44100) {
	// 	c.Close()
	// }
	// rad := (float64(c.sampleCount) * math.Pi) / 100.0
	// // return (float32(math.Sin(rad)) * 0.5) + 0.5
	// return audio.Sample{Value: float32(math.Sin(rad)), Volume: float64(1.0)}
	for {
		wave0 := math.Sin(math.Pi * i) / 100.0
		jack.SendSignal(&audio.Sample{Value: float32(wave0)})
		i += 1.0
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

