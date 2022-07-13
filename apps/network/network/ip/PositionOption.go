package ip

import (
	"fmt"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

const (
	PositionOptionType uint8 = 0x01
)

type PositionOption struct {
	Type        uint8
	Position    Position
	MaxDistance float64
}

const (
	PositionOptionLen = 25
)

func NewPositionOption(position Position, maxDistance float64) *PositionOption {
	return &PositionOption{
		Type:        PositionOptionType,
		Position:    position,
		MaxDistance: maxDistance,
	}
}

func (po *PositionOption) Marshal() []byte {
	bytes := make([]byte, PositionOptionLen)

	bytes[0] = po.Type
	copy(bytes[1:17], po.Position.Marshal())
	copy(bytes[17:25], Float64bytes(po.MaxDistance))

	return bytes
}

func UnmarshalPositionOption(data []byte) (*PositionOption, error) {
	if len(data) < PositionOptionLen {
		return nil, fmt.Errorf("PositionOption length is %d, expected %d", len(data), PositionOptionLen)
	}

	positionOption := &PositionOption{}

	positionOption.Type = uint8(data[0])
	positionOption.Position = UnmarshalPosition(data[1:17])
	positionOption.MaxDistance = Float64FromBytes(data[17:25])

	return positionOption, nil
}
