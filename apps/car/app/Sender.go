package app

import (
	"log"
	"net"
	"time"

	"github.com/aaarafat/vanessa/apps/car/unix"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	"github.com/aaarafat/vanessa/libs/crypto"
	. "github.com/aaarafat/vanessa/libs/vector"
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
			state := a.GetState()
			pos := state.GetPosition()
			speed := state.Speed
			data := NewVZoneMessage(a.ip, pos, uint32(speed)).Marshal()
			log.Printf("Sending zone msg : %f  %f, speed %d  MaxDistance %d m", pos.Lng, pos.Lat, speed, MAX_DIST_METER)
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

func (a *App) sendSpeed(speed uint32) {
	a.sensor.Write(unix.SpeedData{Speed: int(speed)}, unix.ChangeSpeedEvent)

	entries := a.zoneTable.GetBehindMe()
	for _, entry := range entries {
		if entry.Speed > speed {
			// send slow down
			a.sendToRouter(NewVSpeedMessage(a.ip, speed).Marshal(), entry.IP)
		}
	}
}
