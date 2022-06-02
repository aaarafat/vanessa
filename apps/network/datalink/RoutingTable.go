package datalink

import (
	"fmt"
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/scripts/utils"
	"github.com/cornelk/hashmap"
)

type VNeighborTable struct {
	table *hashmap.HashMap
	entryTypes map[string]TypeEnum
}


type VNeighborEntry struct {
	// MAC  net.HardwareAddr
	IP net.IP
}

func NewNeighborTable() *VNeighborTable {
	return &VNeighborTable{
		table: &hashmap.HashMap{},
		entryTypes: map[string]TypeEnum{
			"MAC": String,
			"IP": ByteArray,
		},
	}
}
func NewNeighborEntry(ip net.IP) *VNeighborEntry {
	return &VNeighborEntry{
		IP: ip,
	}
}

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
