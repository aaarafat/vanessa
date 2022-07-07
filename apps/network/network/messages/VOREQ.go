package messages

import (
	"fmt"
	"net"
)




func (VOREQ *VOREQMessage) Marshal() []byte {
	bytes := make([]byte, VOREQMessageLen)

	bytes[0] = byte(VOREQ.Type)
	copy(bytes[1:5], VOREQ.OriginatorIP.To4())



	return bytes
}

//Create a new VOREQ message
func NewVOREQMessage(OriginatorIP net.IP) *VOREQMessage {
	return &VOREQMessage{
		Type: VOREQType,
		OriginatorIP: OriginatorIP,
	}
}



func UnmarshalVOREQ(data []byte) (*VOREQMessage, error) {
	if len(data) < VOREQMessageLen {
		return nil, fmt.Errorf("VOREQ message length is %d, expected %d", len(data), VOREQMessageLen)
	}

	VOREQ := &VOREQMessage{}
	VOREQ.Type = uint8(data[0])
	VOREQ.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	return VOREQ, nil
}


// print the VOREQ message
func (VOREQ *VOREQMessage) String() string {
	return fmt.Sprintf("VOREQ: Type: %d, OriginatorIP: %s", VOREQ.Type, VOREQ.OriginatorIP.String())
}

