package packetfilter

import (
	"fmt"
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

type PacketFilter struct {
	nfq *netfilter.NFQueue
	srcIP net.IP

	// TODO: replace this with router object
	router *aodv.Aodv
}

func NewPacketFilter() (*PacketFilter, error) {
	var err error
	
	AddIPTablesRule()
	if err := RegisterGateway(); err != nil {
		DeleteIPTablesRule()
		log.Panic("done deleting")
		return nil, err
	}
	interfaces, err := net.Interfaces()
	for index , intf := range(interfaces){
		log.Println(index,intf)
	}
	iface := interfaces[1]

	ip, _, err := GetMyIPs(&iface)
	ip = ip.To4()
	if err != nil {
		log.Panicf("failed to get iface ips, err: %s", err)
		return nil, err
	}

	SetMaxMSS(iface.Name, ip, 1400)

	nfq, err := netfilter.NewNFQueue(0, 100, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
		log.Println(err)
		return nil, err
}

	return &PacketFilter{
		nfq: nfq,
		srcIP: ip,
		router: aodv.NewAodv(ip),
	}, nil
}

func (pf *PacketFilter) StealPacket() {
	packets := pf.nfq.GetPackets()
	for {
		select {
		case p := <-packets:
			packet := p.Packet.Data()
			header, err := UnmarshalIPHeader(packet)

			if err != nil {
				p.SetVerdict(netfilter.NF_DROP)
			}

			if pf.srcIP.Equal(header.destIP) {
				fmt.Println(p.Packet)
				p.SetVerdict(netfilter.NF_ACCEPT)
			} else {
				p.SetVerdict(netfilter.NF_DROP)

				go pf.router.SendData(packet, header.destIP)
			}

		}
	}
}

func (pf *PacketFilter) Start() {
	log.Printf("Starting PacketFilter for IP: %s.....\n", pf.srcIP)
	go pf.StealPacket()
}

func (pf *PacketFilter) Close() {
	pf.nfq.Close()
}