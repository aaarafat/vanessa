package network

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func (n *NetworkLayer) listen(channel *DataLinkLayerChannel) {
	log.Printf("Listening for DATA packets on channel: %s....\n", channel.Ifi.Name)
	for {
		packet, addr, err := channel.Read()
		if err != nil {
			continue
		}
		go n.handleMessage(packet, addr)
	}
}

func (n *NetworkLayer) flooderListiner() {
	log.Printf("Listening for DATA packets on Flooder....\n")
	for {
		packet, addr, err := n.broadcastProtocol.Read()
		if err != nil {
			continue
		}
		go n.handleMessage(packet, addr)
	}
}

func (n *NetworkLayer) handleMessage(packet []byte, from net.HardwareAddr) {
	log.Printf("Forwarding message with size %d from %s\n", len(packet), from)
	n.ipConn.Forward(packet)
}
