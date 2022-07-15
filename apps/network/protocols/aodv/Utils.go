package aodv

import (
	"fmt"
	"net"
	"time"
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

func (a *Aodv) connectedTo(destIP net.IP) bool {
	_, exists := a.routingTable.Get(destIP)
	return exists
}

func (a *Aodv) isRREQForMe(rreq *RREQMessage) bool {
	return rreq.DestinationIP.Equal(a.srcIP) 
}

func (a *Aodv) isRREQForNeighbor(rreq *RREQMessage) bool {
	return a.connectedTo(rreq.DestinationIP)
}

func (a *Aodv) inRREQBuffer(destIP net.IP) bool {
	_, ok := a.rreqBuffer.Get(fmt.Sprintf("%s-%d", destIP.String(), a.rreqID))
	return ok
}

func (a *Aodv) addToRREQBuffer(rreq *RREQMessage) {
	callback := func() {
		a.rreqBuffer.Del(fmt.Sprintf("%s-%d", rreq.DestinationIP.String(), rreq.RREQID))
	}
	timer := time.AfterFunc(time.Millisecond * time.Duration(PATH_DISCOVERY_TIME_MS), callback)
	a.rreqBuffer.Set(fmt.Sprintf("%s-%d", rreq.DestinationIP.String(), a.rreqID), *timer)
}