package messages

import "net"

// Heart beat message from car to the associated RSU
type VHBeatMessage struct {
	// Type
	Type uint8  
	// The IP address of the node that originated the RREQ.
	OriginatorIP net.IP
	// The X cooardinates of the car 
	PositionX uint32 
	// The Y cooardinates of the car
	PositionY uint32
}



const (
	VHBeatMessageLen = 13
)

const (
	VHBeatType uint8 = 1
)
