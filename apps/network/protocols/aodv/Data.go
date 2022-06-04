package aodv

import (
	"encoding/binary"
	"fmt"
	"net"
)

type DataMessage struct {
	Type uint8
	Flags uint16
	HopCount uint8
	SeqNumber uint32
	OriginatorIP net.IP
	Data []byte
}

func NewDataMessage(SrcIP net.IP, SeqNum uint32, data []byte) *DataMessage {
	return &DataMessage{
		Type: DataType,
		Flags: 0,
		HopCount: 0,
		SeqNumber: SeqNum,
		OriginatorIP: SrcIP,
		Data: data,
	}
}

func (data *DataMessage) Marshal() []byte {
	bytes := make([]byte, DataMessageLen)

	bytes[0] = byte(data.Type)
	binary.LittleEndian.PutUint16(bytes[1:], data.Flags)
	bytes[3] = byte(data.HopCount)
	binary.LittleEndian.PutUint32(bytes[4:], data.SeqNumber)
	copy(bytes[8:12], data.OriginatorIP.To4())

	bytes = append(bytes, data.Data...)

	return bytes
}

func UnmarshalData(data []byte) (*DataMessage, error) {
	if len(data) < DataMessageLen {
		return nil, fmt.Errorf("Data message length is %d, expected %d", len(data), DataMessageLen)
	}

	dataMsg := &DataMessage{}
	dataMsg.Type = uint8(data[0])
	dataMsg.Flags = binary.LittleEndian.Uint16(data[1:3])
	dataMsg.HopCount = data[3]
	dataMsg.SeqNumber = binary.LittleEndian.Uint32(data[4:8])
	dataMsg.OriginatorIP = net.IPv4(data[8], data[9], data[10], data[11])

	dataMsg.Data = data[12:]

	return dataMsg, nil
}

func (data *DataMessage) String() string {
	return fmt.Sprintf("DataMessage{Type: %d, Flags: %d, HopCount: %d, SeqNumber: %d, OriginatorIP: %s, Data: %s}",
		data.Type, data.Flags, data.HopCount, data.SeqNumber, data.OriginatorIP, string(data.Data))
}