package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

type Logger struct {
	debug bool
	*log.Logger
}

var logger Logger

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
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

func reader(d *json.Decoder) {
	var m map[string]json.RawMessage
	for {
		err := d.Decode(&m)
		if err != nil {
			logger.Log("Error: %v\n", err)
			return
		}

		var eventType string
		err = json.Unmarshal(m["type"], &eventType)
		if err != nil {
			logger.Log("Error: %v\n", err)
			return
		}
		logger.Log("Event type: %s\n", eventType)

		switch eventType {
		case "destination-reached":
			var p DestinationReachedData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				logger.Log("Error decoding destination-reached data: %v", err)
				return
			}
			logger.Log("Destination reached: %v\n", p)

		case "obstacle-detected":
			var p ObstacleDetectedData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				logger.Log("Error decoding obstacle-detected data: %v", err)
				return
			}
			logger.Log("Obstacle detected: %v\n", p)

		case "add-car":
			var p AddCarData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				logger.Log("Error decoding add-car data: %v", err)
				return
			}
			logger.Log("Car added: %v\n", p)

		case "update-location":
			var p UpdateLocationData
			err := json.Unmarshal(m["data"], &p)
			if err != nil {
				logger.Log("Error: %v", err)
				return
			}
			logger.Log("Updated Location: %v\n", p)
		}
	}
}

func testWrite(id int) {
	time.Sleep(time.Millisecond * 100)

	for {
		// sleep for 5 sec
		time.Sleep(time.Second * 5)

		conn, err := initUnixWriteSocket(id)
		if err != nil {
			logger.Log("Error: %v\n", err)
			continue
		}
		err = json.NewEncoder(conn).Encode(map[string]interface{}{
			"id":   id,
			"data": "Hello from car",
		})

		if err != nil {
			logger.Log("Error: %v\n", err)
		}
		logger.Log("Sent message\n")
		conn.Close()
	}
}

func initUnixWriteSocket(id int) (net.Conn, error) {
	addr := fmt.Sprintf("/tmp/car%dwrite.socket", id)
	logger.Log("Connecting to %s\n", addr)
	conn, err := net.Dial("unixgram", addr)
	if err != nil {
		logger.Log("Error: %v\n", err)
		return nil, err
	}
	logger.Log("Connected\n")
	return conn, nil
}

func main() {
	var id int
	flag.IntVar(&id, "id", 0, "ID of the car")
	flag.BoolVar(&logger.debug, "debug", false, "Debug mode")
	flag.Parse()

	l, file := InitLogger(id)
	logger = l
	if file != nil {
		defer file.Close()
	}

	socketAddress := fmt.Sprintf("/tmp/car%d.socket", id)
	err := os.RemoveAll(socketAddress)
	if err != nil {
		logger.Log("Error: %v", err)
		os.Exit(1)
	}

	addr, err := net.ResolveUnixAddr("unixgram", socketAddress)
	if err != nil {
		logger.Log("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		logger.Log("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	d := json.NewDecoder(conn)
	logger.Log("Listening to %s ..\n", socketAddress)

	go testWrite(id)

	reader(d)
}

func InitLogger(id int) (Logger, *os.File) {
	if !logger.debug {
		return Logger{false, log.New(os.Stdout, "", 0)}, nil
	}
	err := os.MkdirAll("/logs", 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating logs directory: %s\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(fmt.Sprintf("/logs/car%d.log", id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}
	return Logger{true, log.New(file, "", log.LstdFlags)}, file
}

func (logger Logger) Log(format string, v ...any) {
	if logger.debug {
		logger.Printf(format, v...)
	}
}
