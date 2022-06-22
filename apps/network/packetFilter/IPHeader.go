package packetfilter

import "net"

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
	Options net.IP
}