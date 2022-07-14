package messages

import (
	"encoding/binary"
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/libs/vector"
)

func NewVZoneMessage(OriginatorIP net.IP, position Position, speed uint32) *VZoneMessage {
	return &VZoneMessage{
		Type:         VZoneType,
		OriginatorIP: OriginatorIP,
		Position:     position,
		Speed:        speed,
	}
}

func (VZone *VZoneMessage) Marshal() []byte {
	bytes := make([]byte, VZoneMessageLen)

	bytes[0] = byte(VZone.Type)
	copy(bytes[1:5], VZone.OriginatorIP.To4())
	// add marashalled position object to bytes
	copy(bytes[5:21], VZone.Position.Marshal())
	binary.LittleEndian.PutUint32(bytes[21:25], VZone.Speed)

	return bytes
}

func UnmarshalVZone(data []byte) (*VZoneMessage, error) {
	if len(data) < VZoneMessageLen {
		return nil, fmt.Errorf("VZone message length is %d, expected %d", len(data), VZoneMessageLen)
	}

	VZone := &VZoneMessage{}
	VZone.Type = uint8(data[0])
	VZone.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VZone.Position = UnmarshalPosition(data[5:21])
	binary.LittleEndian.PutUint32(data[21:25], VZone.Speed)

	return VZone, nil
}

func (VZone *VZoneMessage) String() string {
	return fmt.Sprintf("VZone: Type: %d, OriginatorIP: %s,  Lat: %f, Lng: %f, Speed: %d", VZone.Type, VZone.OriginatorIP.String(), VZone.Position.Lat, VZone.Position.Lng, VZone.Speed)
}
