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
	case VOREQType:
		r.handleVOREQ(payload,from)
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
	log.Println("Recieved HeartBeat from: ", HBeat.OriginatorIP.String())
	r.RARP.Set(HBeat.OriginatorIP.String(), from)

}

func (r* RSU)  handleVObstacle(payload []byte , from net.HardwareAddr) {

	obstacle, err := UnmarshalVObstacle(payload)
	if err != nil {
		log.Println("Failed to unmarshal VObstacle: ", err)
		return
	}
	log.Println("Recieved Obstacle from: ", obstacle.OriginatorIP.String())
	r.RARP.Set(obstacle.OriginatorIP.String(), from)
	r.ethChannel.Broadcast(obstacle.Marshal())
	r.wlanChannel.Broadcast(obstacle.Marshal())
	r.OTable.Set(obstacle.Position,0)
}

// handle VOREQ request 
func (r* RSU) handleVOREQ(payload []byte, from net.HardwareAddr) {
	VOREQ, err := UnmarshalVOREQ(payload)
	if err != nil {
		log.Println("Failed to unmarshal VOREQ: ", err)
		return
	}
	r.RARP.Set(VOREQ.OriginatorIP.String(), from)
	log.Println("Recieved VOREQ from: ", VOREQ.OriginatorIP.String())
	
	VOREP := NewVOREPMessage(r.OTable.GetTable())
	log.Println("Send VOREP to: ", VOREQ.OriginatorIP.String())
	r.wlanChannel.SendTo(VOREP.Marshal(), from)
}
