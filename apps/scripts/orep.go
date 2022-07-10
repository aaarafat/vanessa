package main

import (
	"fmt"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)



func main() {
	// obstacles as a pair of x,y list
	obstacles := []Position{
		{Lat: 1.5, Lng: 2.8},
		{Lat: 3.156, Lng: 4.20},
		{Lat: 5.7, Lng: 6.9},
		{Lat: 325.7, Lng: 69.69},
		{Lat: 85.7298, Lng: 123.45},
	}
	x:= obstacles[0].Marshal()
	fmt.Println(x)
	y:= UnmarshalPosition(x)
	fmt.Println(y)


	fmt.Println("\nCheck VOREP\n")
	VOREP := NewVOREPMessage(obstacles)

	bytes := VOREP.Marshal()
	fmt.Println(bytes)
	fmt.Println(VOREP.Obstacles)
	fmt.Println(VOREP.String())
	VOREP2 := UnmarshalVOREP(bytes)
	fmt.Println(VOREP2.String())
	// convert the obstacles byte array to a list of pair of x,y
	obstacles2 := make([][2]float64, int(VOREP2.Length))
	fmt.Println("Check the obstcales")
	for i := 0; i < len(obstacles2); i++ {
		obstacles2[i][0] = Float64FromBytes(VOREP2.Obstacles[i*16:i*16+8])
		obstacles2[i][1] = Float64FromBytes(VOREP2.Obstacles[i*16+8:i*16+16])
	}
	fmt.Println(obstacles2)
	fmt.Println("\nCheck VOREQ\n")
	VOREQ := NewVOREQMessage(net.IPv4(255,255,255,255),obstacles)
	bytes = VOREQ.Marshal()
	fmt.Println(bytes)
	fmt.Println(VOREQ.Obstacles)
	fmt.Println(VOREQ.String())
	VOREQ2 := UnmarshalVOREQ(bytes)
	fmt.Println(VOREQ2.String())
	// convert the obstacles byte array to a list of pair of x,y
	obstacles2 = make([][2]float64, int(VOREQ2.Length))
	fmt.Println("Check the obstcales")
	for i := 0; i < len(obstacles2); i++ {
		obstacles2[i][0] = Float64FromBytes(VOREQ2.Obstacles[i*16:i*16+8])
		obstacles2[i][1] = Float64FromBytes(VOREQ2.Obstacles[i*16+8:i*16+16])
	}
	fmt.Println(obstacles2)

}