package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	. "github.com/aaarafat/vanessa/apps/network/packetFilter"
	. "github.com/aaarafat/vanessa/apps/network/rsu"
)

func initLogger(debug bool, id int, name string) {
	log.SetPrefix("[vanessa]")
	if !debug {
		log.SetOutput(os.Stdout)
		return // don't do anything if debug is false
	}
	err := os.MkdirAll("/logs", 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating logs directory: %s\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(fmt.Sprintf("/logs/%s%d.log", name, id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}
	log.SetOutput(file)
	log.SetFlags(log.LstdFlags)
}

func main() {
	var id int
	var name string
	var debug bool
	flag.IntVar(&id, "id", 0, "id of the car")
	flag.StringVar(&name, "name", "", "name of the unit (rsu or car)")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	initLogger(debug, id, name)

	if name == "rsu" {
		// create a new RSU
		rsu := NewRSU()

		// start the RSU
		rsu.Start()

		defer rsu.Close()
	} else if name == "car" {
		packetfilter, err := NewPacketFilter(id)

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