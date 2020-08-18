package messages

import (
	"bytes"
	"reflect"
	"testing"
)

func TestMessage_DecodeHeader(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name            string
		fields          fields
		wantErr         bool
		wantVersion     uint8
		wantType        uint8
		wantTokenLength uint8
		wantCode        uint8
		wantMessageID   uint16
	}{
		{
			name: "Version and Type, No Header",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xF5, 0x11, 0x11, 0x11}),
			},
			wantErr:         false,
			wantVersion:     0x1,
			wantType:        0x1,
			wantTokenLength: 0xF,
			wantCode:        0x11,
			wantMessageID:   0x1111,
		},
		{
			name: "Version and Type, No Header",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xA6, 0x11, 0x22, 0x22}),
			},
			wantErr:         false,
			wantVersion:     0x2,
			wantType:        0x1,
			wantTokenLength: 0xA,
			wantCode:        0x11,
			wantMessageID:   0x2222,
		},
		{
			name: "Version and Type, No Header",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xC9, 0x11, 0x22, 0x11}),
			},
			wantErr:         false,
			wantVersion:     0x1,
			wantType:        0x2,
			wantTokenLength: 0xC,
			wantCode:        0x22,
			wantMessageID:   0x1111,
		},
		{
			name: "Empty Buffer",
			fields: fields{
				buff: &bytes.Buffer{},
			},
			wantErr:         true,
			wantVersion:     0x0,
			wantType:        0x0,
			wantTokenLength: 0x0,
			wantCode:        0x0,
			wantMessageID:   0x0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.DecodeHeader(); (err != nil) != tt.wantErr {
				t.Errorf("Message.DecodeHeader() error = %v, wantErr %v", err, tt.wantErr)
			}

			if m.Version != tt.wantVersion {
				t.Errorf("Message.DecodeHeader() version = %v, wantVersion %v", m.Version, tt.wantVersion)
			}

			if m.Type != tt.wantType {
				t.Errorf("Message.DecodeHeader() type = %v, wantType %v", m.Type, tt.wantType)
			}

			if m.TokenLength != tt.wantTokenLength {
				t.Errorf("Message.DecodeHeader() tokenlength = %v, wantTokenLength %v", m.TokenLength, tt.wantTokenLength)
			}

		})
	}
}

func TestMessage_DecodeToken(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name      string
		fields    fields
		wantErr   bool
		wantToken []byte
	}{
		{
			name: "No Token",
			fields: fields{
				TokenLength: 0,
			},
			wantErr: false,
		},
		{
			name: "Valid Token",
			fields: fields{
				TokenLength: 2,
				buff:        bytes.NewBuffer([]byte{0x01, 0x01}),
			},
			wantToken: []byte{0x01, 0x01},
		},
		{
			name: "Invalid Token Length",
			fields: fields{
				TokenLength: 3,
				buff:        bytes.NewBuffer([]byte{0x01, 0x01}),
			},
			wantToken: []byte{},
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.DecodeToken(); (err != nil) != tt.wantErr {
				t.Errorf("Message.DecodeToken() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(m.Token, tt.wantToken) {
				t.Errorf("Message.DecodeToken() Token = %v, wantToken %v", m.Token, tt.wantToken)
			}
		})
	}
}

func TestMessage_DecodeOptions(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantOptions []Option
	}{
		{
			name: "No Options",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xFF}),
			},
		},
		{
			name: "An Option Length Check",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0x21, 0x11, 0x11, 0xFF}),
			},
			wantOptions: []Option{
				Option{
					OptionNumber: 1,
					Length:       2,
					Value:        []byte{0x11, 0x11},
				},
			},
		},
		{
			name: "An Option Delta Check",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0x13, 0x11, 0xFF}),
			},
			wantOptions: []Option{
				Option{
					OptionNumber: 3,
					Length:       1,
					Value:        []byte{0x11},
				},
			},
		},
		{
			name: "Multiple Options",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0x22, 0x11, 0x11, 0x21, 0x11, 0x11, 0x25, 0x11, 0x11, 0xFF}),
			},
			wantOptions: []Option{
				Option{
					OptionNumber: 2,
					Length:       2,
					Value:        []byte{0x11, 0x11},
				},
				Option{
					OptionNumber: 3,
					Length:       2,
					Value:        []byte{0x11, 0x11},
				},
				Option{
					OptionNumber: 8,
					Length:       2,
					Value:        []byte{0x11, 0x11},
				},
			},
		},
		{
			name: "Large Option",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0x2D, 0x0E, 0x11, 0x11, 0xFF}),
			},
			wantOptions: []Option{
				Option{
					OptionNumber: 1,
					Length:       2,
					Value:        []byte{0x11, 0x11},
				},
			},
		},
		{
			name: "Large Length",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xD1, 0x0E, 0x11, 0xFF}),
			},
			wantOptions: []Option{
				Option{
					OptionNumber: 1,
					Length:       1,
					Value:        []byte{0x11},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.DecodeOptions(); (err != nil) != tt.wantErr {
				t.Errorf("Message.DecodeOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(m.Options, tt.wantOptions) {
				t.Errorf("Message.DecodeOptions() Options = %v, wantOptions %v", m.Options, tt.wantOptions)
			}

		})
	}
}

