package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/aaarafat/vanessa/libs/crypto"
	. "github.com/aaarafat/vanessa/libs/vector"
)

func (a *App) getDataFromPacket(packet *ip.IPPacket) ([]byte, error) {
	return crypto.DecryptAES(a.key, packet.Payload)
}

func (a *App) addObstacle(obstacle *Position, from net.IP, obstacleDetectedPacket *ip.IPPacket) {
	if a.state.HasObstacle(obstacle) {
		log.Println("Not Sending as it is not a new obstacle")
		return
	}
	a.state.AddObstacle(obstacle)
	a.sendPacketToETH(obstacleDetectedPacket)
	a.sendToALLWLANInterface(obstacleDetectedPacket.Payload, from.String())
}

func (a *App) addObstacles(obstacles []Position, from net.IP) {
	for _, obstacle := range obstacles {
		if a.state.HasObstacle(&obstacle) {
			log.Println("Not Sending as it is not a new obstacle")
			continue
		}
		a.sendObstacle(&obstacle, from)
		log.Println("Sending as it is a new obstacle")
	}
}
