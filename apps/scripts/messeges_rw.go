package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func read(d *DataLinkLayerChannel, index int) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(payload), addr.String(), index)
		VREP := UnmarshalVOREP(payload)
		log.Printf("Received Obstacles: /n %v", VREP.String())
	}

}

func MyIP(ifi *net.Interface) (net.IP, bool, error) {

	addresses, err := ifi.Addrs()
	if err != nil {
		return nil, false, err
	}
	address := addresses[0]
	s :=strings.Split(address.String(), "/")[0]
	ip := net.ParseIP(s)
	if ip.To4()!= nil {
		return ip,false,nil
	} else if ip.To16()!=nil {
		return ip,true,nil
	}else {
		return nil, false, fmt.Errorf("IP can't be resolved")
	}
}

func main() {

	drsu, err := NewDataLinkLayerChannelWithInterface(VDATAEtherType, 2)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go read(drsu, 2)
	ifis,_ := net.Interfaces()
	ip,_,_ := MyIP(&ifis[2]) 
	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			log.Println("Sending Heartbeat")
			VHB := NewVHBeatMessage(ip , Position{Lat: 0, Lng: 0})
			drsu.Broadcast([]byte(VHB.Marshal()))
		case 1:
			var lng , lat float64
			fmt.Scanf("%f %f", &lat , &lng)
			log.Println("Sending Obstacle Alert")
			VOb := NewVObstacleMessage(ip, Position{Lat: lat, Lng: lng},0)
			drsu.Broadcast([]byte(VOb.Marshal()))
		case 2:
			log.Println("Sending Obstacles REQ")
			VOREQ := NewVOREQMessage(ip)
			drsu.Broadcast([]byte(VOREQ.Marshal()))
		}
		
	}
}
