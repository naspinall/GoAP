package client

import (
	"bytes"
	"crypto/rand"
	"net"

	messages "github.com/naspinall/GoAP/pkg/message"
)

type TokenChannel struct {
	Token   []byte
	Channel chan *messages.Message
}

type Client struct {
	ips      []net.IP
	conn     *net.UDPConn
	channels []TokenChannel
}

func NewClient(address string) (*Client, error) {
	ips, err := net.LookupIP(address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP: ips[0],
	})

	if err != nil {
		return nil, err
	}

	return &Client{
		ips:  ips,
		conn: conn,
	}, nil
}

func (c *Client) Listen() {

}

func (c *Client) generateToken() ([]byte, error) {
	// Generating the random token
	token := make([]byte, 8)
	if _, err := rand.Read(token); err != nil {
		return nil, err
	}

	for _, channel := range c.channels {
		if bytes.Equal(token, channel.Token) {
			return c.generateToken()
		}
	}
	return token, nil

}

func (c *Client) Get() error {
	token, err := c.generateToken()
	if err != nil {
		return err
	}
	m := messages.NewMessage(messages.Get(), messages.WithToken(token))
	m.Write(c.conn)
	return nil
}
