package main

import (
	"log"

	"github.com/naspinall/GoAP/pkg/client"
)

func main() {
	c, err := client.NewClient("coap.me", 5683)
	if err != nil {
		log.Fatal(err)
	}

	m, err := c.Get("coap://coap.me/test")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", m)
}