func TestMessage_DecodePayload(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantPayload []byte
	}{
		{
			name: "No Payload",
			fields: fields{
				buff: &bytes.Buffer{},
			},
		},
		{
			name: "A Payload",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xFF}),
			},
			wantPayload: []byte{0xFF},
		},
		{
			name: "A Bigger Payload",
			fields: fields{
				buff: bytes.NewBuffer([]byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF}),
			},
			wantPayload: []byte{0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.DecodePayload(); (err != nil) != tt.wantErr {
				t.Errorf("Message.DecodePayload() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(m.Payload, tt.wantPayload) {
				t.Errorf("Message.DecodePayload() Payload = %v, wantPayload %v", m.Payload, tt.wantPayload)
			}
		})
	}
}

func TestMessage_EncodePayload(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantPayload []byte
	}{
		{
			name: "No Payload",
			fields: fields{
				buff: &bytes.Buffer{},
			},
			wantErr: false,
		},
		{
			name: "A Payload",
			fields: fields{
				Payload: []byte{0xFF, 0xFF, 0xFF, 0xFF},
				buff:    &bytes.Buffer{},
			},
			wantPayload: []byte{0xFF, 0xFF, 0xFF, 0xFF},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.EncodePayload(); (err != nil) != tt.wantErr {
				t.Errorf("Message.EncodePayload() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(m.Payload, tt.wantPayload) {
				t.Errorf("Message.EncodePayload() Payload = %v, wantPayload %v", m.Payload, tt.wantPayload)
			}
		})
	}
}

func TestMessage_EncodeToken(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name            string
		fields          fields
		wantErr         bool
		wantToken       []byte
		wantTokenLength uint8
	}{
		{
			name: "No Token",
			fields: fields{
				buff: &bytes.Buffer{},
			},
		},
		{
			name: "A Token",
			fields: fields{
				buff:  &bytes.Buffer{},
				Token: []byte{0xFF, 0xFF, 0xFF},
			},
			wantToken:       []byte{0xFF, 0xFF, 0xFF},
			wantTokenLength: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.EncodeToken(); (err != nil) != tt.wantErr {
				t.Errorf("Message.EncodeToken() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !bytes.Equal(m.Token, tt.wantToken) {
				t.Errorf("Message.EncodeToken() Token = %v, wantToken %v", m.Token, tt.wantToken)
			}
			if tt.wantTokenLength != m.TokenLength {
				t.Errorf("Message.EncodeToken() TokenLength = %v, wantTokenLength %v", m.TokenLength, tt.wantTokenLength)
			}
		})
	}
}

func TestMessage_EncodeOptions(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name        string
		fields      fields
		wantErr     bool
		wantOptions []byte
	}{
		{
			name: "No Options",
			fields: fields{
				buff: &bytes.Buffer{},
			},
		},
		{
			name: "An Option",
			fields: fields{
				buff: &bytes.Buffer{},
				Options: []Option{
					Option{
						OptionNumber: 1,
						Length:       2,
						Value:        []byte{0x01, 0x01},
					},
				}},
			wantOptions: []byte{0x21, 0x01, 0x01},
		},
		{
			name: "An Option",
			fields: fields{
				buff: &bytes.Buffer{},
				Options: []Option{
					Option{
						OptionNumber: 1,
						Length:       4,
						Value:        []byte{0x01, 0x01, 0x01, 0x01},
					},
				}},
			wantOptions: []byte{0x41, 0x01, 0x01, 0x01, 0x01},
		},
		{
			name: "Multiple Options",
			fields: fields{
				buff: &bytes.Buffer{},
				Options: []Option{
					Option{
						OptionNumber: 1,
						Length:       4,
						Value:        []byte{0x01, 0x01, 0x01, 0x01},
					},
					Option{
						OptionNumber: 1,
						Length:       4,
						Value:        []byte{0x01, 0x01, 0x01, 0x01},
					},
				}},
			wantOptions: []byte{0x41, 0x01, 0x01, 0x01, 0x01, 0x41, 0x01, 0x01, 0x01, 0x01},
		},
		{
			name: "An Option",
			fields: fields{
				buff: &bytes.Buffer{},
				Options: []Option{
					Option{
						OptionNumber: 13,
						Length:       4,
						Value:        []byte{0x01, 0x01, 0x01, 0x01},
					},
				}},
			wantOptions: []byte{0x4E, 0x01, 0x01, 0x01, 0x01},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.EncodeOptions(); (err != nil) != tt.wantErr {
				t.Errorf("Message.EncodeOptions() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !bytes.Equal(m.buff.Bytes(), tt.wantOptions) {
				t.Errorf("Message.EncodeToken() Written Options = %v, wantOptions %v", m.buff.Bytes(), tt.wantOptions)
			}
		})
	}
}

