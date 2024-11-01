package main

import (
	"flag"
	"log"
	"net"
	"os"
	"strings"
)

var _ = net.Listen
var _ = os.Exit

type method string

var (
	get  method = "GET"
	post method = "POST"
)

func main() {
	tmpdir := flag.String("directory", "/", "Directory to serve files from")
	flag.Parse()

	listener, err := net.Listen("tcp", "0.0.0.0:4221")
	if err != nil {
		log.Fatalln("Failed to bind to port 4221")
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection: ", err.Error())
		}

		go handleRequest(conn, tmpdir)
	}
}

func handleRequest(conn net.Conn, tmpdir *string) {
	buf := make([]byte, 1024)
	conn.Read(buf)

	req := parseRequest(buf)

	if req.method == string(get) && req.path == "/" {
		handleIndex(conn)
		return
	}

	if req.method == string(get) && req.path == "/user-agent" {
		handleUserAgent(conn, req)
		return
	}

	if req.method == string(get) && strings.HasPrefix(req.path, "/echo") {
		handleEcho(conn, req)
		return
	}

	if req.method == string(get) && strings.HasPrefix(req.path, "/files") {
		handleGetFile(conn, req, tmpdir)
		return
	}

	if req.method == string(post) && strings.HasPrefix(req.path, "/files") {
		handleWriteFile(conn, req, tmpdir)
		return
	}

	handleNotFound(conn)
}
