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
	seqNum uint32
	rreqID uint32
}

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
		seqNum: 0,
		rreqID: 0,
	}
}

func (a *Aodv) updateSeqNum(newSeqNum uint32) {
	if newSeqNum > a.seqNum {
		a.seqNum = newSeqNum
	}
}

func (a *Aodv) updateRoutingTable(destIP net.IP, nextHop net.HardwareAddr, hopCount uint8, lifeTime, seqNum uint32) {	
	entry := &VRoutingTableEntry{
		Destination: destIP,
		NextHop: nextHop,
		NoOfHops: hopCount,
		SeqNum: a.seqNum,
		LifeTime: time.Now().Add(time.Duration(lifeTime) * time.Millisecond),
	}

	a.routingTable.Update(entry)
}

func (a *Aodv) sendToAllExcept(payload []byte, addr net.HardwareAddr) {
	for item := range a.neighborTable.Iter() {
		neighborMac := item.MAC
		if neighborMac.String() != addr.String() {
			a.channel.SendTo(payload, neighborMac)
		}
	}
}

func (a *Aodv) Send(payload []byte, dest net.IP) {
	// check if the destination in the routing table
	item, ok := a.routingTable.Get(dest);
	if ok {
		// send the packet
		a.channel.SendTo(payload, item.NextHop)
	} else {
		// send a RREQ or RRER
		a.SendRREQ(dest)
	}
}

func (a *Aodv) Broadcast(payload []byte) {
	a.channel.Broadcast(payload)
}

func (a *Aodv) SendRREQ(destination net.IP) {
	rreq := NewRREQMessage(a.srcIP, destination)
	a.updateSeqNum(a.seqNum + 1)
	a.rreqID = a.rreqID + 1
	rreq.RREQID = a.rreqID
	rreq.OriginatorSequenceNumber = a.seqNum
	item, ok := a.routingTable.Get(destination)
	if ok {
		rreq.DestinationSeqNum = item.SeqNum
		rreq.ClearFlag(RREQFlagU)
	}	

	// broadcast the RREQ
	log.Printf("Sending: %s\n", rreq.String())
	a.Broadcast(rreq.Marshal())
}

func (a *Aodv) SendRREP(destination net.IP) {
	rrep := NewRREPMessage(destination, a.srcIP)
	rrep.DestinationSeqNum = a.seqNum
	
	// broadcast the RREP
	log.Printf("Sending: %s\n", rrep.String())
	a.Send(rrep.Marshal(), destination)
}

func (a *Aodv) handleRREQ(payload []byte, from net.HardwareAddr) {
	rreq, err := UnmarshalRREQ(payload)
	if err != nil {
		log.Printf("Failed to unmarshal RREQ: %v\n", err)
		return
	}
	
	if rreq.HopCount > HopCountLimit || rreq.OriginatorIP.Equal(a.srcIP) {
		// drop the packet
		log.Printf("Dropping %s\n", rreq.String())
		return
	}

	log.Printf("Received: %s\n", rreq.String())
	// update the routing table
	go a.updateRoutingTable(rreq.OriginatorIP, from, rreq.HopCount, ActiveRouteTimeMS, rreq.OriginatorSequenceNumber)

	// check if the RREQ is for me
	if rreq.DestinationIP.Equal(a.srcIP) {
		// update the sequence number if it is not unknown
		if !rreq.HasFlag(RREQFlagU) {
			a.updateSeqNum(rreq.DestinationSeqNum)
		}
		// send a RREP
		a.SendRREP(rreq.OriginatorIP)
	} else {
		// increment hop count
		rreq.HopCount = rreq.HopCount + 1
		// forward the RREQ
		a.sendToAllExcept(rreq.Marshal(), from)
	}
}

func (a *Aodv) handleRREP(payload []byte, from net.HardwareAddr) {
	rrep, err := UnmarshalRREP(payload)
	if err != nil {
		log.Printf("Failed to unmarshal RREP: %v\n", err)
		return
	}

	// check RREP hop count and life time
	if rrep.HopCount > HopCountLimit {
		// drop the packet
		log.Printf("Dropping %s\n", rrep.String())
		return
	}

	log.Printf("Received: %s\n", rrep.String())
	// update the routing table
	go a.updateRoutingTable(rrep.DestinationIP, from, rrep.HopCount, rrep.LifeTime, rrep.DestinationSeqNum)

	// check if the RREP is for me
	if rrep.OriginatorIP.Equal(a.srcIP) {
		// Arrived Successfully
		log.Printf("Path Descovery is successful for ip=%s !!!!", rrep.DestinationIP)
	} else {
		// increment hop count
		rrep.HopCount = rrep.HopCount + 1
		// forward the RREP
		a.Send(rrep.Marshal(), rrep.OriginatorIP)
	}
}

func (a *Aodv) handleMessage(payload []byte, from net.HardwareAddr) {
	msgType := uint8(payload[0])
	// handle the message
	switch msgType {
	case RREQType:
		a.handleRREQ(payload, from)
	case RREPType:
		a.handleRREP(payload, from)
	default:
		log.Println("Unknown message type: ", msgType)
	}
}

func (a *Aodv) Listen() {
	log.Println("Listening for AODV packets...")
	for {
		payload, addr, err := a.channel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		a.handleMessage(payload, addr)
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