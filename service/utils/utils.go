package utils

import (
	"strconv"
	"strings"
)

func MakeAlphaNumeric(s string) string {
	var result strings.Builder
	for i := 0; i < len(s); i++ {
		b := s[i]
		if ('a' <= b && b <= 'z') ||
			('A' <= b && b <= 'Z') ||
			('0' <= b && b <= '9') ||
			b == ' ' {
			result.WriteByte(b)
		}
	}
	return result.String()
}

// Returns an id >= 1, 0 if there is an error
func GetAuthorizationID(authHeader string) uint64 {
	splitted := strings.SplitN(authHeader, " ", 2)
	if len(splitted) != 2 {
		return 0
	}

	id, err := strconv.ParseUint(splitted[1], 10, 64)
	if err != nil {
		return 0
	}
	return id
}
