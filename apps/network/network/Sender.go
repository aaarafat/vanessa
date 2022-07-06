package network

import (
	"log"
	"net"
)

func (n *NetworkLayer) Send(payload []byte, srcIP net.IP, destIP net.IP) {
	n.ipConn.Write(payload, srcIP, destIP)
}


func (n *NetworkLayer) SendUnicast(packet []byte, destIP net.IP) {
	route, found := n.unicastProtocol.GetRoute(destIP)
	if !found {
		log.Printf("No route to %s\n", destIP)
		n.unicastProtocol.BuildRoute(destIP)
		return
	}
	n.forwarders[route.Interface].ForwardTo(packet, route.NextHop)
}

