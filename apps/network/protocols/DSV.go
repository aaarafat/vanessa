package protocols

import (
	"encoding/json"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/tables"
)
type DSV struct {
	etherType Ethertype
	datalink	*DataLinkLayerChannel
	neighborTable		*VNeighborTable
	routingTable 		*VRoutingTable
}	
type DSVMessage struct {
	Source net.IP
	Destination net.IP 
	HopCount int
}

const (
	HopThreshold = 3
)

func (msg *DSVMessage) Print() {
	log.Printf("Source: %s, Destination: %s, HopCount: %d", msg.Source.String(), msg.Destination.String(), msg.HopCount)
}

func NewDSV() *DSV {
	d, err := NewDataLinkLayerChannel(VDSVEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	neighborTable := NewNeighborTable()
	routingTable := NewVRoutingTable()

	return &DSV{
		etherType: VDSVEtherType,
		datalink: d,
		neighborTable: neighborTable,
		routingTable: routingTable,
	}
}


func (d* DSV) Broadcast(payload *DSVMessage) {
	payload_byte, err := json.Marshal(payload)
		if err != nil {
			log.Panic(err)
		}
	d.datalink.Broadcast(payload_byte)
}

func (d* DSV) SendTo(payload *DSVMessage, addr net.HardwareAddr) {
	payload_byte, _ := json.Marshal(payload)
	d.datalink.SendTo(payload_byte, addr)
}

func (d *DSV) Send(source, destination net.IP) {
	log.Println("Sending from: ", source, " to: ", destination)
	payload := DSVMessage{
		Source: source,
		Destination: destination,
		HopCount: 0,	
	}

	// check routing table
	route, exists := d.routingTable.Get(destination)
	if exists {
		log.Println("Found route: ", route)
		d.SendTo(&payload, route.NextHop)
	} else {
		log.Println("No route to destination: ", destination)
		d.Broadcast(&payload)
	}	
}

func (dsv* DSV) Read() {
	for {
		log.Println("waiting for message.....")
		payload_byte, addr, err := dsv.datalink.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}

		payload := DSVMessage{}
		err = json.Unmarshal(payload_byte, &payload)
		if err != nil {
			log.Fatalf("failed to unmarshal payload: %v", err)
		}
		log.Println("Received: ")
		payload.Print()

		// check if hop count equal to HopThreshold => stop
		if payload.HopCount == HopThreshold {
			// Stop
			continue 
		}

		// update Routing Table
		dsv.routingTable.Update(addr, payload.Source, payload.HopCount)


		// check if it is a neighbor message
		if payload.Destination.Equal(net.IP(dsv.datalink.Source)) {
			// send reply
			log.Printf("Sending reply........... to : %s", payload.Source.String())
			dsv.SendTo(&payload, addr)
		} else {
			payload.HopCount = payload.HopCount + 1 // increase hop count
			payload_byte, _ := json.Marshal(payload)
			dsv.datalink.Broadcast(payload_byte)
		}

	}

}
