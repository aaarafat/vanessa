package messages

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/libs/vector"
)

func (VOREQ *VOREQMessage) Marshal() []byte {
	bytes := make([]byte, VOREQMessageLen+16*VOREQ.Length)

	bytes[0] = byte(VOREQ.Type)
	copy(bytes[1:5], VOREQ.OriginatorIP.To4())
	bytes[5] = byte(VOREQ.Length)
	copy(bytes[VOREQMessageLen:VOREQMessageLen+16*VOREQ.Length], VOREQ.Obstacles)

	return bytes
}

//Create a new VOREQ message
func NewVOREQMessage(OriginatorIP net.IP, obstacles []Position) *VOREQMessage {

	bytes := make([]byte, 16*len(obstacles))
	for i, v := range obstacles {
		copy(bytes[i*16:i*16+8], Float64bytes(v.Lat))
		copy(bytes[i*16+8:i*16+16], Float64bytes(v.Lng))
	}
	return &VOREQMessage{
		Type:         VOREQType,
		OriginatorIP: OriginatorIP,
		Length:       uint8(len(obstacles)),
		Obstacles:    bytes,
	}
}

func UnmarshalVOREQ(data []byte) *VOREQMessage {
	VOREQ := &VOREQMessage{}
	VOREQ.Type = uint8(data[0])
	VOREQ.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VOREQ.Length = uint8(data[5])
	VOREQ.Obstacles = data[VOREQMessageLen:]
	return VOREQ
}

// print the VOREQ message
func (VOREQ *VOREQMessage) String() string {
	return fmt.Sprintf("VOREQ: Type: %d, OriginatorIP: %s, Length:%d, Obstacles: %s", VOREQ.Type, VOREQ.OriginatorIP.String(), VOREQ.Length, VOREQ.Obstacles)
}
