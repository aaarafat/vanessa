package rsu

import (
	"log"
	"net"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


func (r* RSU) handleMessage(payload []byte , from net.HardwareAddr) {
	msgType := uint8(payload[0])
	// handle the message
	switch msgType {
	case VHBeatType:
		r.handleVHBeat(payload,from)
	case VObstacleType:
		r.handleVObstacle(payload,from)
	default:
		log.Println("Unknown message type: ", msgType)
	}
}

func (r* RSU) handleVHBeat(payload []byte, from net.HardwareAddr) {

	HBeat, err := UnmarshalVHBeat(payload)
	if err != nil {
		log.Println("Failed to unmarshal VHBeat: ", err)
		return
	}
	r.RARP.Set(HBeat.OriginatorIP.String(), from)

}

func (r* RSU)  handleVObstacle(payload []byte , from net.HardwareAddr) {

	obstacle, err := UnmarshalVObstacle(payload)
	if err != nil {
		log.Println("Failed to unmarshal VObstacle: ", err)
		return
	}
	r.ethChannel.Broadcast(obstacle.Marshal())
	r.wlanChannel.Broadcast(obstacle.Marshal())
	//TODO: Add to Obstacle List
}