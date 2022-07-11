package rsu

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
)

type RSU struct {
	ip          net.IP
	ethChannel  *DataLinkLayerChannel
	wlanChannel *DataLinkLayerChannel
	RARP        *RSUARP
	OTable      *ObstaclesTable
}

const (
	RSUETHInterface  = 1
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
	OTable := NewObstaclesTable()
	return &RSU{
		ip:          net.ParseIP(RsuIP),
		ethChannel:  ethChannel,
		wlanChannel: wlanChannel,
		RARP:        RARP,
		OTable:      OTable,
	}
}

func (r *RSU) readFromETHInterface() {
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

		r.handleEthMessages(*packet, addr)
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

		r.handleMessage(*packet, addr)
	}
}

// get MAC from RSUARP using ip then send to wlan interface
func (r *RSU) sendToWLANInterface(data []byte, ip string) {
	mac := r.RARP.Get(ip)
	r.wlanChannel.SendTo(data, mac)
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
