package rsu

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
)

type RSU struct {
	ip net.IP
	ethChannel *DataLinkLayerChannel
	wlanChannel *DataLinkLayerChannel
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
	return &RSU{
		ip: net.ParseIP(RsuIP),
		ethChannel: ethChannel,
		wlanChannel: wlanChannel,
	}
}


func (r* RSU) readFromETHInterface() {
	for {

		data, addr, err := r.ethChannel.Read()
		if err != nil {
			return
		}
		fmt.Println()
		
		packet, err := UnmarshalPacket(data)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v\n", err)
			continue
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
			return
		}
		// Unmarshalling the data received from the WLAN interface.
		fmt.Println()

		packet, err := UnmarshalPacket(data)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v\n", err)
			continue
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
	log.Printf("Closing RSU....\n")
	r.ethChannel.Close()
	r.wlanChannel.Close()
	log.Printf("RSU closed\n")
}