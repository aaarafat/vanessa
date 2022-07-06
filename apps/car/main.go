package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func initLogger(debug bool, id int) {
	log.SetPrefix("[vanessa]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !debug {
		log.SetOutput(os.Stdout)
		return // don't do anything if debug is false
	}
	err := os.MkdirAll("/logs", 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating logs directory: %s\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(fmt.Sprintf("/logs/car%d.log", id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}
	log.SetOutput(file)
}

func main() {
	var id int
	var debug bool
	flag.IntVar(&id, "id", 0, "id of the car")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.Parse()

	initLogger(debug, id)
}
