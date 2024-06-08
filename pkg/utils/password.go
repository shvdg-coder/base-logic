package utils

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the password.
func HashPassword(plainPassword string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	return string(bytes), err
}

// ArePasswordsEqual compares the plain password with the hashed password.
func ArePasswordsEqual(plainPassword, hashedPassword string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(plainPassword))
	return err == nil
}
