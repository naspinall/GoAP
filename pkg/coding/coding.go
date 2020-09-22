package coding

func EncodeUint16(value uint16) []byte {
	return []byte{byte(value), byte(value >> 8)}
}
func EncodeUint32(value uint32) []byte {
	return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)}
}
func EncodeUint64(value uint64) []byte {
	return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24), byte(value >> 32), byte(value >> 40), byte(value >> 48), byte(value >> 56)}
}
func EncodeUint(value uint) []byte {
	if value > 0xFFFFFFFFFFFFFF {
		return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24), byte(value >> 32), byte(value >> 40), byte(value >> 48), byte(value >> 56)}
	} else if value > 0xFFFFFFFFFFFF {
		return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24), byte(value >> 32), byte(value >> 40), byte(value >> 48)}
	} else if value > 0xFFFFFFFFFF {
		return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24), byte(value >> 32), byte(value >> 40)}
	} else if value > 0xFFFFFFFF {
		return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24), byte(value >> 32)}
	} else if value > 0xFFFFFF {
		return []byte{byte(value), byte(value >> 8), byte(value >> 16), byte(value >> 24)}
	} else if value > 0xFFFF {
		return []byte{byte(value), byte(value >> 8), byte(value >> 16)}
	} else if value > 0xFF {
		return []byte{byte(value), byte(value >> 8)}
	}

	return []byte{byte(value)}
}

func DecodeUint16(input []byte) (value uint16) {
	for index, b := range input {
		if index > 1 {
			return
		}

		// Adding Byte
		value |= uint16(b) << (8 * index)

	}
	return
}

func DecodeUint32(input []byte) (value uint32) {
	for index, b := range input {
		if index > 3 {
			return
		}

		// Adding Byte
		value |= uint32(b) << (8 * index)

	}
	return
}

func DecodeUint64(input []byte) (value uint64) {
	for index, b := range input {
		if index > 7 {
			return
		}

		// Adding Byte
		value |= uint64(b) << (8 * index)

	}
	return
}

func DecodeUint(input []byte) (value uint) {
	return uint(DecodeUint64(input))
}
