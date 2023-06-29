package analog

import (
	"encoding/binary"
	"math"
)

type Float32 float32

func (f Float32) ToBytes() []byte {
	binaryData := make([]byte, 4)
	binary.LittleEndian.PutUint32(binaryData, math.Float32bits(float32(f)))
	return binaryData
}
