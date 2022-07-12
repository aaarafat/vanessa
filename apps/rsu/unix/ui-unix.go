package unix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	"gopkg.in/antage/eventsource.v1"
)

type UiUnix struct {
	id       int
	addr     string
	server   eventsource.EventSource
	getState func() *UiState
}

type UiARPEntry struct {
	IP  string `json:"ip"`
	MAC string `json:"mac"`
}

type UiState struct {
	Id               int          `json:"id"`
	ARP              []UiARPEntry `json:"arp"`
	Obstacles        []Position   `json:"obstacles"`
	ReceivedFromRsus int          `json:"receivedFromRsus"`
	SentToRsus       int          `json:"sentToRsus"`
	ReceivedFromCars int          `json:"receivedFromCars"`
	SentToCars       int          `json:"sentToCars"`
}

func NewUiUnix(id int, getState func() *UiState) *UiUnix {
	server := eventsource.New(nil, func(req *http.Request) [][]byte {
		return [][]byte{[]byte("Access-Control-Allow-Origin: *")}
	})

	addr := fmt.Sprintf("/tmp/rsu%d.ui.socket", id)

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

	http.HandleFunc("/state-rsu", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(*u.getState())
	})
	http.Handle("/", u.server)

	// send refresh message to the ui
	u.Refresh(*u.getState())

	log.Fatal(http.Serve(listener, nil))
}

func (u *UiUnix) Refresh(state UiState) {
	u.Write(state, string(RefreshEvent))
}

func (u *UiUnix) Close() {
	u.server.Close()
}
