package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

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

func main() {

	d, err := NewDataLinkLayerChannel(VEtherType)
	drsu, err := NewDataLinkLayerChannelWithInterface(VEtherType, 2)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go read(d, 1)
	go read(drsu, 2)

	var message string
	var mtype int
	for {

		fmt.Scanf("%d %s", &mtype, &message)
		switch mtype {
		case 0:
			d.Broadcast([]byte(message))
		case 1:
			drsu.Broadcast([]byte(message))
		case 2:
			mac, err := net.ParseMAC(message)
			if err != nil {
				log.Panicln("not a valid mac")
			}

			log.Println(mac.String())
			drsu.SendTo([]byte("Test from car"), mac)
		case 3:
			mac, err := net.ParseMAC(message)
			if err != nil {
				log.Panicln("not a valid mac")
			}

			log.Println(mac.String())
			d.SendTo([]byte("unicast"), mac)
		case 4:
			interfaces, _ := net.Interfaces()
			addresses, _ := interfaces[2].Addrs()
	
			address := addresses[0]
			s :=strings.Split(address.String(), "/")[0]
			log.Println(s)
			drsu.Broadcast([]byte(s))
		}
		
	}
}
