package main

// #cgo LDFLAGS: -lSDL2
// #include <SDL2/SDL.h>
// #include <SDL2/SDL_audio.h>
import "C"
import (
	"flag"
	"fmt"
	"os"
	"unsafe"
)

func main() {
	ret := C.SDL_Init(C.SDL_INIT_AUDIO)
	if ret < 0 {
		_ = fmt.Errorf("Error")
		os.Exit(1)
	}
	defer C.SDL_Quit()
	defer func () {
		err := C.GoString(C.SDL_GetError())
		if err != "" {
			fmt.Println("SDL error: ", err)
		} else {
			fmt.Println("Exiting without errors")
		}
	}()

	var isCapture bool
	flag.BoolVar(&isCapture, "input", true, "Set device type to input")
	toCInt := map[bool]C.int {
		true: C.int(0),
		false: C.int(0),
	}

	flag.Parse()

	nInputDevices := C.SDL_GetNumAudioDevices(toCInt[isCapture])
	for i := C.int(0); i < nInputDevices; i++ {
		deviceName := C.GoString(C.SDL_GetAudioDeviceName(i, toCInt[isCapture]))
		fmt.Println("Device", i, ":", deviceName)
	}

	if nInputDevices == 0 {
		panic("No input devices found")
	}

	deviceName := C.SDL_GetAudioDeviceName(0, toCInt[isCapture])
	var desired, obtained C.SDL_AudioSpec
	var desiredPointer = unsafe.Pointer(&desired)
	// var obtainedPointer = unsafe.Pointer(&obtained)

	C.SDL_memset(desiredPointer, 0, C.sizeof_SDL_AudioSpec)

	desired.freq = 48000
	desired.format = C.AUDIO_F32
	desired.channels = 1
	desired.samples = 2048

	deviceId := C.SDL_OpenAudioDevice(deviceName, toCInt[isCapture], &desired, &obtained, C.SDL_AUDIO_ALLOW_ANY_CHANGE)

	if deviceId == 0 {
		panic("Counldn't open device")
	}

	defer C.SDL_CloseAudioDevice(deviceId)
}
