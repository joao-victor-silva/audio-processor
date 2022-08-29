package main

// #cgo LDFLAGS: -lSDL2
// #include <SDL2/SDL.h>
// #include <SDL2/SDL_audio.h>
import "C"
import (
	"fmt"
	"os"
)

func main() {
	ret := C.SDL_Init(C.SDL_INIT_AUDIO)
	if ret < 0 {
		_ = fmt.Errorf("Error")
		os.Exit(1)
	}

	nInputDevices := C.SDL_GetNumAudioDevices(0)
	for i := C.int(0); i < nInputDevices; i++ {
		deviceName := C.GoString(C.SDL_GetAudioDeviceName(i, 0))
		fmt.Println("Device", i, ":", deviceName)
	}

	C.SDL_Quit()
}
