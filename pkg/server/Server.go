package server

import (
	"log"
	"net"

	messages "github.com/naspinall/GoAP/pkg/message"
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

func Ping() {

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
		log.Println("Message Read")
		m, err := messages.FromBytes(b)
		log.Printf("%+v", b)
		if err != nil {
			log.Fatal(err)
		}

		if m.Type == messages.Confirmable {
			// Encoding ACK
			err := m.AsAcknowledge().Encode()
			if err != nil {
				log.Fatal(err)
			}

			// Sending ACK
			conn.WriteTo(m.Bytes(), raddr)
		}
	}
}
