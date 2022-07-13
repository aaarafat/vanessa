package packetfilter

import (
	"net"
	"testing"

	"github.com/aaarafat/vanessa/apps/network/network/ip"
	"github.com/aaarafat/vanessa/apps/network/network/messages"
)

func CreatePacketFilter() *PacketFilter {
	return &PacketFilter{
		srcIP: net.ParseIP("192.168.1.1"),
		id:    1,
	}
}

func TestValidOptions(t *testing.T) {
	pf := CreatePacketFilter()
	dst := net.ParseIP(ip.BroadcastIP)

	testcases := []struct {
		packet   ip.IPPacket
		expected bool
	}{
		{*ip.NewIPPacket([]byte("test"), pf.srcIP, dst), true},                        // no options
		{*ip.NewIPPacketWithOptions([]byte("test"), pf.srcIP, dst, []byte{0}), false}, // invalid option
		{*ip.NewIPPacketWithOptions([]byte("test"), pf.srcIP, dst, []byte{1}), false}, // invalid position option marshal
		{*ip.NewIPPacketWithOptions([]byte("test"), dst, pf.srcIP,
			ip.NewPositionOption(messages.Position{Lng: 0, Lat: 0}, 1).Marshal()), false}, // no current position
		{*ip.NewIPPacketWithOptions([]byte("test"), pf.srcIP, dst,
			ip.NewPositionOption(messages.Position{Lng: 0, Lat: 0}, 1).Marshal()), true}, // update current position
		{*ip.NewIPPacketWithOptions([]byte("test"), dst, pf.srcIP,
			ip.NewPositionOption(messages.Position{Lng: 0.000001, Lat: 0.000001}, 1).Marshal()), true}, // accepted distance
		{*ip.NewIPPacketWithOptions([]byte("test"), dst, pf.srcIP,
			ip.NewPositionOption(messages.Position{Lng: 10, Lat: 10}, 1).Marshal()), false}, // not accepted distance
	}

	for _, testcase := range testcases {
		result := pf.ValidOptions(&testcase.packet)
		if result != testcase.expected {
			t.Errorf("Expected %v, got %v", testcase.expected, result)
		}
	}
}
