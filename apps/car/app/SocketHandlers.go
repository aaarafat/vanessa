package app

import (
	"encoding/json"
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/car/unix"
)

func (a *App) startSocketHandlers() {
	a.addCarHandler()
	a.obstacleHandler()
	a.destinationReachedHandler()
	a.updateLocationHandler()
	a.checkRouteResponseHandler()
	a.changeStopHandler()
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

				log.Printf("Adding car %v", addCar)

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

func (a *App) checkRouteResponseHandler() {
	checkRouteResponseChannel := make(chan json.RawMessage)
	checkRouteResponseSubscriber := &unix.Subscriber{Messages: &checkRouteResponseChannel}
	a.sensor.Subscribe(unix.CheckRouteResponseEvent, checkRouteResponseSubscriber)

	go func() {
		for {
			select {
			case data := <-*checkRouteResponseSubscriber.Messages:
				var checkRouteResponse unix.CheckRouteResponseData
				err := json.Unmarshal(data, &checkRouteResponse)
				if err != nil {
					log.Printf("Error decoding check-route-response data: %v", err)
					return
				}

				go func() {
					ip, exists := a.checkRouteBuffer.GetStringKey(checkRouteResponse.Coordinates.String())
					if exists {
						a.checkRouteBuffer.Del(checkRouteResponse.Coordinates.String())
						a.zoneTable.Ignore(ip.(net.IP), !checkRouteResponse.InRoute)
					}
				}()
			}
		}
	}()
}

func (a *App) changeStopHandler() {
	changeStopChannel := make(chan json.RawMessage)
	changeStopSubscriber := &unix.Subscriber{Messages: &changeStopChannel}
	a.sensor.Subscribe(unix.ChangeStopEvent, changeStopSubscriber)

	go func() {
		for {
			select {
			case data := <-*changeStopSubscriber.Messages:
				var changeStop unix.ChangeStopData
				err := json.Unmarshal(data, &changeStop)
				if err != nil {
					log.Printf("Error decoding change-stop data: %v", err)
					return
				}

				go a.changeStop(changeStop.Stop)
			}
		}
	}()
}
