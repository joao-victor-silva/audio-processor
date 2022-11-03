package audio

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
#include <SDL2/SDL_audio.h>
*/
import "C"
import (
	"encoding/binary"
	"math"
	"unsafe"
)

//export fillBuffer
func fillBuffer(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	userdata = unsafe.Pointer(uintptr(userdata) ^ 0xFFFFFFFF)
	userdataPointer := (*UserData)(userdata)
	data := C.GoBytes(unsafe.Pointer(stream), length)

	for i := 0; i + 3 < len(data); i = i + 4 {
		buffer := math.Float32frombits(binary.LittleEndian.Uint32(data[i:i+4]))
		select {
			case userdataPointer.Record <- buffer:
			default:
		}
	}

}

//export readBuffer
func readBuffer(userdata unsafe.Pointer, stream *C.Uint8, length C.int) {
	userdata = unsafe.Pointer(uintptr(userdata) ^ 0xFFFFFFFF)
	userdataPointer := (*UserData)(userdata)
	
	streamSlice := CPoiterToSlice(stream, length)
	for i := 0; i + 3 < int(length); i = i + 4 {
		select {
		case data := <- userdataPointer.Playback:
			binaryData := make([]byte, 4)
			binary.LittleEndian.PutUint32(binaryData, math.Float32bits(data))
			streamSlice[i] = (C.Uint8) (binaryData[0])
			streamSlice[i+1] = (C.Uint8) (binaryData[1])
			streamSlice[i+2] = (C.Uint8) (binaryData[2])
			streamSlice[i+3] = (C.Uint8) (binaryData[3])
		default:
			streamSlice[i] = 0
			streamSlice[i+1] = 0
			streamSlice[i+2] = 0
			streamSlice[i+3] = 0
		}
	}
}

func CPoiterToSlice(cArray *C.Uint8, cSize C.int) []C.Uint8 {
    gSlice := (*[1 << 30]C.Uint8)(unsafe.Pointer(cArray))[:cSize:cSize]
    return gSlice
}
