package unix

import (
	"fmt"
	"log"
	"net"
	"os"
)

type Router struct {
	id int

	conn *net.UnixConn
}

func NewRouter(id int) *Router {
	return &Router{id: id}
}

func (r *Router) Read() ([]byte, error) {
	data := make([]byte, 1500)
	n, err := r.conn.Read(data)
	if err != nil {
		log.Printf("Error reading from connection: %v", err)
		return nil, err
	}
	log.Printf("Received from router data with size %d", n)
	return data[:n], nil
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

	r.conn = conn
	log.Printf("Listening to %s ..\n", socketAddress)
}

func (r *Router) Close() {
	r.conn.Close()
}
