package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/AkihiroSuda/go-netfilter-queue"
)

func isIPv4(address string) bool {
	return strings.Count(address, ":") < 2
}

func isIPv6(address string) bool {
	return strings.Count(address, ":") >= 2
}

func RegisterGateway() error {
	exec.Command("route", "del", "default", "gw", "localhost").Run()
	cmd := exec.Command("route", "add", "default", "gw", "localhost")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("couldn't add default gateway, err: %#v, stderr: %#v", err, string(stdoutStderr))
	}
	log.Println("added default gateway")
	return nil
}
func UnregisterGateway() {
	cmd := exec.Command("route", "del", "default", "gw", "localhost")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove default gateway, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove default gateway")
}
func AddIPTablesRule() {
	exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w").Run()
	cmd := exec.Command("iptables", "-t", "filter", "-A", "OUTPUT", "-j", "NFQUEUE", "-w", "--queue-num", "0")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't add iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("added NFQUEUE rule to OUTPUT chain in iptables")
}

func DeleteIPTablesRule() {
	cmd := exec.Command("iptables", "-t", "filter", "-D", "OUTPUT", "-j", "NFQUEUE", "-w")
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("couldn't remove iptables rule, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("remove NFQUEUE rule to OUTPUT chain in iptables")
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

func SetMaxMSS(ifaceName string, ip net.IP, mss int) {
	cmd := exec.Command(
		"ip", "route", "change", "10.0.0.0/8", "dev", ifaceName,
		"proto", "kernel", "scope", "link", "src", ip.String(), "advmss", strconv.Itoa(mss),
	)
	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("Couldn't change MSS, err: ", err, ",stderr: ", string(stdoutStderr))
	}
	log.Println("Changed MSS to: ", mss)
}

func StealPackets() {
	var err error

	nfq, err := netfilter.NewNFQueue(0, 100, netfilter.NF_DEFAULT_PACKET_SIZE)
	if err != nil {
					fmt.Println(err)
					os.Exit(1)
	}
	defer nfq.Close()
	packets := nfq.GetPackets()

	for true {
					select {
					case p := <-packets:
									fmt.Println(p.Packet)
									// drop tha packet
									p.SetVerdict(netfilter.NF_DROP)
					}
	}
}

func main() {
	AddIPTablesRule()
	if err := RegisterGateway(); err != nil {
		DeleteIPTablesRule()
		log.Panic("done deleting")
	}
	interfaces, err := net.Interfaces()
	for index , intf := range(interfaces){
		log.Println(index,intf)
	}
	iface := interfaces[1]

	ip, _, err := GetMyIPs(&iface)
	ip = ip.To4()
	if err != nil {
		log.Panicf("failed to get iface ips, err: %s", err)
	}
	log.Println("iface ipv4: ", ip)
	log.Println("iface: ", iface.Name)

	SetMaxMSS(iface.Name, ip, 1400)

	StealPackets()
}
