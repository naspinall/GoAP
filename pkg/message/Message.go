package messages

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"io/ioutil"
)

type MessageType uint8

// Type Values
const (
	Confirmable     MessageType = 0
	NonConfirmable  MessageType = 1
	Acknowledgement MessageType = 2
	Reset           MessageType = 3
)

type MessagesConfig func(*Message) error

//Code Values
const (
	// Empty Value
	Empty uint8 = 0

	// Request Value
	GET    uint8 = 1
	POST   uint8 = 2
	PUT    uint8 = 3
	DELETE uint8 = 4

	//Response Values
	Created               uint8 = 65
	Deleted               uint8 = 66
	Valid                 uint8 = 67
	Changed               uint8 = 68
	Content               uint8 = 69
	Bad                   uint8 = 128
	Unauthorized          uint8 = 129
	BadOption             uint8 = 130
	Forbidden             uint8 = 131
	NotFound              uint8 = 132
	MethodNotAllowed      uint8 = 133
	NotAcceptable         uint8 = 134
	PreconditionFailed    uint8 = 140
	RequestEntityTooLarge uint8 = 141
	UnsupportedContent    uint8 = 143
	InternalServerError   uint8 = 160
	NotImplemented        uint8 = 161
	BadGateway            uint8 = 162
	ServiceUnavailable    uint8 = 163
	GatewayTimeout        uint8 = 164
	ProxyingNotSupported  uint8 = 165
)

type Message struct {
	Version     uint8       // CoAP Version Number
	Type        MessageType // 2 bit unsigned integer, 0 Confirmable, 1 Non-Confirmable, 2 Acknowledgement (2), or Reset (3).
	TokenLength uint8       // 4 bit unsigned integer, length of the token.
	Code        uint8       // Request type (GET,POST,PUT) or response type.
	MessageID   uint16
	Token       uint64
	Options     []Option
	Payload     []byte
	buff        *bytes.Buffer
}

type Option struct {
	OptionNumber uint16 // Option type
	Length       uint16 // Option length
	Value        []byte // Option Value
}

func (m *Message) SetNonConfirmable() *Message {
	m.Type = 1
	return m
}

// Encoding a message
func (m *Message) EncodeHeader() error {

	// Version, Type and Token Length Encoding
	b := m.Version & 0x03
	b = b | uint8(m.Type)&0x03<<2
	b = b | m.TokenLength&0x0F<<4

	// TODO work on this.
	_, err := m.buff.Write([]byte{b, m.Code})
	if err != nil {
		return err
	}

	// Encoding message id
	err = binary.Write(m.buff, binary.LittleEndian, m.MessageID)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) EncodeToken() error {
	// Writing token to the buffer
	err := binary.Write(m.buff, binary.LittleEndian, m.Token)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) EncodeOptions() error {
	var currentDelta uint16

	// Worry about extended options later
	for _, option := range m.Options {
		delta := option.OptionNumber - currentDelta
		if delta > 13 {
			b := uint8(delta - 13)
			err := m.buff.WriteByte(b)
			if err != nil {
				return err
			}
			delta = 13
		} else if delta > 269 {
			b := delta - 269
			_, err := m.buff.Write([]byte{byte(b >> 8), byte(0x00FF & b)})
			if err != nil {
				return err
			}
			delta = 14
		}
		// Encoding Option header
		header := uint8(delta) | uint8(option.Length&0x0F<<4)
		err := m.buff.WriteByte(header)
		if err != nil {
			return err
		}

		n, err := m.buff.Write(option.Value)
		// Checking option length
		if uint16(n) != option.Length {
			return errors.New("Bad Option Length")
		}

		if err != nil {
			return err
		}

	}
	return nil
}

func (m *Message) EncodePayload() error {

	// Write Payload Marker if Options are present
	if len(m.Options) != 0 {
		m.buff.WriteByte(0xFF)
	}
	// Adding padding byte and writing to buffer
	_, err := m.buff.Write(m.Payload)
	if err != nil {
		return err
	}

	return nil
}

