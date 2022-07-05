package network

import (
	"log"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func (n *NetworkLayer) listen(channel *DataLinkLayerChannel) {
	log.Printf("Listening for AODV packets on channel: %d....\n", channel.IfiIndex)
	for {
		packet, _, err := channel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		go n.handleMessage(packet)
	}
}

func (n *NetworkLayer) handleMessage(packet []byte) {
	n.ipConn.Forward(packet)
}
