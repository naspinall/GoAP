package server

import (
	"log"
	"net"
)

// Echo UDP Server
func Echo() {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		IP:   []byte{0, 0, 0, 0},
		Port: 5000,
		Zone: "",
	})

	if err != nil {
		log.Fatal(err)
	}

	for {
		b := make([]byte, 1024)
		_, raddr, err := conn.ReadFrom(b)

		if err != nil {
			log.Fatal(err)
		}
		_, err = conn.WriteTo(b, raddr)
	}
}
