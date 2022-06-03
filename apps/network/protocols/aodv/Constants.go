package aodv

const (
	RREQType uint8 = 1
	RREPType uint8 = 2
	RERRsType uint8 = 3
	RREPACKType uint8 = 4
)

const (
	RREQMessageLen = 192
	RREPMessageLen = 160
)

const (
	RREPDefaultLifeTimeMS uint32 = 30 * 1000
)