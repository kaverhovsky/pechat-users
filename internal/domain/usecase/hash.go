package usecase

import (
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptPassHashCost = 12 // ~250ms
)

// TODO test
func hashPassword(password string) (string, error) {
	passwordBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcryptPassHashCost)

	return string(passwordBytes), err
}
