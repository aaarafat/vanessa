package packetfilter

import (
	"fmt"
	"log"
	"net"

	"github.com/AkihiroSuda/go-netfilter-queue"
)

type PacketFilter struct {
	nfq *netfilter.NFQueue
	srcIP net.IP
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
	}, nil
}

func (pf *PacketFilter) StealPacket() {
	packets := pf.nfq.GetPackets()
	for {
		select {
		case p := <-packets:
			fmt.Println(p.Packet)
			// drop tha packet
			p.SetVerdict(netfilter.NF_DROP)
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