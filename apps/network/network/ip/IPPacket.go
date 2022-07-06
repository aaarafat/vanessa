package ip

import "net"

type IPPacket struct {
	Header *IPHeader
	Payload   []byte
}

const (
	DefaultIP4HeaderLen = 20
	DefaultTTL = 64
	DefaultProtocol = 144 //! Not Reserverd
)

func NewIPPacket(payload []byte, srcIP, destIp net.IP) *IPPacket {
	header := &IPHeader{
		Version: 4,
		Length: DefaultIP4HeaderLen / 4,
		TypeOfService: 0,
		TotalLength: DefaultIP4HeaderLen + uint16(len(payload)),
		IdentifierFlagsOffset: 0,
		TTL: DefaultTTL,
		Protocol: DefaultProtocol, 
		HeaderChecksum: 0,
		SrcIP: srcIP,
		DestIP: destIp,
	}
	
	return &IPPacket{
		Header: header,
		Payload: payload,
	}
}

func UnmarshalPacket(data []byte) (*IPPacket, error) {
	header, err := UnmarshalIPHeader(data)
	if err != nil {
		return nil, err
	}
	
	return &IPPacket{
		Header: header,
		Payload: data[header.LengthInBytes():],
	}, nil
}

func MarshalIPPacket(packet *IPPacket) []byte {
	data := make([]byte, packet.Header.TotalLength)
	copy(data, MarshalIPHeader(packet.Header))
	copy(data[packet.Header.LengthInBytes():], packet.Payload)
	return data
}