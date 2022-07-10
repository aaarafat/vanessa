package network

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (n *NetworkLayer) Send(payload []byte, srcIP net.IP, destIP net.IP) {
	n.ipConn.Write(payload, srcIP, destIP)
}

func (n *NetworkLayer) SendUnicast(packet []byte, destIP net.IP) {
	log.Printf("Sending unicast to %s\n", destIP)
	route, found := n.unicastProtocol.GetRoute(destIP)
	if !found {
		log.Printf("No route to %s\n", destIP)
		n.addToBuffer(packet, destIP)
		n.unicastProtocol.BuildRoute(destIP)
		return
	}
	n.forwarders[route.Interface].ForwardTo(packet, route.NextHop)
}

func (n *NetworkLayer) SendBroadcast(packet []byte, from net.IP) {
	log.Printf("Sending broadcast from %s\n", from)
	n.forwarders[1].ForwardToAllExceptIP(packet, from)
}

func (n *NetworkLayer) addToBuffer(packet []byte, destIP net.IP) {
	n.packetBuffer.Add(packet, destIP)
	log.Printf("Added to buffer: %s\n", destIP)
}

func (n *NetworkLayer) onPathDiscovery(destIP net.IP) {
	go n.sendPathDiscoveryMessage(destIP)
	data, ok := n.packetBuffer.Get(destIP)
	if ok {
		for _, packet := range data {
			go n.SendUnicast(packet, destIP)
		}
		n.packetBuffer.Del(destIP)
		log.Printf("Removed from buffer: %s\n", destIP)
	}
}

func (n *NetworkLayer) sendPathDiscoveryMessage(destIP net.IP) {
	if destIP.Equal(net.ParseIP(ip.RsuIP)) {
		log.Printf("Sending VOREQ to RSU\n")
		data := NewVOREQMessage(n.ip).Marshal()
		n.ipConn.Write(data, n.ip, net.ParseIP(ip.RsuIP))
	}
}
