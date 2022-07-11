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
		// Calculate distance between my position and the zone
		dist := distancePosition(msg.Position, a.GetPosition())
		// If the distance is less than the max distance, then I am in the zone
		// and forward it again to router
		log.Printf("Distance: %f", dist)
		if dist <= msg.MaxDistance {
			log.Printf("Car with ip: %s  in my zone", msg.OriginatorIP)
			a.ipConn.Forward(bytes)
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
	}
}
