package messages

import (
	"encoding/binary"
	"math"
	"net"
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

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Distance returns the distance between two points on the Earth in meter.
func (p *Position) Distance(p2 *Position) float64 {
	latRad1 := toRadians(p.Lat)
	latRad2 := toRadians(p2.Lat)
	lngRad1 := toRadians(p.Lng)
	lngRad2 := toRadians(p2.Lng)

	dlng := lngRad2 - lngRad1
	dlat := latRad2 - latRad1

	ans := math.Pow(math.Sin(dlat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Sin(dlng/2), 2)

	ans = 2 * math.Asin(math.Sqrt(ans))

	return ans * EARTH_RADIUS_METER
}

// Heart beat message from car to the associated RSU
type VHBeatMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Heart Beat.
	OriginatorIP net.IP
	// Position of the node that originated the Heart Beat.
	Position Position
}

// Obstcale found message from car to the associated RSU
type VObstacleMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Obstcale alert.
	OriginatorIP net.IP
	// Position of the obstacle
	Position Position
	// Set obstacle to true if the obstacle is detected
	Clear uint8
}

type VOREQMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Obstcales request.
	OriginatorIP net.IP
	// length of the list of Obstcales
	Length uint8
	// content of the list of Obstcales
	Obstacles []byte
}

type VOREPMessage struct {
	// Type
	Type uint8
	// length of the list of Obstcales
	Length uint8
	// content of the list of Obstcales
	Obstacles []byte
}

type VPathDiscoveryMessage struct {
	// Type
	Type uint8
	// The IP address of the node that is discovered.
	DestIP net.IP
}

type VZoneMessage struct {
	// Type
	Type uint8
	// The IP address of the node that originated the Obstcale alert.
	OriginatorIP net.IP
	// Position of the obstacle
	Position Position
	// Max Distance from the originator
	MaxDistance float64
}

const (
	VHBeatMessageLen         = 21
	VObstacleMessageLen      = 22
	VOREQMessageLen          = 6
	VOREPMessageLen          = 2
	VPathDiscoveryMessageLen = 5
	VZoneMessageLen          = 29
)

const (
	VHBeatType         uint8 = 1
	VObstacleType      uint8 = 2
	VOREQType          uint8 = 3
	VOREPType          uint8 = 4
	VPathDiscoveryType uint8 = 5
	VZoneType          uint8 = 6
)
