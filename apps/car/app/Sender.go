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
			pos := a.GetPosition()
			data := NewVZoneMessage(a.ip, pos, MAX_DIST_METER).Marshal()
			log.Printf("Sending zone msg : %f  %f, max dist %dm", pos.Lng, pos.Lat, MAX_DIST_METER)
			positionOption := ip.NewPositionOption(pos, MAX_DIST_METER)
			a.sendToRouterWithOptions(data, net.ParseIP(ip.BroadcastIP), positionOption.Marshal())
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
	cipherData, err := crypto.EncryptAES(a.key, data)
	if err != nil {
		return
	}
	a.ipConn.Write(cipherData, a.ip, destIP)
}

func (a *App) sendToRouterWithOptions(data []byte, destIP net.IP, options []byte) {
	cipherData, err := crypto.EncryptAES(a.key, data)
	if err != nil {
		return
	}

	packet := ip.NewIPPacketWithOptions(cipherData, a.ip, destIP, options)
	packetBytes := ip.MarshalIPPacket(packet)

	a.ipConn.Forward(packetBytes)
}
