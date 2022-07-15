package asncf

import (
	"net"
	"testing"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/libs/vector"
)

func CreateAreaSncf() *AreaSNCF {
	return &AreaSNCF{
		srcIP:  net.ParseIP("192.168.1.1"),
		seqNum: 0,
	}
}

func TestValidOptions(t *testing.T) {
	sncf := CreateAreaSncf()
	dst := net.ParseIP(ip.BroadcastIP)

	testcases := []struct {
		packet   ip.IPPacket
		expected bool
	}{
		{*ip.NewIPPacket([]byte("test"), sncf.srcIP, dst), true},                        // no options
		{*ip.NewIPPacketWithOptions([]byte("test"), sncf.srcIP, dst, []byte{0}), false}, // invalid option
		{*ip.NewIPPacketWithOptions([]byte("test"), sncf.srcIP, dst, []byte{1}), false}, // invalid position option marshal
		{*ip.NewIPPacketWithOptions([]byte("test"), dst, sncf.srcIP,
			ip.NewPositionOption(Position{Lng: 0, Lat: 0}, 1).Marshal()), false}, // no current position
		{*ip.NewIPPacketWithOptions([]byte("test"), sncf.srcIP, dst,
			ip.NewPositionOption(Position{Lng: 0, Lat: 0}, 1).Marshal()), true}, // update current position
		{*ip.NewIPPacketWithOptions([]byte("test"), dst, sncf.srcIP,
			ip.NewPositionOption(Position{Lng: 0.000001, Lat: 0.000001}, 1).Marshal()), true}, // accepted distance
		{*ip.NewIPPacketWithOptions([]byte("test"), dst, sncf.srcIP,
			ip.NewPositionOption(Position{Lng: 10, Lat: 10}, 1).Marshal()), false}, // not accepted distance
	}

	for _, testcase := range testcases {
		result := sncf.validOptions(&testcase.packet)
		if result != testcase.expected {
			t.Errorf("Expected %v, got %v", testcase.expected, result)
		}
	}
}
