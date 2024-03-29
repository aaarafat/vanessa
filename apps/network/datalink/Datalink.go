package datalink

import (
	"errors"
	"log"
	"net"
	"strings"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/packet"
)

type DataLinkLayerChannel struct {
	source    net.HardwareAddr
	etherType Ethertype
	channel   *packet.Conn
	IfiIndex  int
	Ifi       net.Interface
}

type Ethertype int

const (
	// VEtherType is the EtherType used by the Vanessa test.
	VEtherType        = 0x7031
	VNDEtherType      = 0x7032 // for neighbor discovery
	VAODVEtherType    = 0x7033 // for AODV protocol
	VIEtherType       = 0x7034 // for contacting infrastructure
	VDATAEtherType    = 0x7035 // for sending and receiving IP data packets
	VRSUARPType       = 0x7036 // for cars to send heartbeats to infrastructure
	VDATAFLDEtherType = 0x7036 // to send data packets to with flooder

)

func newDataLinkLayerChannel(ether Ethertype, ifi net.Interface, ifiIndex int) (*DataLinkLayerChannel, error) {
	log.Printf("connecting to: %s\n", ifi.Name)
	// Open a raw socket using same EtherType as our frame.
	c, err := packet.Listen(&ifi, packet.Raw, int(ether), nil)
	if err != nil {
		log.Printf("failed to listen on channel: %v", err)
		return nil, err
	}

	return &DataLinkLayerChannel{
		etherType: ether, // Set the channel type
		channel:   c,
		source:    ifi.HardwareAddr, // Identify the car as the sender.
		IfiIndex:  ifiIndex,
		Ifi:       ifi,
	}, nil

}
func NewDataLinkLayerChannel(ether Ethertype) (*DataLinkLayerChannel, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Printf("failed to open interface: %v", err)
		return nil, err
	}
	ifi := interfaces[1]
	return newDataLinkLayerChannel(ether, ifi, 1)
}

func NewDataLinkLayerChannelWithInterface(ether Ethertype, index int) (*DataLinkLayerChannel, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return nil, err
	}
	ifi := interfaces[index]
	return newDataLinkLayerChannel(ether, ifi, index)
}

func NewDataLinkLayerChannelWithInterfaceName(ether Ethertype, ifiName string) (*DataLinkLayerChannel, error) {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return nil, err
	}
	for i, ifi := range interfaces {
		if strings.Contains(ifi.Name, ifiName) {
			return newDataLinkLayerChannel(ether, ifi, i)
		}
	}
	return nil, errors.New("interface not found")
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
		log.Printf("failed to write Data link frame: %v", err)
	}
}

func (d *DataLinkLayerChannel) Broadcast(payload []byte) {
	d.SendTo(payload, ethernet.Broadcast)
}

func (d *DataLinkLayerChannel) Read() ([]byte, net.HardwareAddr, error) {
	buf := make([]byte, 1500)
	n, addr, err := d.channel.ReadFrom(buf)
	if err != nil {
		log.Printf("failed to read from channel: %v", err)
		return nil, nil, err
	}
	var f ethernet.Frame

	if err := (&f).UnmarshalBinary(buf[:n]); err != nil {
		log.Fatalf("failed to unmarshal ethernet frame: %v", err)
	}
	mac, err := net.ParseMAC(addr.String())
	if err != nil {
		log.Panic(err)
	}
	return f.Payload, mac, nil
}

func (d *DataLinkLayerChannel) Close() {
	d.channel.Close()
}
