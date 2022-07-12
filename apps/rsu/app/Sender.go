package app

import (
	"log"
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	"github.com/aaarafat/vanessa/libs/crypto"
)

func (a *App) sendToALLWLANInterface(data []byte, ip string) {
	count := a.router.SendToALLWLANInterface(data, ip)
	a.state.SentPacket(1, count)
}

func (a *App) sendPacketToWLAN(packet *ip.IPPacket, to net.IP) {
	bytes := ip.MarshalIPPacket(packet)
	ip.UpdateChecksum(bytes)
	a.router.SendToWLANInterface(bytes, to.String())
	a.state.SentPacket(1, 1)
}

func (a *App) sendPacketToETH(packet *ip.IPPacket) {
	a.router.BroadcastETH(packet)
	a.state.SentPacket(0, 1)
}

func (a *App) sendObstacle(obstacle *Position, from net.IP) {
	VObstacleMessage := NewVObstacleMessage(from, *obstacle, 0)
	cipherData, err := crypto.EncryptAES(a.key, VObstacleMessage.Marshal())
	if err != nil {
		return
	}
	packet := ip.NewIPPacket(cipherData, a.ip, net.ParseIP(ip.RsuIP))

	a.addObstacle(obstacle, from, packet)
}

func (a *App) sendVOREP(to net.IP) {
	VOREP := NewVOREPMessage(a.state.OTable.GetTable())
	cipherData, err := crypto.EncryptAES(a.key, VOREP.Marshal())
	if err != nil {
		return
	}
	log.Println("Send VOREP to: ", to.String())

	packet := ip.NewIPPacket(cipherData, a.ip, to)
	a.sendPacketToWLAN(packet, to)
}
