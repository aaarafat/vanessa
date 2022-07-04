package unix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Position struct {
	Lat float64
	Lng float64
}

type DestinationReachedData struct {
	Coordinates Position
}

type ObstacleDetectedData struct {
	Coordinates         Position
	ObstacleCoordinates Position `json:"obstacle_coordinates"`
}

type AddCarData struct {
	Coordinates Position
}
type UpdateLocationData struct {
	Coordinates Position
}

type UnixSocket struct {
	id  int
}

func NewUnixSocket(id int) *UnixSocket {
	return &UnixSocket{id: id}
}


func (unix *UnixSocket) reader(d *json.Decoder) {
	var m map[string]json.RawMessage
	for {
		err := d.Decode(&m)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}

		var eventType string
		err = json.Unmarshal(m["type"], &eventType)
		if err != nil {
			log.Printf("Error: %v\n", err)
			return
		}
		log.Printf("Event type: %s\n", eventType)

		switch eventType {
		case "destination-reached":
			var p DestinationReachedData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error decoding destination-reached data: %v", err)
				return
			}
			log.Printf("Destination reached: %v\n", p)

		case "obstacle-detected":
			var p ObstacleDetectedData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error decoding obstacle-detected data: %v", err)
				return
			}
			log.Printf("Obstacle detected: %v\n", p)

		case "add-car":
			var p AddCarData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error decoding add-car data: %v", err)
				return
			}
			log.Printf("Car added: %v\n", p)

		case "update-location":
			var p UpdateLocationData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			log.Printf("Updated Location: %v\n", p)
		}
	}
}

func (unix *UnixSocket) testWrite() {
	time.Sleep(time.Millisecond * 100)

	for {
		// sleep for 5 sec
		time.Sleep(time.Second * 5)

		conn, err := unix.initUnixWriteSocket()
		if err != nil {
			log.Printf("Error: %v\n", err)
			continue
		}
		err = json.NewEncoder(conn).Encode(map[string]interface{}{
			"id":   unix.id,
			"data": "Hello from car",
		})

		if err != nil {
			log.Printf("Error: %v\n", err)
		}
		log.Printf("Sent message\n")
		conn.Close()
	}
}

func (unix *UnixSocket) initUnixWriteSocket() (net.Conn, error) {
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

func (unix *UnixSocket) Start() {
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

	go unix.testWrite()

	unix.reader(d)
}