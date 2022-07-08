package unix

import (
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type Event string

const (
	DestinationReachedEvent Event = "destination-reached"
	ObstacleDetectedEvent   Event = "obstacle-detected"
	ObstacleReceivedEvent   Event = "obstacle-received"
	AddCarEvent             Event = "add-car"
	UpdateLocationEvent     Event = "update-location"
)

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
	ObstacleCoordinates []Position `json:"obstacle_coordinates"`
}

func FormatObstacles(pos []Position) ObstacleReceivedData {
	return ObstacleReceivedData{ObstacleCoordinates: pos}
}
