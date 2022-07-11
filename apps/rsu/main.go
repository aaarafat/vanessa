package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	. "github.com/aaarafat/vanessa/apps/rsu/app"
)

func initLogger(debug bool, id int) {
	log.SetPrefix("[vanessa]")
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	if !debug {
		log.SetOutput(os.Stdout)
		return // don't do anything if debug is false
	}
	// delete logs
	err := os.MkdirAll("/var/log/vanessa", 0777)
	if err != nil && !os.IsExist(err) {
		fmt.Printf("Error creating logs directory: %s\n", err)
		os.Exit(1)
	}

	file, err := os.OpenFile(fmt.Sprintf("/var/log/vanessa/rsu%d.log", id), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("Error opening log file: %s\n", err)
		os.Exit(1)
	}
	log.SetOutput(file)
}

func main() {
	var id int
	var debug bool
	var keyStr string
	flag.IntVar(&id, "id", 0, "id of the car")
	flag.BoolVar(&debug, "debug", false, "debug mode")
	flag.StringVar(&keyStr, "key", "", "aes key")
	flag.Parse()

	initLogger(debug, id)

	key := make([]byte, 16)
	n, err := base64.StdEncoding.Decode(key, []byte(keyStr))
	if err != nil {
		log.Fatalf("failed to decode AES key: %v", err)
	}

	log.Printf("My key is %s, from %s written %d  len %d\n", key, keyStr, n, len(key))

	// create a new RSU
	rsu := NewApp(key)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		rsu.Close()
		os.Exit(1)
	}()

	// start the RSU
	go rsu.Start()
	defer rsu.Close()

	select {}
}
