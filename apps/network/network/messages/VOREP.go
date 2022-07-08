package messages

import (
	"encoding/binary"
	"fmt"
	"math"
)

func Float64FromBytes(bytes []byte) float64 {
    bits := binary.LittleEndian.Uint64(bytes)
    float := math.Float64frombits(bits)
    return float
}

func Float64bytes(float float64) []byte {
    bits := math.Float64bits(float)
    bytes := make([]byte, 8)
    binary.LittleEndian.PutUint64(bytes, bits)
    return bytes
}


func (VOREP *VOREPMessage) Marshal() []byte {
	bytes := make([]byte,2+16*VOREP.Length)

	bytes[0] = byte(VOREP.Type)
	bytes[1] = byte(VOREP.Length)
	copy(bytes[2:2+16*VOREP.Length], VOREP.Obstacles)



	return bytes
}


//Create a new VOREP message by convining the obstacles list of x , y to byte array
func NewVOREPMessage(obstacles []Position) *VOREPMessage {
	bytes := make([]byte, 16*len(obstacles)) 
	for i, v := range obstacles {
		copy(bytes[i*16:i*16+8], Float64bytes(v.Lat))
		copy(bytes[i*16+8:i*16+16], Float64bytes(v.Lng))
	}
	return &VOREPMessage{
		Type: VOREPType,
		Length: uint8(len(obstacles)),
		Obstacles: bytes,
	}

}



func UnmarshalVOREP(data []byte) *VOREPMessage {
	VOREP := &VOREPMessage{}
	VOREP.Type = uint8(data[0])
	VOREP.Length = uint8(data[1])
	VOREP.Obstacles = data[2:]
	return VOREP
}


// print the VOREP message
func (VOREP *VOREPMessage) String() string {
	return fmt.Sprintf("VOREP: Type: %d, Length: %d Obstcales: %s", VOREP.Type,VOREP.Length ,VOREP.Obstacles)
}

