package ip

import (
	"fmt"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

type PositionOption struct {
	Position    Position
	MaxDistance float64
}

const (
	PositionOptionLen = 24
)

func NewPositionOption(position Position, maxDistance float64) *PositionOption {
	return &PositionOption{
		Position:    position,
		MaxDistance: maxDistance,
	}
}

func (po *PositionOption) Marshal() []byte {
	bytes := make([]byte, PositionOptionLen)

	copy(bytes[0:16], po.Position.Marshal())
	copy(bytes[16:24], Float64bytes(po.MaxDistance))

	return bytes
}

func UnmarshalPositionOption(data []byte) (*PositionOption, error) {
	if len(data) < PositionOptionLen {
		return nil, fmt.Errorf("PositionOption length is %d, expected %d", len(data), PositionOptionLen)
	}

	positionOption := &PositionOption{}

	positionOption.Position = UnmarshalPosition(data[0:16])
	positionOption.MaxDistance = Float64FromBytes(data[16:24])

	return positionOption, nil
}
