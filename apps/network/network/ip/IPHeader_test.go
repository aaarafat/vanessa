package ip

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func createIPHeader() *IPHeader {
	return &IPHeader{
		Version:               4,
		Length:                5,
		TypeOfService:         0,
		TotalLength:           20,
		IdentifierFlagsOffset: 0,
		TTL:                   10,
		Protocol:              0,
		HeaderChecksum:        0,
		SrcIP:                 net.IPv4(0, 0, 0, 0),
		DestIP:                net.IPv4(0, 0, 0, 0),
	}
}

func TestMarshalUnmarshal(t *testing.T) {
	header := createIPHeader()
	data := MarshalIPHeader(header)
	UpdateChecksum(data, len(data))

	header2, err := UnmarshalIPHeader(data)

	assert.NoError(t, err)

	assert.Equal(t, header.Version, header2.Version)
	assert.Equal(t, header.Length, header2.Length)
	assert.Equal(t, header.TypeOfService, header2.TypeOfService)
	assert.Equal(t, header.TotalLength, header2.TotalLength)
	assert.Equal(t, header.IdentifierFlagsOffset, header2.IdentifierFlagsOffset)
	assert.Equal(t, header.TTL, header2.TTL)
	assert.Equal(t, header.Protocol, header2.Protocol)
	assert.Equal(t, header.SrcIP, header2.SrcIP)
	assert.Equal(t, header.DestIP, header2.DestIP)
}

func TestHeaderChecksum(t *testing.T) {
	header := createIPHeader()
	data := MarshalIPHeader(header)
	UpdateChecksum(data, len(data))

	assert.Zero(t, header.HeaderChecksum)

	UpdateChecksum(data, len(data))
	csm := HeaderChecksum(data, len(data))
	assert.Zero(t, csm)
}

func TestTTL(t *testing.T) {
	header := createIPHeader()
	data := MarshalIPHeader(header)
	UpdateChecksum(data, len(data))

	assert.Equal(t, byte(10), header.TTL)

	header.TTL = 0
	data = MarshalIPHeader(header)
	UpdateChecksum(data, len(data))

	assert.Equal(t, byte(0), header.TTL)

	header.TTL = 1
	data = MarshalIPHeader(header)
	Update(data, len(data))

	header, err := UnmarshalIPHeader(data)
	assert.Error(t, err)
}

func TestLengthInBytes(t *testing.T) {
	header := createIPHeader()
	data := MarshalIPHeader(header)
	UpdateChecksum(data, len(data))

	assert.Equal(t, 20, len(data))
}
