package app

import (
	"fmt"
	"log"
	"math"
	"net"
	"strings"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
	"github.com/aaarafat/vanessa/libs/crypto"
)

func (a *App) getDataFromPacket(bytes []byte) ([]byte, error) {
	packet, err := ip.UnmarshalPacket(bytes)
	if err != nil {
		log.Printf("Error unmarshalling packet: %v", err)
		return nil, err
	}

	data := packet.Payload

	if !packet.Header.SrcIP.Equal(a.ip) {
		// Decrypt data using AES if it is not from self
		data, err = crypto.DecryptAES(a.key, packet.Payload)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func removeDuplicatePositions(positions []Position) []Position {
	keys := make(map[Position]bool)
	list := []Position{}

	for _, entry := range positions {
		if _, value := keys[entry]; !value {
			keys[entry] = true
			list = append(list, entry)
		}
	}
	return list
}

func (a *App) GetPosition() Position {
	state := a.GetState()
	return Position{Lng: state.Lng, Lat: state.Lat}
}

func toRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// Distance returns the distance between two points on the Earth in meter.
func distancePosition(p1 Position, p2 Position) float64 {
	latRad1 := toRadians(p1.Lat)
	latRad2 := toRadians(p2.Lat)
	lngRad1 := toRadians(p1.Lng)
	lngRad2 := toRadians(p2.Lng)

	dlng := lngRad2 - lngRad1
	dlat := latRad2 - latRad1

	ans := math.Pow(math.Sin(dlat/2), 2) + math.Cos(latRad1)*math.Cos(latRad2)*math.Pow(math.Sin(dlng/2), 2)

	ans = 2 * math.Asin(math.Sqrt(ans))

	return ans * EARTH_RADIUS_METER
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
