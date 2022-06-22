package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
)

type Logger struct {
	debug bool
	*log.Logger
}

var logger Logger

type Position struct {
	Lat float64
	Lng float64
}

func reader(d *json.Decoder) {
	var position = Position{Lat: 5, Lng: 5}
	for {
		err := d.Decode(&position)
		if err != nil {
			return
		}
		logger.Log("received: %f %f\n", position.Lat, position.Lng)
	}
}

func main() {
	var id int
	flag.IntVar(&id, "id", 0, "ID of the car")
	flag.BoolVar(&logger.debug, "debug", false, "Debug mode")
	flag.Parse()

	l, file := InitLogger(id)
	logger = l
	if file != nil {
		defer file.Close()
	}

	socketAddress := fmt.Sprintf("/tmp/car%d.socket", id)
	err := os.RemoveAll(socketAddress)
	if err != nil {
		logger.Log("Error: %v", err)
		os.Exit(1)
	}

	addr, err := net.ResolveUnixAddr("unix", socketAddress)
	if err != nil {
		logger.Log("Failed to resolve: %v\n", err)
		os.Exit(1)
	}

	conn, err := net.ListenUnixgram("unixgram", addr)
	if err != nil {
		logger.Log("Failed to resolve: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	d := json.NewDecoder(conn)
	logger.Log("Listening to %s ..\n", socketAddress)
	reader(d)
}

func InitLogger(id int) (Logger, *os.File) {
	if !logger.debug {
		return Logger{false, log.New(os.Stdout, "", 0)}, nil
	}
	err := os.MkdirAll("/logs", 0644)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating logs directory: %s\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(fmt.Sprintf("/logs/car%d.log", id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}
	return Logger{true, log.New(file, "", log.LstdFlags)}, file
}

func (logger Logger) Log(format string, v ...any) {
	if logger.debug {
		logger.Printf(format, v...)
	}
}
