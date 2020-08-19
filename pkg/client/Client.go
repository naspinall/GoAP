package client

import (
	"crypto/rand"
	"encoding/binary"
	"log"
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
	channels map[uint64]chan *messages.Message
}

func NewClient(address string, port int) (*Client, error) {
	ips, err := net.LookupIP(address)
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, &net.UDPAddr{
		IP:   ips[0],
		Port: port,
	})

	if err != nil {
		return nil, err
	}

	c := &Client{
		ips:      ips,
		conn:     conn,
		channels: make(map[uint64]chan *messages.Message),
	}

	// Listener for responses
	go c.listen()

	return c, nil
}

func (c *Client) listen() {
	for {

		// Reading
		b := make([]byte, 1024)
		if _, err := c.conn.Read(b); err != nil {
			log.Fatal(err)
		}

		m, err := messages.FromBytes(b)
		if err != nil {
			log.Fatal(err)
		}

		token, _ := binary.Uvarint(m.Token)
		c, ok := c.channels[token]
		if !ok {
			continue
		}

		// Send the value into the channel
		c <- m
	}
}

func (c *Client) generateToken() ([]byte, chan *messages.Message, error) {
	// Generating the random token
	tokenBytes := make([]byte, 8)
	if _, err := rand.Read(tokenBytes); err != nil {
		return nil, nil, err
	}

	token, _ := binary.Uvarint(tokenBytes)
	_, ok := c.channels[token]
	// Token already exists try again
	if ok {
		return c.generateToken()
	}

	// Creating a new channel
	c.channels[token] = make(chan *messages.Message)

	return tokenBytes, c.channels[token], nil
}

func (c *Client) Get() (*messages.Message, error) {
	token, mc, err := c.generateToken()
	if err != nil {
		return nil, err
	}

	m := messages.NewMessage(messages.Get(), messages.WithToken(token))
	m.Write(c.conn)

	// Wait for response
	resp := <-mc

	return resp, nil
}
