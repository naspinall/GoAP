package messages

import (
	"encoding/binary"
	"errors"
)

const (
	IfMatch       uint = 1
	URIHost       uint = 3
	ETag          uint = 4
	IfNoneMatch   uint = 5
	URIPort       uint = 7
	LocationPath  uint = 8
	URIPath       uint = 11
	ContentFormat uint = 12
	MaxAge        uint = 14
	URIQuery      uint = 15
	Accept        uint = 17
	LocationQuery uint = 20
	ProxyURI      uint = 35
	ProxyScheme   uint = 39
	Size1         uint = 60
)

// Options Data Types
// Empty  zero length sequence of bytes
// Opaque  Bytes
// uint  uint, length is given by option length
// string  UTF8 string

type Options struct {
	ContentFormat *uint
	ETag          [][]byte
	LocationPath  []string
	LocationQuery []string
	MaxAge        *uint
	ProxyURI      *string
	ProxyScheme   *string
	URIHost       *string
	URIPath       []string
	URIPort       *uint
	URIQuery      []string
	Accept        *uint
	IfMatch       [][]byte
	IfNoneMatch   bool
	Size1         *uint
}

func (o *Options) DecodeOption(number uint, b []byte) error {
	switch number {
	// If-Match
	case IfMatch:
		o.IfMatch = append(o.IfMatch, b)

	// URI-Host
	case URIHost:
		host := string(b)
		o.URIHost = &host

	// ETag
	case ETag:
		o.ETag = append(o.ETag, b)

	// If-None-Match
	case IfNoneMatch:
		o.IfNoneMatch = true

	// URI-Port
	case URIPort:
		port, err := ParseUint(b)
		if err != nil {
			return err
		}
		o.URIPort = &port

	// Location Path
	case LocationPath:
		o.LocationPath = append(o.LocationPath, string(b))

	// URI-Path
	case URIPath:
		o.URIPath = append(o.URIPath, string(b))

	// Content Format
	case ContentFormat:
		contentFormat, err := ParseUint(b)
		if err != nil {
			return err
		}
		o.ContentFormat = &contentFormat
	//Max-Age
	case MaxAge:
		maxAge, err := ParseUint(b)
		if err != nil {
			return err
		}
		o.MaxAge = &maxAge

	// URI-Query
	case URIQuery:
		o.URIQuery = append(o.URIQuery, string(b))
	// Accept
	case Accept:
		accept, err := ParseUint(b)
		if err != nil {
			return err
		}
		o.Accept = &accept

	// Location Query
	case LocationQuery:
		o.LocationQuery = append(o.LocationQuery, string(b))

	// Proxy-URI
	case ProxyURI:
		proxyURI := string(b)
		o.ProxyURI = &proxyURI

	// Proxy-Scheme
	case ProxyScheme:
		proxyScheme := string(b)
		o.ProxyScheme = &proxyScheme

	// Size1
	case Size1:
		sizeOne, err := ParseUint(b)
		if err != nil {
			return err
		}
		o.Size1 = &sizeOne
	}
	return nil
}

func EncodeSingleOption(delta uint, b []byte) ([]byte, error) {
	var header byte
	var extendedOptions []byte

	// Encoding Delta
	if delta > 13 {
		header = 0xD0
		extendedDelta := byte(delta - 13)
		extendedOptions = append(b, extendedDelta)
	} else if delta > 269 {
		header = 0xE0
		extendedDelta := uint16(delta - 269)
		extendedOptions = append(b, byte(extendedDelta), byte(extendedDelta>>8))
	} else {
		header = byte(delta)
	}

	length := uint(len(b))
	// Encoding Length
	if length > 13 {
		header = header & 0xFD
		extendedLength := byte(length - 13)
		extendedOptions = append(b, extendedLength)
	} else if length > 269 {
		header = header & 0xFE
		extendedLength := uint16(length - 269)
		extendedOptions = append(b, byte(extendedLength), byte(extendedLength>>8))
	} else {
		header = byte(length >> 4)
	}

	// Adding header to start of slice with extended options or lengths
	return append(append([]byte{header}, extendedOptions...), b...), nil

}

