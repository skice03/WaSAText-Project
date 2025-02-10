package utils

import (
	"regexp"
)

func ValidUsername(username string) bool {
	isValid := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username)
	if isValid && (len(username) >= 3 && len(username) <= 16) {
		return true
	}
	return false
}
