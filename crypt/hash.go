package crypt

import (
	"golang.org/x/crypto/bcrypt"
)

// GenerateSaltedPassword will return a hashed password
func GenerateSaltedPassword(plainPassword string) (saltedPass string, err error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(plainPassword), 10)
	return string(bytes), err
}

// CheckPasswordHash will valid if hash matches a given plaintext password
func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
