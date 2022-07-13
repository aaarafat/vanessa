package sncf

import "encoding/binary"

type SequenceHeader struct {
	SequenceNumber uint32
}

func NewSequenceHeader(sequenceNumber uint32) *SequenceHeader {
	return &SequenceHeader{
		SequenceNumber: sequenceNumber,
	}
}

func (sh *SequenceHeader) Marshal() []byte {
	bytes := make([]byte, 4)

	binary.LittleEndian.PutUint32(bytes, sh.SequenceNumber)

	return bytes
}

func UnmarshalSequenceHeader(bytes []byte) *SequenceHeader {
	sh := &SequenceHeader{}

	sh.SequenceNumber = binary.LittleEndian.Uint32(bytes)

	return sh
}
