package main

import (
	"strings"
)

type request struct {
	method  string
	path    string
	headers map[string]string
	body    string
}

func parseRequest(req []byte) request {
	tmp := strings.Split(string(req), "\r\n")
	status := tmp[0]
	headers := tmp[1 : len(tmp)-2]
	body := tmp[len(tmp)-1]

	method, path := parseMethodAndPath(status)

	return request{
		method:  method,
		path:    path,
		headers: parseHeaders(headers),
		body:    body,
	}
}

func parseMethodAndPath(status string) (string, string) {
	arr := strings.Split(status, " ")
	return arr[0], arr[1]
}

func parseHeaders(headers []string) map[string]string {
	m := make(map[string]string)
	for _, header := range headers {
		tmp := strings.Split(header, ": ")
		m[tmp[0]] = tmp[1]
	}
	return m
}
