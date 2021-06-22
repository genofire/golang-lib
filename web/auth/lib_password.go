package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// HashPassword  - create new hash of password
func HashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(p), nil
}

// ValidatePassword - check if given password is equal to saved hash
func ValidatePassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
