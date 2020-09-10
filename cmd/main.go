package main

import (
	"log"

	"github.com/naspinall/GoAP/pkg/client"
	"github.com/naspinall/GoAP/pkg/server"
)

func main() {
	go server.Ping()
	c, err := client.NewClient("localhost", 5000)
	if err != nil {
		log.Fatal(err)
	}

	c.Get()
}
