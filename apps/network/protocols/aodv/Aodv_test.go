package aodv

import (
	"fmt"
	"net"
	"testing"

	"github.com/cornelk/hashmap"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockFlooder struct {
	mock.Mock
}

func (f *MockFlooder) ForwardTo(payload []byte, addr net.HardwareAddr) {}
func (f *MockFlooder) ForwardToAll(payload []byte) {
	f.Called()
}
func (f *MockFlooder) ForwardToAllExcept(payload []byte, addr net.HardwareAddr) {}
func (f *MockFlooder) ForwardToAllExceptIP(payload []byte, ip net.IP)           {}

func createAodv() *Aodv {
	return &Aodv{
		srcIP:                 net.ParseIP("192.168.1.1"),
		routingTable:          NewVRoutingTable(),
		rreqBuffer:            &hashmap.HashMap{},
		flooder:               &MockFlooder{},
		seqNum:                0,
		rreqID:                0,
		pathDiscoveryCallback: func(net.IP) {},
	}
}

func TestGetRoute(t *testing.T) {
	aodv := createAodv()
	destIP := net.ParseIP("192.168.1.2")
	mac, _ := net.ParseMAC("00:00:00:00:00:02")
	ifi := 10
	hops := 5
	aodv.routingTable.Update(destIP, mac, uint8(hops), 0, 0, ifi)

	route, ok := aodv.GetRoute(destIP)

	assert.True(t, ok)
	assert.NotNil(t, route)
	assert.Equal(t, destIP, route.Destination)
	assert.Equal(t, mac, route.NextHop)
	assert.Equal(t, ifi, route.Interface)
	assert.Equal(t, hops, route.Metric)
}

func TestBuildRoute(t *testing.T) {
	aodv := createAodv()
	oldSeq := aodv.seqNum
	destIP := net.ParseIP("192.168.1.2")

	aodv.flooder.(*MockFlooder).On("ForwardToAll").Return()

	// when RREQ is sent to the destination
	aodv.rreqBuffer.Set(fmt.Sprintf("%s-%d", destIP.String(), aodv.rreqID), "")
	started := aodv.BuildRoute(destIP)
	assert.False(t, started)

	// when RREQ is not in the buffer
	aodv.rreqBuffer.Del(fmt.Sprintf("%s-%d", destIP.String(), aodv.rreqID))
	started = aodv.BuildRoute(destIP)
	assert.True(t, started)

	// check if RREQ is in the buffer
	assert.True(t, aodv.inRREQBuffer(destIP))

	// check sequence number is updated
	assert.Equal(t, oldSeq+1, aodv.seqNum)

	// check RREQ is sent
	aodv.flooder.(*MockFlooder).AssertCalled(t, "ForwardToAll")

}
