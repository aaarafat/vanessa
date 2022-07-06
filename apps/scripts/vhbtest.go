package main

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


const MaxUint = ^uint32(0) 

func main() {
	srcIP := net.ParseIP("192.168.0.1")

	VHBeat := NewVHBeatMessage(srcIP, MaxUint,48)
	bytes := VHBeat.Marshal()
	fmt.Println(bytes)
	VHBeat2, err := UnmarshalVHBeat(bytes)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(VHBeat2.String())
}