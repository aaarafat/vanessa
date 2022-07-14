package messages

import (
	"net"

	. "github.com/aaarafat/vanessa/libs/vector"
)

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
	// length of the list of Obstcales
	Length uint8
	// content of the list of Obstcales
	Obstacles []byte
}

type VOREPMessage struct {
	// Type
	Type uint8
	// length of the list of Obstcales
	Length uint8
	// content of the list of Obstcales
	Obstacles []byte
}

type VPathDiscoveryMessage struct {
	// Type
	Type uint8
	// The IP address of the node that is discovered.
	DestIP net.IP
}

type VZoneMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Obstcale alert.
	OriginatorIP net.IP
	// Position of the obstacle
	Position Position
	// Speed of the vehicle
	Speed uint32
}

const (
	VHBeatMessageLen         = 21
	VObstacleMessageLen      = 22
	VOREQMessageLen          = 6
	VOREPMessageLen          = 2
	VPathDiscoveryMessageLen = 5
	VZoneMessageLen          = 25
)

const (
	VHBeatType         uint8 = 1
	VObstacleType      uint8 = 2
	VOREQType          uint8 = 3
	VOREPType          uint8 = 4
	VPathDiscoveryType uint8 = 5
	VZoneType          uint8 = 6
)
