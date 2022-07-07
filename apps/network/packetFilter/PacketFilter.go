package packetfilter

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/aaarafat/vanessa/apps/network/network"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
)

type PacketFilter struct {
	nfq   *netfilter.NFQueue
	srcIP net.IP
	id    int
	networkLayer *NetworkLayer 
}

func newPacketFilter(id int, ifi net.Interface) (*PacketFilter, error) {
	var err error

	if err := ChainNFQUEUE(); err != nil {
		DeleteIPTablesRule()
		log.Panicf("Reversed chaining NFQUEUE %v\n", err)
		return nil, err
	}

	if err := AddDefaultGateway(); err != nil {
		DeleteDefaultGateway()
		log.Panicf("Removed Default Gatway %v\n", err)
		return nil, err
	}
	ip, _, err := MyIP(&ifi)
	ip = ip.To4()
	if err != nil {
		log.Panicf("failed to get iface ips, err: %s", err)
		return nil, err
	}

	SetMSS(ifi.Name, ip, 1400)

	nfq, err := netfilter.NewNFQueue(0, 100, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pf := &PacketFilter{
		nfq:    nfq,
		srcIP:  ip,
		id:     id,
		networkLayer: NewNetworkLayer(ip),
	}

	return pf, nil
}

func NewPacketFilter(id int) (*PacketFilter, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	iface := interfaces[1]
	return newPacketFilter(id, iface)
}

func NewPacketFilterWithInterface(id int, ifi net.Interface) (*PacketFilter, error) {
	return newPacketFilter(id, ifi)
}

func (pf *PacketFilter) dataCallback(payload []byte) {
	log.Printf("Received: %s\n", payload)
}

func (pf *PacketFilter) StealPacket() {
	packets := pf.nfq.GetPackets()
	for {
		select {
		case p := <-packets:
			packetBytes := p.Packet.Data()
			packet, err := UnmarshalPacket(packetBytes)
			if err != nil {
				log.Printf("Error decoding IP header: %v\n", err)
				p.SetVerdict(netfilter.NF_DROP)
				continue
			}

			if pf.srcIP.Equal(packet.Header.DestIP) {
				pf.dataCallback(packet.Payload)
				p.SetVerdict(netfilter.NF_ACCEPT)

				// TODO : grpc call to the router to process the packet
			} else {
				p.SetVerdict(netfilter.NF_DROP)

				Update(packetBytes)

				log.Printf("Sending packet %s to %s\n", packet.Payload, packet.Header.DestIP)
				go pf.networkLayer.SendUnicast(packetBytes, packet.Header.DestIP)
			}
		}
	}
}

func (pf *PacketFilter) SendHelloToRSU() {
	for {
		time.Sleep(time.Second * 5)
		msg := fmt.Sprintf("Hello From IP: %s\n", pf.srcIP)
		pf.networkLayer.Send([]byte(msg), pf.srcIP, net.ParseIP(RsuIP))
	}
}

func (pf *PacketFilter) Start() {
	log.Printf("Starting PacketFilter for IP: %s.....\n", pf.srcIP)
	go pf.networkLayer.Start()
	// TODO: REMOVE THIS (for testing)
	go pf.SendHelloToRSU()
	
	pf.StealPacket()
}

func (pf *PacketFilter) Close() {
	log.Printf("Closing PacketFilter for IP: %s.....\n", pf.srcIP)
	DeleteIPTablesRule()
	DeleteDefaultGateway()

	pf.networkLayer.Close()

	log.Printf("PacketFilter for IP: %s closed\n", pf.srcIP)
}
