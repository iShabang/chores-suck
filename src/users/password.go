package users

import (
	"golang.org/x/crypto/bcrypt"
)

func checkpword(plain string, hashed string) bool {
	e := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(plain))

	return e == nil
}
