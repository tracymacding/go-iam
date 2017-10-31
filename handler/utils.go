package handler

import (
	"errors"
	"strings"
	"bytes"
	"io"
)

// parse bucket name from host:'bucketName.s3.amazonaws.com'
func ParseBucketFromHost(s3Host string) (string, error) {
	parts := strings.Split(s3Host, ".")
	if len(parts) < 1 {
		return "", errors.New("invalid bucket name")
	}
	return parts[0], nil
}

func StreamToByte(stream io.Reader) []byte {
        buf := new(bytes.Buffer)
        buf.ReadFrom(stream)
        return buf.Bytes()
}
