package vector

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMagnitude(t *testing.T) {
	testCases := []struct {
		lng      float64
		lat      float64
		expected float64
	}{
		{0, 0, 0},
		{1, 0, 1},
		{0, 1, 1},
		{1, 1, math.Sqrt(2)},
		{-1, -1, math.Sqrt(2)},
		{-1, 1, math.Sqrt(2)},
		{1, -1, math.Sqrt(2)},
	}

	for _, tc := range testCases {
		v := Vector{Lng: tc.lng, Lat: tc.lat}
		assert.Equal(t, tc.expected, v.Magnitude())
	}
}

func TestDot(t *testing.T) {
	testCases := []struct {
		lng1     float64
		lat1     float64
		lng2     float64
		lat2     float64
		expected float64
	}{
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 1, 1, 2},
		{1, 1, -1, -1, -2},
		{1, 1, -1, 1, 0},
		{1, 1, 1, -1, 0},
	}

	for _, tc := range testCases {
		v1 := Vector{Lng: tc.lng1, Lat: tc.lat1}
		v2 := Vector{Lng: tc.lng2, Lat: tc.lat2}
		assert.Equal(t, tc.expected, v1.Dot(v2))
	}
}

func TestCross(t *testing.T) {
	testCases := []struct {
		lng1     float64
		lat1     float64
		lng2     float64
		lat2     float64
		expected float64
	}{
		{0, 0, 0, 0, 0},
		{1, 0, 0, 0, 0},
		{0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0},
		{0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0},
		{1, 1, 1, 1, 0},
		{1, 1, -1, -1, 0},
		{1, 1, -1, 1, 2},
		{1, 1, 1, -1, -2},
	}

	for _, tc := range testCases {
		v1 := Vector{Lng: tc.lng1, Lat: tc.lat1}
		v2 := Vector{Lng: tc.lng2, Lat: tc.lat2}
		assert.Equal(t, tc.expected, v1.Cross(v2))
	}
}
