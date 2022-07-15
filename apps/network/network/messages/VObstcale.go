package messages

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/libs/vector"
)

func (VObstacle *VObstacleMessage) Marshal() []byte {
	bytes := make([]byte, VObstacleMessageLen)

	bytes[0] = byte(VObstacle.Type)
	copy(bytes[1:5], VObstacle.OriginatorIP.To4())
	copy(bytes[5:21], VObstacle.Position.Marshal())
	bytes[21] = byte(VObstacle.Clear)

	return bytes
}

//Create a new VObstacle message
func NewVObstacleMessage(OriginatorIP net.IP, position Position, Clear uint8) *VObstacleMessage {
	return &VObstacleMessage{
		Type:         VObstacleType,
		OriginatorIP: OriginatorIP,
		Position:     position,
		Clear:        Clear,
	}
}

func UnmarshalVObstacle(data []byte) (*VObstacleMessage, error) {
	if len(data) < VObstacleMessageLen {
		return nil, fmt.Errorf("VObstacle message length is %d, expected %d", len(data), VObstacleMessageLen)
	}

	VObstacle := &VObstacleMessage{}
	VObstacle.Type = uint8(data[0])
	VObstacle.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VObstacle.Position = UnmarshalPosition(data[5:21])
	VObstacle.Clear = uint8(data[21])
	return VObstacle, nil
}

// print the VObstacle message
func (VObstacle *VObstacleMessage) String() string {
	return fmt.Sprintf("VObstacle: Type: %d, OriginatorIP: %s, Lat: %f, Lng: %f, Clear: %d", VObstacle.Type, VObstacle.OriginatorIP.String(), VObstacle.Position.Lat, VObstacle.Position.Lng, VObstacle.Clear)
}
