package main

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
extern void fillBuffer(void *userdata, Uint8 * stream, int len);
static SDL_AudioCallback get_fn_ptr() {
    return fillBuffer;
}
*/
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
	flag.BoolVar(&isCapture, "input", false, "Set device type to input")
	toCInt := map[bool]C.int {
		true: C.int(1),
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

	deviceName := C.SDL_GetAudioDeviceName(1, toCInt[isCapture])
	var desired, obtained C.SDL_AudioSpec
	var desiredPointer = unsafe.Pointer(&desired)

	C.SDL_memset(desiredPointer, 0, C.sizeof_SDL_AudioSpec)

	desired.freq = 48000
	desired.format = C.AUDIO_F32
	desired.channels = 1
	desired.samples = 2048
	// desired.callback = C.get_fn_ptr()

	deviceId := C.SDL_OpenAudioDevice(deviceName, toCInt[isCapture], &desired, &obtained, C.SDL_AUDIO_ALLOW_ANY_CHANGE)

	if deviceId == 0 {
		panic("Counldn't open device")
	}
	defer C.SDL_CloseAudioDevice(deviceId)

	C.SDL_PauseAudioDevice(deviceId, toCInt[false])

	if isCapture {
		dataWanted := 96000 * 4
		data := make([]byte, dataWanted)
		dataPointer := unsafe.Pointer(&data[0])

		C.SDL_Delay(2020)
		dataSize := C.SDL_DequeueAudio(deviceId, dataPointer, C.Uint32(dataWanted))

		if dataSize == 0 {
			panic("Counldn't retrieve data from device")
		}
		fmt.Println("Got", dataSize, "bytes.")

		err := os.WriteFile("data.bin", data, 0644)
		if err != nil {
			panic("Couldn't write file")
		}

	} else {
		data, err := os.ReadFile("data.bin")
		if err != nil {
			panic("Couldn't read file")
		}

		C.SDL_QueueAudio(deviceId, C.CBytes(data), C.Uint32(len(data)))
		C.SDL_Delay(2020)
	}
}

//export fillBuffer
func fillBuffer(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {

}
