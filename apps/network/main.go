package main

import (
	"log"

	. "github.com/aaarafat/vanessa/apps/network/packetFilter"
)




func main() {
	packetfilter, err := NewPacketFilter()

	if err != nil {
		log.Panicf("failed to create packet filter: %v", err)
	}

	packetfilter.Start()
	

	// wait for the program to exit
	select {}
}