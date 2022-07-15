package vector

import (
	"math"
)

type Vector struct {
	Lng float64 `json:"lat"`
	Lat float64 `json:"lng"`
}

func NewVector(pos1, pos2 Position) Vector {
	return Vector{
		Lng: pos2.Lng - pos1.Lng,
		Lat: pos2.Lat - pos1.Lat,
	}
}

func NewUnitVector(pos1, pos2 Position) Vector {
	v := NewVector(pos1, pos2)
	v.Normalize()
	return v
}

func (v *Vector) Magnitude() float64 {
	return math.Sqrt(float64(v.Lng*v.Lng + v.Lat*v.Lat))
}

func (v *Vector) Normalize() {
	mag := v.Magnitude()
	if mag == 0 {
		return
	}
	v.Lng = float64(v.Lng) / mag
	v.Lat = float64(v.Lat) / mag
}

func (v *Vector) Dot(v2 Vector) float64 {
	return float64(v.Lng*v2.Lng + v.Lat*v2.Lat)
}

func (v *Vector) Cross(v2 Vector) float64 {
	return float64(v.Lng*v2.Lat - v.Lat*v2.Lng)
}

func (v *Vector) Angle(v2 Vector) float64 {
	return math.Atan2(v.Cross(v2), v.Dot(v2))
}
