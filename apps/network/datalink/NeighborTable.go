package datalink

import (
	"log"
	"net"
	"time"

	"github.com/cornelk/hashmap"
)

type VNeighborTable struct {
	table       *hashmap.HashMap
	channel     *DataLinkLayerChannel
	SrcIP       net.IP
	receiveOnly bool
}

type VNeighborEntry struct {
	MAC   net.HardwareAddr
	IP    net.IP
	timer *time.Timer
}

const (
	VNeighborTable_UPDATE_INTERVAL_MS = 6000
)

func NewVNeighborTable(srcIP net.IP, ifiName string, receiveOnly bool) *VNeighborTable {
	d, err := NewDataLinkLayerChannelWithInterfaceName(VNDEtherType, ifiName)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &VNeighborTable{
		table:       &hashmap.HashMap{},
		channel:     d,
		SrcIP:       srcIP,
		receiveOnly: receiveOnly,
	}
}
func NewVNeighborEntry(ip net.IP, mac net.HardwareAddr) *VNeighborEntry {
	return &VNeighborEntry{
		IP:  ip,
		MAC: mac,
	}
}

func (nt *VNeighborTable) Set(MAC string, neighbor VNeighborEntry) {
	entry, exist := nt.Get(MAC)
	if exist {
		entry.timer.Stop()
		nt.table.Del(MAC)
	}

	callback := func() {
		nt.table.Del(MAC)
	}

	timer := time.AfterFunc(VNeighborTable_UPDATE_INTERVAL_MS*time.Millisecond, callback)
	neighbor.timer = timer

	nt.table.Set(MAC, neighbor)
}

func (nt *VNeighborTable) Get(MAC string) (*VNeighborEntry, bool) {
	neighbor, exist := nt.table.Get(MAC)
	if !exist {
		return nil, false
	}
	entry := neighbor.(VNeighborEntry)
	return &entry, true
}

func (nt *VNeighborTable) GetFirst() (*VNeighborEntry, bool) {
	for item := range nt.table.Iter() {
		entry := item.Value.(VNeighborEntry)
		return &entry, true
	}
	return nil, false
}

func (nt *VNeighborTable) Len() int {
	return nt.table.Len()
}

func (nt *VNeighborTable) Print() {

	for item := range nt.table.Iter() {
		itemMAC := item.Key.(string)
		itemIP := item.Value.(VNeighborEntry).IP
		log.Printf("key: %s, Value %s", itemMAC, itemIP.String())
	}
}

func (nt *VNeighborTable) Iter() <-chan VNeighborEntry {
	ch := make(chan VNeighborEntry)
	go func() {
		for item := range nt.table.Iter() {
			value := item.Value.(VNeighborEntry)
			ch <- value
		}
		close(ch)
	}()
	return ch
}
func (nt *VNeighborTable) Update() {
	for {
		payload, addr, err := nt.channel.Read()
		if err != nil {
			continue
		}
		entry := NewVNeighborEntry(net.IPv4(payload[0], payload[1], payload[2], payload[3]), addr)
		nt.Set(addr.String(), *entry)
	}
}

func (nt *VNeighborTable) Run() {
	go nt.Update()
	if nt.receiveOnly {
		return
	}
	for {
		nt.channel.Broadcast(nt.SrcIP.To4())
		time.Sleep((VNeighborTable_UPDATE_INTERVAL_MS / 3) * time.Second)
	}
}

func (nt *VNeighborTable) Close() {
	nt.channel.Close()
}
