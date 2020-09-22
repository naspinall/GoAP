package messages

import (
	"bytes"
	"errors"
)

func (m *Message) AsAcknowledge() *Message {

	// Reusing memory, clearing values
	m.Type = Acknowledgement
	m.Code = Empty
	m.Token = 0
	//m.Options = nil
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

func WithContentType(contentType string) MessagesConfig {
	return func(m *Message) error {

		if contentType == "text/plain" {
			m.Options.ContentFormat = TextPlain
		} else if contentType == "application/link-format" {
			m.Options.ContentFormat = LinkFormat
		} else if contentType == "application/xml" {
			m.Options.ContentFormat = XML
		} else if contentType == "application/octet-stream" {
			m.Options.ContentFormat = OctetStream
		} else if contentType == "application/exi" {
			m.Options.ContentFormat = EXI
		} else if contentType == "application/json" {
			m.Options.ContentFormat = JSON
		} else {
			return errors.New("Bad Content Format Provided")
		}
		return nil
	}

}

func WithURI(URI string) MessagesConfig {
	return func(m *Message) error {
		return m.Options.SetURI(URI)
	}
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
