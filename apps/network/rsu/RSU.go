package rsu

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
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
	c, err := NewDataLinkLayerChannelWithInterface(VAODVEtherType, RSUWLANInterface)
	if err != nil {
		log.Panicln("No interfaces")
	}
	return c
}

func NewRSU() *RSU {
	ethChannel := createETHChannel()
	wlanChannel := createWLANChannel()
	return &RSU{
		ip: net.ParseIP(aodv.RsuIP),
		ethChannel: ethChannel,
		wlanChannel: wlanChannel,
	}
}


func (r* RSU) readFromETHInterface() {
	for {

		payload, addr, err := r.ethChannel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		data, err := aodv.UnmarshalData(payload)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v\n", err)
			return
		}

		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(data.Data), addr.String(), RSUETHInterface)
		r.wlanChannel.Broadcast(payload)
	}
}

func (r *RSU) readFromWLANInterface() {
	for {

		payload, addr, err := r.wlanChannel.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		// Unmarshalling the data received from the WLAN interface.
		fmt.Println()

		data, err := aodv.UnmarshalData(payload)
		if err != nil {
			log.Printf("Failed to unmarshal data: %v\n", err)
			return
		}

		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(data.Data), addr.String(), RSUWLANInterface)
		r.ethChannel.Broadcast(payload)
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