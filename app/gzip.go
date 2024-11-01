package main

import (
	"bytes"
	"compress/gzip"
)

func gzipString(data string) (string, error) {
	var buf bytes.Buffer
	gzipWriter := gzip.NewWriter(&buf)

	_, err := gzipWriter.Write([]byte(data))
	if err != nil {
		return "", err
	}

	err = gzipWriter.Close()
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
