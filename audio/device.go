package audio

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
extern void fillBuffer(void *userdata, Uint8 * stream, int len);
extern void readBuffer(void *userdata, Uint8 * stream, int len);
static SDL_AudioCallback get_fn_writeptr() {
    return fillBuffer;
}
static SDL_AudioCallback get_fn_readptr() {
    return readBuffer;
}
*/
import "C"
import (
	"encoding/binary"
	"math"
	"sync"
)

type audioDevice struct {
	id                 C.uint
	paused             bool
	audioFormat        C.SDL_AudioFormat
	isCapture          bool
	dataChannel        chan Sample
	channelMutex       sync.Mutex
	channelIsOpen      bool
	volumeSamples      []float64
	volumeSamplesIndex int
	data               interface{}
}

type AudioDevice interface {
	Unpause()
	Pause()
	IsPaused() bool
	IsChannelOpen() bool
	TogglePause()
	Close()
	AudioFormat() uint
	ReadData() Sample
	ReadDataUnsafe() float32
	WriteData(Sample)
	WriteSlice([]float32)
}

type Sample struct {
	Value  float32
	Volume float64
}

func (s *Sample) ToBytes() []byte {
	binaryData := make([]byte, 4)
	binary.LittleEndian.PutUint32(binaryData, math.Float32bits(s.Value))
	return binaryData
}

func (device *audioDevice) Pause() {
	C.SDL_PauseAudioDevice(device.id, toCInt(true))
	device.paused = true
}

func (device *audioDevice) Unpause() {
	C.SDL_PauseAudioDevice(device.id, toCInt(false))
	device.paused = false
}

func (device *audioDevice) IsPaused() bool {
	return device.paused
}

func (device *audioDevice) TogglePause() {
	if device.IsPaused() {
		device.Unpause()
	} else {
		device.Pause()
	}
}

func (device *audioDevice) IsChannelOpen() bool {
	return device.channelIsOpen
}

func (device *audioDevice) Close() {
	device.channelMutex.Lock()
	defer device.channelMutex.Unlock()
	device.Pause()
	device.data = nil
	device.channelIsOpen = false
	close(device.dataChannel)
	C.SDL_CloseAudioDevice(device.id)
}

func (device *audioDevice) AudioFormat() uint {
	return uint(device.audioFormat)
}

func (device *audioDevice) WriteData(data Sample) {
	device.channelMutex.Lock()
	defer device.channelMutex.Unlock()

	if !device.channelIsOpen {
		return
	}

	select {
	case device.dataChannel <- data:
	default:
	}
}

func (device *audioDevice) WriteSlice(dataArray []float32) {
	device.channelMutex.Lock()
	defer device.channelMutex.Unlock()

	if !device.channelIsOpen {
		return
	}

	volumeSamplesAmount := float64(len(dataArray))
	average := float32(0.0)
	for _, data := range dataArray {
		average += data
	}
	average64 := float64(average) / volumeSamplesAmount

	// Check if volume is on dB -> 4x can be 16x due to non linear progression
	volume := 0.0
	for _, s := range dataArray {
		volume += math.Pow(float64(s)-average64, 2)
	}
	volume = math.Sqrt(volume) / volumeSamplesAmount

	for _, data := range dataArray {
		select {
		case device.dataChannel <- Sample{Value: data, Volume: volume}:
		default:
		}
	}
}

func (device *audioDevice) ReadData() Sample {
	return <-device.dataChannel
}

func (device *audioDevice) ReadDataUnsafe() float32 {
	select {
	case data := <-device.dataChannel:
		return data.Value
	default:
		return 0
	}
}
