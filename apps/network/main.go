package main

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

func main() {

	rreq := NewRREQMessage(net.IPv4(1,2,3,4), net.IPv4(5,6,7,8))
	bytes := rreq.Marshal()

	fmt.Println(rreq.String())

	rreq2, err := UnmarshalRREQ(bytes)
	if err != nil {
		panic(err)
	}
	fmt.Println(rreq2.String())
}