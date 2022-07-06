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
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_RAW, syscall.IPPROTO_RAW)
	if err != nil {
		return nil, err
	}
	
	return &IPConnection{
		fd:     fd,
	}, nil
}

func (c *IPConnection) Write(payload []byte, srcIp, destIp net.IP) error {
	// 127.0.0.1 for loopback
	packet := NewIPPacket(payload, srcIp, destIp)
	packetBytes := MarshalIPPacket(packet)
	UpdateChecksum(packetBytes)

	return c.Forward(packetBytes)
}

func (c *IPConnection) Forward(packet []byte) error {
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

