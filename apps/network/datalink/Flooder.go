package datalink

import (
	"log"
	"net"
	"sync"
)

type IFlooder interface {
	ForwardTo(payload []byte, addr net.HardwareAddr)
	ForwardToAll(payload []byte)
	ForwardToAllExcept(payload []byte, addr net.HardwareAddr)
	ForwardToAllExceptIP(payload []byte, ip net.IP)
}

type Flooder struct {
	neighborsTable *VNeighborTable
	channel        *DataLinkLayerChannel
	lock           *sync.Mutex
}

func NewFlooder(srcIP net.IP, channel *DataLinkLayerChannel, nt *VNeighborTable) *Flooder {

	return &Flooder{
		neighborsTable: nt,
		channel:        channel,
		lock:           &sync.Mutex{},
	}
}

func (f *Flooder) ForwardToAllExcept(payload []byte, addr net.HardwareAddr) {
	for item := range f.neighborsTable.Iter() {
		neighborMac := item.MAC
		if neighborMac.String() != addr.String() {
			f.ForwardTo(payload, neighborMac)
		}
	}
}

func (f *Flooder) ForwardToAllExceptIP(payload []byte, ip net.IP) {
	for item := range f.neighborsTable.Iter() {
		neighborMac := item.MAC
		neighborIP := item.IP
		if !neighborIP.Equal(ip) {
			f.ForwardTo(payload, neighborMac)
		}
	}
}

func (f *Flooder) ForwardTo(payload []byte, addr net.HardwareAddr) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channel.SendTo(payload, addr)
	log.Printf("Forwarding to %s interface %s\n", addr.String(), f.channel.Ifi.Name)
}

func (f *Flooder) ForwardToAll(payload []byte) {
	f.lock.Lock()
	defer f.lock.Unlock()
	f.channel.Broadcast(payload)
	log.Printf("Broadcasting to interface %s\n", f.channel.Ifi.Name)
}
