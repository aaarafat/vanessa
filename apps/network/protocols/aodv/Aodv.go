package aodv

import (
	"log"
	"net"
	"time"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
	"github.com/cornelk/hashmap"
)

type Aodv struct {
	channel *DataLinkLayerChannel
	flooder IFlooder
	srcIP   net.IP

	// Sequence number
	seqNum uint32
	rreqID uint32

	// tables
	routingTable   *VRoutingTable
	seqTable       *VFloodingSeqTable
	rreqBuffer     *hashmap.HashMap
	neighborTables map[int]*VNeighborTable

	// path discovery callback
	pathDiscoveryCallback func(net.IP)
}

const (
	WLAN0 = 0
	WLAN1 = 1
)

func NewAodv(srcIP net.IP, ifiName string, neighborTables map[int]*VNeighborTable, pathDiscoveryCallback func(net.IP)) *Aodv {
	channel, err := NewDataLinkLayerChannelWithInterfaceName(VAODVEtherType, ifiName)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel:        channel,
		flooder:        NewFlooder(srcIP, channel, neighborTables[WLAN0]),
		routingTable:   NewVRoutingTable(),
		seqTable:       NewVFloodingSeqTable(),
		rreqBuffer:     &hashmap.HashMap{},
		neighborTables: neighborTables,
		srcIP:          srcIP,
		seqNum:         0,
		rreqID:         0,

		pathDiscoveryCallback: pathDiscoveryCallback,
	}
}

func (a *Aodv) GetRoute(destIP net.IP) (*VRoute, bool) {
	item, ok := a.routingTable.Get(destIP)
	if ok {
		// check if hop count is 0
		if item.NoOfHops == 0 {
			// check in neighbor table
			if _, ok := a.neighborTables[item.IfiIndex].Get(item.NextHop.String()); !ok {
				// send RERR
				go a.SendRERR(item.NextHop)
				return nil, false
			}
		}

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
		a.flooder.ForwardTo(payload, item.NextHop)
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
	a.flooder.ForwardToAll(rreq.Marshal())
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

func (a *Aodv) SendRERR(nextHop net.HardwareAddr) {
	// get all unreachable destinations
	unreachable := make([]RERRUnreachableDestination, 0)
	for _, item := range a.routingTable.Items() {
		if item.NextHop.String() == nextHop.String() {
			unreachable = append(unreachable, RERRUnreachableDestination{
				IP:     item.Destination,
				SeqNum: item.SeqNum,
			})
		}
	}

	if len(unreachable) == 0 {
		return
	}

	// delete from routing table
	for _, dest := range unreachable {
		a.routingTable.Del(dest.IP)
		// send RREQ to local repair
		a.SendRREQ(dest.IP)
	}

	rerr := NewRERRMessage(unreachable)
	rerr.SetFlag(RERRFlagN) // set N flag because we performed a local repair
	log.Printf("Sending: %s\n", rerr.String())
	a.flooder.ForwardToAll(rerr.Marshal())
}

func (a *Aodv) updateRSU() {
	for {
		entry, ok := a.neighborTables[WLAN1].GetFirst()
		if ok {
			new := a.routingTable.Set(entry.IP, entry.MAC, 0, VNeighborTable_UPDATE_INTERVAL_MS, 0, WLAN1)
			if new {
				log.Printf("Path Discovery is successful for ip=%s !!!!", entry.IP)
				go a.pathDiscoveryCallback(entry.IP)
			}
		}
		time.Sleep((VNeighborTable_UPDATE_INTERVAL_MS / 3) * time.Millisecond)
	}
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.listen(a.channel)
	go a.updateRSU()
}

func (a *Aodv) Close() {
	a.channel.Close()
}
