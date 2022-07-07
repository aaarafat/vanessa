package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"gopkg.in/antage/eventsource.v1"
)

type Logger struct {
	debug bool
	*log.Logger
}

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type State struct {
	Id               int          `json:"id"`
	Speed            int          `json:"speed"`
	Route            []Coordinate `json:"route"`
	Lat              float64      `json:"lat"`
	Lng              float64      `json:"lng"`
	ObstacleDetected bool         `json:"obstacleDetected"`
	Obstacles        []Coordinate `json:"obstacles"`
}

var logger Logger

func testWrite(server eventsource.EventSource, id int) {
	for {
		// sleep for 5 sec
		time.Sleep(time.Second * 5)

		// send message to all connected users
		logger.Log("Sent message\n")
		payload, err := json.Marshal(map[string]interface{}{
			"id":   id,
			"data": "Hello from car",
		})
		if err != nil {
			logger.Log("Error: %v\n", err)
			continue
		}
		server.SendEventMessage(string(payload), "test", "")

	}
}

func initUnixSocket(addr string) (net.Listener, error) {
	logger.Log("Connecting to %s\n", addr)
	listener, err := net.Listen("unix", addr)
	if err != nil {
		logger.Log("Error: %v\n", err)
		return nil, err
	}
	logger.Log("Connected\n")
	return listener, nil
}

func main() {
	var socketAddress string
	var id int
	flag.IntVar(&id, "id", 0, "ID of the car")
	flag.StringVar(&socketAddress, "s", "", "Unix socket of the car ui")
	flag.Parse()

	logger.debug = false
	l, f := InitLogger(id)
	logger = l
	if f != nil {
		defer f.Close()
	}

	err := os.RemoveAll(socketAddress)
	if err != nil {
		logger.Log("Error: %v", err)
		os.Exit(1)
	}

	state := State{
		Id:    id,
		Speed: 10,
		Route: []Coordinate{
			{
				Lng: 31.21151,
				Lat: 30.02163,
			},
			{
				Lng: 31.21145,
				Lat: 30.02219,
			},
			{
				Lng: 31.21141,
				Lat: 30.02252,
			},
			{
				Lng: 31.21119,
				Lat: 30.02377,
			},
		},
		Lat:              30.02163,
		Lng:              31.21151,
		ObstacleDetected: false,
		Obstacles:        []Coordinate{},
	}

	time.Sleep(time.Millisecond * 100)
	listener, err := initUnixSocket(socketAddress)
	defer listener.Close()

	server := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{[]byte("Access-Control-Allow-Origin: *")}
	})
	defer server.Close()

	go testWrite(server, id)

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(state)
	})
	http.Handle("/", server)

	log.Fatal(http.Serve(listener, nil))
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
