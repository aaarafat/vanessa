package protocols

import "net"

type UnicastProtocol interface {
	GetRoute(destination net.IP) (*VRoute, bool)
	BuildRoute(destination net.IP) (started bool)
	Start()
	Close()
}

type BroadcastProtocol interface {
	Flood(packet []byte)
	Forward(packet []byte, fromIP net.IP)
}
