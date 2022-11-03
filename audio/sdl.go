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
	NewAudioDevice(bool, *UserData) (AudioDevice, error)
}

type audioDevice struct {
	id C.uint
	paused bool
	audioFormat C.SDL_AudioFormat
	isCapture bool
}

type AudioDevice interface {
	Unpause()
	Pause()
	IsPaused() bool
	TogglePause()
	Close()
	AudioFormat() C.SDL_AudioFormat
}

type UserData struct {
	Record chan float32
	Process chan float32
	Playback chan float32
}

func (self *sdl) NewAudioDevice(isCapture bool, userdata *UserData) (AudioDevice, error) {
	if !self.initialized {
		err := fmt.Errorf("SDL isn't initialized")
		return nil, err
	}

	device := audioDevice{}
	device.isCapture = isCapture

	var desired, obtained C.SDL_AudioSpec
	var desiredPointer = unsafe.Pointer(&desired)
	var obtainedPointer = unsafe.Pointer(&obtained)

	C.SDL_memset(desiredPointer, 0, C.sizeof_SDL_AudioSpec)
	C.SDL_memset(obtainedPointer, 0, C.sizeof_SDL_AudioSpec)

	dataPointer := (uintptr)(unsafe.Pointer(userdata)) ^ 0xFFFFFFFF

	desired.freq = 48000
	desired.format = C.AUDIO_F32
	desired.channels = 1
	if isCapture {
		desired.samples = 512
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
		deviceName = C.SDL_GetAudioDeviceName(1, toCInt(isCapture))
	}

	device.id = C.SDL_OpenAudioDevice(deviceName, toCInt(isCapture), &desired, &obtained, C.SDL_AUDIO_ALLOW_ANY_CHANGE)

	var err error
	if (device.id == 0) {
		err = fmt.Errorf("Couldn't open the audio device")
	}

	return &device, err
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

func (device *audioDevice) Close() {
	device.Pause()
	C.SDL_CloseAudioDevice(device.id)
}

func (device *audioDevice) AudioFormat() C.SDL_AudioFormat {
	return device.audioFormat
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

