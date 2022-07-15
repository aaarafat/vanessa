package unix

import (
	. "github.com/aaarafat/vanessa/libs/vector"
)

type Event string

const (
	DestinationReachedEvent Event = "destination-reached"  // from simulator to car ui when destination is reached
	ObstacleDetectedEvent   Event = "obstacle-detected"    // from simulator to car
	ObstacleReceivedEvent   Event = "obstacle-received"    // from car to ui and simulator when obstacle are received from other rsu
	ObstaclesReceivedEvent  Event = "obstacles-received"   // from car to ui and simulator when obstacles are received from other rsu
	RerouteEvent            Event = "reroute"              // from car to simulator when reroute is requested
	ChangeSpeedEvent        Event = "change-speed"         // from car to simulator when speed is changed
	AddCarEvent             Event = "add-car"              // from simulator to car when car is added, rerouted and speed is changed
	UpdateLocationEvent     Event = "update-location"      // from simulator to car when location is changed
	CheckRouteEvent         Event = "check-route"          // from car to simulator when route is checked
	CheckRouteResponseEvent Event = "check-route-response" // from simulator to car when route is checked
	StateEvent              Event = "state"                // from car to ui
	RefreshEvent            Event = "refresh"              // from car to ui
)

type RefreshData = State

type CheckRouteResponseData struct {
	InRoute bool `json:"in_route"`
}

type CheckRouteData struct {
	Coordinate Position `json:"coordinate"`
}

type DestinationReachedData struct {
	Coordinates Position `json:"coordinates"`
}

type ObstacleDetectedData struct {
	Coordinates         Position `json:"coordinates"`
	ObstacleCoordinates Position `json:"obstacle_coordinates"`
}

type AddCarData struct {
	Coordinates Position   `json:"coordinates"`
	Route       []Position `json:"route"`
	Speed       int        `json:"speed"`
}

type UpdateLocationData struct {
	Coordinates Position `json:"coordinates"`
}

type ObstacleReceivedData struct {
	ObstacleCoordinates Position `json:"obstacle_coordinates"`
}

type ObstaclesReceivedData struct {
	ObstacleCoordinates []Position `json:"obstacle_coordinates"`
}

type SpeedData struct {
	Speed int `json:"speed"`
}

type StateData struct {
	Coordinates Position   `json:"coordinates"`
	Speed       int        `json:"speed"`
	Route       []Position `json:"route"`
}

func FormatObstacles(pos []Position) ObstaclesReceivedData {
	return ObstaclesReceivedData{ObstacleCoordinates: pos}
}

func PositionFromBytes(data []byte) Position {
	return UnmarshalPosition(data)
}

func ObstaclesFromBytes(data []byte, len int) ObstaclesReceivedData {
	var obstacles ObstaclesReceivedData
	obstacles.ObstacleCoordinates = make([]Position, len)
	for i := 0; i < len; i++ {
		obstacles.ObstacleCoordinates[i] = PositionFromBytes(data[i*16 : (i+1)*16])
	}
	return obstacles
}