func (o *Options) EncodeOptions() ([]byte, error) {
	var total []byte
	var previousValue uint = 0

	if o.ContentFormat != nil {
		delta := ContentFormat - previousValue
		previousValue = ContentFormat

		value := UintToBytes(*o.ContentFormat)
		b, err := EncodeSingleOption(delta, value)
		if err != nil {
			return nil, err
		}

		total = append(total, b...)
	}

	if o.ETag != nil {
		delta := ETag - previousValue
		previousValue = ETag

		for _, eTag := range o.ETag {

			b, err := EncodeSingleOption(delta, eTag)
			if err != nil {
				return nil, err
			}
			total = append(total, b...)
			// If there are more than one eTag then delta is zero.
			delta = 0
		}

	}

	if o.LocationPath != nil {
		delta := LocationPath - previousValue
		previousValue = LocationPath

		for _, path := range o.LocationPath {

			b, err := EncodeSingleOption(delta, []byte(path))
			if err != nil {
				return nil, err
			}
			total = append(total, b...)
			// If there are more than one eTag then delta is zero.
			delta = 0
		}
	}

	if o.LocationQuery != nil {
		delta := LocationQuery - previousValue
		previousValue = LocationQuery

		for _, query := range o.LocationQuery {

			b, err := EncodeSingleOption(delta, []byte(query))
			if err != nil {
				return nil, err
			}
			total = append(total, b...)
			// If there are more than one eTag then delta is zero.
			delta = 0
		}

	}

	if o.MaxAge != nil {
		delta := MaxAge - previousValue
		previousValue = MaxAge

		value := UintToBytes(*o.MaxAge)
		b, err := EncodeSingleOption(delta, value)
		if err != nil {
			return nil, err
		}

		total = append(total, b...)
	}

	if o.ProxyURI != nil {
		delta := ProxyURI - previousValue
		previousValue = ProxyURI

		b, err := EncodeSingleOption(delta, []byte(*o.ProxyURI))
		if err != nil {
			return nil, err
		}

		total = append(total, b...)
	}

	if o.ProxyScheme != nil {
		delta := ProxyScheme - previousValue
		previousValue = ProxyScheme

		b, err := EncodeSingleOption(delta, []byte(*o.ProxyScheme))
		if err != nil {
			return nil, err
		}

		total = append(total, b...)
	}

	if o.URIHost != nil {
		delta := URIHost - previousValue
		previousValue = URIHost

		b, err := EncodeSingleOption(delta, []byte(*o.URIHost))
		if err != nil {
			return nil, err
		}

		total = append(total, b...)

	}

	if o.URIPath != nil {
		delta := URIPath - previousValue
		previousValue = URIPath

		for _, path := range o.URIPath {

			b, err := EncodeSingleOption(delta, []byte(path))
			if err != nil {
				return nil, err
			}
			total = append(total, b...)
			// If there are more than one eTag then delta is zero.
			delta = 0
		}
	}

	if o.URIPort != nil {
		delta := URIPort - previousValue
		previousValue = URIPort

		value := UintToBytes(*o.URIPort)
		b, err := EncodeSingleOption(delta, value)
		if err != nil {
			return nil, err
		}

		total = append(total, b...)
	}

	if o.URIQuery != nil {
		delta := URIQuery - previousValue
		previousValue = URIQuery

		for _, query := range o.URIQuery {

			b, err := EncodeSingleOption(delta, []byte(query))
			if err != nil {
				return nil, err
			}
			total = append(total, b...)
			// If there are more than one eTag then delta is zero.
			delta = 0
		}
	}

	if o.Accept != nil {
		delta := Accept - previousValue
		previousValue = Accept

		value := UintToBytes(*o.Accept)
		b, err := EncodeSingleOption(delta, value)
		if err != nil {
			return nil, err
		}

		total = append(total, b...)

	}

	if o.IfMatch != nil {
		delta := IfMatch - previousValue
		previousValue = IfMatch

		for _, match := range o.IfMatch {

			b, err := EncodeSingleOption(delta, match)
			if err != nil {
				return nil, err
			}
			total = append(total, b...)
			// If there are more than one eTag then delta is zero.
			delta = 0
		}

	}

	if o.IfNoneMatch {
		delta := IfNoneMatch - previousValue
		previousValue = IfNoneMatch

		b, err := EncodeSingleOption(delta, []byte{})
		if err != nil {
			return nil, err
		}
		total = append(total, b...)
		// If there are more than one eTag then delta is zero.
		delta = 0

	}

	if o.Size1 != nil {
		delta := Size1 - previousValue
		previousValue = Size1

		value := UintToBytes(*o.URIPort)
		b, err := EncodeSingleOption(delta, value)
		if err != nil {
			return nil, err
		}
		total = append(total, b...)
		// If there are more than one eTag then delta is zero.
		delta = 0
	}

	return total, nil
}

func ParseUint(value []byte) (uint, error) {
	val, n := binary.Uvarint(value)
	if n < 0 {
		return 0, errors.New("Bad Option Provided")
	}
	return uint(val), nil
}

// Return minimum number of bytes for uint
func UintToBytes(value uint) (b []byte) {

	if value > 0xFF {
		b = append(b, byte(value))
	}
	if value > 0xFFFF {
		b = append(b, byte(value>>56))
	}
	if value > 0xFFFF {
		b = append(b, byte(value>>48))
	}
	if value > 0xFFFFFF {
		b = append(b, byte(value>>40))
	}
	if value > 0xFFFFFFFF {
		b = append(b, byte(value>>32))
	}
	if value > 0xFFFFFFFFFF {
		b = append(b, byte(value>>24))
	}
	if value > 0xFFFFFFFFFFFF {
		b = append(b, byte(value>>16))
	}
	if value > 0xFFFFFFFFFFFFFF {
		b = append(b, byte(value>>8))
	}

	for i, valueByte := range b {
		if valueByte != 0x00 {
			return b[i:]
		} else if i == len(b)-1 {
			return b[i:]
		}
	}
	return nil
}
