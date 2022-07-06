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
	buf, ok := n.packetBuffer.Get(destIP.String())
	if ok {
		n.packetBuffer.Set(destIP.String(), append(buf.([][]byte), packet))
	} else {
		n.packetBuffer.Set(destIP.String(), [][]byte{packet})
	}
	log.Printf("Added to buffer: %s\n", destIP)
}

func (n *NetworkLayer) onPathDiscovery(destIP net.IP) {
	buf, ok := n.packetBuffer.Get(destIP.String())
	if ok {
		for _, packet := range buf.([][]byte) {
			go n.SendUnicast(packet, destIP)
		}
		n.packetBuffer.Del(destIP.String())
		log.Printf("Removed from buffer: %s\n", destIP)
	}
}