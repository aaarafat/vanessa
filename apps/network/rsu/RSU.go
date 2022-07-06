package rsu

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

type RSU struct {
	ip net.IP
	ethChannel *DataLinkLayerChannel
	wlanChannel *DataLinkLayerChannel
	RARP *RSUARP
}

const (
	RSUETHInterface = 1
	RSUWLANInterface = 3
)

func createETHChannel() *DataLinkLayerChannel {
	c, err := NewDataLinkLayerChannelWithInterface(VEtherType, RSUETHInterface)
	if err != nil {
		log.Panicln("No interfaces")
	}
	return c
}

func createWLANChannel() *DataLinkLayerChannel {
	c, err := NewDataLinkLayerChannelWithInterface(VDATAEtherType, RSUWLANInterface)
	if err != nil {
		log.Panicln("No interfaces")
	}
	return c
}

func NewRSU() *RSU {
	ethChannel := createETHChannel()
	wlanChannel := createWLANChannel()
	RARP := NewRSUARP()
	return &RSU{
		ip: net.ParseIP(aodv.RsuIP), // TODO extract RSUIP out from aodv
		ethChannel: ethChannel,
		wlanChannel: wlanChannel,
		RARP: RARP,
	}
}


func (r* RSU) readFromETHInterface() {
	for {

		data, addr, err := r.ethChannel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		
		packet, err := UnmarshalPacket(data)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v\n", err)
			return
		}

		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(packet.Payload), addr.String(), RSUETHInterface)
		// TODO: Forward up with callback to decide what to do
		r.wlanChannel.Broadcast(data)
	}
}

func (r *RSU) readFromWLANInterface() {
	for {

		data, addr, err := r.wlanChannel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		// Unmarshalling the data received from the WLAN interface.
		fmt.Println()

		packet, err := UnmarshalPacket(data)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v\n", err)
			return
		}

		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(packet.Payload), addr.String(), RSUWLANInterface)
		// TODO: Forward up with callback to decide what to do
		r.ethChannel.Broadcast(data)
	}
}




func (r *RSU) Start() {
	go r.readFromETHInterface()
	go r.readFromWLANInterface()
}

func (r *RSU) Close() {
	r.ethChannel.Close()
	r.wlanChannel.Close()
}

