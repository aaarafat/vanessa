package protocols

import "net"

type UnicastProtocol interface {
	GetRoute(destination net.IP) (*VRoute, bool)
	BuildRoute(destination net.IP) (started bool)
	Start()
	Close()
}
