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
	"encoding/binary"
	"fmt"
	"math"
	"os"
	"os/signal"
	"unsafe"
)

type UserData struct {
	record chan byte
	process chan byte
	playback chan byte
}

type ProcessData interface {
	Process(input <- chan byte, output chan <- byte, dataType C.SDL_AudioFormat)

}

type Copy struct {}
type Effect struct {}

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
	userdata.record = make(chan byte, 1024 * 4)
	userdata.process = make(chan byte, 1024 * 4)
	userdata.playback = make(chan byte, 1024 * 4)

	defer close (userdata.record)

	micId, micAudioFormat := openDevice(true, &userdata)
	headphoneId, headphoneAudioFormat := openDevice(false, &userdata)

	if micAudioFormat != headphoneAudioFormat {
		panic("Couldn't use the same audio format for mic and headphones")
	}

	if micId == 0 {
		panic("Counldn't open the mic device")
	}
	defer C.SDL_CloseAudioDevice(micId)

	if headphoneId == 0 {
		panic("Counldn't open the headphone device")
	}
	defer C.SDL_CloseAudioDevice(headphoneId)

	C.SDL_PauseAudioDevice(micId, toCInt[false])
	defer C.SDL_PauseAudioDevice(micId, toCInt[true])
	C.SDL_PauseAudioDevice(headphoneId, toCInt[false])
	defer C.SDL_PauseAudioDevice(headphoneId, toCInt[true])

	copyFromRecord := Copy{}
	copyToPlayback := Copy{}
	go copyFromRecord.Process(userdata.record, userdata.process, micAudioFormat)
	go copyToPlayback.Process(userdata.process, userdata.playback, headphoneAudioFormat)

	mainThreadSignals := make(chan os.Signal, 1)
	signal.Notify(mainThreadSignals, os.Interrupt)
	_ = <- mainThreadSignals
}

func (*Copy) Process(input <- chan byte, output chan <- byte, audioFormat C.SDL_AudioFormat) {
	for data := range input {
		output <- data
	}
}

func (*Effect) Process(input <- chan byte, output chan <- byte, audioFormat C.SDL_AudioFormat) {
	for true {
		if audioFormat == C.AUDIO_F32 {
			binaryData := make([]byte, 4)
			binaryData[0] = <- input
			binaryData[1] = <- input
			binaryData[2] = <- input
			binaryData[3] = <- input

			buffer := math.Float32frombits(binary.LittleEndian.Uint32(binaryData))
			buffer = buffer / 100
			binary.LittleEndian.PutUint32(binaryData, math.Float32bits(buffer))

			output <- binaryData[0]
			output <- binaryData[1]
			output <- binaryData[2]
			output <- binaryData[3]
		}
	}

	// for data := range input {
	// }
}

func openDevice(isCapture bool, userdata *UserData) (C.uint, C.SDL_AudioFormat) {
	toCInt := map[bool]C.int{
		true:  C.int(1),
		false: C.int(0),
	}

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
		deviceName = C.SDL_GetAudioDeviceName(0, toCInt[isCapture])
	} else {
		deviceName = C.SDL_GetAudioDeviceName(1, toCInt[isCapture])
	}

	deviceId := C.SDL_OpenAudioDevice(deviceName, toCInt[isCapture], &desired, &obtained, C.SDL_AUDIO_ALLOW_ANY_CHANGE)
	return deviceId, obtained.format
}
