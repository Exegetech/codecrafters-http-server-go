package main

import (
	"log"
	"net"
	"os"
	"strconv"
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
	statusLine := getStatusLine(req)
	path := getPath(statusLine)

	if path == "/" {
		msg := []byte("HTTP/1.1 200 OK\r\n\r\n")
		conn.Write(msg)
		conn.Close()
		return
	}

	if strings.HasPrefix(path, "/echo") {
		arr := strings.Split(path, "/")
		needsEcho := arr[2]

		msg := []string{
			"HTTP/1.1 200 OK\r\n",
			"Content-Type: text/plain\r\n",
			"Content-Length: " + strconv.Itoa(len(needsEcho)) + "\r\n",
			"\r\n",
			needsEcho,
		}

		join := strings.Join(msg, "")
		conn.Write([]byte(join))
		conn.Close()
	}

	msg := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	conn.Write(msg)
	conn.Close()
}

func getPath(str string) string {
	arr := strings.Split(str, " ")
	return arr[1]
}

func getStatusLine(str string) string {
	arr := strings.Split(str, "\r\n")
	return arr[0]
}
