package aodv

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createRREQ() *RREQMessage {
	return NewRREQMessage(net.ParseIP("192.168.1.1"), net.ParseIP("192.168.1.2"))
}

func TestRREQMarshalUnmarshal(t *testing.T) {
	t.Log("Testing RREQ marshal and unmarshal")
	rreq := createRREQ()
	data := rreq.Marshal()
	rreq2, err := UnmarshalRREQ(data)

	if err != nil {
		t.Errorf("Error unmarshalling RREQ: %v", err)
	}

	assert.Equal(t, rreq, rreq2)
}

func TestRREQInvalid(t *testing.T) {
	t.Log("Testing RREQ invalid")
	srcIP := net.ParseIP("10.0.0.0")
	rreq := createRREQ()
	seqTable := NewVFloodingSeqTable()

	assert.False(t, rreq.Invalid(seqTable, srcIP))

	// invalidate the RREQ
	rreq.Type = rreq.Type + 1
	assert.True(t, rreq.Invalid(seqTable, srcIP))
	rreq.Type = rreq.Type - 1

	rreq.HopCount = HopCountLimit + 1
	assert.True(t, rreq.Invalid(seqTable, srcIP))
	rreq.HopCount = 0

	rreq.OriginatorIP = srcIP
	assert.True(t, rreq.Invalid(seqTable, srcIP))
	rreq.OriginatorIP = net.ParseIP("192.158.1.6")

	seqTable.Set(rreq.OriginatorIP, rreq.RREQID)
	assert.True(t, rreq.Invalid(seqTable, srcIP))
}


func TestRREQFlag(t *testing.T) {
	t.Log("Testing RREQ flag")
	rreq := createRREQ()
	
	rreq.SetFlag(RREQFlagJ)
	assert.True(t, rreq.HasFlag(RREQFlagJ))

	rreq.SetFlag(RREQFlagR)
	assert.True(t, rreq.HasFlag(RREQFlagR))

	rreq.SetFlag(RREQFlagG)
	assert.True(t, rreq.HasFlag(RREQFlagG))
	
	rreq.SetFlag(RREQFlagD)
	assert.True(t, rreq.HasFlag(RREQFlagD))

	rreq.SetFlag(RREQFlagU)
	assert.True(t, rreq.HasFlag(RREQFlagU))

	rreq.ClearFlag(RREQFlagJ)
	assert.False(t, rreq.HasFlag(RREQFlagJ))

	rreq.ToggleFlag(RREQFlagJ)
	assert.True(t, rreq.HasFlag(RREQFlagJ))

	rreq.ToggleFlag(RREQFlagJ)
	assert.False(t, rreq.HasFlag(RREQFlagJ))
}