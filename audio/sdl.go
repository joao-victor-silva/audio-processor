package audio

// #cgo LDFLAGS: -lSDL2
/*
#include <SDL2/SDL.h>
*/
import "C"
import (
	"fmt"
)

type sdl struct {
	initialized bool
}

type SDL interface {
	Close() error
}

func NewSDL() (SDL, error) {
	sdl := sdl{}

	ret := C.SDL_Init(C.SDL_INIT_AUDIO)
	var err error
	if ret < 0 {
		err = fmt.Errorf("Couldn't initialize SDL")
	}

	return &sdl, err
}

func (*sdl) Close() error {
	err := C.GoString(C.SDL_GetError())
	var retError error
	if err != "" {
		retError = fmt.Errorf("SDL error: %s", err)
	}
	C.SDL_Quit()

	return retError
}

