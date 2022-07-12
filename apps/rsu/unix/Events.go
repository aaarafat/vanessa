package unix

import (
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type Event string

const (
	AddObstacleEvent           Event = "add-obstacle"
	AddARPEntryEvent           Event = "add-arp-entry"
	UpdateReceivedPacketsEvent Event = "update-received-packets"
	UpdateSentPacketsEvent     Event = "update-sent-packets"
)

type AddObstacleData struct {
	Obstacles []Position `json:"obstacles"`
}

type AddARPEntryData = UiARPEntry

type UpdateReceivedPacketsData struct {
	ReceivedFromRsus int `json:"receivedFromRsus"`
	ReceivedFromCars int `json:"receivedFromCars"`
}

type UpdateSentPacketsData struct {
	SentToRsus int `json:"sentToRsus"`
	SentToCars int `json:"sentToCars"`
}
