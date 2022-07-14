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
	pos := Position{Lng: 5, Lat: 5}
	pos2 := Position{Lng: 6, Lat: 6}
	speed := uint32(10)

	zoneTable.Set(ip, speed, pos, myPos)

	entry, exists := zoneTable.Get(ip)
	assert.True(t, exists)
	assert.Equal(t, pos, entry.Position)
	assert.Equal(t, speed, entry.Speed)
	assert.Equal(t, NewUnitVector(pos, pos), entry.Direction)
	assert.Equal(t, NewUnitVector(myPos, pos), entry.DirectionFromMe)
	assert.Equal(t, 0.0, entry.Angle)

	// test adding a duplicate entry
	zoneTable.Set(ip, speed+1, pos2, myPos)

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

func BenchmarkSet(b *testing.B) {
	zoneTable := NewZoneTable()
	ip := net.ParseIP("192.168.1.1")
	myPos := Position{Lng: 1, Lat: 2}
	pos := Position{Lng: 5, Lat: 5}
	speed := uint32(10)

	for i := 0; i < b.N; i++ {
		zoneTable.Set(ip, speed, pos, myPos)
	}
	// 224269	      5370 ns/op	     496 B/op	      14 allocs/op
}
