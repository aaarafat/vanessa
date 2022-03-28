package aodv

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"syscall"
	"time"
)

// Global definitions for the Ethernet IEEE 802.3 interface.
// Source: https://github.com/torvalds/linux/blob/master/include/uapi/linux/if_ether.h
const ROUTE_TTL = 200
const HELLO_INTERVAL = 7
const NEIGHBOR_TTL = 15

const ETH_ALEN = 6;                // Octets in one ethernet addr
const ETH_TLEN = 2;                // Octets in ethernet type field
const ETH_HLEN = 14;               // Total octets in header.
const ETH_ZLEN = 60               // Min. octets in frame sans FCS
const ETH_DATA_LEN = 1500         // Max. octets in payload
const ETH_FRAME_LEN = 1514        // Max. octets in frame sans FCS

const ETH_P_ALL = 0x0003          // Every packet (be careful!!!)
const ETH_P_IP = 0x0800           // Internet Protocol packet
const ETH_P_ARP = 0x0806          // Address Resolution packet
const ETH_P_802_EX1 = 0x88B5      // Local Experimental Ethertype 1
const ETH_P_802_EX2 = 0x88B6      // Local Experimental Ethertype 2


const MetadataPath = "/tmp/mn.metadata.json" // Metadata Path

// AODV Structure
type Aodv struct {
	nodeId int
	broadcastId int
	nodeType string
	numNodes int
	nodeName string
	nodeInterface string
	nodeIp string
	nodeMac string
	nodePort int
	seqNum int
	routingTable map[string]string
	neighborTable map[string]string
}

// Metadata Structure
type Metadata struct {
	Mac map[string] string `json:"mac"`
	Mac2Ip map[string] string `json:"mac2ip"`
}


// Create new AODV based router
func New(station string, numberOfNodes int) (*Aodv, error) {
	aodv := new(Aodv)
	id, err := strconv.Atoi(station[3:])
	if err != nil {
		return nil, err
	}

	aodv.nodeId = id
	aodv.nodeType = station[:3]
	aodv.numNodes = numberOfNodes
	aodv.nodeInterface = station + "-wlan0"
	aodv.nodeName = station
	aodv.broadcastId = 0
	aodv.seqNum = 0
	aodv.nodeMac = ""
	aodv.nodeIp = ""
	aodv.nodePort = 0
	// --------------------------
	// TODO: Add socket and timer
	// ---------------------------
	aodv.routingTable = make(map[string]string)
	aodv.neighborTable = make(map[string]string)
	return aodv, nil
}

// Run AODV Router
func (aodv* Aodv) Run() {
	// TODO: Run AODV in a new go routine
	// get aodv mac address
	aodv.getMacAddress()
	// set aodv port
	aodv.setPort()

	time.Sleep(1 * time.Second)
	// update neighbors
	


	// init Inet socket
	socket, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, syscall.IPPROTO_TCP)
	
	if err != nil {
		panic(err)
	}
	defer syscall.Close(socket)
	fmt.Println("Socket created : ", socket)

	// bind and start listening 
	syscall.Bind(socket, &syscall.SockaddrInet4{Port: aodv.nodePort, Addr: [4]byte{127, 0, 0, 1}})
	syscall.Listen(socket, 1)

	// accept connection
	conn, addr, err := syscall.Accept(socket)

	if err != nil {
		panic(err)
	}

	// read data from socket
	fmt.Println("Connected to: ", addr)
	data := make([]byte, 1024)
	for {
		n, err := syscall.Read(conn, data)
		if err != nil {
			panic(err)
		}
		fmt.Println("Received: ", string(data[:n]))
		syscall.Write(conn, data)
	}
}

// Get Node Mac Address
func (aodv* Aodv) getMacAddress() string {
	// Open metadata json file
	metadataFile, err := os.Open(MetadataPath)
	// handle errors
	if err != nil {
		panic(err)
	}
	fmt.Printf("Read Metadata file successfully in node with ID=%d\n", aodv.nodeId)
	// Close the file after the function is done
	defer metadataFile.Close()

	// read data from metadata file
	byteValue, _ := ioutil.ReadAll(metadataFile)
	var metadata Metadata	
	json.Unmarshal(byteValue, &metadata)
	
	// update aodv
	aodv.nodeMac = metadata.Mac[aodv.nodeName]
	aodv.nodeIp = metadata.Mac2Ip[aodv.nodeMac]

	fmt.Printf("Mac Address : %s, IP : %s for node with ID=%d\n", aodv.nodeMac, aodv.nodeIp, aodv.nodeId)

	return aodv.nodeMac
} 

// Set Node Port
func (aodv* Aodv) setPort() int {
	if aodv.nodeType == "rsu" {
		aodv.nodePort = 3000 + aodv.nodeId
	} else if aodv.nodeType == "car" {
		aodv.nodePort = 4000 + aodv.nodeId
	}

	fmt.Printf("Port Number = %d for node with ID=%d\n", aodv.nodePort, aodv.nodeId)

	return aodv.nodePort
}

// update neighbors
func (aodv* Aodv) updateNeighbors() [] string {
	var sources []string
	
	// open socket
	socket, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, ETH_P_ALL)
	if err != nil {
		panic(err)
	}
	defer syscall.Close(socket)

	// TODO: complete this function

	return sources
}