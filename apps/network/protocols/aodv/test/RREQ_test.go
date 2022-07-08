package test

import (
	"net"
	"testing"

	"github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)

func createRREQ() *aodv.RREQMessage {
	return aodv.NewRREQMessage(net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.2"))
}

func TestRREQMarshalUnmarshal(t *testing.T) {
	t.Log("Testing RREQ marshal and unmarshal")
	rreq := createRREQ()
	data := rreq.Marshal()
	rreq2, err := aodv.UnmarshalRREQ(data)

	if err != nil {
		t.Errorf("Error unmarshalling RREQ: %v", err)
	}

	// check that all attributes are the same
	if rreq.Type != rreq2.Type {
		t.Errorf("Type not the same")
	}
	if rreq.Flags != rreq2.Flags {
		t.Errorf("Flags not the same")
	}
	if rreq.HopCount != rreq2.HopCount {
		t.Errorf("HopCount not the same")
	}
	if rreq.RREQID != rreq2.RREQID {
		t.Errorf("RREQID not the same")
	}
	if rreq.DestinationIP.String() != rreq2.DestinationIP.String() {
		t.Errorf("DestinationIP not the same")
	}
	if rreq.DestinationSeqNum != rreq2.DestinationSeqNum {
		t.Errorf("DestinationSeqNum not the same")
	}
	if rreq.OriginatorIP.String() != rreq2.OriginatorIP.String() {
		t.Errorf("SourceIP not the same")
	}
	if rreq.OriginatorSequenceNumber != rreq2.OriginatorSequenceNumber {
		t.Errorf("SourceSeqNum not the same")
	}
}

func TestRREQInvalid(t *testing.T) {
	t.Log("Testing RREQ invalid")
	srcIP := net.ParseIP("10.0.0.0")
	rreq := createRREQ()
	seqTable := aodv.NewVFloodingSeqTable()

	// invalidate the RREQ
	rreq.Type = rreq.Type + 1
	if !rreq.Invalid(seqTable, srcIP) {
		t.Errorf("RREQ with type %d should be invalid", rreq.Type)
	}
	rreq.Type = rreq.Type - 1

	rreq.HopCount = aodv.HopCountLimit + 1
	if !rreq.Invalid(seqTable, srcIP) {
		t.Errorf("RREQ with hop count %d should be invalid", rreq.HopCount)
	}
	rreq.HopCount = 0

	
	rreq.OriginatorIP = srcIP
	if !rreq.Invalid(seqTable, srcIP) {
		t.Errorf("RREQ with originator IP equal my IP should be invalid")
	}
	rreq.OriginatorIP = net.ParseIP("192.158.1.6")


	seqTable.Set(rreq.OriginatorIP, rreq.RREQID)
	if !rreq.Invalid(seqTable, srcIP) {
		t.Errorf("RREQ in my sequence table should be invalid")
	}
}


func TestRREQFlag(t *testing.T) {
	t.Log("Testing RREQ flag")
	rreq := createRREQ()
	
	rreq.SetFlag(aodv.RREQFlagJ)
	if !rreq.HasFlag(aodv.RREQFlagJ) {
		t.Errorf("RREQ flag J not set")
	}

	rreq.SetFlag(aodv.RREQFlagR)
	if !rreq.HasFlag(aodv.RREQFlagR) {
		t.Errorf("RREQ flag R not set")
	}

	rreq.SetFlag(aodv.RREQFlagG)
	if !rreq.HasFlag(aodv.RREQFlagG) {
		t.Errorf("RREQ flag G not set")
	}
	
	rreq.SetFlag(aodv.RREQFlagD)
	if !rreq.HasFlag(aodv.RREQFlagD) {
		t.Errorf("RREQ flag D not set")
	}

	rreq.SetFlag(aodv.RREQFlagU)
	if !rreq.HasFlag(aodv.RREQFlagU) {
		t.Errorf("RREQ flag U not set")
	}

	rreq.ClearFlag(aodv.RREQFlagJ)
	if rreq.HasFlag(aodv.RREQFlagJ) {
		t.Errorf("RREQ flag J set")
	}

	rreq.ToggleFlag(aodv.RREQFlagJ)
	if !rreq.HasFlag(aodv.RREQFlagJ) {
		t.Errorf("RREQ flag J not set")
	}

	rreq.ToggleFlag(aodv.RREQFlagJ)
	if rreq.HasFlag(aodv.RREQFlagJ) {
		t.Errorf("RREQ flag J set")
	}
}