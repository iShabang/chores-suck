package web

import (
	"errors"
	"regexp"
	"strings"
)

func validateGroupName(name string) error {
	regName := regexp.MustCompile(`^(?:[0-9a-zA-Z]+-)*[0-9a-zA-Z]+$`)

	if strings.TrimSpace(name) == "" {
		return errors.New("Group name cannot be empty")

	} else if !regName.MatchString(name) {
		return errors.New("Name must only consist of alphanumeric characters and hyphens and cannot start or end with a hyphen")
	}
	return nil
}

func validatePassword(password string, confirm string) error {
	if password != confirm {
		return errors.New("Passwords don't match")
	}
	if strings.TrimSpace(password) == "" {
		return errors.New("Password cannot be empty")
	}
	return nil
}

func validateUsername(username string) error {
	regName := regexp.MustCompile(`^(?:[0-9a-zA-Z]+-)*[0-9a-zA-Z]+$`)
	if strings.TrimSpace(username) == "" {
		return errors.New("Username cannot be empty")
	}
	if !regName.MatchString(username) {
		return errors.New("Username must only consist of alphanumeric characters and hyphens and cannot start or end with a hyphen")
	}
	return nil
}

func validateEmail(email string) error {
	regEmail := regexp.MustCompile(`.+@.+\..+`)
	if strings.TrimSpace(email) == "" {
		return errors.New("Email cannot be empty")
	}
	if !regEmail.MatchString(email) {
		return errors.New("Invalid email address")
	}
	return nil
}
