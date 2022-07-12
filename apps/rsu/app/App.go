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

const (
	ETH_IFI  = 0
	WLAN_IFI = 1
)

func NewApp(id int, key []byte) *App {
	router := NewRouter()
	app := &App{
		id:     id,
		ip:     net.ParseIP(ip.RsuIP),
		key:    key,
		router: router,
		state:  NewState(router.RARP),
	}

	app.ui = unix.NewUiUnix(id, app.GetUiState)

	return app
}

func (a *App) GetUiState() *unix.UiState {
	arp := make([]unix.UiARPEntry, a.state.RARP.Len())
	i := 0
	for ip, entry := range a.state.RARP.GetTable() {
		arp[i] = unix.UiARPEntry{
			IP:  ip,
			MAC: entry.MAC.String(),
		}
		i++
	}

	return &unix.UiState{
		ARP:              arp,
		Obstacles:        a.state.OTable.GetTable(),
		ReceivedFromRsus: a.state.RecievedPackets[ETH_IFI],
		SentToRsus:       a.state.SentPackets[ETH_IFI],
		ReceivedFromCars: a.state.RecievedPackets[WLAN_IFI],
		SentToCars:       a.state.SentPackets[WLAN_IFI],
	}
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
