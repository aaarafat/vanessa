package unix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	. "github.com/aaarafat/vanessa/apps/rsu/state"
	"gopkg.in/antage/eventsource.v1"
)

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

	http.HandleFunc("/state", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		json.NewEncoder(w).Encode(*u.getState())
	})
	http.Handle("/", u.server)

	log.Fatal(http.Serve(listener, nil))
}

func (u *UiUnix) Close() {
	u.server.Close()
}
