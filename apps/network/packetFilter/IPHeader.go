package packetfilter

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
	srcIP net.IP
	destIP net.IP
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

	header.Length = (byte(data[0]) << 4) >> 4
	header.TypeOfService = uint8(data[1])
	header.TotalLength = binary.LittleEndian.Uint16(data[2:4])
	header.IdentifierFlagsOffset = binary.LittleEndian.Uint32(data[4:8])
	header.TTL = uint8(data[8])
	header.Protocol = uint8(data[9])
	header.HeaderChecksum = binary.LittleEndian.Uint16(data[10:12])
	header.srcIP = net.IPv4(data[12], data[13], data[14], data[15])
	header.destIP = net.IPv4(data[16], data[17], data[18], data[19])

	if HeaderChecksum(data) != 0 || header.TTL == 0 {
		return nil, fmt.Errorf("IP Packet is invalid or outdated")
	}

	return header, nil
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