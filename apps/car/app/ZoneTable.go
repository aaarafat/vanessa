package app

import (
	"log"
	"net"
	"time"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	"github.com/cornelk/hashmap"
)

type ZoneTable struct {
	table *hashmap.HashMap
}

type ZoneTableEntry struct {
	IP        net.IP
	Speed     uint32
	Position  Position
	Direction Position
	timer     *time.Timer
}

const (
	ZoneTable_UPDATE_INTERVAL_MS = ZONE_MSG_INTERVAL_MS * 3
)

func NewZoneTable() *ZoneTable {
	return &ZoneTable{
		table: &hashmap.HashMap{},
	}
}

func NewZoneTableEntry(ip net.IP, speed uint32, pos Position) *ZoneTableEntry {
	return &ZoneTableEntry{
		IP:        ip,
		Speed:     speed,
		Position:  pos,
		Direction: pos,
		timer:     nil,
	}
}

func (zt *ZoneTable) Set(ip net.IP, speed uint32, pos Position) {
	entry, exists := zt.Get(ip)
	if exists {
		entry.timer.Reset(ZoneTable_UPDATE_INTERVAL_MS * time.Millisecond)

		// get the unit vector
		direction := Position{Lat: entry.Position.Lat - pos.Lat, Lng: entry.Position.Lng - pos.Lng}
		direction.Normalize()

		// update
		entry.Direction = direction
		entry.Position = pos
		entry.Speed = speed

		log.Printf("ZoneTable: updated entry for %s\n", ip.String())
	} else {
		entry := NewZoneTableEntry(ip, speed, pos)
		entry.timer = time.AfterFunc(ZoneTable_UPDATE_INTERVAL_MS*time.Millisecond, func() {
			log.Printf("ZoneTable: entry expired: %s\n", ip.String())
			zt.table.Del(ip)
		})

		zt.table.Set(ip, *entry)

		log.Printf("ZoneTable: added entry for %s\n", ip.String())
	}
}

func (zt *ZoneTable) Get(ip net.IP) (*ZoneTableEntry, bool) {
	entry, exists := zt.table.Get(ip)
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
