package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/car/unix"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
)

type App struct {
	id int
	ip net.IP

	// car data
	position *unix.Position

	// to send messages to the network
	ipConn *ip.IPConnection

	// to connect to the simulator (read sensor data)
	unix *unix.UnixSocket
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

	return &App{id: id, ip: ip, unix: unix.NewUnixSocket(id), ipConn: ipConn}
}

func (a *App) updatePosition(pos *unix.Position) {
	a.position = pos
	log.Printf("Position updated: lng: %f lat: %f", pos.Lng, pos.Lat)
}

func (a *App) Run() {
	log.Printf("App %d starting.....", a.id)
	go a.unix.Start()
	go a.startSocketHandlers()
	log.Printf("App %d started", a.id)
}

func (a *App) Stop() {
	log.Printf("App %d stopping", a.id)
	a.ipConn.Close()
	log.Printf("App %d stopped", a.id)
}
