package aodv

import (
	"fmt"
	"net"
	"time"

	"github.com/cornelk/hashmap"
)

type VFloodingSeqTable struct {
	table *hashmap.HashMap
}

type VSeqEntry struct {
	seqNum uint32
	age *time.Timer
}

const (
	DefaultSeqEntryAge = 5 // seconds
)

func NewVFloodingSeqTable() *VFloodingSeqTable {
	return &VFloodingSeqTable{
		table: &hashmap.HashMap{},
	}
}

func (f *VFloodingSeqTable) Exists(destination net.IP, seq uint32) bool {
	_, exists := f.get(destination, seq)
	return exists
}

func (f *VFloodingSeqTable) key(destination net.IP, seq uint32) string {
	return destination.String() + ":" + fmt.Sprint(seq)
}

func (f *VFloodingSeqTable) get(destination net.IP, seq uint32) (*VSeqEntry, bool) {
	item, exists := f.table.Get(f.key(destination, seq)) 
	if exists {
		entry := item.(VSeqEntry)
		return &entry, true
	}
	return nil, false
}


func (f *VFloodingSeqTable) Set(destination net.IP, seq uint32) {
	item, exists := f.get(destination, seq)
	if exists {
		timer := item.age
		timer.Stop()
	}

	callback := func() {
		f.table.Del(f.key(destination, seq))
	}
	timer := time.AfterFunc(time.Second * DefaultSeqEntryAge, callback)
	entry := VSeqEntry{
		seqNum: seq,
		age: timer,
	}
	f.table.Set(f.key(destination, seq), entry)
}
