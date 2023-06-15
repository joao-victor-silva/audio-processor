package analog

import (
	"bytes"
	"fmt"
	"os"
)

type RecorderEffect struct {
	file          *os.File
	isFileOpen    bool
	buffer        bytes.Buffer
	maxBufferSize int
}

func (e *RecorderEffect) Process(signals []Signal) []Signal {
	if e.isFileOpen {
		for _, signal := range signals {
			e.buffer.Write(signal.ToBytes())
		}

		if e.buffer.Len() > e.maxBufferSize {
			e.buffer.WriteTo(e.file)
		}
	}

	return signals
}

func NewRecorderEffect(filepath string, maxBufferSize int) *RecorderEffect {

	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		panic(fmt.Errorf("Couldn't open %s file", filepath))
	}

	return &RecorderEffect{file: file,
		isFileOpen:    true,
		maxBufferSize: maxBufferSize,
	}
}

func (e *RecorderEffect) Close() {
	if e.isFileOpen {
		if e.buffer.Len() > 0 {
			e.buffer.WriteTo(e.file)
		}

		e.file.Close()
		e.isFileOpen = false
	}
}
