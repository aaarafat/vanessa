package app

import (
	"log"
	"net"
	"time"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	"github.com/aaarafat/vanessa/libs/crypto"
)

func (a *App) sendHeartBeat() {
	for {
		if a.state != nil {
			data := NewVHBeatMessage(a.ip, Position{Lng: a.state.Lng, Lat: a.state.Lat}).Marshal()
			log.Printf("Sending heartbeat : %f  %f", a.state.Lng, a.state.Lat)
			a.sendToRouter(data, net.ParseIP(ip.RsuIP))
		}
		time.Sleep(time.Millisecond * HEART_BEAT_INTERVAL_MS)
	}
}

func (a *App) sendZoneMsg() {
	for {
		if a.state != nil {
			data := NewVZoneMessage(a.ip, Position{Lng: a.state.Lng, Lat: a.state.Lat}, MAX_DIST_METER).Marshal()
			log.Printf("Sending zone msg : %f  %f, max dist %dm", a.state.Lng, a.state.Lat, MAX_DIST_METER)
			a.sendToRouter(data, net.ParseIP(ip.BroadcastIP))
		}
		time.Sleep(time.Millisecond * ZONE_MSG_INTERVAL_MS)
	}
}

func (a *App) sendObstacle(pos Position) {
	log.Printf("Sending obstacle : %f  %f", pos.Lng, pos.Lat)
	data := NewVObstacleMessage(a.ip, pos, 0).Marshal()
	a.sendToRouter(data, net.ParseIP(ip.RsuIP))
}

func (a *App) sendToRouter(data []byte, destIP net.IP) {
	cipherData, err := crypto.EncryptAES(data, a.key)
	if err != nil {
		return
	}
	a.ipConn.Write(cipherData, a.ip, destIP)
}
