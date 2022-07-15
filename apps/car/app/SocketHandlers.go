package app

import (
	"encoding/json"
	"log"

	"github.com/aaarafat/vanessa/apps/car/unix"
)

func (a *App) startSocketHandlers() {
	a.addCarHandler()
	a.obstacleHandler()
	a.destinationReachedHandler()
	a.updateLocationHandler()
}

func (a *App) addCarHandler() {
	addCarChannel := make(chan json.RawMessage)
	addCarSubscriber := &unix.Subscriber{Messages: &addCarChannel}
	a.sensor.Subscribe(unix.AddCarEvent, addCarSubscriber)

	go func() {
		for {
			select {
			case data := <-*addCarSubscriber.Messages:
				var addCar unix.AddCarData
				err := json.Unmarshal(data, &addCar)
				if err != nil {
					log.Printf("Error decoding add-car data: %v", err)
					return
				}

				a.initState(uint32(addCar.Speed), addCar.Route, addCar.Coordinates)
			}
		}
	}()
}

func (a *App) obstacleHandler() {
	obstacleChannel := make(chan json.RawMessage)
	obstableSubscriber := &unix.Subscriber{Messages: &obstacleChannel}
	a.sensor.Subscribe(unix.ObstacleDetectedEvent, obstableSubscriber)

	go func() {
		for {
			select {
			case data := <-*obstableSubscriber.Messages:
				var obstacle unix.ObstacleDetectedData
				err := json.Unmarshal(data, &obstacle)
				if err != nil {
					log.Printf("Error decoding obstacle-detected data: %v", err)
					return
				}

				go a.addObstacle(obstacle.ObstacleCoordinates, true)
			}
		}
	}()
}

func (a *App) destinationReachedHandler() {
	destinationReachedChannel := make(chan json.RawMessage)
	destinationReachedSubscriber := &unix.Subscriber{Messages: &destinationReachedChannel}
	a.sensor.Subscribe(unix.DestinationReachedEvent, destinationReachedSubscriber)

	go func() {
		for {
			select {
			case data := <-*destinationReachedSubscriber.Messages:
				var destinationReached unix.DestinationReachedData
				err := json.Unmarshal(data, &destinationReached)
				if err != nil {
					log.Printf("Error decoding destination-reached data: %v", err)
					return
				}

				go a.destinationReached(destinationReached.Coordinates)
			}
		}
	}()
}

func (a *App) updateLocationHandler() {
	updateLocationChannel := make(chan json.RawMessage)
	updateLocationSubscriber := &unix.Subscriber{Messages: &updateLocationChannel}
	a.sensor.Subscribe(unix.UpdateLocationEvent, updateLocationSubscriber)

	go func() {
		for {
			select {
			case data := <-*updateLocationSubscriber.Messages:
				var updateLocation unix.UpdateLocationData
				err := json.Unmarshal(data, &updateLocation)
				if err != nil {
					log.Printf("Error decoding update-location data: %v", err)
					return
				}

				go a.updatePosition(updateLocation.Coordinates)
			}
		}
	}()
}
