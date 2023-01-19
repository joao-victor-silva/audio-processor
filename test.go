package main

import (
	"os"
	"encoding/binary"
	"fmt"
	"math"
	_ "github.com/joao-victor-silva/audio-processor/audio"
	_ "github.com/joao-victor-silva/audio-processor/effect"
)

type customAudioDevice struct{
	channelIsOpen bool
	sampleCount int
	file *os.File
}

func (c *customAudioDevice) Unpause() {}
func (c *customAudioDevice) Pause() {}
func (c *customAudioDevice) IsPaused() bool {
	return false;
}

func (c *customAudioDevice) IsChannelOpen() bool {
	return c.channelIsOpen;
}

func (c *customAudioDevice) Close() {
	if (c.file == nil) {
		return
	}
	c.channelIsOpen = false
	c.file.Close()
	c.file = nil
}
func (c *customAudioDevice) AudioFormat() uint {
	return 0
}
func (c *customAudioDevice) TogglePause() {}

func (c *customAudioDevice) ReadData() float32 {
	c.sampleCount = c.sampleCount + 1
	if (c.sampleCount >= 44100) {
		c.Close()
	}
	rad := (float64(c.sampleCount) * math.Pi) / 100.0
	// return (float32(math.Sin(rad)) * 0.5) + 0.5
	return float32(math.Sin(rad))
}

func (c *customAudioDevice) ReadDataUnsafe() float32 {
	return 0.0
}

func (c *customAudioDevice) WriteData(data float32) {
	if (!c.channelIsOpen) {
		return
	}

	binaryData := make([]byte, 4)
	binary.LittleEndian.PutUint32(binaryData, math.Float32bits(data))
	c.file.Write(binaryData)
}

func (c *customAudioDevice) WriteSlice(data []float32) {
	for _, d := range data {
		fmt.Println(d)
	}
}

// func main() {
// 	file, _ := os.OpenFile("data.bin", os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
// 	input_output := customAudioDevice{channelIsOpen: true, sampleCount: 0, file: file};
//
// 	effect := effect.Copy{}
// 	effect.Process(&input_output, &input_output);
// }




