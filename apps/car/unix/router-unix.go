package unix

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Router struct {
	id int

	Data *chan []byte
}

func NewRouter(id int) *Router {
	ch := make(chan []byte)
	return &Router{id: id, Data: &ch}
}


func (r *Router) reader(conn net.Conn) {
	for {
		data := make([]byte, 1024)
		n, err := conn.Read(data)
		if err != nil {
			log.Printf("Error reading from connection: %v", err)
			continue
		}
		log.Printf("Received from router: %s", string(data[:n]))
		*r.Data <- data[:n]
	}
}

func (r *Router) Start() {
	socketAddress := fmt.Sprintf("/tmp/car%d-router.socket", r.id)
	err := os.RemoveAll(socketAddress)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	addr, err := net.ResolveUnixAddr("unixgram", socketAddress)
	if err != nil {
		log.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		log.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Printf("Listening to %s ..\n", socketAddress)

	r.reader(conn)
}