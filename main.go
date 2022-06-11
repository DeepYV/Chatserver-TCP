package main

import (
	"log"
	"net"
)

func main() {
	s := Newserver()
	go s.run()
	listener, err := net.Listen("tcp", ":8888")
	if err != nil {
		log.Fatalf("unable to start server: %s", err.Error())
	}

	defer listener.Close()
	log.Printf("server started on :8888")

	for {

		conn, err := listener.Accept()
		if err != nil {
			log.Printf("unable to accept connection")
			continue
		}
		c := s.newClient(conn)
		go c.readinput()
	}

}
