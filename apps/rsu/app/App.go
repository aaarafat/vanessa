package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/rsu/router"
)

type App struct {
	ip     net.IP
	key    []byte
	router *Router
	state  *State
}

func NewApp(key []byte) *App {
	router := NewRouter()
	return &App{
		ip:     net.ParseIP(ip.RsuIP),
		key:    key,
		router: router,
		state:  NewState(router.RARP),
	}
}

func (a *App) Start() {
	log.Printf("RSU is starting.....\n")
	go a.handleEthMessages()
	go a.handleWLANMessgaes()
	log.Printf("RSU is started!!!\n")
}

func (a *App) Close() {
	log.Printf("RSU is closing.....\n")
	a.router.Close()
	log.Printf("RSU is closed!!!\n")
}
