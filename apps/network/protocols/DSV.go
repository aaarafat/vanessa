package protocols

import (
	"fmt"
	"log"

	"github.com/aaarafat/vanessa/apps/network/datalink"
	. "github.com/aaarafat/vanessa/apps/network/datalink"
)

type DSV struct {
	etherType Ethertype
	datalink	*datalink.DataLinkLayerChannel
	neighborTable		*datalink.VNeighborTable
}

type DSVPacketLength int

func (dsv* DSV) read() {
	for {

		payload, addr, err := dsv.datalink.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		fmt.Println()
		log.Printf("Received \"%s\" from: [ %s ]", string(payload), addr.String())
		dsv.datalink.Broadcast(payload)
	}

}

func main() {
	d, err := NewDataLinkLayerChannel(VDSVEtherType)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}

	neighborTable := NewNeighborTable()

	dsv := &DSV{
		etherType: VDSVEtherType,
		datalink: d,
		neighborTable: neighborTable,
	}

	go dsv.read()

}