package app

import (
	"log"
	"reflect"

	"github.com/aaarafat/vanessa/apps/car/unix"
	. "github.com/aaarafat/vanessa/libs/vector"
)

func (a *App) GetState() *unix.State {
	a.stateLock.RLock()
	defer a.stateLock.RUnlock()
	return a.state
}

func (a *App) initState(speed uint32, route []Position, pos Position, stopped bool) {
	log.Printf("Initializing car state......")
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	changedRoute := route
	if a.state != nil && reflect.DeepEqual(a.state.Route, route) {
		changedRoute = nil
	}

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
			Stopped:            stopped,
		}
	} else {
		a.state.Direction = NewUnitVector(a.state.GetPosition(), pos)
		a.state.Speed = speed
		a.state.Route = route
		a.state.Lat = pos.Lat
		a.state.Lng = pos.Lng
		a.state.MaxSpeed = speed
		a.state.Stopped = stopped
	}

	log.Printf("Car state initialized  state:  %s\n", a.state.String())

	go func() {
		a.ui.Write(unix.StateData{Coordinates: pos, Route: changedRoute, Speed: int(speed)}, string(unix.StateEvent))
	}()
}

func (a *App) updateSpeed(speed uint32) {
	a.stateLock.Lock()
	defer a.stateLock.Unlock()
	if a.state == nil || a.state.Speed == speed || speed > a.state.MaxSpeed {
		return
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
	// update direction if the car is moving
	newdirection := NewUnitVector(a.state.GetPosition(), pos)
	if newdirection.Lng != 0 || newdirection.Lat != 0 {
		a.state.Direction = newdirection
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
		a.state.MaxSpeed = 0
	}
	oldLen := len(a.state.Obstacles)
	a.state.Obstacles = removeDuplicatePositions(append(a.state.Obstacles, pos))
	log.Printf("Obstacle added: lng: %f lat: %f", pos.Lng, pos.Lat)

	if oldLen == len(a.state.Obstacles) {
		return
	}

	go func() {
		if fromSensor {
			a.updateSpeed(0)
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
