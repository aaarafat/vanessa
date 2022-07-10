package app

import (
	"log"

	"github.com/aaarafat/vanessa/apps/car/unix"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (a *App) listen() {
	select {
	case data := <-*a.router.Data:
		go a.handleMessage(data)
	}
}

func (a *App) handleMessage(data []byte) {
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
		if dist <= msg.MaxDistance {
			log.Printf("Car with ip: %s  in my zone", msg.OriginatorIP)
			a.ipConn.Forward(data)
		}

	}
}
