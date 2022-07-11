package datalink

import (
	"log"
	"net"
	"os/exec"
	"strings"
)

func getRSU(intfName string) (string, string) {

	out, err := exec.Command("iw", "dev", intfName, "link").Output()
	if err != nil {
		log.Panic(err)
	}
	cmdOut := string(out)
	// println(cmdOut)
	rsuMAC := ""
	SSID := ""
	if strings.Contains(cmdOut, "Not connected") {
		println(intfName, "is not associated")
	} else {
		println(intfName, "is associated")
		arr := strings.Fields(cmdOut)
		rsuMAC = arr[2]
		for ind, v := range arr {
			if strings.Contains(v, "ssid_") {
				println(ind, v)
				SSID = v
				println("mac:", rsuMAC)
				break
			}
		}
	}
	return rsuMAC, SSID
}

func ConnectedToRSU(ifiIndex int) bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return false
	}
	if ifiIndex >= len(interfaces) {
		return false
	}
	mac, _ := getRSU(interfaces[ifiIndex].Name)
	if strings.Compare(mac, "") == 0 {
		return false
	} else {
		return true
	}
}

func GetRSUMac(ifiIndex int) string {
	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatalf("failed to open interface: %v", err)
		return ""
	}
	mac, _ := getRSU(interfaces[ifiIndex].Name)
	return mac
}
