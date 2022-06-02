package main

import (
	"fmt"
	"sort"
)

type TypeEnum string;

const (
	String TypeEnum = "string"
	Int     			  = "int"
	Bool     				= "bool"
	ByteArray 			= "bytearray"
	Frame 					= "frame"
	
)

type PacketMetaData struct {
	// length is -1 if the length is variable (bytearray and string)
	length int
	valueType TypeEnum
}

func ConvertBytesToString(bytes []byte) string {
	return string(bytes)
}

func ConvertStringToBytes(str string) []byte {
	b:= make([]byte, len(str) + 1)
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

func orderedKeys(m map[string]PacketMetaData) []string {
	// get ordered keys
	keys := make([]string, 0)
	for k, _ := range m {
			keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}


func Marshal(packet map[string]any, packetLengths map[string]PacketMetaData) ([]byte, error)  {
	var packetBytes []byte

	for _, key := range orderedKeys(packetLengths) {
		meta := packetLengths[key]
		switch meta.valueType {
		case String:
			packetBytes = append(packetBytes, ConvertStringToBytes(packet[key].(string))...)
		case Int:
			packetBytes = append(packetBytes, ConvertIntToBytes(packet[key].(int), meta.length)...)
		case Bool:
			packetBytes = append(packetBytes, ConvertBoolToBytes(packet[key].(bool))...)
		case ByteArray:
			packetBytes = append(packetBytes, ConvertBytes(packet[key].([]byte))...)
		}
	}
	return packetBytes, nil
}

func Unmarshal(packetBytes []byte, packetLengths map[string]PacketMetaData) (map[string]any, error) {
	packet := make(map[string]any)

	for _, key := range orderedKeys(packetLengths) {
		meta := packetLengths[key]
		switch meta.valueType {
		case String:
			len := GetLength(packetBytes[0])
			packet[key] = ConvertBytesToString(packetBytes[1:len+1])
			packetBytes = packetBytes[len+1:]
		case Int:
			packet[key] = ConvertBytesToInt(packetBytes[:meta.length], meta.length)
			packetBytes = packetBytes[meta.length:]
		case Bool:
			packet[key] = ConvertBytesToBool(packetBytes[:meta.length])
			packetBytes = packetBytes[meta.length:]
		case ByteArray:
			len := GetLength(packetBytes[0])
			packet[key] = packetBytes[1:len+1]
			packetBytes = packetBytes[len+1:]
		}
		p := packet[key]
		println(key)
		fmt.Printf("v: %v\n", p)

	}
	return packet, nil
}


func main()  {
	str:= "hello"
	i32 := 12345
	i16 := 123
	i8 := 123
	b := true
	ba := []byte{1,2,3,4,5}

	println(ba)

	metadata := map[string]PacketMetaData{
		"a": {length: -1, valueType: String},
		"b": {length: 4, valueType: Int},
		"c": {length: 2, valueType: Int},
		"d": {length: 1, valueType: Int},
		"e": {length: 1, valueType: Bool},
		"f": {length: -1, valueType: ByteArray},
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