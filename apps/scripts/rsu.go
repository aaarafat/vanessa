package main

import (
	"fmt"
	"log"

	// "os"
	// "os/exec"
	"time"
	// "strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func listen(d *DataLinkLayerChannel) {
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

	iChannel, err := NewDataLinkLayerChannelWithInterface(VIEtherType, 2)
	// wChannel, err := NewDataLinkLayerChannelWithIntf(VIEtherType, wintf_name)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go read(iChannel)
	// go listen(wChannel)

	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			iChannel.Broadcast([]byte("HI"))
		case 1:
			// wChannel.Broadcast([]byte("HI"))
		}
		time.Sleep(5 * time.Second)
	}
}
