package aodv

import (
	"net"
	"sync"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

type Forwarder struct {
	neighborsTable *VNeighborTable
	channels []*DataLinkLayerChannel
	lock *sync.Mutex
}


func NewForwarder(srcIP net.IP, channels []*DataLinkLayerChannel) *Forwarder {
	return &Forwarder{
		neighborsTable: NewNeighborTable(srcIP),
		channels: channels,
		lock: &sync.Mutex{},
	}
}

func (f *Forwarder) ForwardToAllExcept(payload []byte, addr net.HardwareAddr) {
	for item := range f.neighborsTable.Iter() {
		neighborMac := item.MAC
		neighborIfiIndex := item.IfiIndex
		if neighborMac.String() != addr.String() {
			f.ForwardTo(payload, neighborMac, neighborIfiIndex)
		}
	}
}

func (f *Forwarder) ForwardTo(payload []byte, addr net.HardwareAddr, index int) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channels[index].SendTo(payload, addr)
}

func (f *Forwarder) ForwardToAll(payload []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	for _, channel := range f.channels {
		channel.Broadcast(payload)
	}
}

func (f *Forwarder) Start() {
	go f.neighborsTable.Run()
}