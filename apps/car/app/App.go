package app

import (
	"log"
	"net"
	"sync"
	"time"

	"github.com/aaarafat/vanessa/apps/car/unix"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/libs/vector"
)

type App struct {
	id  int
	ip  net.IP
	key []byte

	// state
	state     *unix.State
	zoneTable *ZoneTable
	stateLock *sync.RWMutex

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
		id:        id,
		ip:        ip,
		key:       key,
		zoneTable: NewZoneTable(),
		ipConn:    ipConn,
		sensor:    unix.NewSensorUnix(id),
		router:    unix.NewRouter(id),
		stateLock: &sync.RWMutex{},
	}

	app.ui = unix.NewUiUnix(id, app.GetState)

	app.initState(0, []Position{}, Position{Lng: 0, Lat: 0})

	return &app
}

func (a *App) printZone() {
	for {
		a.zoneTable.Print()
		time.Sleep(time.Second * 2)
	}
}

func (a *App) checkZone() {
	for {
		// check front
		front := a.zoneTable.GetInFrontOfMe()
		mnSpeed := a.GetState().MaxSpeed
		for _, entry := range front {
			if entry.Speed < mnSpeed {
				mnSpeed = entry.Speed
			}
		}
		a.updateSpeed(mnSpeed)

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
