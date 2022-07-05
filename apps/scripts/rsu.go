package main

import (
	"fmt"
	"log"

	// "os"
	// "os/exec"

	// "strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/rsu"
)



func listenAndBroadcast(d *DataLinkLayerChannel, e *DataLinkLayerChannel) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [ %s ]", string(payload), addr.String())
		e.Broadcast(payload)

	}

}

func read(d *DataLinkLayerChannel, index int) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(payload), addr.String(), index)
	}

}

func test(d *DataLinkLayerChannel, ARP RSUARP) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [%s] for ARP", string(payload), addr.String())
		ARP.Set(string(payload),addr)
	}

}

func main() {

	eChannel, err := NewDataLinkLayerChannelWithInterface(VEtherType, 1)
	wChannel, err := NewDataLinkLayerChannelWithInterface(VEtherType, 3)

	// wChannel, err := NewDataLinkLayerChannelWithIntf(VIEtherType, wintf_name)
	ARP := NewRSUARP()
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go read(eChannel, 1)
	// go listenAndBroadcast(wChannel, eChannel)
	go test(wChannel,*ARP)
	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			wChannel.Broadcast([]byte("HI to Car"))
		case 1:
			eChannel.Broadcast([]byte("HI to RSU"))
		case 2:
			ARP.Print()
		
		}
		
	}
}
