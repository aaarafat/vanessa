package main

import (
	"fmt"
	"log"
	"net"
	// "os"
	"os/exec"
	"time"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func neighborUpdate(d *DataLinkLayerChannel, nt *VNeighborTable) {
	for {
		
		payload, addr, err := d.Read()
		// fmt.Println(net.HardwareAddr(addr.String()))
		entry := NewNeighborEntry(net.IP(payload), addr)
		nt.Set(addr.String() ,entry)
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [ %s ]", string(payload), addr.String())
		// ent, _ := nt.Get(addr.String())
		// fmt.Println("table Entry ", string(ent.IP))
		// fmt.Println("number of enries %d ", nt.Len())
		nt.Print()
	}

}

func main() {
	
	interfaces, err := net.Interfaces()
	adhoc_ifi, err := net.InterfaceByName(interfaces[1].Name)
	net_ifi, err := net.InterfaceByName(interfaces[2].Name)
	addrs , _ := adhoc_ifi.Addrs()
	println(net_ifi.Name, net_ifi.HardwareAddr.String())
	ip := strings.Split(addrs[0].String(), "/")[0]
	println(ip)

	out, err := exec.Command("iw", "dev", net_ifi.Name, "link").Output()
	if err != nil {
		log.Panic(err)
	}
	cmdOut := string(out)
	// println(cmdOut)
	rsuMAC := "" 
	if strings.Contains(cmdOut, "Not connected") {
		println(net_ifi.Name, "is not associated")
	} else {
		println(net_ifi.Name, "is associated")
		arr := strings.Fields(cmdOut) 
		rsuMAC = arr[2] 
		for ind, v := range arr {    
			if strings.Contains(v, "ssid_"){
				println(ind, v)
				println("mac:", rsuMAC)
				break
			}
  		}
	}


	neibourTable := NewNeighborTable(net.IP(ip))

	d, err := NewDataLinkLayerChannel(VNDEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go neighborUpdate(d, neibourTable)

	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			d.Broadcast([]byte(ip))
		case 1:
			log.Print("flooding")
		}
		time.Sleep(5 * time.Second)
	}
}
