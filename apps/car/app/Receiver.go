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
		a.ui.Write(data, string(unix.ObstaclesReceivedEvent))

	case VObstacleType:
		msg, err := UnmarshalVObstacle(data)
		if err != nil {
			log.Printf("Error decoding VObstacle message: %v", err)
			return
		}
		log.Printf("VObstacle message received: %s", msg.String())
		a.addObstacle(msg.Position, false)
		a.ui.Write(unix.ObstacleReceivedData{ObstacleCoordinates: msg.Position}, string(unix.ObstacleReceivedEvent))
	}
}