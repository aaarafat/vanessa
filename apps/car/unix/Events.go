package unix

import (
	. "github.com/aaarafat/vanessa/libs/vector"
)

type Event string

const (
	DestinationReachedEvent Event = "destination-reached"
	ObstacleDetectedEvent   Event = "obstacle-detected"  // from sensor
	ObstacleReceivedEvent   Event = "obstacle-received"  // from router
	ObstaclesReceivedEvent  Event = "obstacles-received" // from RSU
	RerouteEvent            Event = "reroute"
	ChangeSpeedEvent        Event = "change-speed"
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
	ObstacleCoordinates Position `json:"obstacle_coordinates"`
}

type ObstaclesReceivedData struct {
	ObstacleCoordinates []Position `json:"obstacle_coordinates"`
}

type SpeedData struct {
	Speed int `json:"speed"`
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
