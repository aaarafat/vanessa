package aodv

import (
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

type Forwarder struct {
	neighborsTable *VNeighborTable
	channel *DataLinkLayerChannel
}


func NewForwarder(srcIP net.IP, channel *DataLinkLayerChannel) *Forwarder {
	return &Forwarder{
		neighborsTable: NewNeighborTable(srcIP),
		channel: channel,
	}
}

func (f *Forwarder) ForwardToAllExcept(payload []byte, addr net.HardwareAddr) {
	for item := range f.neighborsTable.Iter() {
		neighborMac := item.MAC
		if neighborMac.String() != addr.String() {
			f.ForwardTo(payload, neighborMac)
		}
	}
}

func (f *Forwarder) ForwardTo(payload []byte, addr net.HardwareAddr) {
	f.channel.SendTo(payload, addr)
}

func (f *Forwarder) ForwardToAll(payload []byte) {
	f.channel.Broadcast(payload)
}

func (f *Forwarder) Start() {
	go f.neighborsTable.Run()
}