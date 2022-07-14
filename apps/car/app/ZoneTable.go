package app

import (
	"log"
	"net"
	"time"

	. "github.com/aaarafat/vanessa/libs/vector"
	"github.com/cornelk/hashmap"
)

type ZoneTable struct {
	table *hashmap.HashMap
}

type ZoneTableEntry struct {
	IP              net.IP
	Speed           uint32
	Position        Position
	Direction       Vector  // pos2 - pos1
	DirectionFromMe Vector  // how I'm facing the other vehicle (targetPos - myPos)
	Angle           float64 // angle between Direction and DirectionFromMe
	timer           *time.Timer
}

const (
	ZoneTable_UPDATE_INTERVAL_MS = ZONE_MSG_INTERVAL_MS * 3
)

func NewZoneTable() *ZoneTable {
	return &ZoneTable{
		table: &hashmap.HashMap{},
	}
}

func NewZoneTableEntry(ip net.IP, speed uint32, pos, myPos Position) *ZoneTableEntry {
	entry := &ZoneTableEntry{
		IP:              ip,
		Speed:           speed,
		Position:        pos,
		Direction:       NewUnitVector(pos, pos),
		DirectionFromMe: NewUnitVector(myPos, pos),
		timer:           nil,
	}

	entry.Angle = entry.Direction.Angle(entry.DirectionFromMe)

	return entry
}

func (zt *ZoneTable) Set(ip net.IP, speed uint32, pos, myPos Position) {
	entry, exists := zt.Get(ip)
	if exists {
		log.Printf("ZoneTable: entry already exists for %s\n", ip.String())
		entry.timer.Reset(ZoneTable_UPDATE_INTERVAL_MS * time.Millisecond)

		// update
		entry.Direction = NewUnitVector(entry.Position, pos)
		entry.DirectionFromMe = NewUnitVector(myPos, pos)
		entry.Angle = entry.Direction.Angle(entry.DirectionFromMe)
		entry.Position = pos
		entry.Speed = speed

		zt.table.Set(ip.String(), *entry)
		log.Printf("ZoneTable: updated entry for %s\n", ip.String())
	} else {
		log.Printf("ZoneTable: adding entry for %s\n", ip.String())
		entry := NewZoneTableEntry(ip, speed, pos, myPos)
		entry.timer = time.AfterFunc(ZoneTable_UPDATE_INTERVAL_MS*time.Millisecond, func() {
			log.Printf("ZoneTable: entry expired: %s\n", ip.String())
			zt.table.Del(ip.String())
		})

		zt.table.Set(ip.String(), *entry)

		log.Printf("ZoneTable: added entry for %s\n", ip.String())
	}
}

func (zt *ZoneTable) Get(ip net.IP) (*ZoneTableEntry, bool) {
	entry, exists := zt.table.Get(ip.String())
	if !exists {
		return nil, false
	}
	ztEntry := entry.(ZoneTableEntry)
	return &ztEntry, exists
}

func (zte *ZoneTableEntry) Print() {
	log.Printf("IP: %s, Speed: %d, Position: %v, Direction: %v\n", zte.IP.String(), zte.Speed, zte.Position, zte.Direction)
}

func (zt *ZoneTable) Print() {
	log.Printf("ZoneTable: %d entries\n\n", zt.table.Len())
	for item := range zt.table.Iter() {
		itemEntry := item.Value.(ZoneTableEntry)
		itemEntry.Print()
	}
}
