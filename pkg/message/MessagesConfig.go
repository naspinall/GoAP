package messages

import (
	"bytes"
	"encoding/binary"
)

func (m *Message) AsAcknowledge() *Message {

	// Reusing memory, clearing values
	m.Type = Acknowledgement
	m.Code = Empty
	m.TokenLength = 0
	m.Token = 0
	m.Options = nil
	m.Payload = nil

	// Reseting buffer
	m.buff.Reset()

	return m
}

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

func (m *Message) SetToken(token uint64) *Message {
	m.TokenLength = uint8(binary.Size(token))
	m.Token = token
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

func (m *Message) SetType(Type MessageType) *Message {
	m.Type = Type
	return m
}

func WithType(Type MessageType) MessagesConfig {
	return func(m *Message) error {
		m.SetType(Type)
		return nil
	}
}

func WithToken(token uint64) MessagesConfig {
	return func(m *Message) error {
		m.SetToken(token)
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
