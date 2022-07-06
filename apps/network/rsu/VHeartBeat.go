package rsu

import "net"



type VHeartBeatMessage struct {
	// Type
	Type uint8  
	// The IP address of the node that originated the RREQ.
	OriginatorIP net.IP
	// The X cooardinates of the car 
	PositionX uint32 
	// The Y cooardinates of the car
	PositionY uint32
}


