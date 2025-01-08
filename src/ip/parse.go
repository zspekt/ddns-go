package ip

import (
	"errors"
	"net"
)

func ParseIP(b []byte) (string, error) {
	valid := []byte("1234567890.")

	parsed := parse(b, valid)

	ip := net.ParseIP(parsed)
	if ip == nil {
		return "", errors.New("error parsing IP")
	}

	return ip.String(), nil
}

func parse(input, valid []byte) string {
	buf := make([]byte, len(input))
	var a [256]bool

	for _, b := range valid {
		a[b] = true
	}

	var i int
	for _, b := range input {
		if a[b] {
			buf[i] = b
			i++
		}
	}
	return string(buf[:i])
}
