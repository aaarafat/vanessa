package main

import (
	"log"
	"net"
	"os"

	. "github.com/aaarafat/vanessa/apps/network/packetFilter"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

func print(data []byte) {
	log.Printf("Received data: %s\n", data)
}

func rsu() {
	aodv := aodv.NewAodv(net.ParseIP(aodv.RsuIP), print)	
	aodv.Start()
}


func car() {
	packetfilter, err := NewPacketFilter()

	if err != nil {
		log.Panicf("failed to create packet filter: %v", err)
	}

	packetfilter.Start()
}

func main() {
	name := os.Args[1]

	if name == "rsu" {
		rsu()
	} else if name == "car" {
		car()
	} else {
		log.Panicf("invalid name: %s, Please Enter car or rsu\n", name)
	}
	

	// wait for the program to exit
	select {}
}