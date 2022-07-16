package app

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/aaarafat/vanessa/apps/car/unix"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/libs/vector"
	"github.com/cornelk/hashmap"
)

type App struct {
	id  int
	ip  net.IP
	key []byte

	// state
	state            *unix.State
	zoneTable        *ZoneTable
	stateLock        *sync.RWMutex
	checkRouteBuffer *hashmap.HashMap

	// to send messages to the network
	ipConn *ip.IPConnection

	// to connect to the simulator (read sensor data)
	sensor *unix.SensorUnix

	// to connect to the router
	router *unix.Router

	// to connect to the car ui
	ui *unix.UiUnix
}

func NewApp(id int, key []byte) *App {
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

	app := App{
		id:               id,
		ip:               ip,
		key:              key,
		zoneTable:        NewZoneTable(),
		ipConn:           ipConn,
		sensor:           unix.NewSensorUnix(id),
		router:           unix.NewRouter(id),
		stateLock:        &sync.RWMutex{},
		checkRouteBuffer: &hashmap.HashMap{},
	}

	app.ui = unix.NewUiUnix(id, app.GetState)

	app.initState(0, []Position{}, Position{Lng: 0, Lat: 0}, false)

	return &app
}

func (a *App) printZone() {
	for {
		a.zoneTable.Print()
		time.Sleep(time.Second * 2)
	}
}

func (a *App) checkFront() {
	state := a.GetState()
	pos := state.GetPosition()
	mnSpeed := state.MaxSpeed
	front := a.zoneTable.GetNearestFrontFrom(&pos)
	if front != nil {
		log.Printf("App %d: Front ", a.id)
		front.Print()
	}
	if front != nil && front.Speed < mnSpeed {
		mnSpeed = front.Speed
		if mnSpeed < state.Speed {
			// send check route to simulator
			go a.sendCheckRoute(front.Position, front.IP)
		}
	}
	if mnSpeed != state.Speed {
		a.updateSpeed(mnSpeed)
	}
}

func (a *App) checkZone() {
	for {
		a.checkFront()
		time.Sleep(ZONE_MSG_INTERVAL_MS * time.Millisecond)
	}
}

func (a *App) Run() {
	log.Printf("App %d starting.....", a.id)
	a.startSocketHandlers()
	a.sensor.Start()
	a.router.Start()
	go a.listen()
	go a.sendHeartBeat()
	go a.sendZoneMsg()
	go a.ui.Start()
	go a.printZone()
	go a.checkZone()

	//! Wait 2 seconds until the router is started
	time.Sleep(time.Second * 2)
	a.sensor.Write("", unix.MoveEvent)

	log.Printf("App %d started", a.id)
}

func (a *App) Stop() {
	log.Printf("App %d stopping", a.id)
	a.ipConn.Close()
	a.sensor.Close()
	a.router.Close()
	a.ui.Close()
	log.Printf("App %d stopped", a.id)
}
