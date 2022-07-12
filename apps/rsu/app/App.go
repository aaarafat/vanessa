package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/aaarafat/vanessa/apps/network/network/messages"
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
	app := &App{
		id:  id,
		ip:  net.ParseIP(ip.RsuIP),
		key: key,
	}

	app.ui = unix.NewUiUnix(id, app.GetUiState)
	app.state = NewState(app.ui)
	app.router = NewRouter(app.state.RARP)

	return app
}

func (a *App) GetUiState() *unix.UiState {
	if a.state == nil {
		return &unix.UiState{
			Id:               a.id,
			ARP:              []unix.UiARPEntry{},
			Obstacles:        []messages.Position{},
			ReceivedFromRsus: 0,
			SentToRsus:       0,
			ReceivedFromCars: 0,
			SentToCars:       0,
		}
	}
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
		Id:               a.id,
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

	// send refresh message to the ui
	a.ui.Refresh(*a.GetUiState())

	log.Printf("RSU is started!!!\n")
}

func (a *App) Close() {
	log.Printf("RSU is closing.....\n")
	a.ui.Close()
	a.router.Close()
	log.Printf("RSU is closed!!!\n")
}
