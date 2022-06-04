package aodv

// Types
const (
	RREQType uint8 = 1
	RREPType uint8 = 2
	RERRsType uint8 = 3
	RREPACKType uint8 = 4
)

// Flags
const (
	RREQFlagJ uint16 = 1 << 0
	RREQFlagR uint16 = 1 << 1
	RREQFlagG uint16 = 1 << 2
	RREQFlagD uint16 = 1 << 3
	RREQFlagU uint16 = 1 << 4
)

// Lengths
const (
	RREQMessageLen = 192
	RREPMessageLen = 160
)

// Default values
const (
	RREPDefaultLifeTimeMS uint32 = 5 * 60 * 1000 // 5 mins
	ActiveRouteTimeMS uint32 = 60 * 1000 // 1 min
)

// Limits 
const (
	HopCountLimit = 20
)