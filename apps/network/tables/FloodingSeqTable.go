package tables

import (
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
	DefaultSeqEntryAge = 10 // seconds
)

func NewFloodingSeqTable() *VFloodingSeqTable {
	return &VFloodingSeqTable{
		table: &hashmap.HashMap{},
	}
}

func (f *VFloodingSeqTable) Exists(destination net.IP) bool {
	_, exists := f.get(destination)
	return exists
}

func (f *VFloodingSeqTable) get(destination net.IP) (*VSeqEntry, bool) {
	item, exists := f.table.Get(destination.String())
	if exists {
		entry := item.(VSeqEntry)
		return &entry, true
	}
	return nil, false
}


func (f *VFloodingSeqTable) Set(destination net.IP, seq uint32) {
	item, exists := f.get(destination)
	if exists {
		timer := item.age
		timer.Stop()
	}

	callback := func() {
		f.table.Del(destination.String())
	}
	timer := time.AfterFunc(time.Second * DefaultSeqEntryAge, callback)
	entry := VSeqEntry{
		seqNum: seq,
		age: timer,
	}
	f.table.Set(destination.String(), entry)
}
