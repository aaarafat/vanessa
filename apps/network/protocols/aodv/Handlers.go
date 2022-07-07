package aodv

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
)


func (a *Aodv) listen(channel *DataLinkLayerChannel) {
	log.Printf("Listening for AODV packets on channel: %d....\n", channel.IfiIndex)
	for {
		payload, addr, err := channel.Read()
		if err != nil {
			return
		}
		go a.handleMessage(payload, addr, channel.IfiIndex)
	}
}


func (a *Aodv) handleRREQ(payload []byte, from net.HardwareAddr, IfiIndex int) {
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

	log.Printf("Interface %d: Received: %s\n", IfiIndex, rreq.String())
	// update the routing table
	a.routingTable.Update(rreq.OriginatorIP, from, rreq.HopCount, ActiveRouteTimeMS, rreq.OriginatorSequenceNumber, IfiIndex)

	// update seq table
	a.seqTable.Set(rreq.OriginatorIP, rreq.RREQID)

	// check if the RREQ is for me
	if rreq.DestinationIP.Equal(a.srcIP) || (rreq.DestinationIP.Equal(net.ParseIP(ip.RsuIP)) && ConnectedToRSU(2)) {
		// update the sequence number if it is not unknown
		if !rreq.HasFlag(RREQFlagU) {
			a.updateSeqNum(rreq.DestinationSeqNum)
		}
		// send a RREP
		a.SendRREP(rreq.OriginatorIP, rreq.DestinationIP.Equal(net.ParseIP(ip.RsuIP)))
	} else {
		// increment hop count
		rreq.HopCount = rreq.HopCount + 1
		// forward the RREQ
		a.forwarder.ForwardToAllExcept(rreq.Marshal(), from)
	}
}

func (a *Aodv) handleRREP(payload []byte, from net.HardwareAddr, IfiIndex int) {
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

	log.Printf("Inteface %d: Received: %s\n", IfiIndex, rrep.String())
	// update the routing table
	a.routingTable.Update(rrep.DestinationIP, from, rrep.HopCount, rrep.LifeTime, rrep.DestinationSeqNum, IfiIndex)

	// check if the RREP is for me
	if rrep.OriginatorIP.Equal(a.srcIP) {
		// Arrived Successfully
		log.Printf("Path Descovery is successful for ip=%s !!!!", rrep.DestinationIP)
		a.pathDiscoveryCallback(rrep.DestinationIP)
	} else {
		// increment hop count
		rrep.HopCount = rrep.HopCount + 1
		// forward the RREP
		a.Send(rrep.Marshal(), rrep.OriginatorIP)
	}
}

func (a *Aodv) handleMessage(payload []byte, from net.HardwareAddr, IfiIndex int) {
	msgType := uint8(payload[0])
	// handle the message
	switch msgType {
	case RREQType:
		a.handleRREQ(payload, from, IfiIndex)
	case RREPType:
		a.handleRREP(payload, from, IfiIndex)
	default:
		log.Println("Unknown message type: ", msgType)
	}
}