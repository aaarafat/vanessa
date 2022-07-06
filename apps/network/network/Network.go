package network

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
	"github.com/cornelk/hashmap"
)

type NetworkLayer struct {
	ip net.IP
	channels map[int]*DataLinkLayerChannel
	forwarders map[int]*Forwarder

	// buffer to store packets until path is found
	packetBuffer *hashmap.HashMap

	ipConn *IPConnection

	unicastProtocol Protocol
}

func NewNetworkLayer(ip net.IP) *NetworkLayer {
	ipConn, err := NewIPConnection()
	if err != nil {
		log.Fatalf("failed to open ip connection: %v", err)
	}

	network := &NetworkLayer{
		ip: ip,
		channels: make(map[int]*DataLinkLayerChannel),
		forwarders: make(map[int]*Forwarder),
		packetBuffer: &hashmap.HashMap{},
		ipConn: ipConn,
	}

	network.unicastProtocol = aodv.NewAodv(ip, network.onPathDiscovery)

	return network
}

func (n *NetworkLayer) openChannels()  {
	channels := make(map[int]*DataLinkLayerChannel)
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}

	for index, ifi := range interfaces {
		if ifi.Name == "lo" {
			continue
		}
		channels[index], err = NewDataLinkLayerChannelWithInterface(VDATAEtherType, index)
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
	log.Printf("Closing network layer")

	for _, forwarder := range n.forwarders {
		forwarder.Close()
	}

	for _, channel := range n.channels {
		channel.Close()
	}

	n.ipConn.Close()

	n.unicastProtocol.Close()
	log.Printf("Closed network layer")
}