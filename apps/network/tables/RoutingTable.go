package tables

import (
	"fmt"
	"net"

	"github.com/cornelk/hashmap"
)

type VRoutingTable struct {
	table *hashmap.HashMap
}

type VRoutingTableEntry struct {
	Destination net.IP
	NextHop net.HardwareAddr
	NoOfHops int
}

func NewVRoutingTable() *VRoutingTable {
	return &VRoutingTable{
		table: &hashmap.HashMap{},
	}
}

func (r* VRoutingTable) Update(nextHopMac net.HardwareAddr, destination net.IP, hopCount int) {
	// check if entry exists
	entry, exists := r.table.Get(destination.String())
	if exists {
		// update entry if hop count is less than current entry
		if hopCount < entry.(VRoutingTableEntry).NoOfHops {
			r.table.Set(destination.String(), VRoutingTableEntry{
				Destination: destination,
				NextHop: nextHopMac,
				NoOfHops: hopCount,
			})
		}
	} else {
		r.table.Set(destination.String(), VRoutingTableEntry{
			Destination: destination,
			NextHop: nextHopMac,
			NoOfHops: hopCount,
		})
	}
	r.Print()
}

func (r* VRoutingTable) Get(destination net.IP) (*VRoutingTableEntry, bool) {
	item, exists := r.table.Get(destination)
	if exists {
		return item.(*VRoutingTableEntry), true
	}
	return nil, false
}

func (r* VRoutingTable) Len() int {
	return r.table.Len()
}

func (r *VRoutingTable) Print() {
	
	for item := range r.table.Iter() {
		itemEntry := item.Value.(VRoutingTableEntry)
		fmt.Printf("ip: %s, next hop %s,  no_hops  %d\n",itemEntry.Destination.String(), itemEntry.NextHop.String(), itemEntry.NoOfHops)
	}
}