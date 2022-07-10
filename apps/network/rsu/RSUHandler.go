package rsu

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (r *RSU) handleMessage(packet ip.IPPacket, from net.HardwareAddr) {

	payload := packet.Payload
	msgType := uint8(payload[0])
	// handle the message
	switch msgType {
	case VHBeatType:
		r.handleVHBeat(payload, from)
	case VObstacleType:
		r.handleVObstacle(packet, from)
	case VOREQType:
		r.handleVOREQ(payload, from)
	default:
		log.Println("Unknown message type: ", msgType)
	}
}

func (r *RSU) handleEthMessages(packet ip.IPPacket, from net.HardwareAddr) {
	payload := packet.Payload
	msgType := uint8(payload[0])
	// handle the message
	switch msgType {
	case VObstacleType:
		obstacle, err := UnmarshalVObstacle(payload)
		if err != nil {
			log.Println("Failed to unmarshal VObstacle: ", err)
			return
		}
		log.Println("Recieved Obstacle from: ", obstacle.OriginatorIP.String(), " at: ", obstacle.Position)
		r.sendToALLWLANInterface(payload, obstacle.OriginatorIP.String())
		r.OTable.Set(obstacle.Position, 0)
	default:
		log.Println("Unknown message type: ", msgType)
	}
}

func (r *RSU) handleVHBeat(payload []byte, from net.HardwareAddr) {

	HBeat, err := UnmarshalVHBeat(payload)
	if err != nil {
		log.Println("Failed to unmarshal VHBeat: ", err)
		return
	}
	log.Println("Recieved HeartBeat from: ", HBeat.OriginatorIP.String())
	r.RARP.Set(HBeat.OriginatorIP.String(), from)

}

func (r *RSU) handleVObstacle(packet ip.IPPacket, from net.HardwareAddr) {
	payload := packet.Payload
	obstacle, err := UnmarshalVObstacle(payload)
	if err != nil {
		log.Println("Failed to unmarshal VObstacle: ", err)
		return
	}
	log.Println("Recieved Obstacle from: ", obstacle.OriginatorIP.String(), " at: ", obstacle.Position)
	if _, ok := r.OTable.table[string(obstacle.Position.Marshal())]; ok {
		log.Println("Not Sending as it is not a new obstacle")
	} else {

		r.RARP.Set(obstacle.OriginatorIP.String(), from)
		bytes := ip.MarshalIPPacket(&packet)
		ip.Update(bytes)
		r.ethChannel.Broadcast(bytes)
		r.sendToALLWLANInterface(payload, obstacle.OriginatorIP.String())
		r.OTable.Set(obstacle.Position, 0)
	}
}

// handle VOREQ request
func (r *RSU) handleVOREQ(payload []byte, from net.HardwareAddr) {
	VOREQ := UnmarshalVOREQ(payload)
	r.RARP.Set(VOREQ.OriginatorIP.String(), from)
	newObstacles := UnmarshalPositions(VOREQ.Obstacles, int(VOREQ.Length))
	for _, obs := range newObstacles {
		if _, ok := r.OTable.table[string(obs.Marshal())]; ok {
			log.Println("Not Sending as it is not a new obstacle")
		} else {
			r.OTable.Set(obs, 0)
			VObstacleMessage := NewVObstacleMessage(VOREQ.OriginatorIP, obs, 0)
			packet := ip.NewIPPacket(VObstacleMessage.Marshal(), r.ip, net.ParseIP(ip.RsuIP))
			bytes := ip.MarshalIPPacket(packet)
			ip.UpdateChecksum(bytes)
			r.ethChannel.Broadcast(bytes)
			log.Println("Sending as it is a new obstacle")

		}
	}
	VOREP := NewVOREPMessage(r.OTable.GetTable())
	log.Println("Send VOREP to: ", VOREQ.OriginatorIP.String())

	packet := ip.NewIPPacket(VOREP.Marshal(), r.ip, VOREQ.OriginatorIP)
	bytes := ip.MarshalIPPacket(packet)
	ip.UpdateChecksum(bytes)
	r.wlanChannel.SendTo(bytes, from)
}

// Send to all in RSUARP using wlan exept the one that sent the message
func (r *RSU) sendToALLWLANInterface(data []byte, originatorIP string) {
	for eip, entry := range r.RARP.table {
		if originatorIP == eip {
			continue
		}
		log.Printf("Sending to: %s", eip)
		packet := ip.NewIPPacket(data, r.ip, net.ParseIP(eip))
		bytes := ip.MarshalIPPacket(packet)
		ip.UpdateChecksum(bytes)
		r.wlanChannel.SendTo(bytes, entry.MAC)
	}
}
