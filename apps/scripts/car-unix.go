package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
)

type Position struct {
	Lat float64
	Lng float64
}

func reader(d *json.Decoder) {
	var position = Position{Lat: 5, Lng: 5}
	for {
		err := d.Decode(&position)
		if err != nil {
			return
		}
		fmt.Printf("Client got: %f %f\n", position.Lat, position.Lng)
	}
}

func main() {

	var id int
	flag.IntVar(&id, "id", 0, "ID of the car")
	flag.Parse()

	socketAddress := fmt.Sprintf("/tmp/car%d.socket", id)

	err := os.RemoveAll(socketAddress)
	if err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}

	addr, err := net.ResolveUnixAddr("unix", socketAddress)
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		fmt.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	d := json.NewDecoder(conn)
	fmt.Printf("Listening to %s ..\n", socketAddress)
	reader(d)
}
