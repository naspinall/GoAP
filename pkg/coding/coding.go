package coding

import (
	"encoding/binary"
)

func EncodeUint16(value uint16) (b []byte) {
	b = make([]byte, 2)
	binary.BigEndian.PutUint16(b, value)
	return
}
func EncodeUint32(value uint32) (b []byte) {
	b = make([]byte, 4)
	binary.BigEndian.PutUint32(b, value)
	return
}
func EncodeUint64(value uint64) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, value)
	return
}
func EncodeUint(value uint) (b []byte) {
	b = make([]byte, 8)
	binary.BigEndian.PutUint64(b, uint64(value))

	// Removing zero padding
	for index, val := range b {
		if val != 0 {
			b = b[index:]
			return
		}
	}

	return
}

func DecodeUint16(input []byte) (value uint16) {
	input = ZeroPad(input, 2)
	return binary.BigEndian.Uint16(input)
}

func DecodeUint32(input []byte) (value uint32) {
	input = ZeroPad(input, 4)
	return binary.BigEndian.Uint32(input)
}

func DecodeUint64(input []byte) (value uint64) {
	input = ZeroPad(input, 8)
	return binary.BigEndian.Uint64(input)
}

func DecodeUint(input []byte) (value uint) {

	return uint(DecodeUint64(input))
}

func ZeroPad(input []byte, length int) []byte {
	diff := length - len(input)

	for i := 0; i < diff; i++ {
		input = append(input, 0)
	}

	return input
}
