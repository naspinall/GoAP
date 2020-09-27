package client

import (
	"log"

	messages "github.com/naspinall/GoAP/pkg/message"
)

func (c *Client) Get(URI string) (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}

	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token), messages.WithURI(URI))
	log.Printf("%+v", m)
	log.Printf("%+v", m.Options)
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}

func (c *Client) Post(URI string) (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}

	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token), messages.WithURI(URI))
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}

func (c *Client) Put(URI string) (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}

	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token), messages.WithURI(URI))
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}

func (c *Client) Delete(URI string) (*messages.Message, error) {
	messageID, token, err := c.randomIDs()
	if err != nil {
		return nil, err
	}

	m := messages.NewMessage(messages.Get(), messages.WithMessageID(messageID), messages.WithToken(token), messages.WithURI(URI))
	m, err = c.Do(m)
	if err != nil {
		log.Fatal(err)
	}

	return m, nil
}
