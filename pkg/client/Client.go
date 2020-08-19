package client

import (
	"crypto/rand"
	"encoding/binary"
	"errors"
	"log"
	"net"
	"time"

	messages "github.com/naspinall/GoAP/pkg/message"
)

const (
	MaxTransmitSpan  = 45  // Maximum time from first transmission to it's last retransmission
	MaxTransmitWait  = 93  // Maximum time from first transmission to giving up on recieving an acknowledgement
	MaxLatency       = 100 // Maximum time a datagram is expected to take from send to recieve
	ProcessingDelay  = 2   // Time it takes to send an acknowledgement
	MaxRtt           = 202 // Maximum Round Trip Time
	ExchangeLifetime = 247 // Time for sending to
	NonLifetime      = 145

	AckTimeout      = 2   // Minmum spacing before retransmission
	AckRandomFactor = 1.5 // Random factor used to generate timeout
	MaxRetransmit   = 4   // Maxmimun number of times to do a retransmission
	Nstart          = 1   //
	DefaultLeisure  = 5
	ProbingRate     = 1
)

type TokenChannel struct {
	Token   []byte
	Channel chan *messages.Message
}

type MessageChannel struct {
	Message chan *messages.Message
	Error   chan error
}

type Client struct {
	ips             []net.IP
	conn            *net.UDPConn
	tokenChannels   map[uint64]chan *messages.Message
	messageChannels map[uint16]*MessageChannel
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
		ips:             ips,
		conn:            conn,
		tokenChannels:   make(map[uint64]chan *messages.Message),
		messageChannels: make(map[uint16]*MessageChannel),
	}

	// Listener for responses
	go c.listen()

	return c, nil
}

func (c *Client) listen() {
	for {

		// Read Message.
		b := make([]byte, 1024)
		if _, err := c.conn.Read(b); err != nil {
			log.Fatal(err)
		}

		// Decode Message.
		m, err := messages.FromBytes(b)
		if err != nil {
			log.Fatal(err)
		}

		// Get corresponding message channel.
		mc, ok := c.messageChannels[m.MessageID]
		if !ok {
			// No message ID, message could be a response
			tc, ok := c.tokenChannels[m.Token]
			if ok {
				// Sending response to handler
				tc <- m
				// Send acknowledgement message
				go c.sendAck(m)
			}

			// If no corresponding token, drop message
			continue
		}

		// Send message down corresponding message channel.
		mc.Message <- m

	}
}

func (c *Client) sendAck(m *messages.Message) {
	// Creating ACK Message
	ack := messages.NewMessage(messages.WithMessageID(m.MessageID), messages.WithType(messages.Acknowledgement))

	// Sending ACK
	ack.Write(c.conn)
}

func (c *Client) generateToken() (uint64, error) {
	// Generating the random token
	tokenBytes := make([]byte, 8)
	if _, err := rand.Read(tokenBytes); err != nil {
		return 0, err
	}
	token, _ := binary.Uvarint(tokenBytes)

	// Checking if already in use
	_, ok := c.tokenChannels[token]
	if ok {
		return c.generateToken()
	}

	return token, nil
}

func (c *Client) generateMessageID() (uint16, error) {
	// Generating random message id
	messageBytes := make([]byte, 2)
	if _, err := rand.Read(messageBytes); err != nil {
		return 0, err
	}
	// Converity to 16 bit integer
	bigMessageID, _ := binary.Uvarint(messageBytes)
	messageID := uint16(bigMessageID)

	// Checking if already in use
	_, ok := c.messageChannels[messageID]
	if ok {
		return c.generateMessageID()
	}

	return messageID, nil
}

func (c *Client) transmit(m *messages.Message) {
	// Number of retranmist attemps
	var retransmit int

	// Message timout
	timeout := AckTimeout * AckRandomFactor // TODO make this a random value.

	// Getting MessageID and Token
	messageID := m.MessageID
	token := m.Token

	// Getting Message Channel
	messageChannel := c.messageChannels[messageID]

	// Keep retransmitting until MaxRetransmit
	for retransmit <= MaxRetransmit {

		// Sending Message
		m.Write(c.conn)

		ticker := time.NewTicker(time.Duration(timeout) * time.Second)

		select {
		case m := <-messageChannel.Message:

			// Check message type
			switch m.Type {
			case messages.Acknowledgement:
			case messages.Reset:

			}

		case <-ticker.C:

			// Increase retransmit timmer
			retransmit++

			// Increase timeout
			timeout *= 2
		}
	}

	// Retransmit is done message has timed out.
	messageChannel.Error <- errors.New("Timeout")

	// Removing message and token channels
	delete(c.messageChannels, messageID)
	delete(c.tokenChannels, token)
}

func (c *Client) sendMessage() (*messages.Message, error) {

	// Generating Token
	token, err := c.generateToken()
	if err != nil {
		return nil, err
	}

	// Generating MessageID
	messageID, err := c.generateMessageID()
	if err != nil {
		return nil, err
	}

	// Creating Channels
	tc, mc, ec := make(chan *messages.Message), make(chan *messages.Message), make(chan error)

	// Adding Channels to client map
	c.tokenChannels[token], c.messageChannels[messageID] = tc, &MessageChannel{
		Error:   ec,
		Message: mc,
	}

	// Creating Message
	m := messages.NewMessage(messages.Get(), messages.WithToken(token), messages.WithMessageID(messageID))

	// Transmit the message to the server
	go c.transmit(m)

	// Wait for a response from the server.
	select {
	case <-mc:
		//ACK Recieved, now wait for response
		break
	case err := <-ec:
		// Error Recieved, return error
		return nil, err
	case resp := <-tc:
		// Piggybacked response
		return resp, nil
	}

	// Wait for response
	resp := <-tc

	return resp, nil
}
