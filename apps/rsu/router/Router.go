package router

import (
	"fmt"
	"log"
	"net"
	"os"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/ip"
)

type Router struct {
	ip              net.IP
	ethChannel      *DataLinkLayerChannel
	wlanChannel     *DataLinkLayerChannel
	neighborChannel *DataLinkLayerChannel

	RARP        *RSUARP
	onARPDelete func(ip string, mac net.HardwareAddr)
}

const (
	RSUETHInterface  = 1
	RSUWLANInterface = 3
)

func createETHChannel() *DataLinkLayerChannel {
	c, err := NewDataLinkLayerChannelWithInterface(VEtherType, RSUETHInterface)
	if err != nil {
		log.Println("No interfaces")
		os.Exit(1)
	}
	return c
}

func createWLANChannel() *DataLinkLayerChannel {
	c, err := NewDataLinkLayerChannelWithInterface(VDATAEtherType, RSUWLANInterface)
	if err != nil {
		log.Println("No interfaces")
		os.Exit(1)
	}
	return c
}

func createNeighborChannel() *DataLinkLayerChannel {
	c, err := NewDataLinkLayerChannelWithInterface(VNDEtherType, RSUWLANInterface)
	if err != nil {
		log.Println("No interfaces")
		os.Exit(1)
	}
	return c
}

func NewRouter(arp *RSUARP) *Router {
	ethChannel := createETHChannel()
	wlanChannel := createWLANChannel()
	neighborChannel := createNeighborChannel()
	return &Router{
		ip:              net.ParseIP(RsuIP),
		ethChannel:      ethChannel,
		wlanChannel:     wlanChannel,
		neighborChannel: neighborChannel,
		RARP:            arp,
	}
}

func (r *Router) ReadFromETHInterface() (*IPPacket, net.HardwareAddr, error) {
	data, addr, err := r.ethChannel.Read()
	if err != nil {
		return nil, nil, err
	}
	fmt.Println()

	packet, err := UnmarshalPacket(data)
	if err != nil {
		log.Printf("Failed to unmarshal data: %v\n", err)
		return nil, nil, err
	}

	log.Printf("Received \"%s\" from: [%s] on intf-%d", string(packet.Payload), addr.String(), RSUETHInterface)

	return packet, addr, err
}

func (r *Router) ReadFromWLANInterface() (*IPPacket, net.HardwareAddr, error) {
	data, addr, err := r.wlanChannel.Read()
	if err != nil {
		return nil, nil, err
	}
	// Unmarshalling the data received from the WLAN interface.
	fmt.Println()

	packet, err := UnmarshalPacket(data)
	if err != nil {
		log.Printf("Failed to unmarshal data: %v\n", err)
		return nil, nil, err
	}

	log.Printf("Received \"%s\" from: [%s] on intf-%d", string(packet.Payload), addr.String(), RSUWLANInterface)

	return packet, addr, err
}

// get MAC from RSUARP using ip then send to wlan interface
func (r *Router) SendToWLANInterface(data []byte, ip string) {
	mac := r.RARP.Get(ip)
	r.wlanChannel.SendTo(data, mac)
}

// Send to all in RSUARP using wlan exept the one that sent the message
func (r *Router) SendToALLWLANInterface(data []byte, originatorIP string) int {
	sentPackets := 0
	for eip, entry := range r.RARP.table {
		if originatorIP == eip {
			continue
		}
		log.Printf("Sending to: %s  mac: %s", eip, entry.MAC)
		packet := ip.NewIPPacket(data, r.ip, net.ParseIP(eip))
		bytes := ip.MarshalIPPacket(packet)
		ip.UpdateChecksum(bytes)
		r.wlanChannel.SendTo(bytes, entry.MAC)
		sentPackets++
	}
	return sentPackets
}

func (r *Router) BroadcastETH(packet *IPPacket) {
	bytes := ip.MarshalIPPacket(packet)
	ip.Update(bytes)
	r.ethChannel.Broadcast(bytes)
}

func (r *Router) SendIPToNeighbors() {
	r.neighborChannel.Broadcast(r.ip.To4())
}

func (r *Router) Close() {
	r.ethChannel.Close()
	r.wlanChannel.Close()
}
