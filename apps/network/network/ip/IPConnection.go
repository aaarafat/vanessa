package ip

import (
	"log"
	"net"
	"syscall"
)

type IPConnection struct {
	fd     int
}

func NewIPConnection() (*IPConnection, error) {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW | syscall.IP_HDRINCL)
	if err != nil {
		return nil, err
	}
	
	log.Printf("Opened IP connection with fd: %d\n", fd)

	return &IPConnection{
		fd:     fd,
	}, nil
}

func (c *IPConnection) Write(payload []byte, srcIp, destIp net.IP) error {
	log.Printf("Writing Raw Sockets with payload %s to %s\n", payload, destIp)
	packet := NewIPPacket(payload, srcIp, destIp)
	packetBytes := MarshalIPPacket(packet)

	return c.Forward(packetBytes)
}

func (c *IPConnection) Forward(packet []byte) error {
	// 127.0.0.1 for loopback
	err := syscall.Sendto(c.fd, packet, 0, &syscall.SockaddrInet4{Port: 0, Addr: [4]byte{127, 0, 0, 1}})
	if err != nil {
		log.Printf("failed to send data: %v", err)
		return err
	}
	return nil
}

func (c *IPConnection) Close() {
	syscall.Close(c.fd)
}

