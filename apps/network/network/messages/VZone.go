package messages

import (
	"fmt"
	"net"
)


func NewVZoneMessage(OriginatorIP net.IP, position Position, maxDistance float64) *VZoneMessage {
	return &VZoneMessage{
		Type: VZoneType,
		OriginatorIP: OriginatorIP,
		Position : position,
		MaxDistance : maxDistance,
	}
}

func (VZone *VZoneMessage) Marshal() []byte {
	bytes := make([]byte, VZoneMessageLen)

	bytes[0] = byte(VZone.Type)
	copy(bytes[1:5], VZone.OriginatorIP.To4())
	// add marashalled position object to bytes
	copy(bytes[5:21], VZone.Position.Marshal())
	copy(bytes[21:29],  Float64bytes(VZone.MaxDistance))
	
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
	VZone.MaxDistance = Float64FromBytes(data[21:29])
	
	return VZone, nil
}

func (VZone *VZoneMessage) String() string {
	return fmt.Sprintf("VZone: Type: %d, OriginatorIP: %s,  Lat: %f, Lng: %f, MaxDistance: %f", VZone.Type, VZone.OriginatorIP.String(), VZone.Position.Lat, VZone.Position.Lng, VZone.MaxDistance)
}