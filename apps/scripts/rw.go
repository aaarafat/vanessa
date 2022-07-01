package main

import (
	"fmt"
	"log"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func read(d *DataLinkLayerChannel) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [ %s ]", string(payload), addr.String())
	}

}

func main() {

	d, err := NewDataLinkLayerChannel(VEtherType)
	drsu, err := NewDataLinkLayerChannelWithInterface(VEtherType,2)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go read(d)
	go read(drsu)

	var message string
	var mtype int
	for {

		fmt.Scanf("%d %s", &mtype, &message)
		switch mtype {
		case 0:
			d.Broadcast([]byte(message))
		case 1:
			drsu.Broadcast([]byte(message))
		}
	}
}
