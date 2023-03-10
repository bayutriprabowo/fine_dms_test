package utils

import (
	"regexp"
	"strings"
)

func IsValidTag(tagName string) bool {
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(tagName) {
		return false
	}

	words := strings.Fields(tagName)
	if len(words) > 16 {
		return false
	}

	return true
}
