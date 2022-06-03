package main

import (
	"encoding/json"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/protocols"
)

func main() {
	dsv := NewDSV()

	log.Println("DSV: ", dsv)

	interfaces, err := net.Interfaces()
	iface := interfaces[1]

	ip, _, err := GetMyIPs(&iface)

	log.Println(ip)

	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return
	}

	log.Println("Starting DSV on: ", iface.Name)
	bs, _ := json.Marshal(dsv)
	log.Println()
	log.Printf("Received \"%s\"", bs)

	
	if ip.String() != "10.0.0.3" {
		go dsv.Send(ip, net.ParseIP("10.0.0.3"))
	}
	
	dsv.Read()

}