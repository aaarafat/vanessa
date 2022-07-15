package utils

import (
	"fmt"
	"sort"
	"time"
)

type TypeEnum string

const (
	String    TypeEnum = "string"
	Int32              = "int32"
	Int16              = "int16"
	Int8               = "int8"
	Bool               = "bool"
	ByteArray          = "bytearray"
	Frame              = "frame"
)

// Get the current time in milliseconds
func TimeInMs() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func ConvertBytesToString(bytes []byte) string {
	return string(bytes)
}

func ConvertStringToBytes(str string) []byte {
	b := make([]byte, len(str)+1)
	b[0] = byte(len(str))
	copy(b[1:], str)
	return b
}

func ConvertBytesToInt(bytes []byte, size int) int {
	var number int = 0
	for i := 0; i < size; i++ {
		number = number << 8
		number = number | int(bytes[i])
	}
	return number
}

func ConvertIntToBytes(number int, size int) []byte {
	bytes := make([]byte, size)
	for i := 0; i < size; i++ {
		bytes[i] = byte(number >> int(8*(size-i-1)))
	}
	return bytes
}

func ConvertBytesToBool(bytes []byte) bool {
	return bytes[0] == 1
}

func ConvertBoolToBytes(b bool) []byte {
	if b {
		return []byte{1}
	} else {
		return []byte{0}
	}
}

func ConvertBytes(bytes []byte) []byte {
	var b []byte
	b = append(b, byte(len(bytes)))
	b = append(b, bytes...)
	return b
}

func GetLength(b byte) int {
	return ConvertBytesToInt([]byte{b}, 1)
}

func orderedKeys(m map[string]TypeEnum) []string {
	// get ordered keys
	keys := make([]string, 0)
	for k, _ := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func Marshal(packet map[string]any, packetTypes map[string]TypeEnum) ([]byte, error) {
	var packetBytes []byte

	for _, key := range orderedKeys(packetTypes) {
		valueType := packetTypes[key]
		switch valueType {
		case String:
			packetBytes = append(packetBytes, ConvertStringToBytes(packet[key].(string))...)
		case Int32:
			packetBytes = append(packetBytes, ConvertIntToBytes(packet[key].(int), 4)...)
		case Int16:
			packetBytes = append(packetBytes, ConvertIntToBytes(packet[key].(int), 2)...)
		case Int8:
			packetBytes = append(packetBytes, ConvertIntToBytes(packet[key].(int), 1)...)
		case Bool:
			packetBytes = append(packetBytes, ConvertBoolToBytes(packet[key].(bool))...)
		case ByteArray:
			packetBytes = append(packetBytes, ConvertBytes(packet[key].([]byte))...)
		}
	}
	return packetBytes, nil
}

func Unmarshal(packetBytes []byte, packetTypes map[string]TypeEnum) (map[string]any, error) {
	packet := make(map[string]any)

	for _, key := range orderedKeys(packetTypes) {
		valueType := packetTypes[key]
		switch valueType {
		case String:
			len := GetLength(packetBytes[0])
			packet[key] = ConvertBytesToString(packetBytes[1 : len+1])
			packetBytes = packetBytes[len+1:]
		case Int32:
			packet[key] = ConvertBytesToInt(packetBytes[:4], 4)
			packetBytes = packetBytes[4:]
		case Int16:
			packet[key] = ConvertBytesToInt(packetBytes[:2], 2)
			packetBytes = packetBytes[2:]
		case Int8:
			packet[key] = ConvertBytesToInt(packetBytes[:1], 1)
			packetBytes = packetBytes[1:]
		case Bool:
			packet[key] = ConvertBytesToBool(packetBytes[:1])
			packetBytes = packetBytes[1:]
		case ByteArray:
			len := GetLength(packetBytes[0])
			packet[key] = packetBytes[1 : len+1]
			packetBytes = packetBytes[len+1:]
		}
		p := packet[key]
		println(key)
		fmt.Printf("v: %v\n", p)

	}
	return packet, nil
}

func main() {
	str := "hello"
	i32 := 12345
	i16 := 123
	i8 := 123
	b := true
	ba := []byte{1, 2, 3, 4, 5}

	println(ba)

	metadata := map[string]TypeEnum{
		"a": String,
		"b": Int32,
		"c": Int16,
		"d": Int8,
		"e": Bool,
		"f": ByteArray,
	}

	obj := map[string]any{
		"a": str,
		"b": i32,
		"c": i16,
		"d": i8,
		"e": b,
		"f": ba,
	}

	println(obj)

	bytes, err := Marshal(obj, metadata)

	if err != nil {
		println(err)
	}

	println(bytes)

	packetBytes, err := Unmarshal(bytes, metadata)

	if err != nil {
		print(err)
	}

	print(packetBytes)

}
