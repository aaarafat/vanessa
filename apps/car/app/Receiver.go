package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/car/unix"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (a *App) listen() {
	for {
		data, err := a.router.Read()
		if err != nil {
			continue
		}
		go a.handleMessage(data)
	}
}

func (a *App) handleMessage(bytes []byte) {
	data, err := a.getDataFromPacket(bytes)
	if err != nil {
		return
	}

	mType := data[0]
	switch mType {
	case VOREPType:
		// TODO: Must Request Obstacles from Router then send them to UI
		msg := UnmarshalVOREP(data)
		log.Printf("VOREP message received: %s", msg.String())
		obstacles := unix.ObstaclesFromBytes(msg.Obstacles, int(msg.Length))
		a.updateObstacles(obstacles.ObstacleCoordinates)

	case VObstacleType:
		msg, err := UnmarshalVObstacle(data)
		if err != nil {
			log.Printf("Error decoding VObstacle message: %v", err)
			return
		}
		log.Printf("VObstacle message received: %s", msg.String())
		a.addObstacle(msg.Position, false)

	case VZoneType:
		msg, err := UnmarshalVZone(data)
		if err != nil {
			log.Printf("Error decoding VZone message: %v", err)
			return
		}
		log.Printf("VZone message received: %s", msg.String())
		state := a.GetState()
		entry := a.zoneTable.Set(msg.OriginatorIP, msg.Speed, msg.Position, state.GetPosition(), state.Direction)
		if a.zoneTable.IsFront(entry) && entry.Speed < a.GetState().Speed {
			a.updateSpeed(entry.Speed)
		}

	case VPathDiscoveryType:
		msg, err := UnmarshalVPathDiscovery(data)
		if err != nil {
			log.Printf("Error decoding VPathDiscovery message: %v", err)
			return
		}
		if msg.DestIP.Equal(net.ParseIP(ip.RsuIP)) {
			log.Printf("Sending VOREQ to RSU\n")
			data := NewVOREQMessage(a.ip, a.GetState().Obstacles).Marshal()
			a.sendToRouter(data, net.ParseIP(ip.RsuIP))
		}

	case VSpeedType:
		msg, err := UnmarshalVSpeed(data)
		if err != nil {
			log.Printf("Error decoding VSpeed message: %v", err)
			return
		}
		log.Printf("VSpeed message received: %s", msg.String())
		entry, exists := a.zoneTable.Get(msg.OriginatorIP)
		if exists && a.zoneTable.IsFront(entry) {
			a.updateSpeed(msg.Speed)
		}

	default:
		log.Printf("Unknown message type: %d", mType)
	}
}
