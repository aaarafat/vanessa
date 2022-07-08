package app

import (
	"log"

	"github.com/aaarafat/vanessa/apps/car/unix"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


func (a *App) GetState() *unix.State {
	a.stateLock.RLock()
	defer a.stateLock.RUnlock()
	return a.state
}

func (a *App) initState(speed int, route []Position, pos Position) {
	log.Printf("Initializing car state......")
	a.stateLock.Lock()
	a.state = &unix.State{
		Id: a.id,
		Speed: speed,
		Route: route,
		Lat: pos.Lat,
		Lng: pos.Lng,
		ObstacleDetected: false,
		Obstacles: []Position{},
	}
	a.stateLock.Unlock()
	log.Printf("Car state initialized  state:  %v\n", a.state)
}

func (a *App) updatePosition(pos Position) {
	a.stateLock.Lock()
	a.state.Lat = pos.Lat
	a.state.Lng = pos.Lng
	a.stateLock.Unlock()
	log.Printf("Position updated: lng: %f lat: %f", pos.Lng, pos.Lat)
}

func (a *App) addObstacle(pos Position, fromSensor bool) {
	a.stateLock.Lock()
	if fromSensor {
		a.state.ObstacleDetected = true
	}
	a.state.Obstacles = append(a.state.Obstacles, pos)
	a.stateLock.Unlock()
	log.Printf("Obstacle added: lng: %f lat: %f", pos.Lng, pos.Lat)
}