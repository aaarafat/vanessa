package aodv

import (
	"encoding/binary"
	"fmt"
	"net"
)

type RERRUnreachableDestination struct {
	IP     net.IP
	SeqNum uint32
}

// https://datatracker.ietf.org/doc/html/rfc3561#section-5.3
type RERRMessage struct {
	// Type
	Type uint8
	// Flags
	Flags uint16
	// The number of unreachable destinations included in the	message; MUST be at least 1.
	DestCount uint8
	// Unreachable destinations and their sequence numbers.
	UnreachableDestinations []RERRUnreachableDestination
}

func NewRERRMessage(unreachableDestinations []RERRUnreachableDestination) *RERRMessage {
	return &RERRMessage{
		Type:                    RERRType,
		Flags:                   0,
		DestCount:               uint8(len(unreachableDestinations)),
		UnreachableDestinations: unreachableDestinations,
	}
}

func (d *RERRUnreachableDestination) Marshal() []byte {
	b := make([]byte, 8)
	copy(b[0:4], d.IP.To4())
	binary.LittleEndian.PutUint32(b, uint32(d.SeqNum))
	return b
}

func UnmarshalRERRUnreachableDestination(b []byte) *RERRUnreachableDestination {
	d := RERRUnreachableDestination{}
	d.IP = net.IPv4(b[0], b[1], b[2], b[3])
	d.SeqNum = binary.LittleEndian.Uint32(b[4:8])
	return &d
}

func (RERR *RERRMessage) Marshal() []byte {
	bytes := make([]byte, RERRMessageLen+RERR.DestCount*8)

	bytes[0] = byte(RERR.Type)
	binary.LittleEndian.PutUint16(bytes[1:], RERR.Flags)
	bytes[3] = byte(RERR.DestCount)
	for i := 0; i < int(RERR.DestCount); i++ {
		copy(bytes[4+i*8:], RERR.UnreachableDestinations[i].Marshal())
	}

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
	for i := 0; i < int(RERR.DestCount); i++ {
		RERR.UnreachableDestinations = append(RERR.UnreachableDestinations, *UnmarshalRERRUnreachableDestination(data[4+i*8:]))
	}

	return RERR, nil
}

func (d *RERRUnreachableDestination) String() string {
	return fmt.Sprintf("%s:%d", d.IP, d.SeqNum)
}

func (RERR *RERRMessage) String() string {
	str := fmt.Sprintf("RERR: Type=%d, Flags=%d, DestCount=%d, ", RERR.Type, RERR.Flags, RERR.DestCount)

	for i := 0; i < int(RERR.DestCount); i++ {
		str += fmt.Sprintf("%s ", RERR.UnreachableDestinations[i].String())
	}

	return str
}
