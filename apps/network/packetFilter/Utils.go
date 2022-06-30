package packetfilter

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
)


func AddDefaultGateway() error {
	cmd := exec.Command("route", "add", "default", "gw", "localhost")
	std, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("Couldn't add default gateway, err: %#v, stderr: %#v", err, string(std))
	}
	log.Println("Added default gateway")
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

func MyIP(ifi *net.Interface) (net.IP, bool, error) {

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

func SetMSS(ifaceName string, ip net.IP, mss int) {
	cmd := exec.Command("ip", "route", "change", "10.0.0.0/8", "dev", ifaceName, "proto", "kernel", "scope", "link", "src", ip.String(), "advmss", strconv.Itoa(mss))
	std, err := cmd.CombinedOutput()
	if err != nil {
		log.Panic("Couldn't change MSS, err: ", err, ",stderr: ", string(std))
	}
}