func TestMessage_EncodeHeader(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name       string
		fields     fields
		wantErr    bool
		wantHeader []byte
	}{
		{
			name: "Header Value",
			fields: fields{
				buff:        &bytes.Buffer{},
				Version:     1,
				Type:        1,
				TokenLength: 1,
				Code:        1,
				MessageID:   1,
			},
			wantHeader: []byte{0x15, 0x01, 0x01, 0x00},
		},
		{
			name: "Header Value",
			fields: fields{
				buff:        &bytes.Buffer{},
				Version:     2,
				Type:        3,
				TokenLength: 5,
				Code:        8,
				MessageID:   12,
			},
			wantHeader: []byte{0x5E, 0x08, 0x0C, 0x00},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.EncodeHeader(); (err != nil) != tt.wantErr {
				t.Errorf("Message.EncodeHeader() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !bytes.Equal(m.buff.Bytes(), tt.wantHeader) {
				t.Errorf("Message.EncodeHeader() written Header = %v, wantHeader %v", m.buff.Bytes(), tt.wantHeader)
			}
		})
	}
}

func TestMessage_Encode(t *testing.T) {
	type fields struct {
		Version     uint8
		Type        uint8
		TokenLength uint8
		Code        uint8
		MessageID   uint16
		Token       []byte
		Options     []Option
		Payload     []byte
		buff        *bytes.Buffer
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			fields: fields{
				buff:        &bytes.Buffer{},
				Version:     1,
				Type:        3,
				TokenLength: 2,
				Code:        4,
				MessageID:   3,
				Token:       []byte{0x01, 0x01},
				Options: []Option{
					Option{
						OptionNumber: 1,
						Length:       2,
						Value:        []byte{0x01, 0x01},
					},
					Option{
						OptionNumber: 1,
						Length:       4,
						Value:        []byte{0x02, 0x02, 0x02, 0x02},
					},
				},
				Payload: []byte{0xFF, 0xFF},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Message{
				Version:     tt.fields.Version,
				Type:        tt.fields.Type,
				TokenLength: tt.fields.TokenLength,
				Code:        tt.fields.Code,
				MessageID:   tt.fields.MessageID,
				Token:       tt.fields.Token,
				Options:     tt.fields.Options,
				Payload:     tt.fields.Payload,
				buff:        tt.fields.buff,
			}
			if err := m.Encode(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Encode() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Clearing options
			m.Options = []Option{}

			if err := m.Decode(); (err != nil) != tt.wantErr {
				t.Errorf("Message.Decode() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !reflect.DeepEqual(m.Version, tt.fields.Version) {
				t.Errorf("Message.Decode() Version = %v, wantVersion = %v", m.Version, tt.fields.Version)
			}
			if !reflect.DeepEqual(m.Type, tt.fields.Type) {
				t.Errorf("Message.Decode() Type = %v, wantType = %v", m.Type, tt.fields.Type)
			}
			if !reflect.DeepEqual(m.TokenLength, tt.fields.TokenLength) {
				t.Errorf("Message.Decode() TokenLength = %v, wantTokenLength = %v", m.TokenLength, tt.fields.TokenLength)
			}
			if !reflect.DeepEqual(m.Code, tt.fields.Code) {
				t.Errorf("Message.Decode() Code = %v, wantCode = %v", m.Code, tt.fields.Code)
			}
			if !reflect.DeepEqual(m.MessageID, tt.fields.MessageID) {
				t.Errorf("Message.Decode() MessageID = %v, wantMessageID = %v", m.MessageID, tt.fields.MessageID)
			}
			if !reflect.DeepEqual(m.Token, tt.fields.Token) {
				t.Errorf("Message.Decode() Token = %v, wantToken = %v", m.Token, tt.fields.Token)
			}
			if !reflect.DeepEqual(m.Options, tt.fields.Options) {
				t.Errorf("Message.Decode() Options = %v, wantOptions = %v", m.Options, tt.fields.Options)
			}
			if !reflect.DeepEqual(m.Payload, tt.fields.Payload) {
				t.Errorf("Message.Decode() Payload = %v, wantPayload = %v", m.Payload, tt.fields.Payload)
			}
		})
	}
}
