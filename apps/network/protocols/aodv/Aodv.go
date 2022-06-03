package aodv

import (
	"log"
	"net"
	"time"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/tables"
)

type Aodv struct {
	channel *DataLinkLayerChannel
	neighborTable *VNeighborTable
	routingTable *VRoutingTable
	srcIP net.IP
}

const (
	HopThreshold = 3
)

func NewAodv(srcIP net.IP) *Aodv {
	d, err := NewDataLinkLayerChannel(VAODVEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel: d,
		neighborTable: NewNeighborTable(),
		routingTable: NewVRoutingTable(),
		srcIP: srcIP,
	}
}


func (a *Aodv) SendRREQ(destination net.IP) {
	rreq := NewRREQMessage(a.srcIP, destination)
	log.Printf("Sending: %s\n", rreq.String())

	// broadcast the RREQ
	a.channel.Broadcast(rreq.Marshal())
}

func (a *Aodv) handleRREQ(payload []byte, from net.HardwareAddr) {
	rreq, err := UnmarshalRREQ(payload)
	if err != nil {
		log.Printf("Failed to unmarshal RREQ: %v\n", err)
		return
	}
	
	if rreq.HopCount > HopThreshold || rreq.OriginatorIP.Equal(a.srcIP) {
		// drop the packet
		return
	}

	log.Printf("Received: %s\n", rreq.String())
	// update the routing table
	go a.routingTable.Update(from, rreq.OriginatorIP, rreq.HopCount)

	// check if the RREQ is for me
	if rreq.DestinationIP.Equal(a.srcIP) {
		// send a RREP

		log.Println("Sending RREP to: ", from)
	} else {
		// increment hop count
		rreq.HopCount = rreq.HopCount + 1
		// forward the RREQ
		go a.channel.Broadcast(rreq.Marshal())
	}
}

func (a *Aodv) Listen() {
	log.Println("Listening for AODV packets...")
	for {
		// TODO : Remove this sleep
		time.Sleep(time.Second * 5) 

		payload, addr, err := a.channel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}

		// get the type
		msgType := uint8(payload[0])

		// handle the message
		switch msgType {
		case RREQType:
			go a.handleRREQ(payload, addr)
		default:
			log.Println("Unknown message type: ", msgType)
		}
	}
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.Listen()
}

func (a *Aodv) Close() {
	a.channel.Close()
}