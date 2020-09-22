package messages

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/naspinall/GoAP/pkg/coding"
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
	ContentFormat uint
	ETag          [][]byte
	LocationPath  []string
	LocationQuery []string
	MaxAge        uint
	ProxyURI      *string
	ProxyScheme   *string
	URIHost       *string
	URIPath       []string
	URIPort       uint
	URIQuery      []string
	Accept        uint
	IfMatch       [][]byte
	IfNoneMatch   bool
	Size1         uint
}

func (o *Options) SetURI(rawurl string) error {
	parsedURL, err := url.Parse(rawurl)
	if err != nil {
		return err
	}

	parsedPort := parsedURL.Port()
	if parsedPort == "" {
		parsedPort = "5683"
	}

	// Getting port
	portInt, err := strconv.Atoi(parsedPort)
	if err != nil {
		return err
	}

	o.URIPort = uint(portInt)

	// Getting Host
	host := parsedURL.Hostname()
	o.URIHost = &host

	// Getting and splitting path
	path := strings.Split(parsedURL.Path, "/")
	for index, pathElement := range path {
		path[index] = strings.TrimSpace(pathElement)
	}
	o.URIPath = path[1:len(path)]

	return nil
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
		o.URIPort = coding.DecodeUint(b)

	// Location Path
	case LocationPath:
		o.LocationPath = append(o.LocationPath, string(b))

	// URI-Path
	case URIPath:
		o.URIPath = append(o.URIPath, string(b))

	// Content Format
	case ContentFormat:
		o.ContentFormat = coding.DecodeUint(b)
	//Max-Age
	case MaxAge:
		o.MaxAge = coding.DecodeUint(b)

	// URI-Query
	case URIQuery:
		o.URIQuery = append(o.URIQuery, string(b))
	// Accept
	case Accept:
		o.Accept = coding.DecodeUint(b)

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
		o.Size1 = coding.DecodeUint(b)
	}
	return nil
}

func EncodeSingleOption(delta uint, b []byte) ([]byte, error) {
	var header byte
	var extendedOptions []byte

	// Encoding Delta
	if delta >= 13 {
		header = 0x0D
		extendedDelta := byte(delta - 13)
		extendedOptions = append(extendedOptions, extendedDelta)
	} else if delta >= 269 {
		header = 0x0E
		extendedDelta := uint16(delta - 269)
		extendedOptions = append(extendedOptions, byte(extendedDelta), byte(extendedDelta>>8))
	} else {
		header = byte(delta)
	}

	length := uint(len(b))
	// Encoding Length
	if length >= 13 {
		header = 0xD0 ^ header
		extendedLength := byte(length - 13)
		extendedOptions = append(extendedOptions, extendedLength)
	} else if length >= 269 {
		header = 0xE0 ^ header
		extendedLength := uint16(length - 269)
		extendedOptions = append(extendedOptions, byte(extendedLength), byte(extendedLength>>8))
	} else {
		header = byte(header ^ byte(length<<4))
	}

	// Adding header to start of slice with extended options or lengths
	return append(append([]byte{header}, extendedOptions...), b...), nil
}

func (o *Options) EncodeOptions() ([]byte, error) {
	var total []byte
	var previousValue uint = 0

	{
		delta := ContentFormat - previousValue
		previousValue = ContentFormat

		value := coding.EncodeUint(o.ContentFormat)
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

	{
		delta := MaxAge - previousValue
		previousValue = MaxAge

		value := coding.EncodeUint(o.MaxAge)
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

	{
		delta := URIPort - previousValue
		previousValue = URIPort

		value := coding.EncodeUint(o.URIPort)
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

	{
		delta := Accept - previousValue
		previousValue = Accept

		value := coding.EncodeUint(o.Accept)
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

	{
		delta := Size1 - previousValue
		previousValue = Size1

		value := coding.EncodeUint(o.URIPort)
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
