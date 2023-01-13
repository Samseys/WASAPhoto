package utils

import (
	"regexp"
	"strconv"
	"strings"
)

func CheckName(name string) bool {
	if len(name) < 1 || len(name) > 16 {
		return false
	}

	match, err := regexp.MatchString("^.*?$", name)
	if err != nil {
		return false
	}

	return match
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
