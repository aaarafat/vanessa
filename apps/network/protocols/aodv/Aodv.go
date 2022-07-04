package aodv

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/cornelk/hashmap"
)

type Aodv struct {
	channels []*DataLinkLayerChannel
	forwarder *Forwarder
	routingTable *VRoutingTable
	seqTable *VFloodingSeqTable
	dataSeqTable *VFloodingSeqTable
	dataBuffer *hashmap.HashMap
	srcIP net.IP
	seqNum uint32
	dataSeqNum uint32
	rreqID uint32

	callback func(data[]byte)
}

func createChannels() []*DataLinkLayerChannel {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return nil
	}
	channels := make([]*DataLinkLayerChannel, len(interfaces))
	for i, ifi := range interfaces {
		if ifi.Name == "lo" {
			continue
		}
		log.Printf("Creating channel in AODV for interface: %s Index: %d\n", ifi.Name, i)
		ch, err := NewDataLinkLayerChannelWithInterface(VAODVEtherType, i)
		if err != nil {
			log.Fatalf("failed to create channel: %v", err)
		}
		channels[i] = ch
	}

	return channels
}

func NewAodv(srcIP net.IP, callback func(data[]byte)) *Aodv {
	channels := createChannels()

	return &Aodv{
		channels: channels,
		forwarder: NewForwarder(srcIP, channels),
		routingTable: NewVRoutingTable(),
		seqTable: NewVFloodingSeqTable(),
		dataSeqTable: NewVFloodingSeqTable(),
		dataBuffer: &hashmap.HashMap{},
		srcIP: srcIP,
		seqNum: 0,
		dataSeqNum: 0,
		rreqID: 0,
		callback: callback,
	}
}

func (a *Aodv) updateSeqNum(newSeqNum uint32) {
	if newSeqNum > a.seqNum {
		a.seqNum = newSeqNum
	}
}

func (a *Aodv) Send(payload []byte, dest net.IP) {
	// check if the destination in the routing table
	item, ok := a.routingTable.Get(dest);
	if ok {
		// forward the packet
		a.forwarder.ForwardTo(payload, item.NextHop, item.IfiIndex)
	} else {
		// send a RREQ or RRER
		a.SendRREQ(dest)
	}
}

func (a *Aodv) forwardData(data *DataMessage) {
	// handle rsu connection
	if data.DestenationIP.Equal(net.ParseIP(RsuIP)) && connectedToRSU(2) {
		mac, err := net.ParseMAC(getRSUMac(2))
		if err != nil {
			log.Fatalf("failed to parse MAC: %v", err)
		}
		a.routingTable.Update(data.DestenationIP, mac, data.HopCount, ActiveRouteTimeMS, a.seqNum, 2)
	}
	// check if the destination in the routing table
	item, ok := a.routingTable.Get(data.DestenationIP);
	if ok {
		// forward the packet
		a.forwarder.ForwardTo(data.Marshal(), item.NextHop, item.IfiIndex)
	} else {
		// send a RREQ or RRER
		buf, ok := a.dataBuffer.Get(data.DestenationIP.String())
		if ok {
			a.dataBuffer.Set(data.DestenationIP.String(), append(buf.([]DataMessage), *data))
		} else {
			a.dataBuffer.Set(data.DestenationIP.String(), []DataMessage{*data})
		}
		a.SendRREQ(data.DestenationIP)
	}
}

func (a *Aodv) SendData(payload []byte, dest net.IP) {
	// update the sequence number
	a.dataSeqNum = a.dataSeqNum + 1
	// create the data packet
	data := NewDataMessage(a.srcIP, a.dataSeqNum, payload)
	data.DestenationIP = dest
	// broadcast the data packet
	log.Printf("Sending: %s\n", data.String())
	
	if data.DestenationIP.Equal(net.ParseIP(BroadcastIP)) {
		if connectedToRSU(2) {
			mac, err := net.ParseMAC(getRSUMac(2))
			if err != nil {
				log.Fatalf("failed to parse MAC: %v", err)
			}
			a.forwarder.ForwardTo(payload, mac, 2)
		}
		a.forwarder.ForwardToAll(data.Marshal())
	} else {
		a.forwardData(data)
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
		mac, err := net.ParseMAC(getRSUMac(2))
		if err != nil {
			log.Fatalf("failed to parse MAC: %v", err)
		}
		go a.routingTable.Update(rrep.DestinationIP, mac, rrep.HopCount, rrep.LifeTime, rrep.DestinationSeqNum, 2)
	}
	
	// broadcast the RREP
	log.Printf("Sending: %s\n", rrep.String())
	a.Send(rrep.Marshal(), destination)
}


func (a *Aodv) Listen(channel *DataLinkLayerChannel) {
	log.Println("Listening for AODV packets...")
	for {
		payload, addr, err := channel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		go a.handleMessage(payload, addr, channel.IfiIndex)
	}
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.forwarder.Start()
	for _, channel := range a.channels {
		if channel == nil {
			continue
		}
		go a.Listen(channel)
	}
}

func (a *Aodv) Close() {
	for _, channel := range a.channels {
		if channel == nil {
			continue
		}
		channel.Close()
	}
	a.forwarder.Close()
}