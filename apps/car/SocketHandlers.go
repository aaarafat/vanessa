/*package packetfilter

import (
	"encoding/json"
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
	"github.com/aaarafat/vanessa/apps/network/unix"
)

func (pf *PacketFilter) startSocketHandlers() {
	go pf.obstacleHandler()
	go pf.destinationReachedHandler()
	go pf.updateLocationHandler()
}

func (pf *PacketFilter) obstacleHandler() {
	obstacleChannel := make(chan json.RawMessage)
	obstableSubscriber := &unix.Subscriber{Messages: &obstacleChannel}
	pf.unix.Subscribe(unix.ObstacleDetectedEvent, obstableSubscriber)

	for {
		select {
		case data := <-*obstableSubscriber.Messages:
			var obstacle unix.ObstacleDetectedData
			err := json.Unmarshal(data, &obstacle)
			if err != nil {
				log.Printf("Error decoding obstacle-detected data: %v", err)
				return
			}
			log.Printf("Packet Filter : Obstacle detected: %v\n", data)

			// TODO: send it with loopback interface to the router to be processed by the AODV
			go pf.networkLayer.Send(data, pf.srcIP, net.ParseIP(aodv.RsuIP))
			pf.unix.Write(data)
		}
	}
}

func (pf *PacketFilter) destinationReachedHandler() {
	destinationReachedChannel := make(chan json.RawMessage)
	destinationReachedSubscriber := &unix.Subscriber{Messages: &destinationReachedChannel}
	pf.unix.Subscribe(unix.DestinationReachedEvent, destinationReachedSubscriber)

	for {
		select {
		case data := <-*destinationReachedSubscriber.Messages:
			var destinationReached unix.DestinationReachedData
			err := json.Unmarshal(data, &destinationReached)
			if err != nil {
				log.Printf("Error decoding destination-reached data: %v", err)
				return
			}
			log.Printf("Packet Filter : destination-reached: %v\n", data)

			go pf.networkLayer.Send(data, pf.srcIP, net.ParseIP(aodv.RsuIP))
			pf.unix.Write(data)
		}
	}
}

func (pf *PacketFilter) updateLocationHandler() {
	updateLocationChannel := make(chan json.RawMessage)
	updateLocationSubscriber := &unix.Subscriber{Messages: &updateLocationChannel}
	pf.unix.Subscribe(unix.UpdateLocationEvent, updateLocationSubscriber)

	for {
		select {
		case data := <-*updateLocationSubscriber.Messages:
			var updateLocation unix.UpdateLocationData
			err := json.Unmarshal(data, &updateLocation)
			if err != nil {
				log.Printf("Error decoding update-location data: %v", err)
				return
			}
			log.Printf("Packet Filter : update-location: %v\n", data)

			go pf.networkLayer.Send(data, pf.srcIP, net.ParseIP(aodv.RsuIP))
			pf.unix.Write(data)
		}
	}
}
*/