package auth

import (
	"chores-suck/web/messages"
	"regexp"
	"strings"
)

func validateInput(username string, p1 string, p2 string, email string, msg *messages.RegisterMessage) bool {
	valid := true

	regEmail := regexp.MustCompile(`.+@.+\..+`)
	regName := regexp.MustCompile(`^(?:[0-9a-zA-Z]+-)*[0-9a-zA-Z]+$`)

	if strings.TrimSpace(username) == "" {
		msg.Username = "Must enter a valid username"
		valid = false
	}

	if p1 != p2 {
		msg.Password = "Passwords don't match"
		valid = false
	}

	if strings.TrimSpace(p1) == "" {
		msg.Password = "Must enter a password"
		valid = false
	}

	if strings.TrimSpace(email) == "" {
		msg.Email = "Must enter an email address"
		valid = false
	}

	if !regName.MatchString(username) {
		msg.Username = "Username must only consist of alphanumeric characters and hyphens and cannot start or end with a hyphen"
		valid = false
	}

	if !regEmail.MatchString(email) {
		msg.Email = "Must enter a valid email address"
		valid = false
	}

	return valid
}
