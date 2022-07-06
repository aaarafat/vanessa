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
type VObstacleMessage struct {
	// Type
	Type uint8  
	// The IP address of the node that originated the Obstcale alert.
	OriginatorIP net.IP
	// The X cooardinates of the obstacle 
	PositionX uint32 
	// The Y cooardinates of the obstacle
	PositionY uint32
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
	// length of the list of obstacles
	Length uint8
	// content of the list of obstacles
	Obstacles string

}

const (
	VHBeatMessageLen = 13
	VObstacleMessageLen = 14
	VOREQMessageLen = 5
)

const (
	VHBeatType uint8 = 1
	VObstacleType uint8 = 2
	VOREQType uint8 = 3
	VOREPType uint8 = 4
)
