package utils

import (
	"regexp"
	"strings"
)

// ValidateEmail - Simple email validation
func ValidateEmail(email string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return re.MatchString(email)
}

// SanitizeString - Remove extra spaces
func SanitizeString(s string) string {
	return strings.TrimSpace(s)
}