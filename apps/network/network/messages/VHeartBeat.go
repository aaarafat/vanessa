package messages

import (
	"encoding/binary"
	"fmt"
	"net"
)




func (VHBeat *VHBeatMessage) Marshal() []byte {
	bytes := make([]byte, VHBeatMessageLen)

	bytes[0] = byte(VHBeat.Type)
	copy(bytes[1:5], VHBeat.OriginatorIP.To4())
	binary.LittleEndian.PutUint32(bytes[5:], VHBeat.PositionX)
	binary.LittleEndian.PutUint32(bytes[9:], VHBeat.PositionY)


	return bytes
}

//Create a new VHBeat message
func NewVHBeatMessage(OriginatorIP net.IP, PositionX uint32, PositionY uint32) *VHBeatMessage {
	return &VHBeatMessage{
		Type: VHBeatType,
		OriginatorIP: OriginatorIP,
		PositionX: PositionX,
		PositionY: PositionY,
	}
}



func UnmarshalVHBeat(data []byte) (*VHBeatMessage, error) {
	if len(data) < VHBeatMessageLen {
		return nil, fmt.Errorf("VHBeat message length is %d, expected %d", len(data), VHBeatMessageLen)
	}

	VHBeat := &VHBeatMessage{}
	VHBeat.Type = uint8(data[0])
	VHBeat.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VHBeat.PositionX = binary.LittleEndian.Uint32(data[5:9])
	VHBeat.PositionY = binary.LittleEndian.Uint32(data[9:13])

	return VHBeat, nil
}


// print the VHBeat message
func (VHBeat *VHBeatMessage) String() string {
	return fmt.Sprintf("VHBeat: Type: %d, OriginatorIP: %s, PositionX: %d, PositionY: %d", VHBeat.Type, VHBeat.OriginatorIP.String(), VHBeat.PositionX, VHBeat.PositionY)
}

