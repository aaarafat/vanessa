package ip

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	IPv4HeaderLen  = 20
)

// https://www.techtarget.com/searchnetworking/tutorial/Routing-First-Step-IP-header-format
type IPHeader struct {
	Version uint8
	Length uint8
	TypeOfService uint8
	TotalLength uint16
	IdentifierFlagsOffset uint32
	TTL uint8
	Protocol uint8
	HeaderChecksum uint16
	SrcIP net.IP
	DestIP net.IP
}

func UnmarshalIPHeader(data []byte) (*IPHeader, error) {
	fmt.Printf("UnmarshalIPHeader: %v\n", data)
	if len(data) < IPv4HeaderLen {
		return nil, fmt.Errorf("IPHeader length is %d, expected %d", len(data), IPv4HeaderLen)
	}

	header := &IPHeader{}
	fmt.Println(data[0])
	header.Version = byte(data[0]) >> 4

	if header.Version != 4 {
		return nil, fmt.Errorf("IP Packet is not version 4, it's %d", header.Version)
	}

	header.Length = ((byte(data[0]) << 4) >> 4) 
	header.TypeOfService = uint8(data[1])
	header.TotalLength = binary.LittleEndian.Uint16(data[2:4])
	header.IdentifierFlagsOffset = binary.LittleEndian.Uint32(data[4:8])
	header.TTL = uint8(data[8])
	header.Protocol = uint8(data[9])
	header.HeaderChecksum = binary.LittleEndian.Uint16(data[10:12])
	header.SrcIP = net.IPv4(data[12], data[13], data[14], data[15])
	header.DestIP = net.IPv4(data[16], data[17], data[18], data[19])

	if HeaderChecksum(data) != 0 || header.TTL == 0 {
		return nil, fmt.Errorf("IP Packet is invalid or outdated")
	}

	return header, nil
}

func MarshalIPHeader(header *IPHeader) []byte {
	data := make([]byte, IPv4HeaderLen)

	data[0] = byte(header.Version << 4) | byte(header.Length)
	data[1] = byte(header.TypeOfService)
	binary.LittleEndian.PutUint16(data[2:4], header.TotalLength)
	binary.LittleEndian.PutUint32(data[4:8], header.IdentifierFlagsOffset)
	data[8] = byte(header.TTL)
	data[9] = byte(header.Protocol)
	binary.LittleEndian.PutUint16(data[10:12], header.HeaderChecksum)
	copy(data[12:16], header.SrcIP.To4())
	copy(data[16:20], header.DestIP.To4())

	return data
}

func HeaderChecksum(data []byte) uint16 {
	var sum uint32 = 0
	for i := 0; i < IPv4HeaderLen; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	return ^uint16((sum >> 16) + sum)
}

func UpdateChecksum(data []byte) {
	// update checksum
	data[10] = 0 // MSB
	data[11] = 0 // LSB
	csum := HeaderChecksum(data)
	data[10] = byte(csum >> 8) // MSB
	data[11] = byte(csum)      // LSB
}