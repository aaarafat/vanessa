package aodv

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)


func (a *Aodv) listen(channel *DataLinkLayerChannel) {
	log.Printf("Listening for AODV packets on channel: %d....\n", channel.IfiIndex)
	for {
		payload, addr, err := channel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
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
	go a.routingTable.Update(rreq.OriginatorIP, from, rreq.HopCount, ActiveRouteTimeMS, rreq.OriginatorSequenceNumber, IfiIndex)

	// update seq table
	go a.seqTable.Set(rreq.OriginatorIP, rreq.RREQID)

	// check if the RREQ is for me
	if rreq.DestinationIP.Equal(a.srcIP) || (rreq.DestinationIP.Equal(net.ParseIP(RsuIP)) && ConnectedToRSU(2)) {
		// update the sequence number if it is not unknown
		if !rreq.HasFlag(RREQFlagU) {
			a.updateSeqNum(rreq.DestinationSeqNum)
		}
		// send a RREP
		a.SendRREP(rreq.OriginatorIP, rreq.DestinationIP.Equal(net.ParseIP(RsuIP)))
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
	go a.routingTable.Update(rrep.DestinationIP, from, rrep.HopCount, rrep.LifeTime, rrep.DestinationSeqNum, IfiIndex)

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

	if data.DestenationIP.Equal(a.srcIP) {
		log.Printf("Received: %s\n", data.String())
		go a.callback(data.Data)
		return
	}
	// update seq table
	go a.dataSeqTable.Set(data.OriginatorIP, data.SeqNumber)

	// forward the data
	if data.DestenationIP.Equal(net.ParseIP(BroadcastIP)) {
		log.Printf("Received: %s\n", data.String())
		go a.callback(data.Data)
	
		a.forwarder.ForwardToAllExcept(data.Marshal(), from)
	} else {
		a.forwardData(data)
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
	case DataType:
		a.handleData(payload, from)
	default:
		log.Println("Unknown message type: ", msgType)
	}
}