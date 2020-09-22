package client

import (
	"log"

	messages "github.com/naspinall/GoAP/pkg/message"
)

func (c *Client) Get() (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}
	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token))
	c.setupSession(10, 10)
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}

func (c *Client) Post() (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}
	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token))
	c.setupSession(10, 10)
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}

func (c *Client) Put() (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}
	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token))
	c.setupSession(10, 10)
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}
