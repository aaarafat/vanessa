package main

import (
	"log"
	"net"
	"os"
	"strconv"

	
	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func main() {
	log.Println(os.Args[1])
	index, err := strconv.Atoi(os.Args[1])
	if err == nil {
	}

	interfaces, err := net.Interfaces()
	for index, iface := range interfaces {
		log.Println(index, iface)
	}
	log.Println(interfaces[index].Name)
	ifi, err := net.InterfaceByName(interfaces[index].Name)

	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}

	// Open a raw socket using same EtherType as our frame.
	iChannel, err := NewDataLinkLayerChannelWithIntf(VIEtherType, ifi.Name)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	for {
		
		payload, addr, err := iChannel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		log.Printf("Received \"%s\" from: [ %s ]", string(payload), addr.String())
		
	}
}
