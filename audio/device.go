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
import "sync"

type audioDevice struct {
	id C.uint
	paused bool
	audioFormat C.SDL_AudioFormat
	isCapture bool
	dataChannel chan float32
	channelMutex sync.Mutex
	channelIsOpen bool
	data interface{}
}

type AudioDevice interface {
	Unpause()
	Pause()
	IsPaused() bool
	IsChannelOpen() bool
	TogglePause()
	Close()
	AudioFormat() uint
	ReadData() float32
	ReadDataUnsafe() float32
	WriteData(float32)
	WriteSlice([]float32)
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
	if (device.IsPaused()) {
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

func (device *audioDevice) WriteData(data float32) {
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

	for _, data := range dataArray {
		select {
			case device.dataChannel <- data:
			default:
		}
	}
}

func (device *audioDevice) ReadData() float32 {
	return <- device.dataChannel
}

func (device *audioDevice) ReadDataUnsafe() float32 {
	select {
		case data := <- device.dataChannel:
			return data
		default:
			return 0
	}
}
