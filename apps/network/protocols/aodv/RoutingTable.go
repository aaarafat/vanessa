package aodv

import (
	"log"
	"net"
	"time"

	"github.com/cornelk/hashmap"
)

type VRoutingTable struct {
	table *hashmap.HashMap
}

type VRoutingTableEntry struct {
	Destination net.IP
	NextHop net.HardwareAddr
	NoOfHops uint8
	SeqNum uint32
	IfiIndex int
	LifeTime time.Time
	timer *time.Timer
}

func NewVRoutingTable() *VRoutingTable {
	return &VRoutingTable{
		table: &hashmap.HashMap{},
	}
}

func (r* VRoutingTable) isNewEntry(newEntry *VRoutingTableEntry) bool {
	// https://datatracker.ietf.org/doc/html/rfc3561#section-6.2
	entry, exists := r.Get(newEntry.Destination);
	if exists {
		if entry.LifeTime.Before(time.Now()) {
			return true
		}
		if entry.SeqNum < newEntry.SeqNum {
			return true
		}
		if entry.SeqNum == newEntry.SeqNum && entry.NoOfHops > newEntry.NoOfHops {
			return true
		}
		if entry.SeqNum == newEntry.SeqNum && entry.NoOfHops == newEntry.NoOfHops && entry.LifeTime.Before(newEntry.LifeTime) {
			return true
		}
		return false
	}

	return true
}

func (r* VRoutingTable) Update(destIP net.IP, nextHop net.HardwareAddr, hopCount uint8, lifeTime, seqNum uint32, ifiIndex int) {
	lifeTimeMS := time.Millisecond * time.Duration(lifeTime)
		
	newEntry := &VRoutingTableEntry{
		Destination: destIP,
		NextHop: nextHop,
		NoOfHops: hopCount,
		SeqNum: seqNum,
		LifeTime: time.Now().Add(lifeTimeMS),
		IfiIndex: ifiIndex,
	}

	if r.isNewEntry(newEntry) {
		r.set(newEntry)
	}
}

func (r* VRoutingTable) Get(destination net.IP) (*VRoutingTableEntry, bool) {
	item, exists := r.table.Get(destination.String())
	if exists {
		entry := item.(VRoutingTableEntry)
		return &entry, true
	}
	return nil, false
}

func (r* VRoutingTable) Del(ip net.IP) {
	item, exists := r.table.Get(ip.String())
	if exists {
		// Stop the timer
		itemEntry := item.(VRoutingTableEntry)
		itemEntry.timer.Stop()
		r.table.Del(ip.String())
	}
}

func (r* VRoutingTable) Set(destIP net.IP, nextHop net.HardwareAddr, hopCount uint8, lifeTime, seqNum uint32, ifiIndex int) {
	lifeTimeMS := time.Millisecond * time.Duration(lifeTime)
		
	newEntry := &VRoutingTableEntry{
		Destination: destIP,
		NextHop: nextHop,
		NoOfHops: hopCount,
		SeqNum: seqNum,
		LifeTime: time.Now().Add(lifeTimeMS),
		IfiIndex: ifiIndex,
	}

	r.set(newEntry)
}

func (r* VRoutingTable) set(entry *VRoutingTableEntry) {
	item, exists := r.table.Get(entry.Destination.String())
	if exists {
		// Stop the timer
		itemEntry := item.(VRoutingTableEntry)
		itemEntry.timer.Stop()
		r.table.Del(entry.Destination.String())
	}

	callback := func() {
		r.table.Del(entry.Destination.String())
	}
	timer := time.AfterFunc(entry.LifeTime.Sub(time.Now()), callback)
	entry.timer = timer
	
	r.table.Set(entry.Destination.String(), *entry)
}

func (r* VRoutingTable) Len() int {
	return r.table.Len()
}

func (r *VRoutingTable) Print() {
	log.Println()
	log.Println("Routing Table:")
	for item := range r.table.Iter() {
		itemEntry := item.Value.(VRoutingTableEntry)
		log.Printf("ip: %s, next hop %s, no hops %d, seq num %d, life time %s\n",
		itemEntry.Destination.String(), itemEntry.NextHop.String(), 
		itemEntry.NoOfHops, itemEntry.SeqNum, itemEntry.LifeTime.String())
	}
	log.Println()
}