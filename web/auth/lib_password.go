package auth

import (
	"golang.org/x/crypto/bcrypt"
)

// PasswordHashCost - to set global, for more speed or security
var PasswordHashCost = bcrypt.DefaultCost

// HashPassword  - create new hash of password
func HashPassword(password string) (string, error) {
	p, err := bcrypt.GenerateFromPassword([]byte(password), PasswordHashCost)
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
