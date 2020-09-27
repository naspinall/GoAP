package main

import (
	"log"

	"github.com/naspinall/GoAP/pkg/client"
)

func main() {
	c, err := client.NewClient("localhost", 5688)
	if err != nil {
		log.Fatal(err)
	}

	m, err := c.Get("coap://localhost/a")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("%+v", m)

	log.Println(string(m.Payload))
}
