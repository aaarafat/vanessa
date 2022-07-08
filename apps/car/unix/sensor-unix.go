package unix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type Event string

const (
	DestinationReachedEvent Event = "destination-reached"
	ObstacleDetectedEvent   Event = "obstacle-detected"
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

type Subscriber struct {
	Messages *chan json.RawMessage
}

type SensorUnix struct {
	id     int
	topics map[Event][]*Subscriber
}

func NewSensorUnix(id int) *SensorUnix {
	return &SensorUnix{id: id, topics: make(map[Event][]*Subscriber)}
}

func (u *SensorUnix) Subscribe(topic Event, subscriber *Subscriber) {
	if u.topics[topic] == nil {
		u.topics[topic] = []*Subscriber{}
	}
	u.topics[topic] = append(u.topics[topic], subscriber)
	log.Printf("Subscribed to %s\n", topic)
}

func (u *SensorUnix) publish(topic Event, message json.RawMessage) {
	if u.topics[topic] == nil {
		log.Printf("No subscribers for %s\n", topic)
		return
	}
	for _, subscriber := range u.topics[topic] {
		*subscriber.Messages <- message
	}
	log.Printf("Published to %s\n", topic)
}

func (unix *SensorUnix) reader(d *json.Decoder) {
	var m map[string]json.RawMessage
	for {
		b := make([]byte, 1024)
		d.Buffered().Read(b)
		log.Printf("Reading from socket %s\n", b)
		err := d.Decode(&m)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		var eventType Event
		err = json.Unmarshal(m["type"], &eventType)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		switch eventType {
		case DestinationReachedEvent:
			var p DestinationReachedData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error decoding destination-reached data: %v", err)
				return
			}
			log.Printf("Destination reached: %v\n", p)
			unix.publish(DestinationReachedEvent, m["data"])

		case ObstacleDetectedEvent:
			var p ObstacleDetectedData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error decoding obstacle-detected data: %v", err)
				return
			}
			log.Printf("Obstacle detected: %v\n", p)
			unix.publish(ObstacleDetectedEvent, m["data"])

		case AddCarEvent:
			var p AddCarData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error decoding add-car data: %v", err)
				return
			}
			log.Printf("Car added: %v\n", p)
			unix.publish(AddCarEvent, m["data"])

		case UpdateLocationEvent:
			var p UpdateLocationData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			unix.publish(UpdateLocationEvent, m["data"])
		}
	}
}

func (unix *SensorUnix) Write(message any) {
	log.Printf("Writing to socket...\n")
	conn, err := unix.initUnixWriteSocket()
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	defer conn.Close()

	err = json.NewEncoder(conn).Encode(map[string]interface{}{
		"id":   unix.id,
		"data": message,
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
}

func (unix *SensorUnix) initUnixWriteSocket() (net.Conn, error) {
	addr := fmt.Sprintf("/tmp/car%dwrite.socket", unix.id)
	log.Printf("Connecting to %s\n", addr)
	conn, err := net.Dial("unixgram", addr)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	log.Printf("Connected\n")
	return conn, nil
}

func (unix *SensorUnix) Start() {
	socketAddress := fmt.Sprintf("/tmp/car%d.socket", unix.id)
	err := os.RemoveAll(socketAddress)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	addr, err := net.ResolveUnixAddr("unixgram", socketAddress)
	if err != nil {
		log.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		log.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	d := json.NewDecoder(conn)
	log.Printf("Listening to %s ..\n", socketAddress)

	// go unix.testWrite()

	unix.reader(d)
}
