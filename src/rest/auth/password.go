package auth

import (
	"golang.org/x/crypto/bcrypt"
)

func checkpword(plain string, hashed string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))

	return e == nil
}

func hashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashed), err
}
