package main

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
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
	"os"
	"unsafe"
)

type UserData struct {
	playback chan byte
}

func main() {
	toCInt := map[bool]C.int{
		true:  C.int(1),
		false: C.int(0),
	}

	ret := C.SDL_Init(C.SDL_INIT_AUDIO)
	if ret < 0 {
		_ = fmt.Errorf("Error")
		os.Exit(1)
	}
	defer C.SDL_Quit()
	defer func() {
		err := C.GoString(C.SDL_GetError())
		if err != "" {
			fmt.Println("SDL error: ", err)
		} else {
			fmt.Println("Exiting without errors")
		}
	}()

	var userdata UserData
	userdata.playback = make(chan byte, 1024 * 4)

	micId := openDevice(true, &userdata)
	headphoneId := openDevice(false, &userdata)

	if micId == 0 {
		panic("Counldn't open the mic device")
	}
	defer C.SDL_CloseAudioDevice(micId)

	if headphoneId == 0 {
		panic("Counldn't open the headphone device")
	}
	defer C.SDL_CloseAudioDevice(headphoneId)

	C.SDL_PauseAudioDevice(micId, toCInt[false])
	C.SDL_PauseAudioDevice(headphoneId, toCInt[false])
	C.SDL_Delay(4000);
}

func openDevice(isCapture bool, userdata *UserData) C.uint {
	toCInt := map[bool]C.int{
		true:  C.int(1),
		false: C.int(0),
	}

	var desired, obtained C.SDL_AudioSpec
	var desiredPointer = unsafe.Pointer(&desired)

	C.SDL_memset(desiredPointer, 0, C.sizeof_SDL_AudioSpec)

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
		deviceName = C.SDL_GetAudioDeviceName(0, toCInt[isCapture])
	} else {
		deviceName = C.SDL_GetAudioDeviceName(1, toCInt[isCapture])
	}

	return C.SDL_OpenAudioDevice(deviceName, toCInt[isCapture], &desired, &obtained, C.SDL_AUDIO_ALLOW_ANY_CHANGE)
}
