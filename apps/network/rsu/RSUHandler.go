package rsu

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)


func (r* RSU) handleMessage(packet ip.IPPacket , from net.HardwareAddr) {

	payload := packet.Payload
	msgType := uint8(payload[0])
	// handle the message
	switch msgType {
	case VHBeatType:
		r.handleVHBeat(payload,from)
	case VObstacleType:
		r.handleVObstacle(packet,from)
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

func (r* RSU)  handleVObstacle(packet ip.IPPacket , from net.HardwareAddr) {
	payload := packet.Payload
	obstacle, err := UnmarshalVObstacle(payload)
	if err != nil {
		log.Println("Failed to unmarshal VObstacle: ", err)
		return
	}
	log.Println("Recieved Obstacle from: ", obstacle.OriginatorIP.String(), " at: ", obstacle.Position)
	r.RARP.Set(obstacle.OriginatorIP.String(), from)
	bytes := ip.MarshalIPPacket(&packet)
	ip.Update(bytes)
	r.ethChannel.Broadcast(bytes)
	r.sendToALLWLANInterface(bytes, obstacle.OriginatorIP.String())
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

	packet := ip.NewIPPacket(VOREP.Marshal(), r.ip, net.IPv4(255,255,255,255))
	bytes := ip.MarshalIPPacket(packet)
	ip.UpdateChecksum(bytes)
	r.wlanChannel.SendTo(bytes, from)
}


// Send to all in RSUARP using wlan exept the one that sent the message
func (r* RSU) sendToALLWLANInterface(data []byte, originatorIP string) {
	for ip , entry := range r.RARP.table {
		if originatorIP == ip {
			continue
		}
		r.wlanChannel.SendTo(data, entry.MAC)
	}
}
