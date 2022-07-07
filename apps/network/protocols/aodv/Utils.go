package aodv

import (
	"net"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
)


func (a *Aodv) updateSeqNum(newSeqNum uint32) {
	if newSeqNum > a.seqNum {
		a.seqNum = newSeqNum
	}
}

func (a *Aodv) connectedDirectlyTo(destIP net.IP) bool {
	route, exists := a.routingTable.Get(destIP)
	return exists && route.NoOfHops == 0
}

func (a *Aodv) isRREQForMe(rreq *RREQMessage) bool {
	return rreq.DestinationIP.Equal(a.srcIP) || 
				(rreq.DestinationIP.Equal(net.ParseIP(ip.RsuIP)) && a.connectedDirectlyTo(net.ParseIP(ip.RsuIP))) 
}