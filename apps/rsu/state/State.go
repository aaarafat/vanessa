package state

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	. "github.com/aaarafat/vanessa/apps/rsu/router"
)

type State struct {
	RARP            *RSUARP         `json:"arp"`
	OTable          *ObstaclesTable `json:"obstacles"`
	SentPackets     [2]int          `json:"sentPackets"`     // [0] for eth, [1] for wlan
	RecievedPackets [2]int          `json:"recievedPackets"` // [0] for eth, [1] for wlan
}

func NewState(arp *RSUARP) *State {
	return &State{
		RARP:            arp,
		OTable:          NewObstaclesTable(),
		SentPackets:     [2]int{0, 0},
		RecievedPackets: [2]int{0, 0},
	}
}

func (s *State) ReceivedPacket(interfaceType int) {
	s.RecievedPackets[interfaceType]++
}

func (s *State) SentPacket(interfaceType int, count int) {
	s.SentPackets[interfaceType] += count
}

func (s *State) AddObstacle(obstacle *Position) {
	log.Printf("Adding obstacle: %v\n", obstacle)
	s.OTable.Set(*obstacle, 0)
}

func (s *State) RemoveObstacle(obstacle *Position) {
	log.Printf("Removing obstacle: %v\n", obstacle)
	s.OTable.Set(*obstacle, 1)
}

func (s *State) AddCar(ip string, mac net.HardwareAddr) {
	s.RARP.Set(ip, mac)
}

func (s *State) HasObstacle(obstacle *Position) bool {
	_, ok := s.OTable.table[string(obstacle.Marshal())]
	return ok
}
