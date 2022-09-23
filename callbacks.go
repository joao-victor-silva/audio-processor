package main

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
#include <string.h>
*/
import "C"
import "unsafe"

//export fillBuffer
func fillBuffer(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	userdata = unsafe.Pointer(uintptr(userdata) ^ 0xFFFFFFFF)
	userdataPointer := (*UserData)(userdata)
	data := C.GoBytes(unsafe.Pointer(stream), length)
	for _, b := range data {
		select {
			case userdataPointer.playback <- b:
			default:
		}
	}
}

//export readBuffer
func readBuffer(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	userdata = unsafe.Pointer(uintptr(userdata) ^ 0xFFFFFFFF)
	userdataPointer := (*UserData)(userdata)
	
	streamSlice := CPoiterToSlice(stream, length)
	for i := 0; i < int(length); i++ {
		select {
		case data := <- userdataPointer.playback:
			streamSlice[i] = (C.Uint8) (data)
		default:
			streamSlice[i] = 0
		}
	}
}

func CPoiterToSlice(cArray *C.Uint8, cSize C.int) []C.Uint8 {
    gSlice := (*[1 << 30]C.Uint8)(unsafe.Pointer(cArray))[:cSize:cSize]
    return gSlice
}
