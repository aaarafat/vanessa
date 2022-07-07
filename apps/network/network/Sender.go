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
		n.addToBuffer(packet, destIP)
		n.unicastProtocol.BuildRoute(destIP)
		return
	}
	n.forwarders[route.Interface].ForwardTo(packet, route.NextHop)
}


func (n *NetworkLayer) addToBuffer(packet []byte, destIP net.IP) {
	n.packetBuffer.Add(packet, destIP)
	log.Printf("Added to buffer: %s\n", destIP)
}

func (n *NetworkLayer) onPathDiscovery(destIP net.IP) {
	data, ok := n.packetBuffer.Get(destIP)
	if ok {
		for _, packet := range data {
			go n.SendUnicast(packet, destIP)
		}
		n.packetBuffer.Del(destIP)
		log.Printf("Removed from buffer: %s\n", destIP)
	}
}