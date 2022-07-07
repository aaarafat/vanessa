package app

import (
	"log"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (a *App) listen() {
	select {
	case data := <-*a.router.Data:
		go a.handleMessage(data)
	}
}

func (a *App) handleMessage(data []byte) {
	mType := data[0]
	switch mType {
	case VOREPType:
		msg := UnmarshalVOREP(data)
		log.Printf("VOREP message received: %s", msg.String())
	}
}