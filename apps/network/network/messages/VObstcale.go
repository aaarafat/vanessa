package messages

import (
	"encoding/binary"
	"fmt"
	"net"
)




func (VObstacle *VObstacleMessage) Marshal() []byte {
	bytes := make([]byte, VObstacleMessageLen)

	bytes[0] = byte(VObstacle.Type)
	copy(bytes[1:5], VObstacle.OriginatorIP.To4())
	binary.LittleEndian.PutUint32(bytes[5:], VObstacle.PositionX)
	binary.LittleEndian.PutUint32(bytes[9:], VObstacle.PositionY)
	bytes[13] = byte(VObstacle.Clear)


	return bytes
}

//Create a new VObstacle message
func NewVObstacleMessage(OriginatorIP net.IP, PositionX uint32, PositionY uint32, Clear uint8) *VObstacleMessage {
	return &VObstacleMessage{
		Type: VObstacleType,
		OriginatorIP: OriginatorIP,
		PositionX: PositionX,
		PositionY: PositionY,
		Clear: Clear,
	}
}



func UnmarshalVObstacle(data []byte) (*VObstacleMessage, error) {
	if len(data) < VObstacleMessageLen {
		return nil, fmt.Errorf("VObstacle message length is %d, expected %d", len(data), VObstacleMessageLen)
	}

	VObstacle := &VObstacleMessage{}
	VObstacle.Type = uint8(data[0])
	VObstacle.OriginatorIP = net.IPv4(data[1], data[2], data[3], data[4])
	VObstacle.PositionX = binary.LittleEndian.Uint32(data[5:9])
	VObstacle.PositionY = binary.LittleEndian.Uint32(data[9:13])
	VObstacle.Clear = uint8(data[13])
	return VObstacle, nil
}


// print the VObstacle message
func (VObstacle *VObstacleMessage) String() string {
	return fmt.Sprintf("VObstacle: Type: %d, OriginatorIP: %s, PositionX: %d, PositionY: %d, Clear: %d", VObstacle.Type, VObstacle.OriginatorIP.String(), VObstacle.PositionX, VObstacle.PositionY, VObstacle.Clear)
}

