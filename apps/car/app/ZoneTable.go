package app

import (
	"log"
	"math"
	"net"
	"sync"
	"time"

	. "github.com/aaarafat/vanessa/libs/vector"
	"github.com/cornelk/hashmap"
)

type ZoneTable struct {
	table *hashmap.HashMap

	nearestFrontIP   net.IP
	nearestFrontLock *sync.RWMutex

	nearestBehindIP   net.IP
	nearestBehindLock *sync.RWMutex
}

type ZoneTableEntry struct {
	IP              net.IP
	Speed           uint32
	Position        Position
	Direction       Vector  // pos2 - pos1
	DirectionFromMe Vector  // how I'm facing the other vehicle (targetPos - myPos)
	Angle           float64 // angle between MyDirection and DirectionFromMe
	Ignored         bool
	timer           *time.Timer
}

const (
	ZoneTable_UPDATE_INTERVAL_MS = ZONE_MSG_INTERVAL_MS * 3
	MAX_ANGLE_DEG                = 10
)

type NEAREST_BOOL int

const (
	NEAREST  NEAREST_BOOL = -1
	FARTHEST NEAREST_BOOL = 1
	EQUAL    NEAREST_BOOL = 0
)

func NewZoneTable() *ZoneTable {
	return &ZoneTable{
		table:             &hashmap.HashMap{},
		nearestFrontLock:  &sync.RWMutex{},
		nearestBehindLock: &sync.RWMutex{},
	}
}

func NewZoneTableEntry(ip net.IP, speed uint32, pos, myPos Position, myDir Vector) *ZoneTableEntry {
	entry := &ZoneTableEntry{
		IP:              ip,
		Speed:           speed,
		Position:        pos,
		Direction:       NewUnitVector(pos, pos),
		DirectionFromMe: NewUnitVector(myPos, pos),
		Ignored:         false,
		timer:           nil,
	}

	entry.Angle = myDir.Angle(entry.DirectionFromMe)

	return entry
}

func (zt *ZoneTable) Set(ip net.IP, speed uint32, pos, myPos Position, myDir Vector) *ZoneTableEntry {
	entry, exists := zt.Get(ip)
	if exists {
		log.Printf("ZoneTable: entry already exists for %s\n", ip.String())
		entry.timer.Reset(ZoneTable_UPDATE_INTERVAL_MS * time.Millisecond)

		// update
		entry.Direction = NewUnitVector(entry.Position, pos)
		entry.DirectionFromMe = NewUnitVector(myPos, pos)
		entry.Angle = myDir.Angle(entry.DirectionFromMe)
		entry.Position = pos
		entry.Speed = speed

		zt.table.Set(ip.String(), *entry)
		zt.updateNearest(entry, &myPos)

		log.Printf("ZoneTable: updated entry for %s\n", ip.String())
		return entry
	} else {
		log.Printf("ZoneTable: adding entry for %s\n", ip.String())
		entry := NewZoneTableEntry(ip, speed, pos, myPos, myDir)
		entry.timer = time.AfterFunc(ZoneTable_UPDATE_INTERVAL_MS*time.Millisecond, func() {
			log.Printf("ZoneTable: entry expired: %s\n", ip.String())
			zt.table.Del(ip.String())
			zt.removeNearestIfEqual(entry)
		})

		zt.table.Set(ip.String(), *entry)
		zt.updateNearest(entry, &myPos)

		log.Printf("ZoneTable: added entry for %s\n", ip.String())
		return entry
	}
}

func (zt *ZoneTable) Get(ip net.IP) (*ZoneTableEntry, bool) {
	entry, exists := zt.table.GetStringKey(ip.String())
	if !exists {
		return nil, false
	}
	ztEntry := entry.(ZoneTableEntry)
	return &ztEntry, exists
}

func (zt *ZoneTable) IsFront(entry *ZoneTableEntry) bool {
	return !entry.Ignored && math.Abs(entry.Angle) <= ToRadians(MAX_ANGLE_DEG)
}

func (zt *ZoneTable) IsBehind(entry *ZoneTableEntry) bool {
	return !entry.Ignored && math.Abs(entry.Angle) >= ToRadians(180-MAX_ANGLE_DEG)
}

func (zt *ZoneTable) GetInFrontOfMe() []*ZoneTableEntry {
	var entries []*ZoneTableEntry
	for item := range zt.table.Iter() {
		itemEntry := item.Value.(ZoneTableEntry)
		if zt.IsFront(&itemEntry) {
			entries = append(entries, &itemEntry)
		}
	}
	return entries
}

func (zt *ZoneTable) GetBehindMe() []*ZoneTableEntry {
	var entries []*ZoneTableEntry
	for item := range zt.table.Iter() {
		itemEntry := item.Value.(ZoneTableEntry)
		if zt.IsBehind(&itemEntry) {
			entries = append(entries, &itemEntry)
		}
	}
	return entries
}

func (zt *ZoneTable) GetNearestFrontFrom(pos *Position) *ZoneTableEntry {
	// if stored return it
	zt.nearestFrontLock.RLock()
	if zt.nearestFrontIP != nil {
		entry, exists := zt.Get(zt.nearestFrontIP)
		if exists && zt.IsFront(entry) {
			zt.nearestFrontLock.RUnlock()
			return entry
		}
	}
	zt.nearestFrontLock.RUnlock()

	var nearest *ZoneTableEntry
	var nearestDist float64 = math.MaxFloat64
	for _, entry := range zt.GetInFrontOfMe() {
		dist := entry.Position.Distance(pos)
		if dist < nearestDist {
			nearest = entry
			nearestDist = dist
		}
	}

	// update
	zt.nearestFrontLock.Lock()
	zt.nearestFrontIP = nil
	if nearest != nil {
		zt.nearestFrontIP = nearest.IP
	}
	zt.nearestFrontLock.Unlock()
	return nearest
}

