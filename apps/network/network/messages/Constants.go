package messages

import "net"

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

// Postion to byte array marashalling
func (p Position) Marshal() []byte {
	return append(Float64bytes(p.Lat), Float64bytes(p.Lng)...)
}

// Position from byte array unmarshalling
func UnmarshalPosition(data []byte) Position {
	return Position{
		Lat: Float64frombytes(data[:8]),
		Lng: Float64frombytes(data[8:]),
	}
}

// Heart beat message from car to the associated RSU
type VHBeatMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Heart Beat.
	OriginatorIP net.IP
	// Position of the node that originated the Heart Beat.
	Position Position
}

// Obstcale found message from car to the associated RSU
type VObstacleMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Obstcale alert.
	OriginatorIP net.IP
	// Position of the obstacle
	Position Position
	// Set obstacle to true if the obstacle is detected
	Clear uint8
}

type VOREQMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Obstcales request.
	OriginatorIP net.IP
}

type VOREPMessage struct {
	// Type
	Type uint8
	// length of the list of Obstcales
	Length uint8
	// content of the list of Obstcales
	Obstcales []byte
}

const (
	VHBeatMessageLen    = 21
	VObstacleMessageLen = 22
	VOREQMessageLen     = 5
)

const (
	VHBeatType    uint8 = 1
	VObstacleType uint8 = 2
	VOREQType     uint8 = 3
	VOREPType     uint8 = 4
)
