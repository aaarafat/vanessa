package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/rsu/router"
	. "github.com/aaarafat/vanessa/apps/rsu/state"
	"github.com/aaarafat/vanessa/apps/rsu/unix"
)

type App struct {
	id     int
	ip     net.IP
	key    []byte
	router *Router
	state  *State

	// to connect to the rsu ui
	ui *unix.UiUnix
}

func NewApp(id int, key []byte) *App {
	router := NewRouter()
	app := &App{
		id:     id,
		ip:     net.ParseIP(ip.RsuIP),
		key:    key,
		router: router,
		state:  NewState(router.RARP),
	}

	app.ui = unix.NewUiUnix(id, app.GetState)

	return app
}

func (a *App) GetState() *State {
	return a.state
}

func (a *App) Start() {
	log.Printf("RSU is starting.....\n")
	go a.ui.Start()
	go a.handleEthMessages()
	go a.handleWLANMessgaes()
	log.Printf("RSU is started!!!\n")
}

func (a *App) Close() {
	log.Printf("RSU is closing.....\n")
	a.ui.Close()
	a.router.Close()
	log.Printf("RSU is closed!!!\n")
}
