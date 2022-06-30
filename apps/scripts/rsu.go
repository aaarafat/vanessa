package main

import (
	"fmt"
	"log"
	"net"
	"os"

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
	intf_name := os.Args[1]+"-wlan1"
	// wintf_name := os.Args[1]+"-eth2"
	interfaces, err := net.Interfaces()
	ifi, err := net.InterfaceByName(interfaces[6].Name)
	addrs, _ := ifi.Addrs()

	for index, iface := range interfaces {
		log.Println(index, iface)
	}

	println("number of addresses", len(addrs), intf_name)
	// ip := addrs[0].String()
	


	iChannel, err := NewDataLinkLayerChannelWithIntf(VIEtherType, intf_name)
	// wChannel, err := NewDataLinkLayerChannelWithIntf(VIEtherType, wintf_name)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go listen(iChannel)
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
