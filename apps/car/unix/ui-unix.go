package unix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	. "github.com/aaarafat/vanessa/libs/vector"
	"gopkg.in/antage/eventsource.v1"
)

type Coordinate struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type State struct {
	Id                 int        `json:"id"`
	Speed              uint32     `json:"speed"`
	Route              []Position `json:"route"`
	Lat                float64    `json:"lat"`
	Lng                float64    `json:"lng"`
	ObstacleDetected   bool       `json:"obstacleDetected"`
	Obstacles          []Position `json:"obstacles"`
	MaxSpeed           uint32     `json:"maxSpeed"`
	Direction          Vector     `json:"direction"`
	DestinationReached bool       `json:"destinationReached"`
	Stopped            bool       `json:"stopped"`
}

func (s *State) GetPosition() Position {
	return Position{Lng: s.Lng, Lat: s.Lat}
}

func (s *State) String() string {
	return fmt.Sprintf("State{Id: %d, Speed: %d, Route: %v, Lat: %f, Lng: %f, ObstacleDetected: %t, Obstacles: %v, MaxSpeed: %d, Direction: %v, stopped: %v}", s.Id, s.Speed, s.Route, s.Lat, s.Lng, s.ObstacleDetected, s.Obstacles, s.MaxSpeed, s.Direction, s.Stopped)
}

type UiUnix struct {
	id       int
	addr     string
	server   eventsource.EventSource
	getState func() *State
}

func NewUiUnix(id int, getState func() *State) *UiUnix {
	server := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{[]byte("Access-Control-Allow-Origin: *")}
	})

	addr := fmt.Sprintf("/tmp/car%d.ui.socket", id)

	return &UiUnix{id: id, addr: addr, server: server, getState: getState}
}

func (u *UiUnix) Write(message any, eventName string) {
	payload, err := json.Marshal(map[string]interface{}{
		"id":   u.id,
		"data": message,
	})
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	u.server.SendEventMessage(string(payload), eventName, "")
	log.Printf("Sent event %s to ui\n", eventName)
}

func (u *UiUnix) initUnixSocket() (net.Listener, error) {
	log.Printf("Connecting to %s\n", u.addr)
	listener, err := net.Listen("unix", u.addr)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	log.Printf("Connected\n")
	return listener, nil
}

func (u *UiUnix) Start() {
	err := os.RemoveAll(u.addr)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	time.Sleep(time.Millisecond * 100)
	listener, err := u.initUnixSocket()
	defer listener.Close()

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(*u.getState())
	})
	http.Handle("/", u.server)

	go func() {
		time.Sleep(time.Second * 10)
		// send refresh message to the ui
		u.Refresh(*u.getState())
	}()

	log.Fatal(http.Serve(listener, nil))
}

func (u *UiUnix) Refresh(state State) {
	u.Write(state, string(RefreshEvent))
}

func (u *UiUnix) Close() {
	u.server.Close()
}
