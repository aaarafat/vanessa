package aodv

import (
	"log"
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
		if neighborMac.String() != addr.String() {
			f.ForwardTo(payload, neighborMac, 1)
		}
	}
}

func (f *Forwarder) ForwardTo(payload []byte, addr net.HardwareAddr, index int) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channels[index].SendTo(payload, addr)
	log.Printf("Forwarding to %s interface %d\n", addr.String(), index)
}

func (f *Forwarder) ForwardToAll(payload []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channels[1].Broadcast(payload)
	log.Printf("Broadcasting to interface 1\n")
}

func (f *Forwarder) Start() {
	go f.neighborsTable.Run()
}

func (f *Forwarder) Close() {
	f.neighborsTable.Close()
}