package messages

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/libs/vector"
)

func (VHBeat *VHBeatMessage) Marshal() []byte {
	bytes := make([]byte, VHBeatMessageLen)

	bytes[0] = byte(VHBeat.Type)
	copy(bytes[1:5], VHBeat.OriginatorIP.To4())
	// add marashalled position object to bytes
	copy(bytes[5:], VHBeat.Position.Marshal())

	return bytes
}

//Create a new VHBeat message
func NewVHBeatMessage(OriginatorIP net.IP, position Position) *VHBeatMessage {
	return &VHBeatMessage{
		Type:         VHBeatType,
		OriginatorIP: OriginatorIP,
		Position:     position,
	}
}

func UnmarshalVHBeat(data []byte) (*VHBeatMessage, error) {
	if len(data) < VHBeatMessageLen {
		return nil, fmt.Errorf("VHBeat message length is %d, expected %d", len(data), VHBeatMessageLen)
	}

	VHBeat := &VHBeatMessage{}
	VHBeat.Type = uint8(data[0])
	VHBeat.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VHBeat.Position = UnmarshalPosition(data[5:])

	return VHBeat, nil
}

// print the VHBeat message
func (VHBeat *VHBeatMessage) String() string {
	return fmt.Sprintf("VHBeat: Type: %d, OriginatorIP: %s,  Lat: %f, Lng: %f", VHBeat.Type, VHBeat.OriginatorIP.String(), VHBeat.Position.Lat, VHBeat.Position.Lng)
}
