package unix

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"github.com/cenkalti/backoff/v4"
)

type Subscriber struct {
	Messages *chan json.RawMessage
}

type SensorUnix struct {
	id     int
	server *net.UnixListener
	conn   *net.Conn
	topics map[Event][]*Subscriber
}

func NewSensorUnix(id int) *SensorUnix {
	unix := &SensorUnix{id: id, topics: make(map[Event][]*Subscriber)}
	return unix
}

func (u *SensorUnix) Subscribe(topic Event, subscriber *Subscriber) {
	if u.topics[topic] == nil {
		u.topics[topic] = []*Subscriber{}
	}
	u.topics[topic] = append(u.topics[topic], subscriber)
	log.Printf("Subscribed to %s\n", topic)
}

func (u *SensorUnix) publish(topic Event, message json.RawMessage) {
	if u.topics[topic] == nil {
		log.Printf("No subscribers for %s\n", topic)
		return
	}
	for _, subscriber := range u.topics[topic] {
		*subscriber.Messages <- message
	}
}

func (unix *SensorUnix) reader(d *json.Decoder) {
	var m map[string]json.RawMessage
	err := d.Decode(&m)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	var eventType Event
	err = json.Unmarshal(m["type"], &eventType)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return
	}

	switch eventType {
	case DestinationReachedEvent:
		var p DestinationReachedData
		err := json.Unmarshal(m["data"], &p)
		if err != nil {
			log.Printf("Error decoding destination-reached data: %v", err)
			return
		}
		unix.publish(DestinationReachedEvent, m["data"])

	case ObstacleDetectedEvent:
		var p ObstacleDetectedData
		err := json.Unmarshal(m["data"], &p)
		if err != nil {
			log.Printf("Error decoding obstacle-detected data: %v", err)
			return
		}
		unix.publish(ObstacleDetectedEvent, m["data"])

	case AddCarEvent:
		var p AddCarData
		err := json.Unmarshal(m["data"], &p)
		if err != nil {
			log.Printf("Error decoding add-car data: %v", err)
			return
		}
		unix.publish(AddCarEvent, m["data"])

	case UpdateLocationEvent:
		var p UpdateLocationData
		err := json.Unmarshal(m["data"], &p)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		unix.publish(UpdateLocationEvent, m["data"])

	case ChangeSpeedEvent:
		var p SpeedData
		err := json.Unmarshal(m["data"], &p)
		if err != nil {
			log.Printf("Error: %v", err)
			return
		}
		unix.publish(ChangeSpeedEvent, m["data"])
	}
}

func (unix *SensorUnix) Write(message any, event Event) {
	backoff.Retry(func() error {
		log.Printf("Writing to socket...\n")
		if unix.conn == nil {
			err := unix.initUnixWriteSocket()
			if err != nil {
				return err
			}
		}

		err := json.NewEncoder(*unix.conn).Encode(map[string]interface{}{
			"id":   unix.id,
			"type": event,
			"data": message,
		})
		if err != nil {
			log.Printf("Error: %v\n", err)
			unix.initUnixWriteSocket()
			return err
		}

		return nil
	}, backoff.WithMaxRetries(backoff.NewConstantBackOff(100*time.Millisecond), 10))
}

func (unix *SensorUnix) initUnixWriteSocket() error {
	addr := fmt.Sprintf("/tmp/car%dwrite.socket", unix.id)
	log.Printf("Connecting to %s\n", addr)
	conn, err := net.Dial("unix", addr)
	if err != nil {
		log.Printf("Error: %v\n", err)
		return err
	}
	unix.conn = &conn
	log.Printf("Connected\n")
	return nil
}

func (unix *SensorUnix) Start() {
	socketAddress := fmt.Sprintf("/tmp/car%d.socket", unix.id)
	err := os.RemoveAll(socketAddress)
	if err != nil {
		log.Printf("Error: %v", err)
		os.Exit(1)
	}

	addr, err := net.ResolveUnixAddr("unix", socketAddress)
	if err != nil {
		log.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	server, err := net.ListenUnix("unix", addr)
	if err != nil {
		log.Printf("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	unix.server = server

	log.Printf("Listening to %s ..\n", socketAddress)

	go func() {
		for {
			conn, err := server.Accept()
			if err != nil {
				log.Printf("Error: %v\n", err)
				continue
			}
			d := json.NewDecoder(conn)
			unix.reader(d)
			conn.Close()
		}
	}()
}

func (u *SensorUnix) Close() {
	u.server.Close()
	if u.conn != nil {
		(*u.conn).Close()
	}
}
