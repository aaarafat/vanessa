package protocols

import "net"

type Protocol interface {
	GetRoute(destination net.IP) (*VRoute, bool)
	BuildRoute(destination net.IP) 
	Start()
	Close()
}