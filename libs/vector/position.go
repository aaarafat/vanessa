package vector

import (
	"encoding/binary"
	"math"
)

const (
	EARTH_RADIUS_METER = 6371000
)

type Position struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

func Float64FromBytes(bytes []byte) float64 {
	bits := binary.LittleEndian.Uint64(bytes)
	float := math.Float64frombits(bits)
	return float
}

func Float64bytes(float float64) []byte {
	bits := math.Float64bits(float)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}

func ToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Postion to byte array marashalling
func (p Position) Marshal() []byte {
	return append(Float64bytes(p.Lat), Float64bytes(p.Lng)...)
}

// Position from byte array unmarshalling
func UnmarshalPosition(data []byte) Position {
	return Position{
		Lat: Float64FromBytes(data[:8]),
		Lng: Float64FromBytes(data[8:]),
	}
}

// Unmarshal the payload into a list of positions
func UnmarshalPositions(payload []byte, len int) []Position {
	var list []Position
	for i := 0; i < len; i++ {
		pos := Position{
			Lat: Float64FromBytes(payload[i*16 : i*16+8]),
			Lng: Float64FromBytes(payload[i*16+8 : i*16+16]),
		}
		list = append(list, pos)
	}
	return list
}

// Distance returns the distance between two points on the Earth in meter.
func (p *Position) Distance(p2 *Position) float64 {
	latRad1 := ToRadians(p.Lat)
	latRad2 := ToRadians(p2.Lat)
	lngRad1 := ToRadians(p.Lng)
	lngRad2 := ToRadians(p2.Lng)

	dlng := lngRad2 - lngRad1
	dlat := latRad2 - latRad1

	ans := math.Pow(math.Sin(dlat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Sin(dlng/2), 2)

	ans = 2 * math.Asin(math.Sqrt(ans))

	return ans * EARTH_RADIUS_METER
}