func (m *Message) Encode() error {
	if err := m.EncodeHeader(); err != nil {
		return err
	}
	if err := m.EncodeToken(); err != nil {
		return err
	}
	if err := m.EncodeOptions(); err != nil {
		return err
	}
	if err := m.EncodePayload(); err != nil {
		return err
	}
	return nil
}

func (m *Message) Write(w io.Writer) error {
	// Encoding
	if err := m.Encode(); err != nil {
		return err
	}

	// Writing to writer
	if _, err := w.Write(m.buff.Bytes()); err != nil {
		return err
	}
	return nil
}

// Decoding Message
func (m *Message) DecodeHeader() error {
	// Version, Type byte

	b := make([]byte, 4)

	n, err := m.buff.Read(b)

	// Bad Read
	if err != nil {
		return err
	}

	// Malformed Packet
	if n != 4 {
		return errors.New("Malformed Packet")
	}

	m.Version = uint8(b[0] & 0x03)
	m.Type = MessageType(b[0] >> 2 & 0x03)
	m.TokenLength = uint8(b[0] >> 4)

	m.Code = b[1]

	messageID, n := binary.Uvarint(b[2:])
	m.MessageID = uint16(messageID)

	return nil

}

func (m *Message) DecodeToken() error {
	// If no token just skip
	if m.TokenLength == 0 {
		return nil
	}

	b := make([]byte, m.TokenLength)
	n, err := m.buff.Read(b)
	if err != nil {
		return err
	}

	if uint8(n) != m.TokenLength {
		return errors.New("Malformed Packet")
	}

	// Reading the token
	m.Token, _ = binary.Uvarint(b)
	return nil
}

func (m *Message) OneByteOption() (uint16, error) {

	b, err := m.buff.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint16(b - 13), nil

}

func (m *Message) TwoByteOption() (uint16, error) {
	b := make([]byte, 2)
	_, err := m.buff.Read(b)
	if err != nil {
		return 0, err
	}
	option, _ := binary.Uvarint(b)
	return uint16(option - 269), nil
}

func (m *Message) DecodeOptions() error {
	var prevDelta uint16

	// Holds the option header
	b, err := m.buff.ReadByte()
	if err != nil {
		return err
	}

	// 0xFF is the payload indicator
	for b != 0xFF {
		delta := uint16(b & 0xF)

		// Taking into account extended options
		switch delta {
		case 13:
			delta, err = m.OneByteOption()
			if err != nil {
				return nil
			}

		case 14:
			delta, err = m.TwoByteOption()
			if err != nil {
				return nil
			}
		}

		// Getting the length
		length := uint16(b >> 4)
		switch length {
		case 13:
			length, err = m.OneByteOption()
			if err != nil {
				return nil
			}

		case 14:
			length, err = m.TwoByteOption()
			if err != nil {
				return nil
			}
		}

		// Getting option data
		val := make([]byte, length)
		_, err := m.buff.Read(val)
		if err != nil {
			return err
		}

		// Adding option
		m.Options = append(m.Options, Option{
			OptionNumber: delta + prevDelta,
			Length:       length,
			Value:        val,
		})

		prevDelta += delta
		// Reading next header or payload indicator byte
		b, err = m.buff.ReadByte()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *Message) DecodePayload() error {
	b, err := ioutil.ReadAll(m.buff)
	if err != nil {
		return err
	}
	m.Payload = b
	return nil
}

func (m *Message) Decode() error {
	if err := m.DecodeHeader(); err != nil {
		return err
	}
	if err := m.DecodeToken(); err != nil {
		return err
	}
	if err := m.DecodeOptions(); err != nil {
		return err
	}
	if err := m.DecodePayload(); err != nil {
		return err
	}

	// Won't be reading from the buffer anymore so reseting
	m.buff.Reset()
	return nil
}

func FromReader(r io.Reader) (*Message, error) {
	b, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, err
	}
	// Creating new message with bytes as buffer.
	m := &Message{buff: bytes.NewBuffer(b)}

	// Decoding message
	if err := m.Decode(); err != nil {
		return nil, err
	}
	return m, nil
}

func FromBytes(b []byte) (*Message, error) {
	m := NewMessage()
	m.buff = bytes.NewBuffer(b)
	if err := m.Decode(); err != nil {
		return nil, err
	}
	return m, nil
}
