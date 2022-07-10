package app

import (
	"fmt"
	"log"
	"math"
	"net"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func (a *App) GetPosition() Position {
	state := a.GetState()
	return Position{Lng: state.Lng, Lat: state.Lat}
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

func distancePosition(p1 Position, p2 Position) float64 {
	latRad1 := toRadians(p1.Lat)
	latRad2 := toRadians(p2.Lat)
	lngRad1 := toRadians(p1.Lng)
	lngRad2 := toRadians(p2.Lng)

	dlng := lngRad2 - lngRad1
	dlat := latRad2 - latRad1

	ans := math.Pow(math.Sin(dlat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Cos(dlng/2), 2)

	ans = 2 * math.Asin(math.Sqrt(ans))

	return ans * 6371
}

func MyInterface() net.Interface {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("Error getting interfaces: %v", err)
	}
	iface := interfaces[1]
	return iface
}

func MyIP() (net.IP, bool, error) {
	ifi := MyInterface()
	addresses, err := ifi.Addrs()
	if err != nil {
		return nil, false, err
	}
	address := addresses[0]
	s := strings.Split(address.String(), "/")[0]
	ip := net.ParseIP(s)
	if ip.To4() != nil {
		return ip, false, nil
	} else if ip.To16() != nil {
		return ip, true, nil
	} else {
		return nil, false, fmt.Errorf("IP can't be resolved")
	}
}
