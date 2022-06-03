package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"time"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

func neighborUpdate(d *DataLinkLayerChannel, nt *VNeighborTable) {
	for {
		
		payload, addr, err := d.Read()
		// fmt.Println(net.HardwareAddr(addr.String()))
		entry := NewNeighborEntry(net.IP(payload))
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
	args := os.Args
	if len(args) < 3 {
			fmt.Println("Please put station name and stations length");
			os.Exit(1);
	}
	
	station := args[1]
	fmt.Println(station)
	intf := station+"-wlan1"
	out, err := exec.Command("iw", "dev", intf, "link").Output()
	if err != nil {
		log.Panic(err)
	}
	cmdOut := string(out)
	fmt.Println(cmdOut)
	rsuMAC := "" 
	if strings.Contains(cmdOut, "Not connected") {
		println(station, "is not associated")
	} else {
		println(station, "is associated")
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


	ip := args[2]
	neibourTable := NewNeighborTable()

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
