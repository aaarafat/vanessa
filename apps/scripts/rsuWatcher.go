package main

import (
	"log"
	"os"
	"os/exec"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"strings"
	"github.com/juju/fslock"
	// "time"
)

func main() {
	fileName := "./"+os.Args[1]+".log"

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
				getLastLine(fileName)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			}
		}

	}()

	err = watcher.Add(fileName)
	if err != nil {
		log.Fatal("Add failed:", err)
	}
	<-done
	//todo
	//delete files when closing
}

func getLastLine(path string) {
    c := exec.Command("tail", "-1" , path)
    output, _ := c.Output()
	out := strings.ReplaceAll(string(output), "\n", "")
    fmt.Println(string(out))
}


func writeToRSUs(fileName, msg string) {
    lock := fslock.New(fileName)
    lockErr := lock.TryLock()
    if lockErr != nil {
        fmt.Println("falied to acquire lock > " + lockErr.Error())
        return
    }

    

    // release the lock
    lock.Unlock()
}