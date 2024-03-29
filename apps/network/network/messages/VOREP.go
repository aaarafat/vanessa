package messages

import (
	"fmt"

	. "github.com/aaarafat/vanessa/libs/vector"
)

func (VOREP *VOREPMessage) Marshal() []byte {
	bytes := make([]byte, VOREPMessageLen+16*VOREP.Length)

	bytes[0] = byte(VOREP.Type)
	bytes[1] = byte(VOREP.Length)

	copy(bytes[VOREPMessageLen:VOREPMessageLen+16*VOREP.Length], VOREP.Obstacles)

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
		Type:      VOREPType,
		Length:    uint8(len(obstacles)),
		Obstacles: bytes,
	}

}

func UnmarshalVOREP(data []byte) *VOREPMessage {
	VOREP := &VOREPMessage{}
	VOREP.Type = uint8(data[0])
	VOREP.Length = uint8(data[1])
	VOREP.Obstacles = data[VOREPMessageLen:]
	return VOREP
}

// print the VOREP message
func (VOREP *VOREPMessage) String() string {
	return fmt.Sprintf("VOREP: Type: %d, Length: %d Obstacles: %s", VOREP.Type, VOREP.Length, VOREP.Obstacles)
}
