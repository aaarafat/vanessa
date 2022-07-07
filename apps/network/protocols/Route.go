package protocols

import (
	"fmt"
	"net"
)


type VRoute struct {
	Destination net.IP
	NextHop net.HardwareAddr
	Interface int
	Metric int
}

func NewVRoute(destination net.IP, nextHop net.HardwareAddr, ifi int, metric int) *VRoute {
	return &VRoute{
		Destination: destination,
		NextHop: nextHop,
		Interface: ifi,
		Metric: metric,
	}
}

func (v *VRoute) String() string {
	return fmt.Sprintf("%s %s %d %d", v.Destination, v.NextHop, v.Interface, v.Metric)
}