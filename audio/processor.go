package audio

import (
	"bufio"
	"encoding/binary"
	"math"
	"os"
	"sync"
)

type AudioProcessor interface {
	IsChannelOpen() bool
	Close()
	ReadData() Sample
	WriteData(Sample)
}

type processor struct {
	File       *os.File
	isFileOpen bool
	writer *bufio.Writer
	// Input      chan float32
	// Output     chan float32
	bypass chan Sample
	mu sync.Mutex
}

func (p *processor) IsChannelOpen() bool {
	return p.isFileOpen
}

func (p *processor) Close() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.writer.Flush()
	p.File.Close()
	p.File = nil
	p.isFileOpen = false
	close(p.bypass)
}

func (p *processor) ReadData() Sample {
	// data := make([]byte, 4)
	// p.File.Read(data)
	//
	// return math.Float32frombits(binary.LittleEndian.Uint32(data))
	return <- p.bypass
}

func (p *processor) WriteData(data Sample) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if !p.isFileOpen {
		return
	}

	binaryData := make([]byte, 4)
	binary.LittleEndian.PutUint32(binaryData, math.Float32bits(data.Value))
	p.writer.Write(binaryData)
	p.bypass <- data
}

func NewProcessor(filePath string) *processor {
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		panic(err)
	}
	return &processor{File: file, isFileOpen: true, bypass: make(chan Sample, 2048), writer: bufio.NewWriter(file)}
}
