package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"strings"
)

func handleIndex(conn net.Conn) {
	msg := []byte("HTTP/1.1 200 OK\r\n\r\n")
	conn.Write(msg)
	conn.Close()
}

func handleUserAgent(conn net.Conn, req request) {
	userAgent, ok := req.headers["User-Agent"]
	if !ok {
		log.Println("Error getting User-Agent")
	}

	msg := []string{
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: text/plain\r\n",
		"Content-Length: " + strconv.Itoa(len(userAgent)) + "\r\n",
		"\r\n",
		userAgent,
	}

	join := strings.Join(msg, "")
	conn.Write([]byte(join))
	conn.Close()
}

func handleEcho(conn net.Conn, req request) {
	arr := strings.Split(req.path, "/")
	value := arr[2]

	acceptEncoding, ok := req.headers["Accept-Encoding"]
	if !ok || !strings.Contains(acceptEncoding, "gzip") {
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
		return
	}

	compressed, err := gzipString(value)
	if err != nil {
		log.Println("Error compressing value")
	}

	msg := []string{
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: text/plain\r\n",
		"Content-Encoding: gzip\r\n",
		"Content-Length: " + strconv.Itoa(len(compressed)) + "\r\n",
		"\r\n",
		compressed,
	}

	join := strings.Join(msg, "")
	conn.Write([]byte(join))
	conn.Close()

}

func handleGetFile(conn net.Conn, req request, tmpdir *string) {
	arr := strings.Split(req.path, "/")
	filename := arr[2]

	content, err := os.ReadFile(*tmpdir + filename)
	if err != nil {
		handleNotFound(conn)
		return
	}

	msg := []string{
		"HTTP/1.1 200 OK\r\n",
		"Content-Type: application/octet-stream\r\n",
		"Content-Length: " + strconv.Itoa(len(content)) + "\r\n",
		"\r\n",
		string(content),
	}

	join := strings.Join(msg, "")
	conn.Write([]byte(join))
	conn.Close()
}

func handleWriteFile(conn net.Conn, req request, tmpdir *string) {
	arr := strings.Split(req.path, "/")
	filename := arr[2]

	cleaned := strings.Trim(req.body, "\x00")
	err := os.WriteFile(*tmpdir+filename, []byte(cleaned), 0644)
	if err != nil {
		log.Println("Error writing file")
	}

	msg := []string{
		"HTTP/1.1 201 Created\r\n",
		"\r\n",
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
