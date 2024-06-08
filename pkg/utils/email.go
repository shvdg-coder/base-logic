package utils

import "net/mail"

// IsValidEmail checks if the given email string is a valid email address.
func IsValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return email != "" && err == nil
}
