package app

import (
	"log"

	"github.com/aaarafat/vanessa/apps/car/unix"
	. "github.com/aaarafat/vanessa/libs/vector"
)

func (a *App) GetState() *unix.State {
	a.stateLock.RLock()
	defer a.stateLock.RUnlock()
	return a.state
}

func (a *App) initState(speed uint32, route []Position, pos Position) {
	log.Printf("Initializing car state......")
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil {
		a.state = &unix.State{
			Id:                 a.id,
			Speed:              speed,
			Route:              route,
			Lat:                pos.Lat,
			Lng:                pos.Lng,
			ObstacleDetected:   false,
			Obstacles:          []Position{},
			MaxSpeed:           speed,
			Direction:          NewUnitVector(pos, pos),
			DestinationReached: false,
		}
	} else {
		a.state.Direction = NewUnitVector(a.state.GetPosition(), pos)
		a.state.Speed = speed
		a.state.Route = route
		a.state.Lat = pos.Lat
		a.state.Lng = pos.Lng
		a.state.MaxSpeed = speed
	}

	log.Printf("Car state initialized  state:  %s\n", a.state.String())
}

func (a *App) updateSpeed(speed uint32) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil || a.state.Speed == speed {
		return
	}
	if a.state.MaxSpeed < speed {
		a.state.MaxSpeed = speed
	}
	a.state.Speed = speed

	go func() {
		a.sendSpeed(speed)
	}()
	log.Printf("Speed updated: %d", speed)
}

func (a *App) updatePosition(pos Position) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil {
		return
	}
	a.state.Direction = NewUnitVector(a.state.GetPosition(), pos)
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
	oldLen := len(a.state.Obstacles)
	a.state.Obstacles = removeDuplicatePositions(append(a.state.Obstacles, pos))
	log.Printf("Obstacle added: lng: %f lat: %f", pos.Lng, pos.Lat)

	if oldLen == len(a.state.Obstacles) {
		return
	}

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
	if a.state == nil || len(obstacles) == 0 {
		return
	}
	oldLen := len(a.state.Obstacles)
	a.state.Obstacles = removeDuplicatePositions(append(a.state.Obstacles, obstacles...))
	log.Printf("Obstacles updated: %v", a.state.Obstacles)

	if oldLen == len(a.state.Obstacles) {
		return
	}

	go func() {
		data := unix.FormatObstacles(obstacles)
		a.ui.Write(data, string(unix.ObstaclesReceivedEvent))
		a.sensor.Write(data, unix.RerouteEvent)
	}()
}

func (a *App) destinationReached(pos Position) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil {
		return
	}
	a.state.DestinationReached = true
	a.state.MaxSpeed = 0
	log.Printf("Destination reached")

	go func() {
		a.updateSpeed(0)
		a.ui.Write(unix.DestinationReachedData{Coordinates: pos}, string(unix.DestinationReachedEvent))
	}()
}
