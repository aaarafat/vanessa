package network

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

type NetworkLayer struct {
	ip net.IP
	channels map[int]*DataLinkLayerChannel
	forwarders map[int]*Forwarder

	ipConn *IPConnection

	unicastProtocol Protocol
}

func callback(packet []byte) {
	log.Println("Received packet")
}

func NewNetworkLayer(ip net.IP) *NetworkLayer {
	ipConn, err := NewIPConnection()
	if err != nil {
		log.Fatalf("failed to open ip connection: %v", err)
	}

	return &NetworkLayer{
		ip: ip,
		channels: make(map[int]*DataLinkLayerChannel),
		forwarders: make(map[int]*Forwarder),
		ipConn: ipConn,
		unicastProtocol: aodv.NewAodv(ip, callback),
	}
}

func (n *NetworkLayer) openChannels()  {
	channels := make(map[int]*DataLinkLayerChannel)
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}

	for _, ifi := range interfaces {
		if ifi.Name == "lo" {
			continue
		}
		channels[ifi.Index], err = NewDataLinkLayerChannelWithInterface(VDATAEtherType, ifi.Index)
		if err != nil {
			log.Fatalf("failed to open interface: %v", err)
		}
	}

	n.channels = channels
}

func (n *NetworkLayer) openForwarders() {
	for ifiIndex, channel := range n.channels {
		n.forwarders[ifiIndex] = NewForwarder(n.ip, channel)
	}
}

func (n *NetworkLayer) openListeners() {
	for _, channel := range n.channels {
		go n.listen(channel)
	}
}

func (n *NetworkLayer) Start() {
	n.openChannels()
	n.openForwarders()
	n.openListeners()
	
	n.unicastProtocol.Start()
}

func (n *NetworkLayer) Close() {
	for _, channel := range n.channels {
		channel.Close()
	}

	n.ipConn.Close()

	for _, forwarder := range n.forwarders {
		forwarder.Close()
	}

	n.unicastProtocol.Close()
}