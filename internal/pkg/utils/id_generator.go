package utils

import (
	"crypto/rand"
	"io"
)

var idTable = []byte{
	'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm',
	'n', 'o', 'p', 'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z',
	'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M',
	'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
	'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}

// GenerateRandomString generates a secure random string of given length using alphanumeric characters
func GenerateRandomString(length int) string {
	b := make([]byte, length)
	n, err := io.ReadAtLeast(rand.Reader, b, length)
	if n != length || err != nil {
		return "0Aj3kL"[:length] // Fallback
	}
	for i := 0; i < len(b); i++ {
		b[i] = idTable[int(b[i])%len(idTable)]
	}
	return string(b)
}

// GenerateUserID generates a custom ID for users with prefix USR- (e.g., USR-1a2B3c)
func GenerateUserID() string {
	return "USR-" + GenerateRandomString(6)
}

// GenerateDeviceID generates a custom ID for devices with prefix DEVC- (e.g., DEVC-1a2B3c)
func GenerateDeviceID() string {
	return "DEVC-" + GenerateRandomString(6)
}

// GenerateSessionID generates a custom ID for sessions with prefix SESS- (e.g., SESS-1a2B3c4D)
func GenerateSessionID() string {
	return "SESS-" + GenerateRandomString(8)
}
