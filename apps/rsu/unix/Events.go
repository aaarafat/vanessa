package unix

import (
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type Event string

const (
	AddObstacleEvent           Event = "add-obstacle"
	AddARPEntryEvent           Event = "add-arp-entry"
	RemoveARPEntryEvent        Event = "remove-arp-entry"
	UpdateReceivedPacketsEvent Event = "update-received-packets"
	UpdateSentPacketsEvent     Event = "update-sent-packets"
	RefreshEvent               Event = "refresh"
)

type AddObstacleData struct {
	Obstacle Position `json:"obstacle"`
}

type AddARPEntryData = UiARPEntry
type RemoveARPEntryData = UiARPEntry

type UpdateReceivedPacketsData struct {
	ReceivedFromRsus int `json:"receivedFromRsus"`
	ReceivedFromCars int `json:"receivedFromCars"`
}

type UpdateSentPacketsData struct {
	SentToRsus int `json:"sentToRsus"`
	SentToCars int `json:"sentToCars"`
}

type RefreshData = UiState
