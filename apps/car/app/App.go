package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/car/unix"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type App struct {
	id int
	ip net.IP

	// state
	state *unix.State

	// to send messages to the network
	ipConn *ip.IPConnection

	// to connect to the simulator (read sensor data)
	sensor *unix.SensorUnix

	// to connect to the router 
	router *unix.Router

	// to connect to the car ui
	ui *unix.UiUnix
}

func NewApp(id int) *App {
	ipConn, err := ip.NewIPConnection()
	if err != nil {
		log.Fatalf("Error creating IP connection: %v", err)
		return nil
	}

	ip, _, err := MyIP()
	ip = ip.To4()
	if err != nil {
		log.Fatalf("Error getting IP: %v", err)
		return nil
	}

	return &App{
		id: id, 
		ip: ip, 
		ipConn: ipConn, 
		sensor: unix.NewSensorUnix(id), 
		router: unix.NewRouter(id),
	}
}

func (a *App) initState() {
	log.Printf("Initializing car state......")
	a.state = &unix.State{
		Id: a.id,
		Speed: 0,
		Route: []unix.Coordinate{},
		Lat: 0,
		Lng: 0,
		ObstacleDetected: false,
		Obstacles: []unix.Coordinate{},
	}
	log.Printf("Car state initialized  state:  %v\n", a.state)
}

func (a *App) updatePosition(pos *Position) {
	a.state.Lat = pos.Lat
	a.state.Lng = pos.Lng
	log.Printf("Position updated: lng: %f lat: %f", pos.Lng, pos.Lat)
}

func (a *App) Run() {
	log.Printf("App %d starting.....", a.id)
	a.initState()
	go a.sensor.Start()
	go a.router.Start()
	go a.ui.Start()
	go a.startSocketHandlers()
	go a.sendHeartBeat()
	go a.listen()

	log.Printf("App %d started", a.id)
}

func (a *App) Stop() {
	log.Printf("App %d stopping", a.id)
	a.ipConn.Close()
	a.ui.Close()
	log.Printf("App %d stopped", a.id)
}
