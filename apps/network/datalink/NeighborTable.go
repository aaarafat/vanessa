package datalink

import (
	"fmt"
	"log"
	"net"
	"time"

	"github.com/cornelk/hashmap"
)

type VNeighborTable struct {
	table *hashmap.HashMap
	channel *DataLinkLayerChannel
	SrcIP net.IP
}

type VNeighborEntry struct {
	MAC  net.HardwareAddr
	IP net.IP
}

const (
	VNeighborTable_UPDATE_INTERVAL = 5
)

func NewNeighborTable(srcIP  net.IP) *VNeighborTable {
	d, err := NewDataLinkLayerChannel(VNDEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	return &VNeighborTable{
		table: &hashmap.HashMap{},
		channel: d,
		SrcIP: srcIP,
	}
}
func NewNeighborEntry(ip net.IP, mac net.HardwareAddr) *VNeighborEntry {
	return &VNeighborEntry{
		IP: ip,
		MAC: mac,
	}
}

/*
func (nt *VNeighborTable) MarshalBinary() []byte {

	var payload []byte
	for item := range nt.table.Iter() {
		itemMAC := item.Key.(string)
		itemIP := item.Value.(*VNeighborEntry)

		data := map[string]any{
			"MAC": itemMAC,
			"IP": itemIP.IP,
		}

		b, err := Marshal(data, nt.entryTypes)
		if err != nil {
			log.Panic(err)
		}
		payload = append(payload, b...)
	}

	return payload
}

func (nt *VNeighborTable) UnmarshalBinary(data []byte) {
	for len(data) > 0 {
		item, err := Unmarshal(data, nt.entryTypes)
		if err != nil {
			log.Panic(err)
		}

		itemMAC := item["MAC"].(string)
		itemIP := item["IP"].(net.IP)

		nt.Set(itemMAC, NewNeighborEntry(itemIP))
	}
}
*/

func (nt *VNeighborTable) Set(MAC string, neighbor *VNeighborEntry) {
	if neighbor == nil {
		log.Panic("You are trying to add null neighbor")
	}
	
	nt.table.Set(MAC, neighbor)
}

func (nt *VNeighborTable) Get(MAC string) (*VNeighborEntry, bool) {
	neighbor, exist := nt.table.Get(MAC)
	if !exist {
		log.Panic("Neighbor doesn't exist")
		return nil, false
	}
	return neighbor.(*VNeighborEntry), true
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
		entry := NewNeighborEntry(net.IP(payload), addr)
		nt.Set(addr.String() ,entry)
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
	}
}

func (nt *VNeighborTable) Run() {
	go nt.Update()
	for {
		nt.channel.Broadcast([]byte(nt.SrcIP))
		time.Sleep(VNeighborTable_UPDATE_INTERVAL * time.Second)
	}
}