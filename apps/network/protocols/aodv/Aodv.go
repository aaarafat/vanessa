package aodv

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
)

type Aodv struct {
	channel *DataLinkLayerChannel
	forwarder *Forwarder
	routingTable *VRoutingTable
	seqTable *VFloodingSeqTable
	dataSeqTable *VFloodingSeqTable
	srcIP net.IP
	seqNum uint32
	dataSeqNum uint32
	rreqID uint32

	// path discovery callback
	pathDiscoveryCallback func(net.IP)
}


func NewAodv(srcIP net.IP, pathDiscoveryCallback func(net.IP)) *Aodv {
	channel, err := CreateChannel(VAODVEtherType, 1)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel: channel,
		forwarder: NewForwarder(srcIP, channel),
		routingTable: NewVRoutingTable(),
		seqTable: NewVFloodingSeqTable(),
		dataSeqTable: NewVFloodingSeqTable(),
		srcIP: srcIP,
		seqNum: 0,
		dataSeqNum: 0,
		rreqID: 0,

		pathDiscoveryCallback: pathDiscoveryCallback,
	}
}

func (a *Aodv) GetRoute(destIP net.IP) (*VRoute, bool) {
	item, ok := a.routingTable.Get(destIP)
	if ok {
		return NewVRoute(item.Destination, item.NextHop, item.IfiIndex, int(item.NoOfHops)), true
	} 

	if destIP.Equal(net.ParseIP(RsuIP)) && ConnectedToRSU(2) {
		mac, err := net.ParseMAC(GetRSUMac(2))
		if err != nil {
			log.Fatalf("failed to parse mac: %v", err)
		}
		return NewVRoute(destIP, mac, 2, 0), true
	}

	return nil, false
}

func (a *Aodv) BuildRoute(destIP net.IP) {
	// send a RREQ
	a.SendRREQ(destIP)
}

func (a *Aodv) Send(payload []byte, dest net.IP) {
	// check if the destination in the routing table
	item, ok := a.routingTable.Get(dest);
	if ok {
		// forward the packet
		a.forwarder.ForwardTo(payload, item.NextHop)
	} else {
		// send a RREQ or RRER
		a.SendRREQ(dest)
	}
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
	a.forwarder.ForwardToAll(rreq.Marshal())
}

func (a *Aodv) SendRREP(destination net.IP, forRSU bool) {
	rrep := NewRREPMessage(destination, a.srcIP)
	rrep.DestinationSeqNum = a.seqNum
	if forRSU {
		rrep.DestinationIP = net.ParseIP(RsuIP)
		rrep.HopCount = 1
		mac, err := net.ParseMAC(GetRSUMac(2))
		if err != nil {
			log.Fatalf("failed to parse MAC: %v", err)
		}
		go a.routingTable.Update(rrep.DestinationIP, mac, rrep.HopCount, rrep.LifeTime, rrep.DestinationSeqNum, 2)
	}
	
	// broadcast the RREP
	log.Printf("Sending: %s\n", rrep.String())
	a.Send(rrep.Marshal(), destination)
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.forwarder.Start()
	go a.listen(a.channel)
}

func (a *Aodv) Close() {
	a.channel.Close()
	a.forwarder.Close()
}