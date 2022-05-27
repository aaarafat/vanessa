package datalink

import (
	"log"
	"net"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/packet"
)

type DataLinkLayerChannel struct {
	source    net.HardwareAddr
	etherType Ethertype
	channel   *packet.Conn
}

type Ethertype int

const (
	// VEtherType is the EtherType used by the Vanessa test.
	VEtherType = 0x7031
)

func NewDataLinkLayerChannel(ether Ethertype) (*DataLinkLayerChannel, error) {
	interfaces, err := net.Interfaces()
	ifi, err := net.InterfaceByName(interfaces[1].Name)

	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return nil, err
	}

	// Open a raw socket using same EtherType as our frame.
	c, err := packet.Listen(ifi, packet.Raw, int(ether), nil)
	if err != nil {
		log.Fatalf("failed to listen on channel: %v", err)
		return nil, err
	}

	return &DataLinkLayerChannel{
		etherType: ether, // Set the channel type
		channel:   c,
		source:    ifi.HardwareAddr, // Identify the car as the sender.
	}, nil

}

func (d *DataLinkLayerChannel) SendTo(payload []byte, destination net.HardwareAddr) {
	frame := &ethernet.Frame{
		// Set the destination MAC address
		Destination: destination,
		// Identify the car as the sender.
		Source: d.source,
		// Identify frame with the same type as channel
		EtherType: ethernet.EtherType(d.etherType),
		Payload:   payload,
	}
	frame_binary, err := frame.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal Data link frame: %v", err)
	}
	// Broadcast the frame to all devices on our network segment.
	addr := &packet.Addr{HardwareAddr: destination}
	if _, err := d.channel.WriteTo(frame_binary, addr); err != nil {
		log.Fatalf("failed to write Data link frame: %v", err)
	}
}

func (d *DataLinkLayerChannel) Broadcast(payload []byte) {
	d.SendTo(payload, ethernet.Broadcast)
}

func (d *DataLinkLayerChannel) Read() ([]byte, net.Addr, error) {
	buf := make([]byte, 1500)
	n, addr, err := d.channel.ReadFrom(buf)
	if err != nil {
		log.Fatalf("failed to read from channel: %v", err)
		return nil, nil, err
	}
	var f ethernet.Frame

	if err := (&f).UnmarshalBinary(buf[:n]); err != nil {
		log.Fatalf("failed to unmarshal ethernet frame: %v", err)
	}
	return f.Payload, addr, nil
}

func (d *DataLinkLayerChannel) Close() {
	d.channel.Close()
}
