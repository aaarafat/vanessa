package main

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

const MaxUint = ^uint32(0)

func main() {
	srcIP := net.ParseIP("192.168.0.1")

	VHBeat := NewVHBeatMessage(srcIP, Position{Lat: 1.5, Lng: 2.8})
	bytes := VHBeat.Marshal()
	fmt.Println(bytes)
	VHBeat2, err := UnmarshalVHBeat(bytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(VHBeat2.String())

	VObstacle := NewVObstacleMessage(srcIP, Position{Lat: 89.554, Lng: 260.15664}, 0)
	bytes = VObstacle.Marshal()
	fmt.Println(bytes)
	VObstacle2, err := UnmarshalVObstacle(bytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(VObstacle2.String())

}
