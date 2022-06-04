package aodv

import (
	"encoding/binary"
	"fmt"
	"net"
)

// https://datatracker.ietf.org/doc/html/rfc3561#section-5.2
type RREPMessage struct {
	// Type
	Type uint8
	// Flags
	Flags uint16
	// The number of hops that the RREP has already made.
	HopCount uint8
	// The IP address of the node that is the destination of the RREQ.
	DestinationIP net.IP
	// Destination sequence number
	DestinationSeqNum uint32 
	// The IP address of the node that originated the RREQ.
	OriginatorIP net.IP
	// The time in milliseconds for which nodes receiving the RREP consider the route to be valid.
	LifeTime uint32
}

func NewRREPMessage(SrcIP, DestIP net.IP) *RREPMessage {
	return &RREPMessage{
		Type: RREPType,
		Flags: 0,
		HopCount: 0,
		DestinationIP: DestIP,
		OriginatorIP: SrcIP,
		DestinationSeqNum: 0,
		LifeTime: RREPDefaultLifeTimeMS,
	}
}

func (rrep *RREPMessage) Marshal() []byte {
	bytes := make([]byte, RREPMessageLen)

	bytes[0] = byte(rrep.Type)
	binary.LittleEndian.PutUint16(bytes[1:], rrep.Flags)
	bytes[3] = byte(rrep.HopCount)
	binary.LittleEndian.PutUint32(bytes[4:], rrep.DestinationSeqNum)
	copy(bytes[8:12], rrep.DestinationIP.To4())
	copy(bytes[12:16], rrep.OriginatorIP.To4())
	binary.LittleEndian.PutUint32(bytes[16:], rrep.LifeTime)

	return bytes
}

func UnmarshalRREP(data []byte) (*RREPMessage, error) {
	if len(data) != RREPMessageLen {
		return nil, fmt.Errorf("RREP message length is %d, expected %d", len(data), RREPMessageLen)
	}

	rrep := &RREPMessage{}
	rrep.Type = uint8(data[0])
	rrep.Flags = binary.LittleEndian.Uint16(data[1:3])
	rrep.HopCount = data[3]
	rrep.DestinationSeqNum = binary.LittleEndian.Uint32(data[4:8])
	rrep.DestinationIP = net.IPv4(data[8], data[9], data[10], data[11])
	rrep.OriginatorIP = net.IPv4(data[12], data[13], data[14], data[15])
	rrep.LifeTime = binary.LittleEndian.Uint32(data[16:20])

	return rrep, nil
}

func (rrep *RREPMessage) Invalid() bool {
	return rrep.Type != RREPType ||
				rrep.HopCount > HopCountLimit
}

func (rrep *RREPMessage) String() string {
	return fmt.Sprintf("RREP: Type=%d, Flags=%d, HopCount=%d, DestinationSeqNum=%d, DestinationIP=%s, OriginatorIP=%s, LifeTime=%d",
		rrep.Type, rrep.Flags, rrep.HopCount, rrep.DestinationSeqNum, rrep.DestinationIP.String(), rrep.OriginatorIP.String(), rrep.LifeTime)
}