package main

import (
	"fmt"
	"log"
	"net"
	"os/exec"
	"strings"

	. "github.com/aaarafat/vanessa/apps/network/datalink"
	"github.com/aaarafat/vanessa/apps/network/network/ip"
	. "github.com/aaarafat/vanessa/apps/network/network/messages"
)

func read(d *DataLinkLayerChannel, index int) {
	for {

		payload, addr, err := d.Read()
		if err != nil {
			log.Fatalf("failed to read from channel: %v", err)
		}
		p, err := ip.UnmarshalPacket(payload)
		if err != nil {
			log.Fatalf("failed to Unmarshal IP: %v", err)
		}
		payload = p.Payload
	
		fmt.Println()
		log.Printf("Received \"%s\" from: [%s] on intf-%d", string(payload), addr.String(), index)


		msgType := uint8(payload[0])
		// handle the message
		switch msgType {
		case VObstacleType:
			VOb,_:= UnmarshalVObstacle(payload)
			fmt.Println(VOb.String())

		case VOREPType:
			VOREP2 := UnmarshalVOREP(payload)
			obstacles2 := make([][2]float64, int(VOREP2.Length))
			fmt.Println("Check the obstcales")
			for i := 0; i < len(obstacles2); i++ {
				obstacles2[i][0] = Float64frombytes(VOREP2.Obstcales[i*16:i*16+8])
				obstacles2[i][1] = Float64frombytes(VOREP2.Obstcales[i*16+8:i*16+16])
			}
			fmt.Println(obstacles2)
		default:
			log.Println("Unknown message type: ", msgType)
		}
		
		
		}

}

func getRSU(intfName string) (string, string){

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
			if strings.Contains(v, "ssid_"){
				println(ind, v)
				SSID = v
				println("mac:", rsuMAC)
				break
			}
  		}
	}
	return rsuMAC, SSID
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

func main() {
	
	drsu, err := NewDataLinkLayerChannelWithInterface(VDATAEtherType, 2)
	if err != nil {
		log.Fatalf("failed to create channel: %v", err)
	}
	go read(drsu, 2)
	ifis,_ := net.Interfaces()
	myip,_,_ := MyIP(&ifis[2]) 
	macstr , _ := getRSU(ifis[2].Name)
	mac, _ := net.ParseMAC(macstr)
	var mtype int
	for {

		fmt.Scanf("%d", &mtype)
		switch mtype {
		case 0:
			log.Println("Sending Heartbeat")
			VHB := NewVHBeatMessage(myip , Position{Lat: 0, Lng: 0})
			packet := ip.NewIPPacket(VHB.Marshal(),myip,myip)
			bytes := ip.MarshalIPPacket(packet)
			ip.UpdateChecksum(bytes)
			drsu.SendTo([]byte(bytes), mac)
		case 1:
			var lng , lat float64
			fmt.Scanf("%f %f", &lat , &lng)
			log.Println("Sending Obstacle Alert")
			VOb := NewVObstacleMessage(myip, Position{Lat: lat, Lng: lng},0)
			packet := ip.NewIPPacket(VOb.Marshal(),myip,myip)
			bytes := ip.MarshalIPPacket(packet)
			ip.UpdateChecksum(bytes)
			drsu.SendTo([]byte(bytes), mac)
		case 2:
			log.Println("Sending Obstacles REQ")
			VOREQ := NewVOREQMessage(myip)
			packet := ip.NewIPPacket(VOREQ.Marshal(),myip,myip)
			bytes := ip.MarshalIPPacket(packet)
			ip.UpdateChecksum(bytes)
			drsu.SendTo([]byte(bytes), mac)
		}
		
	}
}
