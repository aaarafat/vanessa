package packetfilter

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

type PacketFilter struct {
	nfq *netfilter.NFQueue
	srcIP net.IP

	// TODO: replace this with router object
	router *aodv.Aodv
}

func newPacketFilter(ifi net.Interface) (*PacketFilter, error) {
	var err error
	
	if err := ChainNFQUEUE(); err != nil {
		DeleteIPTablesRule()
		log.Panic("Reversed chaining NFQUEUE")
		return nil, err
	}

	if err := AddDefaultGateway(); err != nil {
		DeleteDefaultGateway()
		log.Panic("Removed Default Gatway")
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
		nfq: nfq,
		srcIP: ip,
		router: nil,
	}

	pf.router = aodv.NewAodv(ip, pf.DataCallback)

	return pf, nil
}

func NewPacketFilter() (*PacketFilter, error) {

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Println(err)
		return nil, err
	}
	iface := interfaces[1]
	return newPacketFilter(iface)
}

func NewPacketFilterWithInterface(ifi net.Interface) (*PacketFilter, error) {
	return newPacketFilter(ifi)
}


func (pf *PacketFilter) DataCallback(data []byte) {
	header, err := UnmarshalIPHeader(data)

	if err != nil {
		log.Println(err)
		return	
	}

	if header.destIP.Equal(pf.srcIP) {
		log.Println("packet is for me")
	}

	payload := data[header.Length:]
	log.Printf("PacketFilter received data: %v\n", payload)
}

func (pf *PacketFilter) StealPacket() {
	packets := pf.nfq.GetPackets()
	for {
		select {
		case p := <-packets:
			packet := p.Packet.Data()
			log.Printf("PacketFilter received packet: %v\n", packet)

			header, err := UnmarshalIPHeader(packet)

			if err != nil {
				p.SetVerdict(netfilter.NF_DROP)
			}

			if pf.srcIP.Equal(header.destIP) {
				fmt.Println(p.Packet)
				p.SetVerdict(netfilter.NF_ACCEPT)
			} else {
				p.SetVerdict(netfilter.NF_DROP)

				UpdateChecksum(packet)


				log.Println(header.Version)
				go pf.router.SendData(packet, header.destIP)
			}

		}
	}
}

func (pf *PacketFilter) Start() {
	log.Printf("Starting PacketFilter for IP: %s.....\n", pf.srcIP)
	go pf.StealPacket()
	go pf.router.Start()

	// TODO: REMOVE THIS (for testing)
	for {
		time.Sleep(time.Second * 5)
		msg := fmt.Sprintf("Hello From IP: %s\n", pf.srcIP)
		pf.router.SendData([]byte(msg), net.ParseIP(aodv.RsuIP))
	}
}

func (pf *PacketFilter) Close() {
	pf.nfq.Close()
}