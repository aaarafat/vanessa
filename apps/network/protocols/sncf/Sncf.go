package sncf

import (
	"net"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
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

func (sncf *SNCF) Forward(data []byte, from net.HardwareAddr) {
	seqHeader := UnmarshalSequenceHeader(data[0:4])
	packet, err := ip.UnmarshalPacket(data[4:])
	if err != nil || sncf.seqTable.Exists(packet.Header.SrcIP, seqHeader.SequenceNumber) {
		return
	}

	sncf.seqTable.Set(packet.Header.SrcIP, seqHeader.SequenceNumber)

	sncf.flooder.ForwardToAllExcept(data[4:], from)
}
