package packetfilter

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
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

func GetMyIPs(ifi *net.Interface) (net.IP, bool, error) {

	addresses, err := ifi.Addrs()
	if err != nil {
		return nil, false, err
	}
	address := addresses[0]
	s :=strings.Split(address.String(), "/")[0]
	ip := net.ParseIP(s)
	if ip.To4()!= nil {
		return ip,false,nil
	} else if ip.To16()!=nil {
		return ip,true,nil
	}else {
		return nil, false, fmt.Errorf("IP can't be resolved")
	}
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
