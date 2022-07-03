package main

import (
	"log"
	"net"
	"net/netip"
	"os"
	"strconv"
	"strings"

	"github.com/mdlayher/ethernet"
	"github.com/mdlayher/packet"
)

func main() {
	type EtherType ethernet.EtherType
	const VEtherType = 0x7031
	index, err := strconv.Atoi(os.Args[1])
	if err == nil {
	}
	interfaces, err := net.Interfaces()

	for index, iface := range interfaces {
		log.Println(index, iface)
	}
	log.Println(interfaces[index].Name)
	log.Println("=================================")
	ifi, err := net.InterfaceByName(interfaces[index].Name)
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
	}
	addresses, err := ifi.Addrs()
	address := addresses[0]
	s := strings.Split(address.String(), "/")[0]
	ip := net.ParseIP(s)
	ip4, _ := netip.ParseAddr(s)
	ip6 := net.ParseIP(netip.AddrFrom16(ip4.As16()).String())
	if ip.To4() != nil {
		log.Printf("IP4: %s, Mapped IP6: %s , %s", ip.String(), netip.AddrFrom16(ip4.As16()).String(), ip6.String())
	} else if ip.To16() != nil {
		log.Printf("IP6: %s", ip.String())
	} else {
		log.Println("Not found")
	}

	log.Printf("MAC Add: %v", ifi.HardwareAddr)
	// Open a raw socket using same EtherType as our frame.
	c, err := packet.Listen(ifi, packet.Raw, VEtherType, nil)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer c.Close()

	// Marshal a frame to its binary format.
	f := &ethernet.Frame{
		// Broadcast frame to all machines on same network segment.
		Destination: ethernet.Broadcast,
		// Identify our machine as the sender.
		Source: ifi.HardwareAddr,
		// Identify frame with an unused EtherType.
		EtherType: VEtherType,
		// Send a simple message.
		Payload: []byte("Hello from :" + s),
	}
	b, err := f.MarshalBinary()
	if err != nil {
		log.Fatalf("failed to marshal frame: %v", err)
	}

	// Broadcast the frame to all devices on our network segment.
	addr := &packet.Addr{HardwareAddr: ethernet.Broadcast}
	if _, err := c.WriteTo(b, addr); err != nil {
		log.Fatalf("failed to write frame: %v", err)
	}

}
