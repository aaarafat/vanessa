package aodv

import (
	"net"
	"sync"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

type Forwarder struct {
	neighborsTable *VNeighborTable
	channel *DataLinkLayerChannel
	lock *sync.Mutex
}


func NewForwarder(srcIP net.IP, channel *DataLinkLayerChannel) *Forwarder {
	return &Forwarder{
		neighborsTable: NewNeighborTable(srcIP),
		channel: channel,
		lock: &sync.Mutex{},
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
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channel.SendTo(payload, addr)
}

func (f *Forwarder) ForwardToAll(payload []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channel.Broadcast(payload)
}

func (f *Forwarder) Start() {
	go f.neighborsTable.Run()
}