func (zt *ZoneTable) GetNearestBehindFrom(pos *Position) *ZoneTableEntry {
	// if stored return it
	zt.nearestBehindLock.RLock()
	if zt.nearestBehindIP != nil {
		entry, exists := zt.Get(zt.nearestBehindIP)
		if exists && zt.IsBehind(entry) {
			zt.nearestBehindLock.RUnlock()
			return entry
		}
	}
	zt.nearestBehindLock.RUnlock()

	var nearest *ZoneTableEntry
	var nearestDist float64 = math.MaxFloat64
	for _, entry := range zt.GetBehindMe() {
		dist := entry.Position.Distance(pos)
		if dist < nearestDist {
			nearest = entry
			nearestDist = dist
		}
	}

	// update
	zt.nearestBehindLock.Lock()
	zt.nearestBehindIP = nil
	if nearest != nil {
		zt.nearestBehindIP = nearest.IP
	}
	zt.nearestBehindLock.Unlock()
	return nearest
}

func (zt *ZoneTable) Ignore(ip net.IP) {
	entry, exists := zt.Get(ip)
	if exists {
		entry.Ignored = true
		zt.table.Set(ip.String(), *entry)
		zt.removeNearestIfEqual(entry)
		log.Printf("ZoneTable: ignored entry for %s\n", ip.String())
	}
}

func (zt *ZoneTable) updateNearest(entry *ZoneTableEntry, pos *Position) {
	zt.updateNearestFront(entry, pos)
	zt.updateNearestBehind(entry, pos)
}

func (zt *ZoneTable) checkNearestBehind(entry *ZoneTableEntry, pos *Position) NEAREST_BOOL {
	if !zt.IsBehind(entry) {
		return FARTHEST
	}
	zt.nearestBehindLock.RLock()
	defer zt.nearestBehindLock.RUnlock()
	if zt.nearestBehindIP == nil {
		return NEAREST
	}

	nearest, exists := zt.Get(zt.nearestBehindIP)
	if !exists {
		return NEAREST
	}
	if zt.nearestBehindIP.Equal(entry.IP) {
		return EQUAL
	}

	if nearest.Position.Distance(pos) < entry.Position.Distance(pos) {
		return FARTHEST
	}
	return NEAREST
}

func (zt *ZoneTable) updateNearestBehind(entry *ZoneTableEntry, pos *Position) {
	check := zt.checkNearestBehind(entry, pos)
	if check == NEAREST || check == EQUAL {
		zt.nearestBehindLock.Lock()
		zt.nearestBehindIP = entry.IP
		zt.nearestBehindLock.Unlock()
	}
}

func (zt *ZoneTable) checkNearestFront(entry *ZoneTableEntry, pos *Position) NEAREST_BOOL {
	if !zt.IsFront(entry) {
		return FARTHEST
	}
	zt.nearestFrontLock.RLock()
	defer zt.nearestFrontLock.RUnlock()
	if zt.nearestFrontIP == nil {
		return NEAREST
	}
	nearest, exists := zt.Get(zt.nearestFrontIP)
	if !exists {
		return NEAREST
	}
	if zt.nearestFrontIP.Equal(entry.IP) {
		return EQUAL
	}

	if nearest.Position.Distance(pos) < entry.Position.Distance(pos) {
		return FARTHEST
	}
	return NEAREST
}

func (zt *ZoneTable) updateNearestFront(entry *ZoneTableEntry, pos *Position) {
	check := zt.checkNearestFront(entry, pos)
	if check == NEAREST || check == EQUAL {
		zt.nearestFrontLock.Lock()
		zt.nearestFrontIP = entry.IP
		zt.nearestFrontLock.Unlock()
	}
}

func (zt *ZoneTable) removeNearestIfEqual(entry *ZoneTableEntry) {
	// check if this is the nearest front
	zt.nearestFrontLock.RLock()
	if zt.IsFront(entry) && zt.nearestFrontIP != nil && zt.nearestFrontIP.Equal(entry.IP) {
		zt.nearestFrontLock.RUnlock()
		zt.nearestFrontLock.Lock()
		zt.nearestFrontIP = nil
		zt.nearestFrontLock.Unlock()
	} else {
		zt.nearestFrontLock.RUnlock()
	}

	// check if this is the nearest behind
	zt.nearestBehindLock.RLock()
	if zt.IsBehind(entry) && zt.nearestBehindIP != nil && zt.nearestBehindIP.Equal(entry.IP) {
		zt.nearestBehindLock.RUnlock()
		zt.nearestBehindLock.Lock()
		zt.nearestBehindIP = nil
		zt.nearestBehindLock.Unlock()
	} else {
		zt.nearestBehindLock.RUnlock()
	}
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
	log.Println()
	log.Println("ZoneTable: front of me:")
	for _, entry := range zt.GetInFrontOfMe() {
		entry.Print()
	}
	log.Println("ZoneTable: behind me:")
	for _, entry := range zt.GetBehindMe() {
		entry.Print()
	}
	log.Println()
}
