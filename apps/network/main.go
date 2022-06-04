package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/protocols/aodv"
)


func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}


func GetMyIPs(iface *net.Interface) (net.IP, net.IP, error) {
	addrs, err := iface.Addrs()
	if err != nil {
		return nil, nil, err
	}

	var ip4, ip6 net.IP
	for _, addr := range addrs {
		var ip net.IP

		switch v := addr.(type) {
		case *net.IPNet:
			ip = v.IP
		case *net.IPAddr:
			ip = v.IP
		}

		if isIPv4(ip.String()) {
			ip4 = ip
		} else if isIPv6(ip.String()) {
			ip6 = ip
		} else {
			return nil, nil, fmt.Errorf("ip is not ip4 or ip6")
		}
	}

	return ip4, ip6, nil
}


func main() {
	interfaces, err := net.Interfaces()
	iface := interfaces[1]

	ip, _, err := GetMyIPs(&iface)

	log.Printf("My IP is %s", ip.String())

	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return
	}

	aodv := NewAodv(ip)

	aodv.Start()

	// wait for ever
	var ipString, msg string
	for {
		fmt.Scanf("%s %s", &ipString, &msg)

		ip := net.ParseIP(ipString)

		aodv.SendData([]byte(msg), ip)
	}
}