package app

import (
	"net"
	"testing"
	"time"

	. "github.com/aaarafat/vanessa/libs/vector"
	"github.com/stretchr/testify/assert"
)

func TestSet(t *testing.T) {
	zoneTable := NewZoneTable()
	ip := net.ParseIP("192.168.1.1")
	myPos := Position{Lng: 1, Lat: 2}
	myPos2 := Position{Lng: 2, Lat: 3}
	myDir := NewUnitVector(myPos, myPos2)
	pos := Position{Lng: 5, Lat: 5}
	pos2 := Position{Lng: 6, Lat: 6}
	speed := uint32(10)

	zoneTable.Set(ip, speed, pos, myPos, myDir)

	entry, exists := zoneTable.Get(ip)
	assert.True(t, exists)
	assert.Equal(t, pos, entry.Position)
	assert.Equal(t, speed, entry.Speed)
	assert.Equal(t, NewUnitVector(pos, pos), entry.Direction)
	assert.Equal(t, NewUnitVector(myPos, pos), entry.DirectionFromMe)
	assert.Equal(t, -0.141897054604164, entry.Angle)

	// test adding a duplicate entry
	zoneTable.Set(ip, speed+1, pos2, myPos, myDir)

	expectedDirection := NewUnitVector(pos, pos2)

	entry, exists = zoneTable.Get(ip)
	assert.True(t, exists)
	assert.Equal(t, expectedDirection, entry.Direction)
	assert.Equal(t, pos2, entry.Position)
	assert.Equal(t, speed+1, entry.Speed)
	assert.Equal(t, NewUnitVector(myPos, pos2), entry.DirectionFromMe)
	assert.Equal(t, -0.11065722117389568, entry.Angle)

	// test delete entry
	time.Sleep(ZoneTable_UPDATE_INTERVAL_MS * time.Millisecond)
	time.Sleep(100 * time.Millisecond)
	_, exists = zoneTable.Get(ip)
	assert.False(t, exists)
}

func TestGetInFrontOfMe(t *testing.T) {
	zoneTable := NewZoneTable()
	ip := net.ParseIP("192.168.1.1")
	myPos := Position{Lng: 1, Lat: 2}
	myPos2 := Position{Lng: 2, Lat: 3}
	myDir := NewUnitVector(myPos, myPos2)
	pos := Position{Lng: 5, Lat: 5}
	speed := uint32(10)

	zoneTable.Set(ip, speed, pos, myPos, myDir)

	front := zoneTable.GetInFrontOfMe()

	assert.Len(t, front, 1)
	assert.Equal(t, front[0].IP, ip)

	// set one behind me
	ip2 := net.ParseIP("192.168.1.2")
	pos2 := Position{Lng: -5, Lat: -5}

	zoneTable.Set(ip2, speed, pos2, myPos, myDir)

	front = zoneTable.GetInFrontOfMe()

	assert.Len(t, front, 1)
	assert.Equal(t, front[0].IP, ip)
}

func TestGetBehindMe(t *testing.T) {
	zoneTable := NewZoneTable()
	ip := net.ParseIP("192.168.1.1")
	myPos := Position{Lng: 1, Lat: 2}
	myPos2 := Position{Lng: 2, Lat: 3}
	myDir := NewUnitVector(myPos, myPos2)
	pos := Position{Lng: 5, Lat: 5}
	speed := uint32(10)

	zoneTable.Set(ip, speed, pos, myPos, myDir)

	back := zoneTable.GetBehindMe()

	assert.Len(t, back, 0)

	// set one behind me
	ip2 := net.ParseIP("192.168.1.2")
	pos2 := Position{Lng: -5, Lat: -5}

	zoneTable.Set(ip2, speed, pos2, myPos, myDir)

	back = zoneTable.GetBehindMe()

	assert.Len(t, back, 1)
	assert.Equal(t, back[0].IP, ip2)
}

func BenchmarkSet(b *testing.B) {
	zoneTable := NewZoneTable()
	ip := net.ParseIP("192.168.1.1")
	myPos := Position{Lng: 1, Lat: 2}
	myPos2 := Position{Lng: 2, Lat: 3}
	myDir := NewUnitVector(myPos, myPos2)
	pos := Position{Lng: 5, Lat: 5}
	speed := uint32(10)

	for i := 0; i < b.N; i++ {
		zoneTable.Set(ip, speed, pos, myPos, myDir)
	}
	// 224269	      5370 ns/op	     496 B/op	      14 allocs/op
}
