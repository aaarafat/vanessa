package aodv


func (a *Aodv) updateSeqNum(newSeqNum uint32) {
	if newSeqNum > a.seqNum {
		a.seqNum = newSeqNum
	}
}
