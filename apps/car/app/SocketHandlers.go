package app

import (
	"encoding/json"
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/car/unix"
)

func (a *App) startSocketHandlers() {
	go a.obstacleHandler()
	go a.destinationReachedHandler()
	go a.updateLocationHandler()
}

func (a *App) obstacleHandler() {
	obstacleChannel := make(chan json.RawMessage)
	obstableSubscriber := &unix.Subscriber{Messages: &obstacleChannel}
	a.unix.Subscribe(unix.ObstacleDetectedEvent, obstableSubscriber)

	for {
		select {
		case data := <-*obstableSubscriber.Messages:
			var obstacle unix.ObstacleDetectedData
			err := json.Unmarshal(data, &obstacle)
			if err != nil {
				log.Printf("Error decoding obstacle-detected data: %v", err)
				return
			}
			log.Printf("Sockets : Obstacle detected: %v\n", data)

			// TODO: send it with loopback interface to the router to be processed by the AODV
			go a.ipConn.Write(data, a.ip, net.ParseIP(RSUIP))
			a.unix.Write(data)
		}
	}
}

func (a *App) destinationReachedHandler() {
	destinationReachedChannel := make(chan json.RawMessage)
	destinationReachedSubscriber := &unix.Subscriber{Messages: &destinationReachedChannel}
	a.unix.Subscribe(unix.DestinationReachedEvent, destinationReachedSubscriber)

	for {
		select {
		case data := <-*destinationReachedSubscriber.Messages:
			var destinationReached unix.DestinationReachedData
			err := json.Unmarshal(data, &destinationReached)
			if err != nil {
				log.Printf("Error decoding destination-reached data: %v", err)
				return
			}
			log.Printf("Sockets : destination-reached: %v\n", data)

			go a.ipConn.Write(data, a.ip, net.ParseIP(RSUIP))
			a.unix.Write(data)
		}
	}
}

func (a *App) updateLocationHandler() {
	updateLocationChannel := make(chan json.RawMessage)
	updateLocationSubscriber := &unix.Subscriber{Messages: &updateLocationChannel}
	a.unix.Subscribe(unix.UpdateLocationEvent, updateLocationSubscriber)

	for {
		select {
		case data := <-*updateLocationSubscriber.Messages:
			var updateLocation unix.UpdateLocationData
			err := json.Unmarshal(data, &updateLocation)
			if err != nil {
				log.Printf("Error decoding update-location data: %v", err)
				return
			}

			a.updatePosition(&updateLocation.Coordinates)

			go a.ipConn.Write(data, a.ip, net.ParseIP(RSUIP))
			a.unix.Write(data)
		}
	}
}
