package aodv

// Types
const (
	RREQType    uint8 = 1
	RREPType    uint8 = 2
	RERRType    uint8 = 3
	RREPACKType uint8 = 4
	DataType    uint8 = 5
)

// Flags
const (
	RREQFlagJ uint16 = 1 << 0
	RREQFlagR uint16 = 1 << 1
	RREQFlagG uint16 = 1 << 2
	RREQFlagD uint16 = 1 << 3
	RREQFlagU uint16 = 1 << 4

	RREPFlagR uint16 = 1 << 0
	RREPFlagA uint16 = 1 << 1
)

// Lengths
const (
	RREQMessageLen = 24
	RREPMessageLen = 20
	RERRMessageLen = 12
	DataMessageLen = 16
)

// Default values
const (
	RREPDefaultLifeTimeMS  uint32 = 5 * 60 * 1000 // 5 mins
	ActiveRouteTimeMS      uint32 = 60 * 1000     // 1 min
	RSUActiveRouteTimeMS   uint32 = 3000          // 3 seconds
	PATH_DISCOVERY_TIME_MS uint32 = 500           // 500 ms
)

// Limits
const (
	HopCountLimit = 20
)
