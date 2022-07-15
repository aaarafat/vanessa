package messages

import (
	"fmt"
	"net"
)

func NewVPathDiscoveryMessage(destIP net.IP) *VPathDiscoveryMessage {
	return &VPathDiscoveryMessage{
		Type:   VPathDiscoveryType,
		DestIP: destIP,
	}
}

func (p *VPathDiscoveryMessage) Marshal() []byte {
	bytes := make([]byte, VPathDiscoveryMessageLen)
	bytes[0] = byte(p.Type)
	copy(bytes[1:5], p.DestIP.To4())
	return bytes
}

func UnmarshalVPathDiscovery(data []byte) (*VPathDiscoveryMessage, error) {
	if len(data) < VPathDiscoveryMessageLen {
		return nil, fmt.Errorf("VPathDiscovery message length is %d, expected %d", len(data), VPathDiscoveryMessageLen)
	}

	p := &VPathDiscoveryMessage{}
	p.Type = uint8(data[0])
	p.DestIP = net.IPv4(data[1], data[2], data[3], data[4])
	return p, nil
}

func (p *VPathDiscoveryMessage) String() string {
	return fmt.Sprintf("VPathDiscovery: Type: %d, DestIP: %s", p.Type, p.DestIP.String())
}
