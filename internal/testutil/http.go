package testutil

import "strings"

func LinesToRawHeader(lines ...string) []byte {
	s := strings.Join(lines, "\r\n")
	s += "\r\n\r\n"
	return []byte(s)
}

func LinesToRawPayload(lines ...string) []byte {
	s := strings.Join(lines, "\r\n")
	return []byte(s)
}
