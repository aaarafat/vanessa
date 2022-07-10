package packetfilter

import (
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/aaarafat/vanessa/apps/network/network"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/aaarafat/vanessa/apps/network/unix"
)

type PacketFilter struct {
	nfq          *netfilter.NFQueue
	srcIP        net.IP
	id           int
	networkLayer *NetworkLayer
	routerSocket *unix.RouterSocket
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
		nfq:          nfq,
		srcIP:        ip,
		id:           id,
		networkLayer: NewNetworkLayer(ip),
		routerSocket: unix.NewRouterSocket(id),
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

func (pf *PacketFilter) dataCallback(packet []byte, from net.IP) {
	log.Printf("Received packet size %d from %s\n", len(packet), from)
	pf.routerSocket.Write(packet)
}

func (pf *PacketFilter) StealPacket() {
	packets := pf.nfq.GetPackets()
	for {
		select {
		case p := <-packets:
			go func() {
				p.SetVerdict(netfilter.NF_DROP)
				packetBytes := p.Packet.Data()
				packet, err := UnmarshalPacket(packetBytes)
				if err != nil {
					log.Printf("Error decoding IP header: %v\n", err)
					return
				}

				if pf.srcIP.Equal(packet.Header.DestIP) {
					pf.dataCallback(packetBytes, packet.Header.SrcIP)

					// TODO : grpc call to the router to process the packet
				} else if packet.Header.DestIP.Equal(net.ParseIP(ip.BroadcastIP)) {
					if !pf.srcIP.Equal(packet.Header.SrcIP) {
						pf.dataCallback(packetBytes, packet.Header.SrcIP)
					}

					Update(packetBytes)

					log.Printf("Sending packet size %d to %s\n", len(packetBytes), packet.Header.DestIP)
					pf.networkLayer.SendBroadcast(packetBytes, packet.Header.SrcIP)
				} else {
					Update(packetBytes)

					log.Printf("Sending packet size %d to %s\n", len(packetBytes), packet.Header.DestIP)
					pf.networkLayer.SendUnicast(packetBytes, packet.Header.DestIP)
				}
			}()
		}
	}
}

func (pf *PacketFilter) Start() {
	log.Printf("Starting PacketFilter for IP: %s.....\n", pf.srcIP)
	go pf.networkLayer.Start()

	pf.StealPacket()
}

func (pf *PacketFilter) Close() {
	log.Printf("Closing PacketFilter for IP: %s.....\n", pf.srcIP)
	DeleteIPTablesRule()
	DeleteDefaultGateway()

	pf.routerSocket.Close()
	pf.networkLayer.Close()

	log.Printf("PacketFilter for IP: %s closed\n", pf.srcIP)
}
