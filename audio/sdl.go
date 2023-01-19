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
	"fmt"
	"sync"
	"unsafe"
)

func toCInt(value bool) C.int {
	if value {
		return C.int(1)
	}
	return C.int(0)
}

type sdl struct {
	initialized bool
}

type SDL interface {
	Close() error
	NewAudioDevice(bool) (AudioDevice, error)
	ListAudioDevice(bool) error
}

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

// TODO: Create sink and source interface as a subset of audio device, read-only
// and write-only

func (self *sdl) NewAudioDevice(isCapture bool) (AudioDevice, error) {
	if !self.initialized {
		err := fmt.Errorf("SDL isn't initialized")
		return nil, err
	}

	device := audioDevice{}
	device.isCapture = isCapture
	device.dataChannel = make(chan float32, 1024)
	device.channelIsOpen = true

	var desired, obtained C.SDL_AudioSpec
	var desiredPointer = unsafe.Pointer(&desired)
	var obtainedPointer = unsafe.Pointer(&obtained)

	C.SDL_memset(desiredPointer, 0, C.sizeof_SDL_AudioSpec)
	C.SDL_memset(obtainedPointer, 0, C.sizeof_SDL_AudioSpec)

	var data AudioDevice
	data = &device
device.data = &data
	dataPointer := (uintptr)(unsafe.Pointer(&data)) ^ 0xFFFFFFFF

	desired.freq = 44100
	desired.format = C.AUDIO_F32
	desired.channels = 1
	if isCapture {
		desired.samples = 1024
	} else {
		desired.samples = 1024
	}

	desired.userdata = (unsafe.Pointer)(dataPointer)

	if isCapture {
		desired.callback = C.get_fn_writeptr()
	} else {
		desired.callback = C.get_fn_readptr()
	}

	var deviceName *C.char
	if isCapture {
		deviceName = C.SDL_GetAudioDeviceName(0, toCInt(isCapture))
	} else {
		deviceName = C.SDL_GetAudioDeviceName(0, toCInt(isCapture))
	}

	device.id = C.SDL_OpenAudioDevice(deviceName, toCInt(isCapture), &desired, &obtained, C.SDL_AUDIO_ALLOW_ANY_CHANGE)

	var err error
	if (device.id == 0) {
		err = fmt.Errorf("Couldn't open the audio device")
	}

	return &device, err
}

func (s *sdl) ListAudioDevice(isCapture bool) error {
	numberOfAudioDevices := C.SDL_GetNumAudioDevices(toCInt(isCapture))

	for i := C.int(0); i < numberOfAudioDevices; i++ {
		deviceName := C.GoString(C.SDL_GetAudioDeviceName(i, toCInt(isCapture)))
		fmt.Println("Device:", deviceName)
	}

	return nil
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

func NewSDL() (SDL, error) {
	sdl := sdl{}
	sdl.initialized = true

	ret := C.SDL_Init(C.SDL_INIT_AUDIO)
	var err error
	if ret < 0 {
		err = fmt.Errorf("Couldn't initialize SDL")
	}

	return &sdl, err
}

func (self *sdl) Close() error {
	err := C.GoString(C.SDL_GetError())
	var retError error
	if err != "" {
		retError = fmt.Errorf("SDL error: %s", err)
	}
	C.SDL_Quit()
	self.initialized = false

	return retError
}

