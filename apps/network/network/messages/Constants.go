package messages

import "net"

// Heart beat message from car to the associated RSU
type VHBeatMessage struct {
	// Type
	Type uint8  
	// The IP address of the node that originated the Heart Beat.
	OriginatorIP net.IP
	// The X cooardinates of the car 
	PositionX uint32 
	// The Y cooardinates of the car
	PositionY uint32
}

// Obstcale found message from car to the associated RSU
type VObstcaleMessage struct {
	// Type
	Type uint8  
	// The IP address of the node that originated the Obstcale alert.
	OriginatorIP net.IP
	// The X cooardinates of the obstacle 
	PositionX uint32 
	// The Y cooardinates of the obstacle
	PositionY uint32
}

const (
	VHBeatMessageLen = 13
	VObstcaleMessageLen = 13
)

const (
	VHBeatType uint8 = 1
	VObstcaleType uint8 = 2
)
