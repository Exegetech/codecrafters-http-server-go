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
	statusLine, headers := separateRequest(req)
	path := getPath(statusLine)

	if path == "/" {
		handleIndex(conn)
		return
	}

	if path == "/user-agent" {
		handleUserAgent(conn, headers)
		return
	}

	if strings.HasPrefix(path, "/echo") {
		handleEcho(conn, path)
		return
	}

	handleNotFound(conn)
}

func handleIndex(conn net.Conn) {
	msg := []byte("HTTP/1.1 200 OK\r\n\r\n")
	conn.Write(msg)
	conn.Close()
}

func handleEcho(conn net.Conn, path string) {
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

func handleUserAgent(conn net.Conn, headers []string) {
	var userAgentHeader string
	for _, header := range headers {
		if strings.HasPrefix(header, "User-Agent") {
			userAgentHeader = header
			break
		}
	}

	value := strings.Split(userAgentHeader, ": ")[1]

	msg := []string{
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: text/plain\r\n",
		"Content-Length: " + strconv.Itoa(len(value)) + "\r\n",
		"\r\n",
		value,
	}

	join := strings.Join(msg, "")
	conn.Write([]byte(join))
	conn.Close()
}

func handleNotFound(conn net.Conn) {
	msg := []byte("HTTP/1.1 404 Not Found\r\n\r\n")
	conn.Write(msg)
	conn.Close()
}

func separateRequest(req string) (string, []string) {
	arr := strings.Split(req, "\r\n")
	statusLine := arr[0]
	headers := arr[1:]

	return statusLine, headers
}

func getPath(str string) string {
	arr := strings.Split(str, " ")
	return arr[1]
}
