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
	defer a.stateLock.Unlock()
	if a.state == nil {
		a.state = &unix.State{
			Id:               a.id,
			Speed:            speed,
			Route:            route,
			Lat:              pos.Lat,
			Lng:              pos.Lng,
			ObstacleDetected: false,
			Obstacles:        []Position{},
		}
	} else {
		a.state.Speed = speed
		a.state.Route = route
		a.state.Lat = pos.Lat
		a.state.Lng = pos.Lng
	}

	log.Printf("Car state initialized  state:  %v\n", a.state)
}

func (a *App) updatePosition(pos Position) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil {
		return
	}
	a.state.Lat = pos.Lat
	a.state.Lng = pos.Lng
	log.Printf("Position updated: lng: %f lat: %f", pos.Lng, pos.Lat)

	go a.ui.Write(unix.UpdateLocationData{Coordinates: pos}, string(unix.UpdateLocationEvent))
}

func (a *App) addObstacle(pos Position, fromSensor bool) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if fromSensor {
		a.state.ObstacleDetected = true
	}
	a.state.Obstacles = removeDuplicatePositions(append(a.state.Obstacles, pos))
	log.Printf("Obstacle added: lng: %f lat: %f", pos.Lng, pos.Lat)

	go func() {
		if fromSensor {
			a.sendObstacle(pos)
			a.ui.Write(unix.ObstacleDetectedData{ObstacleCoordinates: pos}, string(unix.ObstacleDetectedEvent))
		} else {
			a.ui.Write(unix.ObstacleReceivedData{ObstacleCoordinates: pos}, string(unix.ObstacleReceivedEvent))
		}
		data := unix.FormatObstacles(a.GetState().Obstacles)
		a.sensor.Write(data, unix.RerouteEvent)
	}()
}

func (a *App) updateObstacles(obstacles []Position) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil {
		return
	}
	a.state.Obstacles = removeDuplicatePositions(append(a.state.Obstacles, obstacles...))
	log.Printf("Obstacles updated: %v", a.state.Obstacles)

	go func() {
		data := unix.FormatObstacles(obstacles)
		a.ui.Write(data, string(unix.ObstaclesReceivedEvent))
		a.sensor.Write(data, unix.RerouteEvent)
	}()
}
