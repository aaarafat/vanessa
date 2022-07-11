package app

import (
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type State struct {
	RARP            *RSUARP
	OTable          *ObstaclesTable
	SentPackets     [2]int // [0] for eth, [1] for wlan
	RecievedPackets [2]int // [0] for eth, [1] for wlan
}

func NewState() *State {
	return &State{
		RARP:            NewRSUARP(),
		OTable:          NewObstaclesTable(),
		SentPackets:     [2]int{0, 0},
		RecievedPackets: [2]int{0, 0},
	}
}

func (s *State) ReceivedPacket(interfaceType int) {
	s.RecievedPackets[interfaceType]++
}

func (s *State) SentPacket(interfaceType int) {
	s.SentPackets[interfaceType]++
}

func (s *State) AddObstacle(obstacle *Position) {
	s.OTable.Set(*obstacle, 0)
}

func (s *State) RemoveObstacle(obstacle *Position) {
	s.OTable.Set(*obstacle, 1)
}

func (s *State) AddCar(ip string, mac net.HardwareAddr) {
	s.RARP.Set(ip, mac)
}
