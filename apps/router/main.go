package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/aaarafat/vanessa/apps/router/aodv"
)

func main() {
	args := os.Args
	if len(args) < 3 {
			fmt.Println("Please put station name and stations length");
			os.Exit(1);
	}

	station := args[1]
	numberOfNodes, err := strconv.Atoi(args[2])

	if err != nil {
		fmt.Println("Number of nodes must be an integer");
		os.Exit(1);
	}

	a, err := aodv.New(station, numberOfNodes)

	if err != nil {
		fmt.Println("Error in ceating aodv router", err);
		os.Exit(1);
	}

	a.Run()
}
