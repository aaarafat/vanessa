package aodv

import (
	"fmt"
	"net"

	"encoding/binary"

	. "github.com/aaarafat/vanessa/apps/network/tables"
)

// https://datatracker.ietf.org/doc/html/rfc3561#section-5.1
type RREQMessage struct {
	// Type
	Type uint8  
	// Flags J|R|G|D|U|   Reserved
	Flags uint16 
	// The number of hops that the RREQ has already made.
	HopCount uint8 
	// RREQ ID
	RREQID uint32 
	// The IP address of the node that is the destination of the RREQ.
	DestinationIP net.IP
	// Destination sequence number
	DestinationSeqNum uint32 
	// The IP address of the node that originated the RREQ.
	OriginatorIP net.IP
	// The sequence number of the RREQ.
	OriginatorSequenceNumber uint32 
}


func NewRREQMessage(SrcIP, DestIP net.IP) *RREQMessage {
	rreq := &RREQMessage{
		Type: RREQType,
		Flags: 0,
		HopCount: 0,
		RREQID: 0,
		DestinationIP: DestIP,
		OriginatorIP: SrcIP,
		DestinationSeqNum: 0,
		OriginatorSequenceNumber: 0,
	}

	// Destination sequence number is not defined
	rreq.SetFlag(RREQFlagU)

	return rreq
}

func (rreq* RREQMessage) Marshal() []byte {
	bytes := make([]byte, RREQMessageLen)

	bytes[0] = byte(rreq.Type)
	binary.LittleEndian.PutUint16(bytes[1:], rreq.Flags)
	bytes[3] = byte(rreq.HopCount)
	binary.LittleEndian.PutUint32(bytes[4:], rreq.RREQID)
	copy(bytes[8:12], rreq.DestinationIP.To4())
	binary.LittleEndian.PutUint32(bytes[12:], rreq.DestinationSeqNum)
	copy(bytes[16:20], rreq.OriginatorIP.To4())
	binary.LittleEndian.PutUint32(bytes[20:], rreq.OriginatorSequenceNumber)

	return bytes
}

func UnmarshalRREQ(data []byte) (*RREQMessage, error) {
	if len(data) != RREQMessageLen {
		return nil, fmt.Errorf("RREQ message length is %d, expected %d", len(data), RREQMessageLen)
	}
	
	rreq := &RREQMessage{}
	rreq.Type = uint8(data[0])
	rreq.Flags = binary.LittleEndian.Uint16(data[1:3])
	rreq.HopCount = data[3]
	rreq.RREQID = binary.LittleEndian.Uint32(data[4:8])
	rreq.DestinationIP = net.IPv4(data[8], data[9], data[10], data[11])
	rreq.DestinationSeqNum = binary.LittleEndian.Uint32(data[12:16])
	rreq.OriginatorIP = net.IPv4(data[16], data[17], data[18], data[19])
	rreq.OriginatorSequenceNumber = binary.LittleEndian.Uint32(data[20:24])

	return rreq, nil
}

func (rreq* RREQMessage) SetFlag(flag uint16) *RREQMessage {
	rreq.Flags |= flag
	return rreq
}

func (rreq* RREQMessage) ClearFlag(flag uint16) *RREQMessage {
	rreq.Flags &= ^flag
	return rreq
}

func (rreq* RREQMessage) ToggleFlag(flag uint16) *RREQMessage {
	rreq.Flags ^= flag
	return rreq
}

func (rreq* RREQMessage) HasFlag(flag uint16) bool {
	return rreq.Flags & flag != 0
}


func (rreq *RREQMessage) Invalid(seqTable *VFloodingSeqTable, srcIP net.IP) bool {
	return rreq.Type != RREQType ||
				rreq.HopCount > HopCountLimit || 
				rreq.OriginatorIP.Equal(srcIP) || 
				seqTable.Exists(rreq.OriginatorIP, rreq.RREQID) 
}

func (rreq* RREQMessage) String() string {
	return fmt.Sprintf("RREQ: Type: %d, Flags: %d, HopCount: %d, RREQID: %d, DestinationIP: %s, DestinationSeqNum: %d, OriginatorIP: %s, OriginatorSequenceNumber: %d",
		rreq.Type, rreq.Flags, rreq.HopCount, rreq.RREQID, rreq.DestinationIP.String(), rreq.DestinationSeqNum, rreq.OriginatorIP.String(), rreq.OriginatorSequenceNumber)
}
