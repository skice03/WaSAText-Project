package utils

import (
	"regexp"
)

// Error handling
/*
var ErrUserDoesNotExist = errors.New("user doesn't exist")
var ErrUserNotFound = errors.New("user not found")
var ErrInternalServerError = errors.New("internal server error")
var ErrPermissionDenied = errors.New("permission denied")
var ErrCommentNotFound = errors.New("comment not found")

func ErrorTranslate(w http.ResponseWriter, err error) {

	if errors.Is(err, ErrUserDoesNotExist) {
		w.WriteHeader(http.StatusNotFound)
	} else if errors.Is(err, ErrUserNotFound) {
		w.WriteHeader(http.StatusForbidden)
	} else if errors.Is(err, ErrInternalServerError) {
		w.WriteHeader(http.StatusInternalServerError)
	} else if errors.Is(err, ErrPermissionDenied) {
		w.WriteHeader(http.StatusForbidden)
	} else if errors.Is(err, ErrCommentNotFound) {
		w.WriteHeader(http.StatusNotFound)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
*/

func ValidUsername(username string) bool {
	isValid := regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString(username)
	if isValid && (len(username) >= 3 && len(username) <= 16) {
		return true
	}
	return false
}
