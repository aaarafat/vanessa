package datalink

import (
	"log"
	"net"
)


func CreateChannel(eth Ethertype, ifiIndex int) (*DataLinkLayerChannel, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return nil, err
	}
	ifi := interfaces[ifiIndex]
	return newDataLinkLayerChannel(eth,ifi,ifiIndex)
}