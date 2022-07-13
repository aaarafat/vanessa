package sncf

import (
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/protocols"
)

type SNCF struct {
	seqNum   uint32
	seqTable *VFloodingSeqTable

	flooder IFlooder
}

func NewSNCF(flooder IFlooder) *SNCF {
	return &SNCF{
		seqNum:   0,
		seqTable: NewVFloodingSeqTable(),
		flooder:  flooder,
	}
}

func (sncf *SNCF) updateSeqNum() {
	sncf.seqNum++
}

func (sncf *SNCF) Flood(packet []byte) {
	// Add flooding header
	seqHeader := NewSequenceHeader(sncf.seqNum)
	sncf.updateSeqNum()

	data := make([]byte, len(packet)+4)

	copy(data[0:4], seqHeader.Marshal())
	copy(data[4:], packet)

	sncf.flooder.ForwardToAll(data)
}

func (sncf *SNCF) Forward(data []byte, fromIP net.IP) {
	seqHeader := UnmarshalSequenceHeader(data[0:4])
	if sncf.seqTable.Exists(fromIP, seqHeader.SequenceNumber) {
		return
	}

	sncf.seqTable.Set(fromIP, seqHeader.SequenceNumber)

	sncf.flooder.ForwardToAllExceptIP(data, fromIP)
}
