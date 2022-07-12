package datalink

import (
	"fmt"
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
	VNeighborTable_UPDATE_INTERVAL = 6
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

func (nt *VNeighborTable) Set(MAC string, neighbor *VNeighborEntry) {
	if neighbor == nil {
		log.Panic("You are trying to add null neighbor")
	}

	entry, exist := nt.Get(MAC)
	if exist {
		entry.timer.Stop()
		nt.table.Del(MAC)
	}

	callback := func() {
		nt.table.Del(MAC)
	}

	timer := time.AfterFunc(VNeighborTable_UPDATE_INTERVAL*time.Second, callback)
	neighbor.timer = timer

	nt.table.Set(MAC, neighbor)
}

func (nt *VNeighborTable) Get(MAC string) (*VNeighborEntry, bool) {
	neighbor, exist := nt.table.Get(MAC)
	if !exist {
		return nil, false
	}
	return neighbor.(*VNeighborEntry), true
}

func (nt *VNeighborTable) GetFirst() (*VNeighborEntry, bool) {
	for item := range nt.table.Iter() {
		return item.Value.(*VNeighborEntry), true
	}
	return nil, false
}

func (nt *VNeighborTable) Len() int {
	return nt.table.Len()
}

func (nt *VNeighborTable) Print() {

	for item := range nt.table.Iter() {
		itemMAC := item.Key.(string)
		itemIP := item.Value.(*VNeighborEntry)
		fmt.Printf("key: %s, Value %s", itemMAC, string(itemIP.IP))
	}
}

func (nt *VNeighborTable) Iter() <-chan *VNeighborEntry {
	ch := make(chan *VNeighborEntry)
	go func() {
		for item := range nt.table.Iter() {
			value := item.Value.(*VNeighborEntry)
			ch <- value
		}
		close(ch)
	}()
	return ch
}
func (nt *VNeighborTable) Update() {
	for {
		payload, addr, err := nt.channel.Read()
		entry := NewVNeighborEntry(net.ParseIP(string(payload)), addr)
		nt.Set(addr.String(), entry)
		if err != nil {
			return
		}
	}
}

func (nt *VNeighborTable) Run() {
	go nt.Update()
	if nt.receiveOnly {
		return
	}
	for {
		nt.channel.Broadcast([]byte(nt.SrcIP.String()))
		time.Sleep((VNeighborTable_UPDATE_INTERVAL / 3) * time.Second)
	}
}

func (nt *VNeighborTable) Close() {
	nt.channel.Close()
}
