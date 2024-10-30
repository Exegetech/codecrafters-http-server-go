package main

import (
	"log"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalln("Failed to bind to port 4221")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln("Error accepting connection: ", err.Error())
		}

		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	buf := make([]byte, 1024)
	conn.Read(buf)

	req := string(buf)

	msg := []byte("HTTP/1.1 200 OK\r\n\r\n")
	if !strings.HasPrefix(req, "GET / HTTP/1.1") {
		msg = []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	}

	conn.Write(msg)

	conn.Close()
}
