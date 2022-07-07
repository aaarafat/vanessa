package unix

import (
	"fmt"
	"log"
	"net"
	"time"
)

type RouterSocket struct {
	id int
	conn net.Conn
}

func NewRouterSocket(id int) *RouterSocket {
	conn, err := initSocket(id)
	if err != nil {
		return nil
	}
	return &RouterSocket{id: id, conn: conn}
}

func initSocket(id int) (net.Conn, error) {
	time.Sleep(time.Second * 1)
	addr := fmt.Sprintf("/tmp/car%d-router.socket", id)
	log.Printf("Connecting to %s\n", addr)
	conn, err := net.Dial("unixgram", addr)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return nil, err
	}
	log.Printf("Connected to %s\n", addr)
	return conn, nil
}

func (a *RouterSocket) Write(data []byte) {
	n, err := a.conn.Write(data)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}
	log.Printf("Sent %d bytes to the Car Router\n", n)
}

func (a *RouterSocket) Close() {
	a.conn.Close()
}