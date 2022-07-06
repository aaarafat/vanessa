package packetfilter

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AkihiroSuda/go-netfilter-queue"
	. "github.com/aaarafat/vanessa/apps/network/network"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
	"github.com/aaarafat/vanessa/apps/network/unix"
)

type PacketFilter struct {
	nfq   *netfilter.NFQueue
	srcIP net.IP
	id    int
	networkLayer *NetworkLayer 

	unix   *unix.UnixSocket
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
		unix:   unix.NewUnixSocket(id),
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

func (pf *PacketFilter) DataCallback(dataByte []byte) {
	packet, err := UnmarshalPacket(dataByte)

	if err != nil {
		log.Println(err)
		return
	}

	if packet.Header.DestIP.Equal(pf.srcIP) {
		log.Println("packet is for me")
	}

	log.Printf("PacketFilter received data: %s\n", packet.Payload)
	pf.unix.Write(packet.Payload)
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
				log.Printf("Error decoding IP header: %v\n", err)
				p.SetVerdict(netfilter.NF_DROP)
				continue
			}

			log.Printf("PacketFilter received DestIP: %s, SrcIP: %s\n", header.DestIP, header.SrcIP)

			if pf.srcIP.Equal(header.DestIP) {
				fmt.Println(p.Packet)
				p.SetVerdict(netfilter.NF_ACCEPT)

				// TODO : grpc call to the router to process the packet
			} else {
				p.SetVerdict(netfilter.NF_DROP)

				Update(packet)

				go pf.networkLayer.SendUnicast(packet, header.DestIP)
			}
		}
	}
}

func (pf *PacketFilter) ObstacleHandler() {
	obstacleChannel := make(chan json.RawMessage)
	obstableSubscriber := &unix.Subscriber{Messages: &obstacleChannel}
	pf.unix.Subscribe(unix.ObstacleDetectedEvent, obstableSubscriber)

	for {
		select {
		case data := <-*obstableSubscriber.Messages:
			var obstacle unix.ObstacleDetectedData
			err := json.Unmarshal(data, &obstacle)
			if err != nil {
				log.Printf("Error decoding obstacle-detected data: %v", err)
				return
			}
			log.Printf("Packet Filter : Obstacle detected: %v\n", data)

			// TODO: send it with loopback interface to the router to be processed by the AODV
			go pf.networkLayer.Send(data, pf.srcIP, net.ParseIP(aodv.BroadcastIP))
			pf.unix.Write(data)
		}
	}
}

func (pf *PacketFilter) DestinationReachedHandler() {
	destinationReachedChannel := make(chan json.RawMessage)
	destinationReachedSubscriber := &unix.Subscriber{Messages: &destinationReachedChannel}
	pf.unix.Subscribe(unix.DestinationReachedEvent, destinationReachedSubscriber)

	for {
		select {
		case data := <-*destinationReachedSubscriber.Messages:
			var destinationReached unix.DestinationReachedData
			err := json.Unmarshal(data, &destinationReached)
			if err != nil {
				log.Printf("Error decoding destination-reached data: %v", err)
				return
			}
			log.Printf("Packet Filter : destination-reached: %v\n", data)

			pf.unix.Write(data)
		}
	}
}

func (pf *PacketFilter) UpdateLocationHandler() {
	updateLocationChannel := make(chan json.RawMessage)
	updateLocationSubscriber := &unix.Subscriber{Messages: &updateLocationChannel}
	pf.unix.Subscribe(unix.UpdateLocationEvent, updateLocationSubscriber)

	for {
		select {
		case data := <-*updateLocationSubscriber.Messages:
			var updateLocation unix.UpdateLocationData
			err := json.Unmarshal(data, &updateLocation)
			if err != nil {
				log.Printf("Error decoding update-location data: %v", err)
				return
			}
			log.Printf("Packet Filter : update-location: %v\n", data)

			pf.unix.Write(data)
		}
	}
}

func (pf *PacketFilter) SendHelloToRSU() {
	for {
		time.Sleep(time.Second * 5)
		msg := fmt.Sprintf("Hello From IP: %s\n", pf.srcIP)
		pf.networkLayer.Send([]byte(msg), pf.srcIP, net.ParseIP(aodv.RsuIP))
	}
}

func (pf *PacketFilter) Start() {
	log.Printf("Starting PacketFilter for IP: %s.....\n", pf.srcIP)
	go pf.unix.Start()
	go pf.networkLayer.Start()
	// go pf.ObstacleHandler()
	// go pf.DestinationReachedHandler()
	// go pf.UpdateLocationHandler()
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
