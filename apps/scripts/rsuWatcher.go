package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/fsnotify/fsnotify"

	// "github.com/juju/fslock"
	"database/sql"
	// "time"
)

func main() {
	const create string = `
		CREATE TABLE IF NOT EXISTS RSUs (
		sender TEXT NOT NULL PRIMARY KEY,
		msg TEXT
		);`
	lFileName := "./" + os.Args[1] + ".log"
	gFileName := "./rsus.db"
	db, err := sql.Open("sqlite3", gFileName)
	db.Exec(create)
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal("NewWatcher failed: ", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		defer close(done)

		for {
			select {
			case _, ok := <-watcher.Events:
				if !ok {
					return
				}
				msg := getLastLine(lFileName)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	err = watcher.Add(lFileName)
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
	//todo
	//delete files when closing
}

func getLastLine(path string) (msg string) {
	c := exec.Command("tail", "-1", path)
	output, _ := c.Output()
	msg = strings.ReplaceAll(string(output), "\n", "")
	fmt.Println(msg)
	return msg
}
