package aodv

import (
	"log"
	"net"
	"sync"
	"time"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/tables"
)

type Aodv struct {
	channel *DataLinkLayerChannel
	neighborTable *VNeighborTable
	routingTable *VRoutingTable
	srcIP net.IP
	chLock sync.RWMutex
}

const (
	HopThreshold = 20
)

func NewAodv(srcIP net.IP) *Aodv {
	d, err := NewDataLinkLayerChannel(VAODVEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel: d,
		neighborTable: NewNeighborTable(srcIP),
		routingTable: NewVRoutingTable(),
		srcIP: srcIP,
	}
}

func (a *Aodv) sendToAllExcept(payload []byte, addr net.HardwareAddr) {
	a.neighborTable.Print()
	for item := range a.neighborTable.Iter() {
		stringMac := item.MAC.String()
		log.Println("Mac = ", stringMac)
		neighborMac := item.MAC

		if stringMac != addr.String() {
			log.Println("Sending to: ", stringMac)
			go a.channel.SendTo(payload, neighborMac)
		}
	}
}

func (a *Aodv) Send(payload []byte, dest net.IP) {
	// check if the destination in the routing table
	item, ok := a.routingTable.Get(dest);
	if ok {
		// send the packet
		go a.channel.SendTo(payload, item.NextHop)
	} else {
		// send a RREQ or RRER
		go a.SendRREQ(dest)
	}
}

func (a *Aodv) Broadcast(payload []byte) {
	a.channel.Broadcast(payload)
}

func (a *Aodv) SendRREQ(destination net.IP) {
	rreq := NewRREQMessage(a.srcIP, destination)
	log.Printf("Sending: %s\n", rreq.String())

	// broadcast the RREQ
	go a.Broadcast(rreq.Marshal())
}

func (a *Aodv) SendRREP(destination net.IP) {
	rrep := NewRREPMessage(destination, a.srcIP)
	log.Printf("Sending: %s\n", rrep.String())

	// broadcast the RREP
	go a.Send(rrep.Marshal(), destination)
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
		go a.SendRREP(rreq.OriginatorIP)
	} else {
		// increment hop count
		rreq.HopCount = rreq.HopCount + 1
		// forward the RREQ
		go a.sendToAllExcept(rreq.Marshal(), from)
	}
}

func (a *Aodv) handleRREP(payload []byte, from net.HardwareAddr) {
	rrep, err := UnmarshalRREP(payload)
	if err != nil {
		log.Printf("Failed to unmarshal RREP: %v\n", err)
		return
	}

	if rrep.HopCount > HopThreshold {
		// drop the packet
		return
	}

	log.Printf("Received: %s\n", rrep.String())
	// update the routing table
	go a.routingTable.Update(from, rrep.DestinationIP, rrep.HopCount)

	// check if the RREP is for me
	if rrep.OriginatorIP.Equal(a.srcIP) {
		// Arrived Successfully
		log.Printf("Path Descovery is successful for ip=%s !!!!", rrep.DestinationIP)
	} else {
		// increment hop count
		rrep.HopCount = rrep.HopCount + 1
		// forward the RREP
		go a.Send(rrep.Marshal(), rrep.OriginatorIP)
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
		case RREPType:
			go a.handleRREP(payload, addr)
		default:
			log.Println("Unknown message type: ", msgType)
		}
	}
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.neighborTable.Run()
	go a.Listen()
}

func (a *Aodv) Close() {
	a.channel.Close()
}