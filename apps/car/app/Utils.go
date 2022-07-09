package app

import (
	"fmt"
	"log"
	"net"
	"strings"
)

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
