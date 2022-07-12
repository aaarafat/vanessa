package router

import (
	"log"
	"net"
	"time"
)

type RARPEntry struct {
	MAC   net.HardwareAddr
	timer *time.Timer
}

const lifeTimeMS = 10000

type RSUARP struct {
	table map[string]RARPEntry
}

func NewRSUARP() *RSUARP {
	return &RSUARP{
		table: make(map[string]RARPEntry),
	}
}

func (RARP *RSUARP) Set(ip string, mac net.HardwareAddr) {

	if mac == nil {
		log.Panic("You are trying to add null neighbor")
	}
	callback := func() {
		delete(RARP.table, ip)
	}
	if val, ok := RARP.table[ip]; ok {
		val.timer.Reset(lifeTimeMS * time.Millisecond)
	} else {

		entry := &RARPEntry{
			MAC:   mac,
			timer: time.AfterFunc(lifeTimeMS*time.Millisecond, callback),
		}

		RARP.table[ip] = *entry
		log.Printf("adding %s with %s", ip, mac.String())
	}
}

func (RARP *RSUARP) GetTable() map[string]RARPEntry {
	return RARP.table
}

func (RSUARP *RSUARP) Len() int {
	return len(RSUARP.table)
}

func (RARP *RSUARP) Get(ip string) net.HardwareAddr {
	return RARP.table[ip].MAC
}

func (RARP *RSUARP) Print() {
	log.Println("Printing ARP Table")
	for k, v := range RARP.table {
		log.Println(k, v)
	}
}
