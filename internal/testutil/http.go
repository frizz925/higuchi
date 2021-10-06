package testutil

import "strings"

func LinesToRawRequestHeader(lines ...string) []byte {
	s := strings.Join(lines, "\r\n")
	s += "\r\n\r\n"
	return []byte(s)
}

func LinesToRawRequest(lines ...string) []byte {
	s := strings.Join(lines, "\r\n")
	return []byte(s)
}
