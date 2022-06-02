package datalink

import (
	"fmt"
	"log"
	"net"

	"github.com/cornelk/hashmap"
)

type VNeighborTable struct {
	table *hashmap.HashMap
}


type VNeighborEntry struct {
	// MAC  net.HardwareAddr
	IP net.IP
}

func NewNeighborTable() *VNeighborTable {
	return &VNeighborTable{
		table: &hashmap.HashMap{},
	}
}
func NewNeighborEntry(ip net.IP) *VNeighborEntry {
	return &VNeighborEntry{
		IP: ip,
	}
}
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
