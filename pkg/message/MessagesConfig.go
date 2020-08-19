package messages

import (
	"bytes"
	"encoding/binary"
)

// Create a new message with sane defaults
func NewMessage(cfgs ...MessagesConfig) *Message {
	m := &Message{
		Version: 1,
		buff:    &bytes.Buffer{},
	}

	for _, cfg := range cfgs {
		if err := cfg(m); err != nil {
			return nil
		}
	}
	return m
}

func (m *Message) GET() *Message {
	m.Code = GET
	return m
}

func (m *Message) POST() *Message {
	m.Code = POST
	return m
}

func (m *Message) PUT() *Message {
	m.Code = PUT
	return m
}

func (m *Message) DELETE() *Message {
	m.Code = DELETE
	return m
}

func (m *Message) SetToken(b []byte) *Message {
	m.TokenLength = uint8(len(b))
	m.Token, _ = binary.Uvarint(b)
	return m
}

func (m *Message) SetMessageID(id uint16) *Message {
	m.MessageID = id
	return m
}

func (m *Message) SetPayload(b []byte) *Message {
	m.Payload = b
	return m
}

func WithToken(b []byte) MessagesConfig {
	return func(m *Message) error {
		m.SetToken(b)
		return nil
	}
}

func WithMessageID(id uint16) MessagesConfig {
	return func(m *Message) error {
		m.SetMessageID(id)
		return nil
	}
}

func WithPayload(b []byte) MessagesConfig {
	return func(m *Message) error {
		m.SetPayload(b)
		return nil
	}
}

func Get() MessagesConfig {
	return func(m *Message) error {
		m.GET()
		return nil
	}
}
func Post() MessagesConfig {
	return func(m *Message) error {
		m.POST()
		return nil
	}
}
func Put() MessagesConfig {
	return func(m *Message) error {
		m.PUT()
		return nil
	}
}
func Delete() MessagesConfig {
	return func(m *Message) error {
		m.DELETE()
		return nil
	}
}
