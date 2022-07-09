package app

import (
	"log"
	"net"
	"time"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


func (a *App) sendHeartBeat() {
	for {
		if a.state != nil {
			data := NewVHBeatMessage(a.ip, Position{Lng: a.state.Lng, Lat: a.state.Lat}).Marshal()
			log.Printf("Sending heartbeat : %f  %f", a.state.Lng, a.state.Lat)
			a.ipConn.Write(data, a.ip, net.ParseIP(ip.RsuIP))
		}
		time.Sleep(time.Millisecond * DATA_SENDING_INTERVAL_MS)
	}
}


func (a *App) sendObstacle(pos Position) {
	log.Printf("Sending obstacle : %f  %f", pos.Lng, pos.Lat)
	data := NewVObstacleMessage(a.ip, pos, 0).Marshal()
	a.ipConn.Write(data, a.ip, net.ParseIP(ip.RsuIP))
}