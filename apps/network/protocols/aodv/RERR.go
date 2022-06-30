package aodv

import (
	"encoding/binary"
	"fmt"
	"net"
)

// https://datatracker.ietf.org/doc/html/rfc3561#section-5.3
type RERRMessage struct {
	// Type
	Type uint8
	// Flags
	Flags uint16
	// The number of unreachable destinations included in the	message; MUST be at least 1.
	DestCount uint8
	// The IP address of the node that is the destination of the RREQ.
	UnreachableDestinationIP net.IP
	// Destination sequence number
	UnreachableDestinationSeqNum uint32 
}

func NewRERRMessage(SrcIP, DestIP net.IP) *RERRMessage {
	return &RERRMessage{
		Type: RERRType,
		Flags: 0,
		DestCount: 0,
		UnreachableDestinationIP: DestIP,
		UnreachableDestinationSeqNum: 0,
	}
}

func (RERR *RERRMessage) Marshal() []byte {
	bytes := make([]byte, RERRMessageLen)

	bytes[0] = byte(RERR.Type)
	binary.LittleEndian.PutUint16(bytes[1:], RERR.Flags)
	bytes[3] = byte(RERR.DestCount)
	copy(bytes[4:8], RERR.UnreachableDestinationIP.To4())
	binary.LittleEndian.PutUint32(bytes[8:], RERR.UnreachableDestinationSeqNum)

	return bytes
}

func UnmarshalRERR(data []byte) (*RERRMessage, error) {
	if len(data) < RERRMessageLen {
		return nil, fmt.Errorf("RERR message length is %d, expected %d", len(data), RERRMessageLen)
	}

	RERR := &RERRMessage{}
	RERR.Type = uint8(data[0])
	RERR.Flags = binary.LittleEndian.Uint16(data[1:3])
	RERR.DestCount = data[3]
	RERR.UnreachableDestinationIP = net.IPv4(data[4], data[5], data[6], data[7])
	RERR.UnreachableDestinationSeqNum = binary.LittleEndian.Uint32(data[8:12])

	return RERR, nil
}


func (RERR *RERRMessage) String() string {
	return fmt.Sprintf("RERR: Type=%d, Flags=%d, DestCount=%d, DestinationSeqNum=%d, DestinationIP=%s",
		RERR.Type, RERR.Flags, RERR.DestCount, RERR.UnreachableDestinationSeqNum, RERR.UnreachableDestinationIP.String())
}