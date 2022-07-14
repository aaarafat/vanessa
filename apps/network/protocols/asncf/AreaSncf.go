package asncf

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
)

type AreaSNCF struct {
	srcIP    net.IP
	position *Position

	seqNum   uint32
	seqTable *VFloodingSeqTable
	channel  *DataLinkLayerChannel

	flooder IFlooder
}

func NewAreaSNCF(ip net.IP, ifiName string, nt *VNeighborTable) *AreaSNCF {
	channel, err := NewDataLinkLayerChannelWithInterfaceName(VDATAFLDEtherType, ifiName)
	if err != nil {
		log.Fatalf("failed to open channel: %v", err)
	}

	return &AreaSNCF{
		srcIP:    ip,
		seqNum:   0,
		channel:  channel,
		seqTable: NewVFloodingSeqTable(),
		flooder:  NewFlooder(ip, channel, nt),
	}
}

func (sncf *AreaSNCF) updateSeqNum() {
	sncf.seqNum++
}

func (sncf *AreaSNCF) Flood(bytes []byte) {
	packet, err := ip.UnmarshalPacket(bytes)
	if err != nil {
		log.Printf("Error unmarshaling packet: %v\n", err)
		return
	}
	if !sncf.validOptions(packet) {
		return
	}

	// Add flooding header
	seqHeader := NewSequenceHeader(sncf.seqNum)
	sncf.updateSeqNum()

	data := make([]byte, len(bytes)+4)

	copy(data[0:4], seqHeader.Marshal())
	copy(data[4:], bytes)

	sncf.flooder.ForwardToAll(data)
}

func (sncf *AreaSNCF) forward(data []byte, fromIP net.IP, fromMac net.HardwareAddr) {
	seqHeader := UnmarshalSequenceHeader(data[0:4])
	if sncf.seqTable.Exists(fromIP, seqHeader.SequenceNumber) || sncf.srcIP.Equal(fromIP) {
		log.Printf("SNCF: duplicate packet from %v", fromIP)
		return
	}

	sncf.seqTable.Set(fromIP, seqHeader.SequenceNumber)

	sncf.flooder.ForwardToAllExcept(data, fromMac)
}

func (sncf *AreaSNCF) validOptions(packet *ip.IPPacket) bool {
	if !packet.HasOptions() {
		log.Printf("Packet has no options\n")
		return true
	}

	if packet.Header.Options[0] == ip.PositionOptionType {
		positionOption, err := ip.UnmarshalPositionOption(packet.Header.Options)
		if err != nil {
			log.Printf("Error decoding Position Option: %v\n", err)
			return false
		}

		if sncf.srcIP.Equal(packet.Header.SrcIP) {
			// update my position
			sncf.position = &positionOption.Position
			log.Printf("Updated position to lng: %f  lat: %f\n", sncf.position.Lng, sncf.position.Lat)
			return true
		}

		if sncf.position == nil {
			// have not received my position yet
			log.Printf("Have not received my position yet\n")
			return false
		}

		// check distance
		distance := sncf.position.Distance(&positionOption.Position)

		log.Printf("Distance: %f from %s\n", distance, packet.Header.SrcIP)

		if distance < positionOption.MaxDistance {
			log.Printf("Distance is less than %f from %s\n", positionOption.MaxDistance, packet.Header.SrcIP)
			return true
		}
		log.Printf("Distance is greater than %f from %s\n", positionOption.MaxDistance, packet.Header.SrcIP)
	}

	log.Printf("Packet has invalid options %d\n", packet.Header.Options[0])

	return false
}

func (sncf *AreaSNCF) Read() ([]byte, net.HardwareAddr, error) {
	data, addr, err := sncf.channel.Read()
	if err != nil {
		log.Printf("Error reading from channel: %v\n", err)
		return nil, nil, err
	}

	log.Printf("SNCF: received packet from %v\n", addr)

	packet, err := ip.UnmarshalPacket(data[4:])
	if err != nil {
		log.Printf("Error unmarshaling packet: %v\n", err)
		return nil, nil, err
	}

	seqHeader := UnmarshalSequenceHeader(data[0:4])
	if sncf.seqTable.Exists(packet.Header.SrcIP, seqHeader.SequenceNumber) || sncf.srcIP.Equal(packet.Header.SrcIP) {
		log.Printf("SNCF: duplicate packet from %v", packet.Header.SrcIP)
		return nil, nil, fmt.Errorf("duplicate packet")
	}

	// add options check
	if !sncf.validOptions(packet) {
		return nil, nil, fmt.Errorf("Packet has invalid options")
	}

	ip.Update(data[4:], packet.Header.LengthInBytes())

	go sncf.forward(data, packet.Header.SrcIP, addr)

	return data[4:], addr, nil
}

func (sncf *AreaSNCF) Close() {
	sncf.channel.Close()
}
