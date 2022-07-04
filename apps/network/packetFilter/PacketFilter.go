package packetfilter

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/AkihiroSuda/go-netfilter-queue"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
	"github.com/aaarafat/vanessa/apps/network/unix"
)

type PacketFilter struct {
	nfq *netfilter.NFQueue
	srcIP net.IP
	id int

	// TODO: replace this with router object
	router *aodv.Aodv
	unix *unix.UnixSocket
}

func newPacketFilter(id int, ifi net.Interface) (*PacketFilter, error) {
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
		id: id,
		router: nil,
		unix: unix.NewUnixSocket(id),
	}

	pf.router = aodv.NewAodv(ip, pf.DataCallback)

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
	header, err := UnmarshalIPHeader(dataByte)

	if err != nil {
		log.Println(err)
		return	
	}

	if header.destIP.Equal(pf.srcIP) {
		log.Println("packet is for me")
	}

	payload := dataByte[header.Length:]
	data, err := aodv.UnmarshalData(payload)
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("PacketFilter received data: %v\n", data.Data)
	pf.unix.Write(data.Data)
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

func (pf *PacketFilter) ObstacleHandler() {
	obstacleChannel := make(chan json.RawMessage)
	obstableSubscriber := &unix.Subscriber{Messages: &obstacleChannel}
	pf.unix.Subscribe(unix.ObstacleDetectedEvent, obstableSubscriber)

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
		go pf.router.SendData(data, net.ParseIP(aodv.BroadcastIP))
		pf.unix.Write(data)
	}
}


func (pf *PacketFilter) DestinationReachedHandler() {
	destinationReachedChannel := make(chan json.RawMessage)
	destinationReachedSubscriber := &unix.Subscriber{Messages: &destinationReachedChannel}
	pf.unix.Subscribe(unix.DestinationReachedEvent, destinationReachedSubscriber)

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

func (pf *PacketFilter) UpdateLocationHandler() {
	updateLocationChannel := make(chan json.RawMessage)
	updateLocationSubscriber := &unix.Subscriber{Messages: &updateLocationChannel}
	pf.unix.Subscribe(unix.UpdateLocationEvent, updateLocationSubscriber)

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


func (pf *PacketFilter) Start() {
	log.Printf("Starting PacketFilter for IP: %s.....\n", pf.srcIP)
	go pf.StealPacket()
	go pf.router.Start()
	go pf.unix.Start()
	go pf.ObstacleHandler()
	go pf.DestinationReachedHandler()
	go pf.UpdateLocationHandler()

	// TODO: REMOVE THIS (for testing)
	for {
		time.Sleep(time.Second * 5)
		msg := fmt.Sprintf("Hello From IP: %s\n", pf.srcIP)
		pf.router.SendData([]byte(msg), net.ParseIP(aodv.RsuIP))
	}
}

func (pf *PacketFilter) Close() {
	pf.nfq.Close()
	pf.router.Close()
}