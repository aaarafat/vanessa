package ip

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	IPv4HeaderLen = 20
)

// https://www.techtarget.com/searchnetworking/tutorial/Routing-First-Step-IP-header-format
type IPHeader struct {
	Version               uint8
	Length                uint8
	TypeOfService         uint8
	TotalLength           uint16
	IdentifierFlagsOffset uint32
	TTL                   uint8
	Protocol              uint8
	HeaderChecksum        uint16
	SrcIP                 net.IP
	DestIP                net.IP
	Options               []byte
}

func UnmarshalIPHeader(data []byte) (*IPHeader, error) {
	if len(data) < IPv4HeaderLen {
		return nil, fmt.Errorf("IPHeader length is %d, expected %d", len(data), IPv4HeaderLen)
	}

	header := &IPHeader{}
	header.Version = byte(data[0]&0xf0) / 16

	if header.Version != 4 {
		return nil, fmt.Errorf("IP Packet is not version 4, it's %d", header.Version)
	}

	header.Length = data[0] & 0x0f
	header.TypeOfService = uint8(data[1])
	header.TotalLength = binary.LittleEndian.Uint16(data[2:4])
	bytes := make([]byte, 2)
	binary.LittleEndian.PutUint16(bytes[:], header.TotalLength)
	header.IdentifierFlagsOffset = binary.LittleEndian.Uint32(data[4:8])
	header.TTL = uint8(data[8])
	header.Protocol = uint8(data[9])
	header.HeaderChecksum = binary.LittleEndian.Uint16(data[10:12])
	header.SrcIP = net.IPv4(data[12], data[13], data[14], data[15])
	header.DestIP = net.IPv4(data[16], data[17], data[18], data[19])

	header.Options = data[20:header.LengthInBytes()]

	if csm := HeaderChecksum(data, header.LengthInBytes()); csm != 0 || header.TTL == 0 {
		return nil, fmt.Errorf("IP Packet is invalid or outdated csm: %d, ttl: %d", csm, header.TTL)
	}

	return header, nil
}

func MarshalIPHeader(header *IPHeader) []byte {
	data := make([]byte, header.LengthInBytes())

	data[0] = byte(header.Version<<4) | byte(header.Length)
	data[1] = byte(header.TypeOfService)
	binary.LittleEndian.PutUint16(data[2:], header.TotalLength)
	binary.LittleEndian.PutUint32(data[4:], header.IdentifierFlagsOffset)
	data[8] = byte(header.TTL)
	data[9] = byte(header.Protocol)
	binary.LittleEndian.PutUint16(data[10:], header.HeaderChecksum)
	copy(data[12:16], header.SrcIP.To4())
	copy(data[16:20], header.DestIP.To4())

	copy(data[20:header.LengthInBytes()], header.Options)

	return data
}

func (h *IPHeader) LengthInBytes() int {
	return int(h.Length * 4)
}

func HeaderChecksum(data []byte, len int) uint16 {
	sum := uint32(0)
	for i := 0; i < len; i += CHECKSUM_BLOCK_SIZE {
		cur := uint32(0)
		for j := i; j < i+CHECKSUM_BLOCK_SIZE && j < len; j++ {
			cur |= uint32(data[j]) << uint(8*(CHECKSUM_BLOCK_SIZE-1)-8*(j-i))
		}
		sum += cur
	}
	check := uint16(sum >> 16)
	carry := uint16(sum & 0xffff)
	check += carry
	// do one's complement
	check = ^check
	return check
}

func UpdateChecksum(data []byte, len int) {
	data[10] = 0
	data[11] = 0
	csum := HeaderChecksum(data, len)
	data[10] = byte(csum >> 8)
	data[11] = byte(csum)
}

func UpdateTTL(data []byte) {
	data[8] = byte(uint8(data[8]) - 1)
}

func Update(data []byte, len int) {
	UpdateTTL(data)
	UpdateChecksum(data, len)
}

func (h *IPHeader) String() string {
	return fmt.Sprintf("IPv4 Header: %s -> %s", h.SrcIP, h.DestIP)
}
