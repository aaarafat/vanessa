package aodv

import (
	"log"
	"net"
	"time"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
	"github.com/cornelk/hashmap"
)

type Aodv struct {
	channel   *DataLinkLayerChannel
	forwarder IForwarder
	srcIP     net.IP

	// Sequence number
	seqNum uint32
	rreqID uint32

	// tables
	routingTable *VRoutingTable
	seqTable     *VFloodingSeqTable
	rreqBuffer   *hashmap.HashMap

	// path discovery callback
	pathDiscoveryCallback func(net.IP)
}

func NewAodv(srcIP net.IP, IfiIndex int, pathDiscoveryCallback func(net.IP)) *Aodv {
	channel, err := NewDataLinkLayerChannelWithInterfaceName(VAODVEtherType, "wlan0")
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel:      channel,
		forwarder:    NewForwarder(srcIP, channel),
		routingTable: NewVRoutingTable(),
		seqTable:     NewVFloodingSeqTable(),
		rreqBuffer:   &hashmap.HashMap{},
		srcIP:        srcIP,
		seqNum:       0,
		rreqID:       0,

		pathDiscoveryCallback: pathDiscoveryCallback,
	}
}

func (a *Aodv) GetRoute(destIP net.IP) (*VRoute, bool) {
	item, ok := a.routingTable.Get(destIP)
	if ok {
		return NewVRoute(item.Destination, item.NextHop, item.IfiIndex, int(item.NoOfHops)), true
	}

	return nil, false
}

func (a *Aodv) BuildRoute(destIP net.IP) (started bool) {
	if a.inRREQBuffer(destIP) {
		return false
	}

	log.Printf("Building route for IP: %s.....\n", destIP)
	a.SendRREQ(destIP)
	return true
}

func (a *Aodv) Send(payload []byte, dest net.IP) {
	item, ok := a.routingTable.Get(dest)
	if ok {
		a.forwarder.ForwardTo(payload, item.NextHop)
	} else {
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

	a.addToRREQBuffer(rreq)

	// broadcast the RREQ
	log.Printf("Sending: %s\n", rreq.String())
	a.forwarder.ForwardToAll(rreq.Marshal())
}

func (a *Aodv) SendRREPFor(rreq *RREQMessage) {
	rrep := NewRREPMessage(rreq.OriginatorIP, rreq.DestinationIP)
	if a.isRREQForMe(rreq) {
		if !rreq.HasFlag(RREQFlagU) {
			a.updateSeqNum(rreq.DestinationSeqNum)
		}
		rrep.DestinationSeqNum = a.seqNum
	} else {
		route, _ := a.routingTable.Get(rreq.DestinationIP)
		rrep.DestinationSeqNum = route.SeqNum
		rrep.HopCount = route.NoOfHops + 1
		rrep.LifeTime = uint32(route.LifeTime.Sub(time.Now()).Milliseconds())
		// TODO: Send Gratious RREP to the RREQ originator if the RREQ has the G flag set
	}
	// broadcast the RREP
	log.Printf("Sending: %s\n", rrep.String())
	a.Send(rrep.Marshal(), rrep.OriginatorIP)
}

func (a *Aodv) updateRSU() {
	for {
		if ConnectedToRSU(2) {
			mac, err := net.ParseMAC(GetRSUMac(2))
			if err != nil {
				log.Printf("failed to parse MAC: %v", err)
				continue
			}
			new := a.routingTable.Set(net.ParseIP(ip.RsuIP), mac, 0, RSUActiveRouteTimeMS, 0, 2)
			if new {
				log.Printf("Path Discovery is successful for ip=%s !!!!", ip.RsuIP)
				go a.pathDiscoveryCallback(net.ParseIP(ip.RsuIP))
			}
		}
		time.Sleep(time.Millisecond * time.Duration(RSUActiveRouteTimeMS/3))
	}
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.forwarder.Start()
	go a.listen(a.channel)
	go a.updateRSU()
}

func (a *Aodv) Close() {
	a.channel.Close()
	a.forwarder.Close()
}
