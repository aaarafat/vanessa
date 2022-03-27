package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
			fmt.Println("Please put station name and stations length");
			os.Exit(1);
	}

	station := args[1]
	numberOfNodes := args[2]

	fmt.Println("Station: ", station)
	fmt.Println("numberOfNodes: ", numberOfNodes)

	/*
	aodv := aodv.New(station, numberOfNodes)

	go aodv.Run()
	*/
}
