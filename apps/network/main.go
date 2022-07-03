package main

import (
	"log"
	"os"

	. "github.com/aaarafat/vanessa/apps/network/packetFilter"
	. "github.com/aaarafat/vanessa/apps/network/rsu"
)

func main() {
	name := os.Args[1]

	if name == "rsu" {
		// create a new RSU
		rsu := NewRSU()

		// start the RSU
		rsu.Start()

		defer rsu.Close()
	} else if name == "car" {
		packetfilter, err := NewPacketFilter()

		if err != nil {
			log.Panicf("failed to create packet filter: %v", err)
		}

		packetfilter.Start()

		defer packetfilter.Close()
	} else {
		log.Panicf("invalid name: %s, Please Enter car or rsu\n", name)
	}
	

	// wait for the program to exit
	select {}
}