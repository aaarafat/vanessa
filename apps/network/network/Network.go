package network

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

const (
	UNICAST_IFI = "wlan0"
	RSU_IFI     = "wlan1"
	WLAN0       = 0
	WLAN1       = 1
)

type NetworkLayer struct {
	ip       net.IP
	channels map[int]*DataLinkLayerChannel
	flooders map[int]*Flooder

	neighborTables map[int]*VNeighborTable

	// buffer to store packets until path is found
	packetBuffer *PacketBuffer

	ipConn *IPConnection

	unicastProtocol Protocol
}

func NewNetworkLayer(ip net.IP) *NetworkLayer {
	ipConn, err := NewIPConnection()
	if err != nil {
		log.Fatalf("failed to open ip connection: %v", err)
	}

	network := &NetworkLayer{
		ip:             ip,
		channels:       make(map[int]*DataLinkLayerChannel),
		flooders:       make(map[int]*Flooder),
		neighborTables: make(map[int]*VNeighborTable),
		packetBuffer:   NewPacketBuffer(),
		ipConn:         ipConn,
	}

	return network
}

func (n *NetworkLayer) openChannels() {
	// open WLAN0 channel
	ch0, err := NewDataLinkLayerChannelWithInterfaceName(VDATAEtherType, UNICAST_IFI)
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	n.channels[WLAN0] = ch0
	n.neighborTables[WLAN0] = NewVNeighborTable(n.ip, UNICAST_IFI, false)
	n.flooders[WLAN0] = NewFlooder(n.ip, n.channels[WLAN0], n.neighborTables[WLAN0])
	go n.neighborTables[WLAN0].Run()

	// open WLAN1 channel
	ch1, err := NewDataLinkLayerChannelWithInterfaceName(VDATAEtherType, RSU_IFI)
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}
	n.channels[WLAN1] = ch1
	n.neighborTables[WLAN1] = NewVNeighborTable(n.ip, RSU_IFI, true)
	n.flooders[WLAN1] = NewFlooder(n.ip, n.channels[WLAN1], n.neighborTables[WLAN1])
	go n.neighborTables[WLAN1].Run()
}

func (n *NetworkLayer) openListeners() {
	for _, channel := range n.channels {
		go n.listen(channel)
	}
}

func (n *NetworkLayer) Start() {
	n.unicastProtocol = aodv.NewAodv(n.ip, UNICAST_IFI, n.neighborTables, n.onPathDiscovery)
	n.openChannels()
	n.openListeners()

	n.unicastProtocol.Start()
}

func (n *NetworkLayer) Close() {
	log.Printf("Closing network layer")
	n.unicastProtocol.Close()

	for _, channel := range n.channels {
		channel.Close()
	}

	for _, nt := range n.neighborTables {
		nt.Close()
	}

	n.ipConn.Close()

	log.Printf("Closed network layer")
}
