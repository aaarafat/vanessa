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
		entry := NewNeighborEntry(net.IP(payload), addr)
		nt.Set(addr.String() ,entry)
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [ %s ]", string(payload), addr.String())
		
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

	neibourTable := NewNeighborTable(net.IP(ip))

	nChannel, err := NewDataLinkLayerChannel(VNDEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go neighborUpdate(nChannel, neibourTable)

	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			nChannel.Broadcast([]byte(ip))
			rsuMAC, SSID := getRSU(net_ifi.Name)
			println(rsuMAC)
			if  strings.Compare(SSID, "") != 0{
				fileName := "./"+SSID+".log"
				appendToFile(fileName, ip)
			}
			
		case 1:
			log.Print("flooding")
		}
		time.Sleep(5 * time.Second)
	}
}

//returns mac, ssid of assiciated rsu
func getRSU(intfName string) (string, string){

	out, err := exec.Command("iw", "dev", intfName, "link").Output()
	if err != nil {
		log.Panic(err)
	}
	cmdOut := string(out)
	// println(cmdOut)
	rsuMAC := "" 
	SSID := ""
	if strings.Contains(cmdOut, "Not connected") {
		println(intfName, "is not associated")
	} else {
		println(intfName, "is associated")
		arr := strings.Fields(cmdOut) 
		rsuMAC = arr[2] 
		for ind, v := range arr {    
			if strings.Contains(v, "ssid_"){
				println(ind, v)
				SSID = v
				println("mac:", rsuMAC)
				break
			}
  		}
	}
	return rsuMAC, SSID
}

func appendToFile(fileName, msg string) {
	f, err := os.OpenFile(fileName,
	os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString("Hi I am "+msg+"\n"); err != nil {
		log.Println(err)
	}
}