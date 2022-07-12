package state

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	. "github.com/aaarafat/vanessa/apps/rsu/router"
	"github.com/aaarafat/vanessa/apps/rsu/unix"
)

type State struct {
	RARP            *RSUARP         `json:"arp"`
	OTable          *ObstaclesTable `json:"obstacles"`
	SentPackets     [2]int          `json:"sentPackets"`     // [0] for eth, [1] for wlan
	RecievedPackets [2]int          `json:"recievedPackets"` // [0] for eth, [1] for wlan

	ui *unix.UiUnix
}

func NewState(ui *unix.UiUnix) *State {
	state := &State{
		OTable:          NewObstaclesTable(),
		SentPackets:     [2]int{0, 0},
		RecievedPackets: [2]int{0, 0},
		ui:              ui,
	}

	state.RARP = NewRSUARP(state.RemoveCar)

	return state

}

func (s *State) ReceivedPacket(interfaceType int) {
	s.RecievedPackets[interfaceType]++

	s.ui.Write(unix.UpdateReceivedPacketsData{
		ReceivedFromRsus: s.RecievedPackets[0],
		ReceivedFromCars: s.RecievedPackets[1],
	}, string(unix.UpdateReceivedPacketsEvent))
}

func (s *State) SentPacket(interfaceType int, count int) {
	s.SentPackets[interfaceType] += count

	s.ui.Write(unix.UpdateSentPacketsData{
		SentToRsus: s.SentPackets[0],
		SentToCars: s.SentPackets[1],
	}, string(unix.UpdateSentPacketsEvent))
}

func (s *State) AddObstacle(obstacle *Position) {
	log.Printf("Adding obstacle: %v\n", obstacle)
	s.OTable.Set(*obstacle, 0)

	s.ui.Write(unix.AddObstacleData{
		Obstacles: s.OTable.GetTable(),
	}, string(unix.AddObstacleEvent))
}

func (s *State) RemoveObstacle(obstacle *Position) {
	log.Printf("Removing obstacle: %v\n", obstacle)
	s.OTable.Set(*obstacle, 1)
}

func (s *State) AddCar(ip string, mac net.HardwareAddr) {
	if s.RARP.Set(ip, mac) {
		s.ui.Write(unix.AddARPEntryData{
			IP:  ip,
			MAC: mac.String(),
		}, string(unix.AddARPEntryEvent))
	}
}

func (s *State) RemoveCar(ip string, mac net.HardwareAddr) {
	s.RARP.Del(ip)

	s.ui.Write(unix.RemoveARPEntryData{
		IP:  ip,
		MAC: mac.String(),
	}, string(unix.RemoveARPEntryEvent))
}

func (s *State) HasObstacle(obstacle *Position) bool {
	_, ok := s.OTable.table[string(obstacle.Marshal())]
	return ok
}
