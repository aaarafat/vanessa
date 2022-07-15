package messages

import (
	"encoding/binary"
	"fmt"
	"net"
)

func NewVSpeedMessage(OriginatorIP net.IP, speed uint32) *VSpeedMessage {
	return &VSpeedMessage{
		Type:         VSpeedType,
		OriginatorIP: OriginatorIP,
		Speed:        speed,
	}
}

func (VSpeed *VSpeedMessage) Marshal() []byte {
	bytes := make([]byte, VSpeedMessageLen)

	bytes[0] = byte(VSpeed.Type)
	copy(bytes[1:5], VSpeed.OriginatorIP.To4())
	binary.LittleEndian.PutUint32(bytes[5:9], VSpeed.Speed)

	return bytes
}

func UnmarshalVSpeed(data []byte) (*VSpeedMessage, error) {
	if len(data) < VSpeedMessageLen {
		return nil, fmt.Errorf("VSpeed message length is %d, expected %d", len(data), VSpeedMessageLen)
	}

	VSpeed := &VSpeedMessage{}
	VSpeed.Type = uint8(data[0])
	VSpeed.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VSpeed.Speed = binary.LittleEndian.Uint32(data[5:9])

	return VSpeed, nil
}

func (VSpeed *VSpeedMessage) String() string {
	return fmt.Sprintf("VSpeed: Type: %d, OriginatorIP: %s, Speed: %d", VSpeed.Type, VSpeed.OriginatorIP.String(), VSpeed.Speed)
}
