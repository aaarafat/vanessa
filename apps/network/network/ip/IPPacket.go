package ip

import (
	"fmt"
	"net"
)

type IPPacket struct {
	Header  *IPHeader
	Payload []byte
}

const (
	DefaultIP4HeaderLen = 20
	DefaultTTL          = 20
	DefaultProtocol     = 144 //! Not Reserverd
)

func NewIPPacket(payload []byte, srcIP, destIp net.IP) *IPPacket {
	header := &IPHeader{
		Version:               4,
		Length:                DefaultIP4HeaderLen / 4,
		TypeOfService:         0,
		TotalLength:           DefaultIP4HeaderLen + uint16(len(payload)),
		IdentifierFlagsOffset: 0,
		TTL:                   DefaultTTL,
		Protocol:              DefaultProtocol,
		HeaderChecksum:        0,
		SrcIP:                 srcIP,
		DestIP:                destIp,
	}

	return &IPPacket{
		Header:  header,
		Payload: payload,
	}
}

func NewIPPacketWithOptions(payload []byte, srcIP, destIp net.IP, options []byte) *IPPacket {
	optionLen := len(options) + (4-len(options)%4)%4 // to make it multiple of 32 bits
	header := &IPHeader{
		Version:               4,
		Length:                (DefaultIP4HeaderLen + uint8(optionLen)) / 4,
		TypeOfService:         0,
		TotalLength:           (DefaultIP4HeaderLen + uint16(optionLen)) + uint16(len(payload)),
		IdentifierFlagsOffset: 0,
		TTL:                   DefaultTTL,
		Protocol:              DefaultProtocol,
		HeaderChecksum:        0,
		SrcIP:                 srcIP,
		DestIP:                destIp,
		Options:               options,
	}

	return &IPPacket{
		Header:  header,
		Payload: payload,
	}
}

func (p *IPPacket) HasOptions() bool {
	return p.Header.LengthInBytes() > DefaultIP4HeaderLen
}

func UnmarshalPacket(data []byte) (*IPPacket, error) {
	header, err := UnmarshalIPHeader(data)
	if err != nil {
		return nil, err
	}

	return &IPPacket{
		Header:  header,
		Payload: data[header.LengthInBytes():],
	}, nil
}

func MarshalIPPacket(packet *IPPacket) []byte {
	data := make([]byte, packet.Header.LengthInBytes()+len(packet.Payload))
	copy(data, MarshalIPHeader(packet.Header))
	copy(data[packet.Header.LengthInBytes():], packet.Payload)
	return data
}

func (packet *IPPacket) String() string {
	return fmt.Sprintf("%s   %s", packet.Payload, packet.Header.String())
}
