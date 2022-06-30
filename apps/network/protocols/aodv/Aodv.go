package aodv

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/cornelk/hashmap"
)

type Aodv struct {
	channel *DataLinkLayerChannel
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

func NewAodvWithInterface(srcIP net.IP, callback func(data[]byte),ifiname string) *Aodv {
	d, err := NewDataLinkLayerChannelWithInterface(VAODVEtherType, ifiname)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel: d,
		forwarder: NewForwarder(srcIP, d),
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

func NewAodv(srcIP net.IP, callback func(data[]byte)) *Aodv {
	d, err := NewDataLinkLayerChannel(VAODVEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &Aodv{
		channel: d,
		forwarder: NewForwarder(srcIP, d),
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
		a.forwarder.ForwardTo(payload, item.NextHop)
	} else {
		// send a RREQ or RRER
		a.SendRREQ(dest)
	}
}

func (a *Aodv) forwardData(data *DataMessage) {
	// check if the destination in the routing table
	item, ok := a.routingTable.Get(data.DestenationIP);
	if ok {
		// forward the packet
		a.forwarder.ForwardTo(data.Marshal(), item.NextHop)
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
	
	if rreq.Invalid(a.seqTable, a.srcIP) {
		// drop the packet
		log.Printf("Dropping %s\n", rreq.String())
		return
	}

	log.Printf("Received: %s\n", rreq.String())
	// update the routing table
	go a.routingTable.Update(rreq.OriginatorIP, from, rreq.HopCount, ActiveRouteTimeMS, rreq.OriginatorSequenceNumber)
	// update seq table
	go a.seqTable.Set(rreq.OriginatorIP, rreq.RREQID)

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
		a.forwarder.ForwardToAllExcept(rreq.Marshal(), from)
	}
}

func (a *Aodv) handleRREP(payload []byte, from net.HardwareAddr) {
	rrep, err := UnmarshalRREP(payload)
	if err != nil {
		log.Printf("Failed to unmarshal RREP: %v\n", err)
		return
	}

	// check RREP hop count
	if rrep.Invalid() {
		// drop the packet
		log.Printf("Dropping %s\n", rrep.String())
		return
	}

	log.Printf("Received: %s\n", rrep.String())
	// update the routing table
	go a.routingTable.Update(rrep.DestinationIP, from, rrep.HopCount, rrep.LifeTime, rrep.DestinationSeqNum)

	// check if the RREP is for me
	if rrep.OriginatorIP.Equal(a.srcIP) {
		// Arrived Successfully
		log.Printf("Path Descovery is successful for ip=%s !!!!", rrep.DestinationIP)
		// handle data in the buffer
		data, ok := a.dataBuffer.Get(rrep.DestinationIP.String())
		if ok {
			// send the data
			msgs := data.([]DataMessage)
			for _, msg := range msgs {
				go a.SendData(msg.Marshal(), msg.DestenationIP)
			}
			// remove the data from the buffer
			a.dataBuffer.Del(rrep.DestinationIP.String())
		}
	} else {
		// increment hop count
		rrep.HopCount = rrep.HopCount + 1
		// forward the RREP
		a.Send(rrep.Marshal(), rrep.OriginatorIP)
	}
}

func (a *Aodv) handleData(payload []byte, from net.HardwareAddr) {
	data, err := UnmarshalData(payload)
	if err != nil {
		log.Printf("Failed to unmarshal data: %v\n", err)
		return
	}

	if data.OriginatorIP.Equal(a.srcIP) || a.dataSeqTable.Exists(data.OriginatorIP, data.SeqNumber) {
		// drop the packet
		log.Printf("Dropping %s\n", data.String())
		return
	} 

	if data.DestenationIP.Equal(a.srcIP) || data.DestenationIP.Equal(net.ParseIP(BroadcastIP)) {
		log.Printf("Received: %s\n", data.String())
		go a.callback(data.Data)
	}
	// update seq table
	go a.dataSeqTable.Set(data.OriginatorIP, data.SeqNumber)

	// forward the data
	if data.DestenationIP.Equal(net.ParseIP(BroadcastIP)) {
		a.forwarder.ForwardToAllExcept(data.Marshal(), from)
	} else {
		a.forwardData(data)
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
	case DataType:
		a.handleData(payload, from)
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
		go a.handleMessage(payload, addr)
	}
}

func (a *Aodv) Start() {
	log.Printf("Starting AODV for IP: %s.....\n", a.srcIP)
	go a.forwarder.Start()
	go a.Listen()
}

func (a *Aodv) Close() {
	a.channel.Close()
}