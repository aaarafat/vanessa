package main

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

var channels []*DataLinkLayerChannel

func readFromInterface(d *DataLinkLayerChannel, index int) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(payload), addr.String(), index)
		BroadcastToInterfaces(payload, index)

	}

}
func BroadcastToInterfaces(payload []byte, from int) {
	for index, c := range channels {
		if index+1 == from {
			continue
		}
		c.Broadcast(payload)
	}
}

func main() {

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Panicln("No interfaces")
	}

	for index, iface := range interfaces {
		log.Println(index, iface)
		if index > 0 {

			c, err := NewDataLinkLayerChannelWithInterface(VEtherType, index)
			if err != nil {
				log.Panicln("No interfaces")
			}
			channels = append(channels, c)
			log.Println("Created channel for ", index)
			go readFromInterface(c, index)
		}
	}

	select {}

}
