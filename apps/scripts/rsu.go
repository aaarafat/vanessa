package main

import (
	"fmt"
	"log"
	"net"
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
	
	interfaces, err := net.Interfaces()
	ifi, err := net.InterfaceByName(interfaces[6].Name)
	addrs, _ := ifi.Addrs()

	for _, intf := range interfaces {
		println(intf.Name)
		addrs, _ := ifi.Addrs()
		for _, addr := range addrs {
			println(addr.String())
	
		}
	}
	println("number of addresses", len(addrs), ifi.Name)
	// ip := addrs[0].String()
	


	iChannel, err := NewDataLinkLayerChannel(VIEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go listen(iChannel)

	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			iChannel.Broadcast([]byte("HI"))
		case 1:
			log.Print("flooding")
		}
		time.Sleep(5 * time.Second)
	}
}
