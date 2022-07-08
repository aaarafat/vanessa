package unix

import (
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


type Event string

const (
	DestinationReachedEvent Event = "destination-reached"
	ObstacleDetectedEvent   Event = "obstacle-detected" // from sensor
	ObstacleReceivedEvent   Event = "obstacle-received" // from router
	ObstaclesReceivedEvent  Event = "obstacles-received" // from RSU
	AddCarEvent             Event = "add-car"
	UpdateLocationEvent     Event = "update-location"
)


type DestinationReachedData struct {
	Coordinates Position
}

type ObstacleDetectedData struct {
	Coordinates         Position `json:"coordinates"`
	ObstacleCoordinates Position `json:"obstacle_coordinates"`
}

type AddCarData struct {
	Coordinates Position `json:"coordinates"`
	Route       []Position `json:"route"`
	Speed       int 			`json:"speed"`
}

type UpdateLocationData struct {
	Coordinates Position
}

type ObstacleReceivedData struct {
	ObstacleCoordinates []Position `json:"obstacle_coordinates"`
}


func FormatObstacles(pos []Position) ObstacleReceivedData {
	return ObstacleReceivedData{ObstacleCoordinates: pos}
}