package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (a *App) handleEthMessages() {
	for {
		packet, _, err := a.router.ReadFromETHInterface()
		if err != nil {
			continue
		}
		data, err := a.getDataFromPacket(packet)
		if err != nil {
			continue
		}
		a.state.ReceivedPacket(0)
		msgType := uint8(data[0])
		switch msgType {
		case VObstacleType:
			obstacle, err := UnmarshalVObstacle(data)
			if err != nil {
				log.Println("Failed to unmarshal VObstacle: ", err)
				continue
			}
			log.Println("Recieved Obstacle from: ", obstacle.OriginatorIP.String(), " at: ", obstacle.Position)
			a.sendToALLWLANInterface(packet.Payload, obstacle.OriginatorIP.String())
			a.state.AddObstacle(&obstacle.Position)

		default:
			log.Println("Unknown message type: ", msgType)
		}
	}
}

func (a *App) handleWLANMessgaes() {
	for {
		packet, from, err := a.router.ReadFromWLANInterface()
		if err != nil {
			continue
		}
		data, err := a.getDataFromPacket(packet)
		if err != nil {
			continue
		}
		a.state.ReceivedPacket(1)
		msgType := uint8(data[0])
		switch msgType {
		case VObstacleType:
			a.handleVObstacle(data, from, packet)
		case VHBeatType:
			a.handleVHBeat(data, from)
		case VOREQType:
			a.handleVOREQ(data, from)
		default:
			log.Println("Unknown message type: ", msgType)
		}
	}
}

func (a *App) handleVHBeat(payload []byte, from net.HardwareAddr) {
	HBeat, err := UnmarshalVHBeat(payload)
	if err != nil {
		log.Println("Failed to unmarshal VHBeat: ", err)
		return
	}
	log.Println("Recieved VHBeat from: ", from, " at: ", HBeat.Position)
	a.router.RARP.Set(HBeat.OriginatorIP.String(), from)
}

func (a *App) handleVObstacle(payload []byte, from net.HardwareAddr, packet *ip.IPPacket) {
	obstacle, err := UnmarshalVObstacle(payload)
	if err != nil {
		log.Println("Failed to unmarshal VObstacle: ", err)
		return
	}
	log.Println("Recieved Obstacle from: ", obstacle.OriginatorIP.String(), " at: ", obstacle.Position)
	a.router.RARP.Set(obstacle.OriginatorIP.String(), from)
	a.addObstacle(&obstacle.Position, obstacle.OriginatorIP, packet)
}

func (a *App) handleVOREQ(payload []byte, from net.HardwareAddr) {
	VOREQ := UnmarshalVOREQ(payload)
	log.Printf("Received VOREQ from %s\n", VOREQ.OriginatorIP.String())
	a.router.RARP.Set(VOREQ.OriginatorIP.String(), from)
	obstacles := UnmarshalPositions(VOREQ.Obstacles, int(VOREQ.Length))
	a.addObstacles(obstacles, VOREQ.OriginatorIP)
	a.sendVOREP(VOREQ.OriginatorIP)
}
