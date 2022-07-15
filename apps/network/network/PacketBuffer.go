package network

import (
	"net"
	"sync"
	"time"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/cornelk/hashmap"
)

type PacketBuffer struct {
	buf  *hashmap.HashMap
	lock *sync.RWMutex
}

type PacketBufferEntry struct {
	Packet []byte
	DestIP net.IP
	timer  *time.Timer
}

func NewPacketBuffer() *PacketBuffer {
	return &PacketBuffer{
		buf:  &hashmap.HashMap{},
		lock: &sync.RWMutex{},
	}
}

func (p *PacketBuffer) Add(packet []byte, destIP net.IP) {
	p.lock.Lock()
	defer p.lock.Unlock()
	m, ok := p.buf.Get(destIP.String())
	if !ok {
		m = &hashmap.HashMap{}
	}

	callback := func() {
		m.(*hashmap.HashMap).Del(destIP.String())
	}
	timer := time.AfterFunc(ip.PacketTimeoutMS, callback)
	m.(*hashmap.HashMap).Set(destIP.String(), PacketBufferEntry{
		Packet: packet,
		DestIP: destIP,
		timer:  timer,
	})

	p.buf.Set(destIP.String(), m)
}

func (p *PacketBuffer) Get(destIP net.IP) ([][]byte, bool) {
	p.lock.RLock()
	defer p.lock.RUnlock()
	m, ok := p.buf.Get(destIP.String())
	if !ok {
		return nil, false
	}

	data := make([][]byte, m.(*hashmap.HashMap).Len())
	i := 0
	for item := range m.(*hashmap.HashMap).Iter() {
		value := item.Value.(PacketBufferEntry)
		data[i] = value.Packet
		i++
	}

	return data, true
}

func (p *PacketBuffer) Del(destIP net.IP) {
	p.lock.Lock()
	defer p.lock.Unlock()
	m, ok := p.buf.Get(destIP.String())
	if ok {
		// Stop the timer
		for item := range m.(*hashmap.HashMap).Iter() {
			value := item.Value.(PacketBufferEntry)
			value.timer.Stop()
		}
	}
	p.buf.Del(destIP.String())
}
