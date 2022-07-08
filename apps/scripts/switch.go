package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

var channels []*DataLinkLayerChannel

func readFromInterface(d *DataLinkLayerChannel, index int) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(payload), addr.String(), index)
		BroadcastToInterfaces(payload, index)

	}

}
func BroadcastToInterfaces(payload []byte, from int) {
	for index, c := range channels {
		if index+1 == from {
			continue
		}
		c.Broadcast(payload)
	}
}


func initLogger(debug bool) {
	log.SetPrefix("[vanessa]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !debug {
		log.SetOutput(os.Stdout)
		return // don't do anything if debug is false
	}
	// delete logs
	err := os.MkdirAll("/var/log/vanessa", 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating logs directory: %s\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(fmt.Sprintf("/var/log/vanessa/s0.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}
	log.SetOutput(file)
}

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	initLogger(debug)

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Panicln("No interfaces")
	}

	for index, iface := range interfaces {
		log.Println(index, iface)
		if index > 0 {

			c, err := NewDataLinkLayerChannelWithInterface(VEtherType, index)
			if err != nil {
				log.Panicln("No interfaces")
			}
			channels = append(channels, c)
			log.Println("Created channel for ", index)
			go readFromInterface(c, index)
		}
	}

	select {}

}
